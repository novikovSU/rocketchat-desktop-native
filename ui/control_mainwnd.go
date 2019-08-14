package ui

import (
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/novikovSU/rocketchat-desktop-native/utils"
)

var (
	mainWindowIsFocused bool
)

func OpenMainWindow(app *gtk.Application) {
	wnd := CreateWindow("main_window")

	utils.Safe(wnd.Connect("focus-in-event", func() {
		log.Printf("DEBUG: Main window is focused\n")
		mainWindowIsFocused = true
	}))

	utils.Safe(wnd.Connect("focus-out-event", func() {
		log.Printf("DEBUG: Main window is UNfocused\n")
		mainWindowIsFocused = false
	}))

	createMenuBar(app, wnd)

	InitContactListControl()
	InitChatListControl()
	InitSendMsgControl()

	wnd.ShowAll()
	app.AddWindow(wnd)
}

func createMenuBar(app *gtk.Application, wnd *gtk.Window) {
	// Set up the menu model for the button
	menu := glib.MenuNew()

	// Actions with the prefix 'app' reference actions on the application
	// Actions with the prefix 'win' reference actions on the current window (specific to ApplicationWindow)
	// Other prefixes can be added to widgets via InsertActionGroup
	menu.Append("Connect", "custom.connect")
	menu.Append("Disconnect", "custom.disconnect")
	menu.Append("Quit", "app.quit")

	wnd.InsertActionGroup("custom", createCustomActionGroup(app))

	// Create a new menu button
	menuBtn := GetMenuButton("main_menu_button")
	menuBtn.SetMenuModel(&menu.MenuModel)

	createTitleBar(wnd, menuBtn)

	// Add Quit action to menu
	quitAction := glib.SimpleActionNew("quit", nil)
	utils.Safe(quitAction.Connect("activate", app.Quit))
	app.AddAction(quitAction)

	// Add action for X-button
	utils.Safe(wnd.Connect("destroy", app.Quit))
}

func createCustomActionGroup(app *gtk.Application) *glib.SimpleActionGroup {
	actionGroup := glib.SimpleActionGroupNew()

	// Create actions in the custom action group
	addAction(app, actionGroup, "connect", onConnectMenuAction)
	addAction(app, actionGroup, "disconnect", onDisconnectMenuAction)

	return actionGroup
}

func onConnectMenuAction(app *gtk.Application) {
	log.Println("CONNECTED")
	OpenConnectionWindow(app)
}

func onDisconnectMenuAction(app *gtk.Application) {
	log.Println("DISCONNECTED")
}

func addAction(app *gtk.Application, group *glib.SimpleActionGroup, name string, fn func(app *gtk.Application)) {
	action := glib.SimpleActionNew(name, nil)
	group.AddAction(action)

	utils.Safe(action.Connect("activate", func() { fn(app) }))
}

func createTitleBar(wnd *gtk.Window, menuBtn *gtk.MenuButton) {
	header := GetHeaderBar("main_header")
	header.SetShowCloseButton(true)

	// add the menu button to the header
	header.PackStart(menuBtn)

	// Assemble the window
	wnd.SetTitlebar(header)
}
