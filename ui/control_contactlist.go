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
	contactList   *gtk.TreeView
	contactsStore *gtk.ListStore
)

func InitContactListControl() {
	contactList, contactsStore = createContactListTreeView()
	chatCaption := GetLabel("chat_caption")

	sel, err := contactList.GetSelection()
	utils.AssertErr(err)

	utils.Safe(sel.Connect("changed", func() {
		selVal := GetTreeViewSelectionVal(contactList, ContactListNameColumn)
		chatCaption.SetText(selVal)
		refreshChatHistory(chatStore, selVal)
		clearContactUnreadCount(contactsStore, selVal)
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

	mdl := model.Chat.GetModelById(strings.Replace(currID, rocket.Me.ID, "", 1))
	mdl.ClearUnreadCount()

	if mdl != nil {
		iter, exists := cs.GetIterFirst()
		if exists {
			//TODO can we don't use infinite loop syntax?
			for {
				val, err := cs.GetValue(iter, ContactListNameColumn)
				if err == nil {
					strVal, err := val.GetString()
					if err == nil {
						if strings.Compare(strVal, mdl.String()) == 0 {
							utils.AssertErr(cs.SetValue(iter, ContactListUnreadCountColumn, getUnreadCount(&mdl)))
							break
						} else {
							if !contactsStore.IterNext(iter) {
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
