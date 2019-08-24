package ui

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/novikovSU/rocketchat-desktop-native/bus"
	"strconv"

	"github.com/novikovSU/rocketchat-desktop-native/model"
	"github.com/novikovSU/rocketchat-desktop-native/utils"
)

const (
	ContactListNameColumn        = 0
	ContactListUnreadCountColumn = 1
)

var (
	contactList   *gtk.TreeView
	contactsStore *gtk.ListStore
)

func InitContactListControl() {
	contactList, contactsStore = createContactListTreeView()
	chatCaption := GetLabel("chat_caption")

	sel, err := contactList.GetSelection()
	utils.AssertErr(err)

	utils.Safe(sel.Connect("changed", func() {
		contactName := GetTreeViewSelectionVal(contactList, ContactListNameColumn)
		chatCaption.SetText(contactName)

		bus.Pub(bus.Ui_contacts_selected, contactName)
		refreshChatHistory(chatStore, contactName)
	}))
}

func createContactListTreeView() (*gtk.TreeView, *gtk.ListStore) {
	contactList := GetTreeView("contact_list")

	contactList.AppendColumn(createTextViewTextColumn("name", ContactListNameColumn))
	contactList.AppendColumn(createTextViewTextColumn("unreadCount", ContactListUnreadCountColumn))

	contactListStore := createListStore(glib.TYPE_STRING, glib.TYPE_STRING)
	contactList.SetModel(contactListStore)

	return contactList, contactListStore
}

func getUnreadCount(contactModel *model.IContactModel) string {
	count := (*contactModel).GetUnreadCount()
	if count > 0 {
		return strconv.Itoa(count)
	}
	return ""
}
