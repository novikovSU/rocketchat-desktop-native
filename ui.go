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

// NameSorter sorts users by name.
type NameSorter []model.UserModel

func (a NameSorter) Len() int           { return len(a) }
func (a NameSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a NameSorter) Less(i, j int) bool { return a[i].User.Name < a[j].User.Name }

// DateSorter sorts messages by timestamp.
type DateSorter []api.Message

func (a DateSorter) Len() int           { return len(a) }
func (a DateSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DateSorter) Less(i, j int) bool { return a[i].Timestamp.Before(*a[j].Timestamp) }

func initContactListSubscribers(cs *gtk.ListStore) {
	//TODO react to remove users/groups/channels

	bus.Sub(bus.Contacts_update_finished, func() {
		cs.Clear()
		for _, channel := range model.Chat.Channels {
			addToList(cs, hashSign+channel.Channel.Name)
		}

		for _, group := range model.Chat.Groups {
			addToList(cs, lockSign+group.Group.Name)
		}

		for _, user := range getSortedUsers() {
			if user.User.Name != "" {
				addToList(cs, user.User.Name)
			}
		}
	})
}

func fillChat(cs *gtk.ListStore, name string) {
	msgs, err := getHistoryByName(name)
	if err != nil {
		log.Printf("ERROR: can't get history by name %s: %s\n", name, err)
		return
	}

	currentChatID, _ = getIDByName(name)

	sort.Sort(DateSorter(msgs))

	cs.Clear()
	for _, msg := range msgs {
		text := strings.Replace(msg.Text, "&nbsp;", "", -1)
		text = strings.Replace(text, "<", "", -1)
		text = strings.Replace(text, ">", "", -1)
		//log.Printf("Text: %s\n", text)
		text = fmt.Sprintf("<b>%s</b> <i>%s</i>\n%s", msg.User.Name, msg.Timestamp.Format("2006-01-02 15:04:05"), text)
		addToList(cs, text)
	}
}

func showUpdates(cs *gtk.ListStore) {

}

func getSortedUsers() []model.UserModel {
	users := make([]model.UserModel, 0, len(model.Chat.Users))
	for _, u := range model.Chat.Users {
		users = append(users, u)
	}
	sort.Sort(NameSorter(users))

	return users
}
