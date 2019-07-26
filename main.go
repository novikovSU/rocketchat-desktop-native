package main

import (
	"errors"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
)

const appID = "com.github.novikovSU.rocketchat-desktop-native"

const (
	iistItem = iota
	nColumns
	keyEnter = 65293
)

var (
	GtkApplication gtk.Application
	GtkBuilder     gtk.Builder
	ctrlPressed    = false
	chatStore      *gtk.ListStore
)

func initList(list *gtk.TreeView) *gtk.ListStore {
	cellRenderer, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Fatal(err)
	}
	cellRenderer.SetProperty("wrap-mode", pango.WRAP_WORD_CHAR)
	cellRenderer.SetProperty("wrap-width", 530)
	cellRenderer.SetProperty("ypad", 5)
	cellRenderer.SetProperty("xpad", 3)

	column, err := gtk.TreeViewColumnNewWithAttribute("List Items", cellRenderer, "markup", 0)
	if err != nil {
		log.Fatal(err)
	}
	column.SetSizing(gtk.TREE_VIEW_COLUMN_AUTOSIZE)
	list.AppendColumn(column)

	store, err := gtk.ListStoreNew(glib.TYPE_STRING)
	if err != nil {
		log.Fatal(err)
	}
	list.SetModel(store)

	return store
}

func addToList(store *gtk.ListStore, text string) {
	iter := store.Append()
	store.SetValue(iter, 0, text)
}

func onChanged(selection *gtk.TreeSelection, label *gtk.Label) {
	var iter *gtk.TreeIter
	var model gtk.ITreeModel
	var ok bool
	model, iter, ok = selection.GetSelected()

	if ok {
		value, err := model.(*gtk.TreeModel).GetValue(iter, 0)
		if err != nil {
			log.Fatal(err)
		}

		text, err := value.GetString()
		if err != nil {
			log.Fatal(err)
		}

		label.SetText(text)
	}
}

func getSelectionText(selection *gtk.TreeSelection) (selectionText string) {
	var iter *gtk.TreeIter
	var model gtk.ITreeModel
	var ok bool
	model, iter, ok = selection.GetSelected()

	if ok {
		value, err := model.(*gtk.TreeModel).GetValue(iter, 0)
		if err != nil {
			log.Fatal(err)
		}

		text, err := value.GetString()
		if err != nil {
			log.Fatal(err)
		}

		selectionText = text
	}

	return
}

func isWindow(obj glib.IObject) (*gtk.Window, error) {
	// Make type assertion (as per gtk.go).
	if win, ok := obj.(*gtk.Window); ok {
		return win, nil
	}
	return nil, errors.New("not a *gtk.Window")
}

// onMainWindowDestory is the callback that is linked to the
// on_main_window_destroy handler. It is not required to map this,
// and is here to simply demo how to hook-up custom callbacks.
func onMainWindowDestroy() {
	log.Println("onMainWindowDestroy")
}

func main() {
	// Get application config
	config, _ = getConfig()

	// Get Rocket.Chat connection
	getConnection()

	// Init chat history
	allHistory = make(map[string]chatHistory)
	_ = getNewMessages(client)

	// Create a new application.
	app, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		log.Panic(err)
	}
	GtkApplication = *app

	// Connect function to application startup event, this is not required.
	app.Connect("startup", func() {
		log.Println("application startup")
	})

	// Connect function to application activate event
	app.Connect("activate", func() {
		GtkBuilder = *createGtkBuilder()
		// Get application config
		config, err = getConfig()
		if err == nil {
			// Get Rocket.Chat connection
			getConnection()
			openMainWindow(app)
		} else {
			OpenConnectionWindow()
		}
	})

	// Connect function to application shutdown event, this is not required.
	app.Connect("shutdown", func() {
		log.Println("application shutdown")
	})

	// Launch the application
	os.Exit(app.Run(os.Args))
}

func openMainWindow(app *gtk.Application) {
	log.Println("application activate")

	notif := glib.NotificationNew("Rocket.Chat Desktop native")
	notif.SetBody("application activate")
	app.SendNotification(appID, notif)

	win := CreateWindow("main_window")

	/* DISABLE custom header and menu */
	// Create menu
	// Get a headerbar
	obj, err := GtkBuilder.GetObject("main_header")
	if err != nil {
		log.Panic(err)
	}
	header, ok := obj.(*gtk.HeaderBar)
	if ok != true {
		log.Fatal("Could not create header bar:", err)
	}
	header.SetShowCloseButton(true)

	// Create a new menu button
	obj, err = GtkBuilder.GetObject("main_menu_button")
	if err != nil {
		log.Panic(err)
	}
	mbtn, ok := obj.(*gtk.MenuButton)
	if ok != true {
		log.Fatal("Could not create menu button:", err)
	}

	// Set up the menu model for the button
	menu := glib.MenuNew()
	if menu == nil {
		log.Fatal("Could not create menu (nil)")
	}
	// Actions with the prefix 'app' reference actions on the application
	// Actions with the prefix 'win' reference actions on the current window (specific to ApplicationWindow)
	// Other prefixes can be added to widgets via InsertActionGroup
	menu.Append("Connect", "custom.connect")
	menu.Append("Disconnect", "custom.disconnect")
	menu.Append("Quit", "app.quit")

	customActionGroup := glib.SimpleActionGroupNew()
	win.InsertActionGroup("custom", customActionGroup)

	// Create an action in the custom action group
	aConnect := glib.SimpleActionNew("connect", nil)
	aConnect.Connect("activate", func() {
		log.Println("CONNECTED")
		OpenConnectionWindow()
	})
	customActionGroup.AddAction(aConnect)
	app.AddAction(aConnect)

	aDisconnect := glib.SimpleActionNew("disconnect", nil)
	aDisconnect.Connect("activate", func() {
		log.Println("DISCONNECTED")
	})
	customActionGroup.AddAction(aDisconnect)
	app.AddAction(aDisconnect)

	mbtn.SetMenuModel(&menu.MenuModel)

	// add the menu button to the header
	header.PackStart(mbtn)

	// Assemble the window
	win.SetTitlebar(header)

	// Add Quit action to menu
	aQuit := glib.SimpleActionNew("quit", nil)
	aQuit.Connect("activate", func() {
		app.Quit()
	})
	app.AddAction(aQuit)

	// Add action for X-button
	win.Connect("destroy", app.Quit)

	// END creating menu

	obj, err = GtkBuilder.GetObject("contact_list")
	if err != nil {
		log.Fatal(err)
	}
	contactList, ok := obj.(*gtk.TreeView)
	if ok != true {
		log.Fatal(err)
	}

	// -------------

	obj, err = GtkBuilder.GetObject("chat_caption")
	if err != nil {
		log.Fatal(err)
	}
	chatCaption, ok := obj.(*gtk.Label)
	if ok != true {
		log.Fatal(err)
	}

	obj, err = GtkBuilder.GetObject("chat_list")
	if err != nil {
		log.Fatal(err)
	}
	chatList, ok := obj.(*gtk.TreeView)
	if ok != true {
		log.Fatal(err)
	}

	// Autoscroll of chatList
	obj, err = GtkBuilder.GetObject("right_scrolled_window")
	if err != nil {
		log.Fatal(err)
	}
	rightScrolledWindow, ok := obj.(*gtk.ScrolledWindow)
	if ok != true {
		log.Fatal(err)
	}
	chatList.Connect("size-allocate", func() {
		adj := rightScrolledWindow.GetVAdjustment()
		adj.SetValue(adj.GetUpper() - adj.GetPageSize())
	})
	chatList.ConnectAfter("size-allocate", func() {

	})

	obj, err = GtkBuilder.GetObject("text_input")
	if err != nil {
		log.Fatal(err)
	}
	textInput, ok := obj.(*gtk.TextView)
	if ok != true {
		log.Fatal(err)
	}

	// ------------------

	contactsStore := initList(contactList)
	fillContactList(contactsStore)

	chatStore = initList(chatList)

	contactsSelection, err := contactList.GetSelection()
	if err != nil {
		log.Fatal(err)
	}

	contactsSelection.Connect("changed", func() {
		//onChanged(contactsSelection, chatCaption)
		selectionText := getSelectionText(contactsSelection)
		//header.SetSubtitle(selectionText)
		chatCaption.SetText(selectionText)
		fillChat(chatStore, selectionText)
	})

	textInput.Connect("key-press-event", func(tv *gtk.TextView, ev *gdk.Event) {
		keyEvent := &gdk.EventKey{ev}
		switch keyEvent.KeyVal() {
		case gdk.KEY_Control_L, gdk.KEY_Control_R:
			ctrlPressed = true
		}
	})

	textInput.Connect("key-release-event", func(tv *gtk.TextView, ev *gdk.Event) {
		keyEvent := &gdk.EventKey{ev}
		switch keyEvent.KeyVal() {
		case keyEnter:
			if !ctrlPressed {
				buffer, err := tv.GetBuffer()
				if err != nil {
					log.Fatal("Unable to get buffer:", err)
				}
				start, end := buffer.GetBounds()

				inputText, err := buffer.GetText(start, end, true)
				if err != nil {
					log.Fatal("Unable to get text:", err)
				}
				inputText = strings.TrimSuffix(inputText, "\n")

				match, _ := regexp.MatchString("^(\\s*)$", inputText)
				if !match {
					selectionText := getSelectionText(contactsSelection)
					postByName(selectionText, inputText)
					buffer.SetText("")
				}
			}
			break
		case gdk.KEY_Control_L, gdk.KEY_Control_R:
			ctrlPressed = false
		default:
			//log.Printf("Keycode: %d\n", keyEvent.KeyVal())
		}
	})

	pullChan = subscribeToUpdates(client, 500)

	win.ShowAll()
	app.AddWindow(win)
}

// CreateWindow AAA
func CreateWindow(id string) *gtk.Window {
	obj, err := GtkBuilder.GetObject(id)
	if err != nil {
		log.Panic(err)
	}

	wnd, err := isWindow(obj)
	if err != nil {
		log.Panic(err)
	}

	// Create the action "wnd.close"
	wndCloseAction := glib.SimpleActionNew("close", nil)
	wndCloseAction.Connect("activate", func() {
		wnd.Close()
	})
	GtkApplication.AddAction(wndCloseAction)

	return wnd
}

func ttt(obj *glib.IObject) *gtk.Window {
	return (*obj).(*gtk.Window)
}

func createGtkBuilder() *gtk.Builder {
	// Get the GtkBuilder UI definition in the glade file.
	builder, err := gtk.BuilderNewFromFile("main.glade")
	if err != nil {
		log.Panic(err)
	}

	// Map the handlers to callback functions, and connect the signals to the Builder.
	signals := map[string]interface{}{
		"on_main_window_destroy": onMainWindowDestroy,
	}
	builder.ConnectSignals(signals)

	return builder
}
