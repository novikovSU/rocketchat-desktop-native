package model

import (
	"github.com/novikovSU/rocketchat-desktop-native/bus"
)

func Init(meUserName string) {
	Chat.meUserName = meUserName

	bus.Sub(bus.Rocket_users_load, Chat.loadUsers)
	bus.Sub(bus.Rocket_channels_load, Chat.loadChannels)
	bus.Sub(bus.Rocket_groups_load, Chat.loadGroups)
	bus.Sub(bus.Rocket_messages_new, Chat.addMessage)
	bus.Sub(bus.Ui_contacts_selected, onContactsSelected)
}

func onContactsSelected(name string) {
	id := Chat.GetIdByName(name)
	if len(id) > 0 {
		Chat.clearUnreadCount(id)
	}
}
