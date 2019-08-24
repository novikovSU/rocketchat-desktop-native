package ui

import (
	"github.com/gotk3/gotk3/gtk"
	"sort"

	"github.com/novikovSU/rocketchat-desktop-native/model"
	"github.com/novikovSU/rocketchat-desktop-native/rocket"
)

func refreshChatHistory(cs *gtk.ListStore, name string) {
	model.Chat.ActiveContactId, _ = rocket.GetRIDByName(name)

	msgs, err := rocket.GetHistoryByID(model.Chat.ActiveContactId)
	if err != nil {
		logger.Error("can't get history by name %s: %s\n", name, err)
		return
	}

	sort.Sort(MessageSorter(msgs))

	cs.Clear()
	for _, msg := range msgs {
		addTextMessageToActiveChat(&msg)
	}
}
