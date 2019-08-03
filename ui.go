package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/novikovSU/rocketchat-desktop-native/bus"
	"github.com/novikovSU/rocketchat-desktop-native/model"

	"github.com/gotk3/gotk3/gtk"
	"github.com/novikovSU/gorocket/api"
)

const (
	hashSign = "\u0023"     // Hash sign for channels
	lockSign = "\U0001F512" // Lock sign for private groups
)

// ContactsSorter sorts by name.
type ContactsSorter []model.IContactModel

func (c ContactsSorter) Len() int           { return len(c) }
func (c ContactsSorter) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ContactsSorter) Less(i, j int) bool { return c[i].GetName() < c[j].GetName() }

// DateSorter sorts messages by timestamp.
type DateSorter []api.Message

func (a DateSorter) Len() int           { return len(a) }
func (a DateSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DateSorter) Less(i, j int) bool { return a[i].Timestamp.Before(*a[j].Timestamp) }

func initContactListSubscribers(cs *gtk.ListStore) {
	//TODO react to remove users/groups/channels

}

func initUI() {
	bus.Sub(bus.Contacts_update_finished, func() {
		cs := contactsStore
		cs.Clear()

		addContactsToList(cs, getSortedChannels())
		addContactsToList(cs, getSortedGroups())
		addContactsToList(cs, getSortedUsers())
	})

	bus.Sub(bus.Messages_new, func(msg api.Message) {
		cs := chatStore
		if msg.ChannelID == currentChatID || msg.ChannelID == me.ID+currentChatID {
			text := strings.Replace(msg.Text, "&nbsp;", "", -1)
			text = strings.Replace(text, "<", "", -1)
			text = strings.Replace(text, ">", "", -1)
			//log.Printf("Text: %s\n", text)
			text = fmt.Sprintf("<b>%s</b> <i>%s</i>\n%s", msg.User.Name, msg.Timestamp.Format("2006-01-02 15:04:05"), text)
			addToList(cs, text)
		}
	})

	bus.Sub(bus.Messages_new, func(msg api.Message) {
		log.Printf("DEBUG: Prepare to notificate")
		notif.SetBody("application activate")
		//GtkApplication.SendNotification(appID, notif)
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
	currentChatID, _ = getRIDByName(name)

	msgs, err := getHistoryByID(currentChatID)
	if err != nil {
		log.Printf("ERROR: can't get history by name %s: %s\n", name, err)
		return
	}

	sort.Sort(DateSorter(msgs))

	cs.Clear()
	for _, msg := range msgs {
		drawMessage(cs, msg)
	}
}

func getSortedChannels() []model.IContactModel {
	return getSortedContacts(func(contacts []model.IContactModel) []model.IContactModel {
		for _, u := range model.Chat.Channels {
			contacts = append(contacts, u)
		}
		return contacts
	})
}

func getSortedGroups() []model.IContactModel {
	return getSortedContacts(func(contacts []model.IContactModel) []model.IContactModel {
		for _, u := range model.Chat.Groups {
			contacts = append(contacts, u)
		}
		return contacts
	})
}

func getSortedUsers() []model.IContactModel {
	return getSortedContacts(func(contacts []model.IContactModel) []model.IContactModel {
		for _, u := range model.Chat.Users {
			contacts = append(contacts, u)
		}
		return contacts
	})
}

func getSortedContacts(appender func([]model.IContactModel) []model.IContactModel) []model.IContactModel {
	contacts := appender(make([]model.IContactModel, 0))
	sort.Sort(ContactsSorter(contacts))
	return contacts
}

func addContactsToList(cs *gtk.ListStore, contacts []model.IContactModel) {
	for _, contact := range contacts {
		if contact.GetName() != "" {
			addToList(cs, contact.GetDisplayName())
		}
	}
}
