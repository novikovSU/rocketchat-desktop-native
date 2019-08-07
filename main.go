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

	"github.com/novikovSU/rocketchat-desktop-native/ui"
)

const appID = "com.github.novikovSU.rocketchat-desktop-native"

const (
	iistItem = iota
	nColumns
	keyEnter = 65293
)

var (
	GtkApplication *gtk.Application
	MainWindow     *gtk.Window

	contactsStore *gtk.ListStore
	chatStore     *gtk.ListStore

	ctrlPressed  = false
	shiftPressed = false

	mainWindowIsFocused = false
)

//deprecated
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

func main() {

	// Get application config
	config, _ = getConfig()

	// Create a new application.
	app, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		log.Panic(err)
	}
	GtkApplication = app

	// Connect function to application startup event, this is not required.
	app.Connect("startup", func() {
		log.Println("application startup")
	})

	// Connect function to application activate event
	app.Connect("activate", func() {
		log.Println("application activate")

		notif := glib.NotificationNew("Rocket.Chat Desktop native")
		notif.SetBody("application activated")
		GtkApplication.SendNotification(appID, notif)

		ui.InitUI()
		// Get application config
		config, err = getConfig()
		if err == nil {
			// Get Rocket.Chat connection
			getConnection()
			openMainWindow(app)
		} else {
			OpenConnectionWindow()
		}

		initUI()

		// Get Rocket.Chat connection
		initRocket()

		//		systray.Run(onSysTrayReady, onSysTrayExit)
	})

	// Connect function to application shutdown event, this is not required.
	app.Connect("shutdown", func() {
		log.Println("application shutdown")
	})

	// Launch the application
	os.Exit(app.Run(os.Args))
}

func openMainWindow(app *gtk.Application) {
	win := CreateWindow("main_window")
	MainWindow = win

	win.Connect("focus-in-event", func() {
		log.Printf("DEBUG: Main window is focused\n")
		mainWindowIsFocused = true
	})

	win.Connect("focus-out-event", func() {
		log.Printf("DEBUG: Main window is UNfocused\n")
		mainWindowIsFocused = false
	})

	/* DISABLE custom header and menu */
	// Create menu
	// Get a headerbar
	obj, err := ui.GtkBuilder.GetObject("main_header")
	if err != nil {
		log.Panic(err)
	}
	header, ok := obj.(*gtk.HeaderBar)
	if ok != true {
		log.Fatal("Could not create header bar:", err)
	}
	header.SetShowCloseButton(true)

	// Create a new menu button
	obj, err = ui.GtkBuilder.GetObject("main_menu_button")
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

	contactList, cs := ui.CreateContactListTreeView()
	chatCaption := ui.GetLabel("chat_caption")
	chatList := ui.GetTreeView("chat_list")
	rightScrolledWindow := ui.GetScrolledWindow("right_scrolled_window")

	// Autoscroll of chatList
	chatList.Connect("size-allocate", func() {
		adj := rightScrolledWindow.GetVAdjustment()
		adj.SetValue(adj.GetUpper() - adj.GetPageSize())
	})
	chatList.ConnectAfter("size-allocate", func() {

	})

	textInput := ui.GetTextView("text_input")
	// ------------------

	contactsStore = cs

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
					postByNameRT(selectionText, inputText)
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

	win.ShowAll()
	app.AddWindow(win)
}

// CreateWindow AAA
func CreateWindow(id string) *gtk.Window {
	obj, err := ui.GtkBuilder.GetObject(id)
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
