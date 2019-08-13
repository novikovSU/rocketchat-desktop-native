package ui

import (
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/novikovSU/rocketchat-desktop-native/model"
	"github.com/novikovSU/rocketchat-desktop-native/rocket"
	"github.com/novikovSU/rocketchat-desktop-native/utils"
)

const (
	ContactListNameColumn        = 0
	ContactListUnreadCountColumn = 1
)

var (
	ContactList   *gtk.TreeView
	ContactsStore *gtk.ListStore
)

func InitContactListControl() {
	ContactList, ContactsStore = createContactListTreeView()
	chatCaption := GetLabel("chat_caption")

	sel, err := ContactList.GetSelection()
	utils.AssertErr(err)

	utils.Safe(sel.Connect("changed", func() {
		selVal := GetTreeViewSelectionVal(ContactList, ContactListNameColumn)
		chatCaption.SetText(selVal)
		refreshChatHistory(chatStore, selVal)
		clearContactUnreadCount(ContactsStore, selVal)
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

func clearContactUnreadCount(cs *gtk.ListStore, name string) {
	currID, _ := rocket.GetRIDByName(name)

	model := model.Chat.GetModelById(strings.Replace(currID, rocket.Me.ID, "", 1))
	model.ClearUnreadCount()

	if model != nil {
		iter, exists := ContactsStore.GetIterFirst()
		if exists {
			for {
				val, err := ContactsStore.GetValue(iter, ContactListNameColumn)
				if err == nil {
					strVal, err := val.GetString()
					if err == nil {
						if strings.Compare(strVal, model.String()) == 0 {
							ContactsStore.SetValue(iter, ContactListUnreadCountColumn, getUnreadCount(&model))
							break
						} else {
							if !ContactsStore.IterNext(iter) {
								break
							}
						}
					}
				}
			}
		}
	}
}

func getUnreadCount(contactModel *model.IContactModel) string {
	count := (*contactModel).GetUnreadCount()
	if count > 0 {
		return strconv.Itoa(count)
	}
	return ""
}
