package main

import (
	"github.com/novikovSU/rocketchat-desktop-native/rocket"
	"github.com/novikovSU/rocketchat-desktop-native/settings"
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/novikovSU/rocketchat-desktop-native/ui"
)

func OpenConnectionWindow() {
	wnd := createModal("connection_window")

	okBtn := ui.GetGtkButton("connection_ok_button")
	_, _ = okBtn.Connect("clicked", func() {
		log.Println("connection_ok_button clicked")
		newConf := createConfig()
		if testConfig(newConf) {
			log.Println("connection settings are correct")

			settings.Conf = newConf
			settings.StoreConfig(settings.Conf)
			openMainWindow(GtkApplication)

			wnd.Close()
		} else {
			log.Println("connection settings are incorrect!")
			//TODO
		}
	})

	wnd.ShowAll()
	GtkApplication.AddWindow(wnd)
}

func createModal(id string) *gtk.Dialog {
	obj, err := ui.GtkBuilder.GetObject(id)
	if err != nil {
		log.Panic(err)
	}

	wnd := obj.(*gtk.Dialog)

	// Create the action "wnd.close"
	wndCloseAction := glib.SimpleActionNew("close", nil)
	wndCloseAction.Connect("activate", func() {
		wnd.Close()
	})
	GtkApplication.AddAction(wndCloseAction)

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
	ctrl := ui.GetGtkInputText(name)
	val, _ := ctrl.GetText()

	return val
}

func testConfig(config *settings.Config) bool {
	err := rocket.GetConnectionSafe(config)
	return err == nil
}
