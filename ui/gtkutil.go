package ui

import (
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func GetGtkButton(name string) *gtk.Button {
	return (*getGtkObjectSafe(name)).(*gtk.Button)
}

func GetGtkInputText(name string) *gtk.Entry {
	return (*getGtkObjectSafe(name)).(*gtk.Entry)
}

func GetTreeView(name string) *gtk.TreeView {
	return (*getGtkObjectSafe(name)).(*gtk.TreeView)
}

func GetTextView(name string) *gtk.TextView {
	return (*getGtkObjectSafe(name)).(*gtk.TextView)
}

func GetLabel(name string) *gtk.Label {
	return (*getGtkObjectSafe(name)).(*gtk.Label)
}

func GetScrolledWindow(name string) *gtk.ScrolledWindow {
	return (*getGtkObjectSafe(name)).(*gtk.ScrolledWindow)
}

func getGtkObjectSafe(name string) *glib.IObject {
	obj, err := GtkBuilder.GetObject(name)
	if err != nil {
		log.Panic(err)
	}

	return &obj
}
