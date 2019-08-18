package main

import (
	"log"
	"os"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	cfg "github.com/novikovSU/rocketchat-desktop-native/config"
	"github.com/novikovSU/rocketchat-desktop-native/model"
	"github.com/novikovSU/rocketchat-desktop-native/rocket"
	"github.com/novikovSU/rocketchat-desktop-native/settings"
	"github.com/novikovSU/rocketchat-desktop-native/ui"
	"github.com/novikovSU/rocketchat-desktop-native/utils"
)

var (
	GtkApplication *gtk.Application
	MainWindow     *gtk.Window
)

//deprecated
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
	utils.Safe(app.Connect("startup", func() {
		log.Println("application startup")
	}))

	// Connect function to application activate event
	utils.Safe(app.Connect("activate", func() {
		log.Println("application activate")

		ui.InitUI(app)

		// Get application config
		settings.Conf, err = settings.GetConfig()
		if err == nil {
			ui.OpenMainWindow(app)
		} else {
			//TODO handle situation properly: try to connect in goroutine and open window again if fails
			ui.OpenConnectionWindow(app)
		}

		model.Init(settings.Conf.User)
		ui.InitSubscribers()

		// Get Rocket.Chat connection
		rocket.InitRocket()

		//		systray.Run(onSysTrayReady, onSysTrayExit)
	}))

	// Connect function to application shutdown event, this is not required.
	utils.Safe(app.Connect("shutdown", func() {
		log.Println("application shutdown")
	}))

	// Launch the application
	os.Exit(app.Run(os.Args))
}
