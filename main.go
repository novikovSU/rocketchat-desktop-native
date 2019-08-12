package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
	cfg "github.com/novikovSU/rocketchat-desktop-native/config"
	"github.com/novikovSU/rocketchat-desktop-native/rocket"
	"github.com/novikovSU/rocketchat-desktop-native/settings"
	"github.com/novikovSU/rocketchat-desktop-native/ui"
	"log"
	"os"
)

var (
	GtkApplication *gtk.Application
	MainWindow     *gtk.Window

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

func main() {
	// Create a new application.
	app, err := gtk.ApplicationNew(cfg.AppID, glib.APPLICATION_FLAGS_NONE)
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

		ui.InitUI(app)
		ui.SendNotification("Rocket.Chat Desktop native", "application activated")

		// Get application config
		settings.Conf, err = settings.GetConfig()
		if err == nil {
			openMainWindow(app)
		} else {
			//TODO handle situation properly: try to connect in goroutine and open window again if fails
			OpenConnectionWindow()
		}

		initUI()

		// Get Rocket.Chat connection
		rocket.InitRocket()

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
	MainWindow = ui.CreateWindow("main_window")

	MainWindow.Connect("focus-in-event", func() {
		log.Printf("DEBUG: Main window is focused\n")
		mainWindowIsFocused = true
	})

	MainWindow.Connect("focus-out-event", func() {
		log.Printf("DEBUG: Main window is UNfocused\n")
		mainWindowIsFocused = false
	})

	createMenuBar()

	ui.InitContactListControl()

	chatList := ui.GetTreeView("chat_list")
	rightScrolledWindow := ui.GetScrolledWindow("right_scrolled_window")

	// Autoscroll of chatList
	chatList.Connect("size-allocate", func() {
		adj := rightScrolledWindow.GetVAdjustment()
		adj.SetValue(adj.GetUpper() - adj.GetPageSize())
	})
	chatList.ConnectAfter("size-allocate", func() {

	})

	// ------------------

	ui.ChatStore = initList(chatList)

	ui.InitSendMsgControl()

	MainWindow.ShowAll()
	app.AddWindow(MainWindow)
}

func createMenuBar() {
	/* DISABLE custom header and menu */
	// Get a headerbar

	// Create a new menu button
	mbtn := ui.GetMenuButton("main_menu_button")

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
	MainWindow.InsertActionGroup("custom", customActionGroup)

	// Create an action in the custom action group
	aConnect := glib.SimpleActionNew("connect", nil)
	aConnect.Connect("activate", func() {
		log.Println("CONNECTED")
		OpenConnectionWindow()
	})
	customActionGroup.AddAction(aConnect)
	GtkApplication.AddAction(aConnect)

	aDisconnect := glib.SimpleActionNew("disconnect", nil)
	aDisconnect.Connect("activate", func() {
		log.Println("DISCONNECTED")
	})
	customActionGroup.AddAction(aDisconnect)
	GtkApplication.AddAction(aDisconnect)

	mbtn.SetMenuModel(&menu.MenuModel)

	createTitleBar(mbtn)

	// Add Quit action to menu
	aQuit := glib.SimpleActionNew("quit", nil)
	aQuit.Connect("activate", GtkApplication.Quit)
	GtkApplication.AddAction(aQuit)

	// Add action for X-button
	MainWindow.Connect("destroy", GtkApplication.Quit)
}

func createTitleBar(menuBtn *gtk.MenuButton) {
	header := ui.GetHeaderBar("main_header")
	header.SetShowCloseButton(true)

	// add the menu button to the header
	header.PackStart(menuBtn)

	// Assemble the window
	MainWindow.SetTitlebar(header)
}
