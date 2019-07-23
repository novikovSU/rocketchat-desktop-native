package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/novikovSU/rocketchat-desktop-native/appconf"
)

const appID = "com.github.novikovSU.rocketchat-desktop-native"

const (
	iistItem = iota
	nColumns
	keyEnter    = 65293
	senderConst = "Полиграфов Полиграф"
)

var (
	ctrlPressed = false
)

func initList(list *gtk.TreeView) *gtk.ListStore {
	cellRenderer, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Fatal(err)
	}

	column, err := gtk.TreeViewColumnNewWithAttribute("List Items", cellRenderer, "markup", 0)
	if err != nil {
		log.Fatal(err)
	}
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
	// Create a new application.
	app, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		log.Panic(err)
	}

	// Connect function to application startup event, this is not required.
	app.Connect("startup", func() {
		log.Println("application startup")
	})

	// Connect function to application activate event
	app.Connect("activate", func() {
		log.Println("application activate")

		notif := glib.NotificationNew("Rocket.Chat Desktop native")
		notif.SetBody("application activate")
		app.SendNotification(appID, notif)

		// Get the GtkBuilder UI definition in the glade file.
		builder, err := gtk.BuilderNewFromFile("main.glade")
		if err != nil {
			log.Panic(err)
		}

		// Map the handlers to callback functions, and connect the signals
		// to the Builder.
		signals := map[string]interface{}{
			"on_main_window_destroy": onMainWindowDestroy,
		}
		builder.ConnectSignals(signals)

		// Get the object with the id of "main_window".
		obj, err := builder.GetObject("main_window")
		if err != nil {
			log.Panic(err)
		}

		// Verify that the object is a pointer to a gtk.ApplicationWindow.
		win, err := isWindow(obj)
		if err != nil {
			log.Panic(err)
		}

		// Create menu
		// Create a header bar
		header, err := gtk.HeaderBarNew()
		if err != nil {
			log.Fatal("Could not create header bar:", err)
		}
		header.SetShowCloseButton(true)
		header.SetTitle("Rocket.Chat Desktop Native")
		header.SetSubtitle("Do we need a subtitle?")

		// Create a new menu button
		mbtn, err := gtk.MenuButtonNew()
		if err != nil {
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

		// Create the action "win.close"
		aClose := glib.SimpleActionNew("close", nil)
		aClose.Connect("activate", func() {
			win.Close()
		})
		app.AddAction(aClose)

		customActionGroup := glib.SimpleActionGroupNew()
		win.InsertActionGroup("custom", customActionGroup)

		// Create an action in the custom action group
		aConnect := glib.SimpleActionNew("connect", nil)
		aConnect.Connect("activate", func() {
			log.Println("CONNECTED")
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

		obj, err = builder.GetObject("contact_list")
		if err != nil {
			log.Fatal(err)
		}
		contactList, ok := obj.(*gtk.TreeView)
		if ok != true {
			log.Fatal(err)
		}

		// -------------

		obj, err = builder.GetObject("chat_caption")
		if err != nil {
			log.Fatal(err)
		}
		chatCaption, ok := obj.(*gtk.Label)
		if ok != true {
			log.Fatal(err)
		}

		obj, err = builder.GetObject("chat_list")
		if err != nil {
			log.Fatal(err)
		}
		chatList, ok := obj.(*gtk.TreeView)
		if ok != true {
			log.Fatal(err)
		}

		// Autoscroll of chatList
		obj, err = builder.GetObject("right_scrolled_window")
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

		obj, err = builder.GetObject("text_input")
		if err != nil {
			log.Fatal(err)
		}
		textInput, ok := obj.(*gtk.TextView)
		if ok != true {
			log.Fatal(err)
		}

		// ------------------

		contactsStore := initList(contactList)
		addToList(contactsStore, "#general")
		addToList(contactsStore, "#rocket-native")
		addToList(contactsStore, "Иванов Иван")
		addToList(contactsStore, "Петров Петр")
		addToList(contactsStore, "Сидоров Сидор")
		addToList(contactsStore, "Иванов Иван")
		addToList(contactsStore, "Петров Петр")

		chatStore := initList(chatList)
		addToList(chatStore, "<i>--- Вы вошли в чат. ---</i>")

		contactsSelection, err := contactList.GetSelection()
		if err != nil {
			log.Fatal(err)
		}

		contactsSelection.Connect("changed", func() {
			//onChanged(contactsSelection, chatCaption)
			selectionText := getSelectionText(contactsSelection)
			header.SetSubtitle(selectionText)
			chatCaption.SetText(selectionText)
		})

		/*chatSelection, err := chatList.GetSelection()
		if err != nil {
			log.Fatal(err)
		}

		chatSelection.Connect("changed", func() {
			onChanged(chatSelection, rightLabel)
		}) */

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
						now := time.Now()
						timeText := now.Format("2006-01-02 15:04:05")
						text := fmt.Sprintf("<b>%s</b> <i>%s</i>\n%s", senderConst, timeText, inputText)
						addToList(chatStore, text)
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
	})

	// Connect function to application shutdown event, this is not required.
	app.Connect("shutdown", func() {
		log.Println("application shutdown")
	})

	// Launch the application
	os.Exit(app.Run(os.Args))
}
