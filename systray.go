package main

import (
	"github.com/novikovSU/gorocket/api"
	"github.com/novikovSU/rocketchat-desktop-native/bus"
	"io/ioutil"
	"log"

	"github.com/getlantern/systray"
)

var (
	showMenuItem, quitMenuItem *systray.MenuItem
)

func onSysTrayReady() {
	err := setIcon("icon.ico")
	if err == nil {
		systray.SetTooltip("Rocket.Chat Desktop native")

		showMenuItem = systray.AddMenuItem("Show", "Show chat window")
		quitMenuItem = systray.AddMenuItem("Quit", "Quit")

		bus.SubscribeAsync("messages.new", onNewMessage)
		go handleMenuItemEvents()
	}
}

func onSysTrayExit() {
	// Cleaning stuff here.
}

func handleMenuItemEvents() {
	for {
		select {
		case <-showMenuItem.ClickedCh:
			onShowMenuItemClicked()
		case <-quitMenuItem.ClickedCh:
			onQuitMenuItemClicked()
			return
		}
	}
}

func onShowMenuItemClicked() {
	//TODO recreate window, if it has been closed
	MainWindow.Present()
}

func onQuitMenuItemClicked() {
	GtkApplication.Quit()
	systray.Quit()
}

func onNewMessage(msg api.Message) {
	//TODO we need event handler to reset icon to normal state then all unread messages will be read
	_ = setIcon("iconExcited.ico")
}

//--------------------------------------------------------

func setIcon(name string) error {
	icondData, err := ioutil.ReadFile("resources/" + name)
	if err != nil {
		log.Println("Could not load icon %s. err: %s\n", name, err)
	} else {
		systray.SetIcon(icondData)
	}

	return err
}
