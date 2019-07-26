package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/gotk3/gotk3/gtk"
	"github.com/novikovSU/gorocket/api"
)

var (
	hashSign = "\u0023"     // Hash sign for channels
	lockSign = "\U0001F512" // Lock sign for private groups
)

// NameSorter sorts users by name.
type NameSorter []api.User

func (a NameSorter) Len() int           { return len(a) }
func (a NameSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a NameSorter) Less(i, j int) bool { return a[i].Name < a[j].Name }

// DateSorter sorts messages by timestand.
type DateSorter []api.Message

func (a DateSorter) Len() int           { return len(a) }
func (a DateSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DateSorter) Less(i, j int) bool { return a[i].Timestamp.Before(*a[j].Timestamp) }

func fillContactList(cs *gtk.ListStore) {
	cs.Clear()
	channels, err := client.Channel().List()
	if err != nil {
		log.Printf("Can't get channels: %s\n", err)
	}
	for _, channel := range channels {
		addToList(cs, "#"+channel.Name)
	}

	groups, err := client.Groups().ListGroups()
	if err != nil {
		log.Printf("Can't get groups: %s\n", err)
	}
	for _, group := range groups {
		addToList(cs, lockSign+group.Name)
	}

	users, err := client.Users().List()
	if err != nil {
		log.Printf("Can't get ims: %s\n", err)
	}
	sort.Sort(NameSorter(users))

	for _, user := range users {
		if user.Name != "" {
			addToList(cs, user.Name)
		}
	}
}

func fillChat(cs *gtk.ListStore, name string) {
	msgs, err := getHistoryByName(name)
	if err != nil {
		log.Printf("ERROR: can't get history by name %s: %s\n", name, err)
		return
	}

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
