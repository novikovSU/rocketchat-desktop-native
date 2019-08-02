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

// ChannelsSorter sorts users by name.
type ChannelsSorter []model.ChannelModel

func (a ChannelsSorter) Len() int           { return len(a) }
func (a ChannelsSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ChannelsSorter) Less(i, j int) bool { return a[i].Channel.Name < a[j].Channel.Name }

// GroupsSorter sorts users by name.
type GroupsSorter []model.GroupModel

func (a GroupsSorter) Len() int           { return len(a) }
func (a GroupsSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a GroupsSorter) Less(i, j int) bool { return a[i].Group.Name < a[j].Group.Name }

// UsersSorter sorts users by name.
type UsersSorter []model.UserModel

func (a UsersSorter) Len() int           { return len(a) }
func (a UsersSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a UsersSorter) Less(i, j int) bool { return a[i].User.Name < a[j].User.Name }

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
		for _, channel := range getSortedChannels() {
			addToList(cs, hashSign+channel.Channel.Name)
		}

		for _, group := range getSortedGroups() {
			addToList(cs, lockSign+group.Group.Name)
		}

		for _, user := range getSortedUsers() {
			if user.User.Name != "" {
				addToList(cs, user.User.Name)
			}
		}
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

func getSortedChannels() []model.ChannelModel {
	channels := make([]model.ChannelModel, 0, len(model.Chat.Channels))
	for _, u := range model.Chat.Channels {
		channels = append(channels, u)
	}
	sort.Sort(ChannelsSorter(channels))

	return channels
}

func getSortedGroups() []model.GroupModel {
	groups := make([]model.GroupModel, 0, len(model.Chat.Groups))
	for _, u := range model.Chat.Groups {
		groups = append(groups, u)
	}
	sort.Sort(GroupsSorter(groups))

	return groups
}

func getSortedUsers() []model.UserModel {
	users := make([]model.UserModel, 0, len(model.Chat.Users))
	for _, u := range model.Chat.Users {
		users = append(users, u)
	}
	sort.Sort(UsersSorter(users))

	return users
}
