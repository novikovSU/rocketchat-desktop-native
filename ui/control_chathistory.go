package ui

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/gotk3/gotk3/gtk"

	"github.com/novikovSU/gorocket/api"
	"github.com/novikovSU/rocketchat-desktop-native/model"
	"github.com/novikovSU/rocketchat-desktop-native/rocket"
)

func refreshChatHistory(cs *gtk.ListStore, name string) {
	model.Chat.ActiveContactId, _ = rocket.GetRIDByName(name)

	msgs, err := rocket.GetHistoryByID(model.Chat.ActiveContactId)
	if err != nil {
		log.Printf("ERROR: can't get history by name %s: %s\n", name, err)
		return
	}

	sort.Sort(MessageSorter(msgs))

	cs.Clear()
	for _, msg := range msgs {
		drawMessage(cs, msg)
	}
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

func addToList(store *gtk.ListStore, text string) {
	iter := store.Append()
	store.SetValue(iter, 0, text)
}
