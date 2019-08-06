package ui

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	ContactListNameColumn        = 0
	ContactListUnreadCountColumn = 1
)

func CreateContactListTreeView() (*gtk.TreeView, *gtk.ListStore) {
	contactList := GetTreeView("contact_list")

	contactList.AppendColumn(createTextViewTextColumn("name", ContactListNameColumn))
	contactList.AppendColumn(createTextViewTextColumn("unreadCount", ContactListUnreadCountColumn))

	contactListStore := createListStore(glib.TYPE_STRING, glib.TYPE_STRING)
	contactList.SetModel(contactListStore)

	return contactList, contactListStore
}
