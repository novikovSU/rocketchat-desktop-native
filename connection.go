package main

import (
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func OpenConnectionWindow() {
	wnd := createModal("connection_window")
	wnd.ShowAll()
	GtkApplication.AddWindow(wnd)
}

func createModal(id string) *gtk.Dialog {
	obj, err := GtkBuilder.GetObject(id)
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
