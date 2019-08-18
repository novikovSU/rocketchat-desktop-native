package ui

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/novikovSU/rocketchat-desktop-native/rocket"
	"github.com/novikovSU/rocketchat-desktop-native/settings"
	"github.com/novikovSU/rocketchat-desktop-native/utils"
)

func OpenConnectionWindow(app *gtk.Application) {
	wnd := createModal(app, "connection_window")

	okBtn := GetGtkButton("connection_ok_button")
	_, _ = okBtn.Connect("clicked", func() {
		logger.Debug("connection_ok_button clicked")
		newConf := createConfig()
		if testConfig(newConf) {
			logger.Debug("connection settings are correct")

			settings.Conf = newConf
			settings.StoreConfig(settings.Conf)
			OpenMainWindow(app)

			wnd.Close()
		} else {
			logger.Debug("connection settings are incorrect!")
			//TODO
		}
	})

	wnd.ShowAll()
	app.AddWindow(wnd)
}

func createModal(app *gtk.Application, id string) *gtk.Dialog {
	obj, err := GtkBuilder.GetObject(id)
	utils.AssertErr(err)

	wnd := obj.(*gtk.Dialog)

	// Create the action "wnd.close"
	wndCloseAction := glib.SimpleActionNew("close", nil)
	wndCloseAction.Connect("activate", func() {
		wnd.Close()
	})
	app.AddAction(wndCloseAction)

	return wnd
}

func createConfig() *settings.Config {
	//TODO validation
	config := settings.CreateDefaultConfig()
	config.Server = getInputTextValue("server_input_text")
	config.User = getInputTextValue("login_input_text")
	config.Email = getInputTextValue("e_mail_input_text")
	config.Password = getInputTextValue("password_input_text")

	return config
}

func getInputTextValue(name string) string {
	ctrl := GetGtkInputText(name)
	val, _ := ctrl.GetText()

	return val
}

func testConfig(config *settings.Config) bool {
	err := rocket.GetConnectionSafe(config)
	return err == nil
}
