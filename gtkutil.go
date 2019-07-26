package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"log"
)

func GetGtkButton(name string) *gtk.Button {
	return (*getGtkObjectSafe(name)).(*gtk.Button)
}

func GetGtkInputText(name string) *gtk.Entry {
	return (*getGtkObjectSafe(name)).(*gtk.Entry)
}

func getGtkObjectSafe(name string) *glib.IObject {
	obj, err := GtkBuilder.GetObject(name)
	if err != nil {
		log.Panic(err)
	}

	return &obj
}
