package main

//TODO move all of them to ui package

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/gotk3/gotk3/gtk"

	"github.com/novikovSU/gorocket/api"
	"github.com/novikovSU/rocketchat-desktop-native/bus"
	"github.com/novikovSU/rocketchat-desktop-native/config"
	"github.com/novikovSU/rocketchat-desktop-native/model"
	"github.com/novikovSU/rocketchat-desktop-native/rocket"
	"github.com/novikovSU/rocketchat-desktop-native/ui"
	"github.com/novikovSU/rocketchat-desktop-native/utils"
)

func initUI() {
	bus.Sub(bus.Contacts_update_finished, func() {
		//TODO react to remove users/groups/channels: if it happens, what should we do with selection, for example?
		cs := ui.ContactsStore
		cs.Clear()

		addContactsToList(cs, getSortedChannels())
		addContactsToList(cs, getSortedGroups())
		addContactsToList(cs, getSortedUsers())
	})

	bus.Sub(bus.Messages_new, func(msg api.Message) {
		cs := ui.ChatStore
		if msg.ChannelID == model.Chat.ActiveContactId || msg.ChannelID == rocket.Me.ID+model.Chat.ActiveContactId {
			text := strings.Replace(msg.Text, "&nbsp;", "", -1)
			text = strings.Replace(text, "<", "", -1)
			text = strings.Replace(text, ">", "", -1)
			//log.Printf("Text: %s\n", text)
			text = fmt.Sprintf("<b>%s</b> <i>%s</i>\n%s", msg.User.Name, msg.Timestamp.Format("2006-01-02 15:04:05"), text)
			ui.AddToList(cs, text)
		}

		//TODO create function for get contactId by message
		model := model.Chat.GetModelById(strings.Replace(msg.ChannelID, rocket.Me.ID, "", 1))
		if model != nil {
			iter, exists := ui.ContactsStore.GetIterFirst()
			if exists {
				for {
					val, err := ui.ContactsStore.GetValue(iter, ui.ContactListNameColumn)
					if err == nil {
						strVal, err := val.GetString()
						if err == nil {
							if strings.Compare(strVal, model.String()) == 0 {
								ui.ContactsStore.SetValue(iter, ui.ContactListUnreadCountColumn, ui.GetUnreadCount(&model))
								break
							} else {
								if !ui.ContactsStore.IterNext(iter) {
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
		if rocket.OwnMessage(msg) {
			return
		}
		if mainWindowIsFocused && msg.ChannelID == model.Chat.ActiveContactId {
			return
		}

		notifTitle := fmt.Sprintf("%s (%s)", msg.User.Name, msg.User.UserName)
		ui.SendNotification(notifTitle, msg.Text)
	})
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
			utils.AssertErr(cs.SetValue(iter, ui.ContactListNameColumn, contact.String()))
			utils.AssertErr(cs.SetValue(iter, ui.ContactListUnreadCountColumn, ui.GetUnreadCount(&contact)))
		}
	}
}

func updateUnreadCount(msg api.Message) {
	/*msg.ChannelID
	model.Chat.*/
}
