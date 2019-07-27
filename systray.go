package main

import (
	"io/ioutil"
	"log"

	"github.com/getlantern/systray"
)

var (
	showMenuItem, quitMenuItem *systray.MenuItem
)

func onSysTrayReady() {
	icondData, err := ioutil.ReadFile("resources/icon.ico")
	if err != nil {
		log.Println("Could not load icon. err: %s\n", err)
	} else {
		systray.SetIcon(icondData)
		systray.SetTooltip("Rocket.Chat Desktop native")

		showMenuItem = systray.AddMenuItem("Show", "Show chat window")
		quitMenuItem = systray.AddMenuItem("Quit", "Quit")

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
