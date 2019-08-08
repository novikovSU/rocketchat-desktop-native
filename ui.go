package main

//TODO move all of them to ui package

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/novikovSU/gorocket/api"

	"github.com/novikovSU/rocketchat-desktop-native/bus"
	"github.com/novikovSU/rocketchat-desktop-native/model"
	"github.com/novikovSU/rocketchat-desktop-native/ui"
)

const (
	hashSign = "\u0023"     // Hash sign for channels
	lockSign = "\U0001F512" // Lock sign for private groups
)

func initUI() {
	bus.Sub(bus.Contacts_update_finished, func() {
		//TODO react to remove users/groups/channels: if it happens, what should we do with selection, for example?
		cs := contactsStore
		cs.Clear()

		addContactsToList(cs, getSortedChannels())
		addContactsToList(cs, getSortedGroups())
		addContactsToList(cs, getSortedUsers())
	})

	bus.Sub(bus.Messages_new, func(msg api.Message) {
		cs := chatStore
		if msg.ChannelID == model.Chat.ActiveContactId || msg.ChannelID == me.ID+model.Chat.ActiveContactId {
			text := strings.Replace(msg.Text, "&nbsp;", "", -1)
			text = strings.Replace(text, "<", "", -1)
			text = strings.Replace(text, ">", "", -1)
			//log.Printf("Text: %s\n", text)
			text = fmt.Sprintf("<b>%s</b> <i>%s</i>\n%s", msg.User.Name, msg.Timestamp.Format("2006-01-02 15:04:05"), text)
			addToList(cs, text)
		}

		//TODO create function for get contactId by message
		model := model.Chat.GetModelById(strings.Replace(msg.ChannelID, me.ID, "", 1))
		if model != nil {
			iter, exists := contactsStore.GetIterFirst()
			if exists {
				for {
					val, err := contactsStore.GetValue(iter, ui.ContactListNameColumn)
					if err == nil {
						strVal, err := val.GetString()
						if err == nil {
							if strings.Compare(strVal, model.GetDisplayName()) == 0 {
								contactsStore.SetValue(iter, ui.ContactListUnreadCountColumn, getUnreadCount(&model))
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
	})

	bus.Sub(bus.Messages_new, func(msg api.Message) {
		if config.Debug {
			log.Printf("DEBUG: Prepare to notificate")
		}
		if ownMessage(msg) {
			return
		}
		if mainWindowIsFocused && msg.ChannelID == model.Chat.ActiveContactId {
			return
		}
		notif := glib.NotificationNew(fmt.Sprintf("%s (%s)", msg.User.Name, msg.User.UserName))
		notif.SetBody(msg.Text)
		GtkApplication.SendNotification(appID, notif)
	})
}

func drawMessage(cs *gtk.ListStore, msg api.Message) error {
	text := strings.Replace(msg.Text, "&nbsp;", "", -1)
	text = strings.Replace(text, "<", "", -1)
	text = strings.Replace(text, ">", "", -1)
	//log.Printf("Text: %s\n", text)
	text = fmt.Sprintf("<b>%s</b> <i>%s</i>\n%s", msg.User.Name, msg.Timestamp.Format("2006-01-02 15:04:05"), text)
	addToList(cs, text)

	return nil
}

func fillChat(cs *gtk.ListStore, name string) {
	model.Chat.ActiveContactId, _ = getRIDByName(name)

	msgs, err := getHistoryByID(model.Chat.ActiveContactId)
	if err != nil {
		log.Printf("ERROR: can't get history by name %s: %s\n", name, err)
		return
	}

	sort.Sort(ui.MessageSorter(msgs))

	cs.Clear()
	for _, msg := range msgs {
		drawMessage(cs, msg)
	}
}

func clearContactUnreadCount(cs *gtk.ListStore, name string) {
	currID, _ := getRIDByName(name)

	model := model.Chat.GetModelById(strings.Replace(currID, me.ID, "", 1))
	model.ClearUnreadCount()

	if model != nil {
		iter, exists := contactsStore.GetIterFirst()
		if exists {
			for {
				val, err := contactsStore.GetValue(iter, ui.ContactListNameColumn)
				if err == nil {
					strVal, err := val.GetString()
					if err == nil {
						if strings.Compare(strVal, model.GetDisplayName()) == 0 {
							contactsStore.SetValue(iter, ui.ContactListUnreadCountColumn, getUnreadCount(&model))
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

func getSortedChannels() []model.IContactModel {
	contacts := model.ChannelsMapToModels(model.Chat.Channels)
	sort.Sort(ui.ContactsSorter(contacts))

	return contacts
}

func getSortedGroups() []model.IContactModel {
	contacts := model.GroupsMapToModels(model.Chat.Groups)
	sort.Sort(ui.ContactsSorter(contacts))

	return contacts
}

func getSortedUsers() []model.IContactModel {
	contacts := model.UsersMapToModels(model.Chat.Users)
	sort.Sort(ui.ContactsSorter(contacts))

	return contacts
}

func addContactsToList(cs *gtk.ListStore, contacts []model.IContactModel) {
	for _, contact := range contacts {
		if contact.GetName() != "" {
			iter := cs.Append()
			cs.SetValue(iter, ui.ContactListNameColumn, contact.GetDisplayName())
			cs.SetValue(iter, ui.ContactListUnreadCountColumn, getUnreadCount(&contact))
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

func updateUnreadCount(msg api.Message) {
	/*msg.ChannelID
	model.Chat.*/
}
