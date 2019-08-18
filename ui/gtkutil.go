package ui

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/novikovSU/rocketchat-desktop-native/utils"
)

func GetTextViewBuffer(tv *gtk.TextView) *gtk.TextBuffer {
	buf, err := tv.GetBuffer()
	utils.AssertErr(err)

	return buf
}

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

func GetHeaderBar(name string) *gtk.HeaderBar {
	return (*getGtkObjectSafe(name)).(*gtk.HeaderBar)
}

func GetMenuButton(name string) *gtk.MenuButton {
	return (*getGtkObjectSafe(name)).(*gtk.MenuButton)
}

func GetWindow(name string) *gtk.Window {
	return (*getGtkObjectSafe(name)).(*gtk.Window)
}

func GetScrolledWindow(name string) *gtk.ScrolledWindow {
	return (*getGtkObjectSafe(name)).(*gtk.ScrolledWindow)
}

func getGtkObjectSafe(name string) *glib.IObject {
	obj, err := GtkBuilder.GetObject(name)
	utils.AssertErr(err)

	return &obj
}
