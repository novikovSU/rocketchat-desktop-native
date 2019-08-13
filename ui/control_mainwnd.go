package ui

import (
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/novikovSU/rocketchat-desktop-native/utils"
)

var (
	MainWindowIsFocused bool
)

func OpenMainWindow(app *gtk.Application) {
	wnd := CreateWindow("main_window")

	utils.Safe(wnd.Connect("focus-in-event", func() {
		log.Printf("DEBUG: Main window is focused\n")
		MainWindowIsFocused = true
	}))

	utils.Safe(wnd.Connect("focus-out-event", func() {
		log.Printf("DEBUG: Main window is UNfocused\n")
		MainWindowIsFocused = false
	}))

	createMenuBar(app, wnd)

	InitContactListControl()
	InitChatListControl()
	InitSendMsgControl()

	wnd.ShowAll()
	app.AddWindow(wnd)
}

func createMenuBar(app *gtk.Application, wnd *gtk.Window) {
	/* DISABLE custom header and menu */
	// Get a headerbar

	// Create a new menu button
	mbtn := GetMenuButton("main_menu_button")

	// Set up the menu model for the button
	menu := glib.MenuNew()

	// Actions with the prefix 'app' reference actions on the application
	// Actions with the prefix 'win' reference actions on the current window (specific to ApplicationWindow)
	// Other prefixes can be added to widgets via InsertActionGroup
	menu.Append("Connect", "custom.connect")
	menu.Append("Disconnect", "custom.disconnect")
	menu.Append("Quit", "app.quit")

	customActionGroup := glib.SimpleActionGroupNew()
	wnd.InsertActionGroup("custom", customActionGroup)

	// Create an action in the custom action group
	aConnect := glib.SimpleActionNew("connect", nil)
	aConnect.Connect("activate", func() {
		log.Println("CONNECTED")
		OpenConnectionWindow(app)
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

	createTitleBar(wnd, mbtn)

	// Add Quit action to menu
	aQuit := glib.SimpleActionNew("quit", nil)
	aQuit.Connect("activate", app.Quit)
	app.AddAction(aQuit)

	// Add action for X-button
	wnd.Connect("destroy", app.Quit)
}

func createTitleBar(wnd *gtk.Window, menuBtn *gtk.MenuButton) {
	header := GetHeaderBar("main_header")
	header.SetShowCloseButton(true)

	// add the menu button to the header
	header.PackStart(menuBtn)

	// Assemble the window
	wnd.SetTitlebar(header)
}
