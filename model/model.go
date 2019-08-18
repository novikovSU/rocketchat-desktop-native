package model

import (
	"github.com/novikovSU/rocketchat-desktop-native/bus"
)

func Init(meUserName string) {
	Chat.meUserName = meUserName

	bus.Sub(bus.Rocket_users_load, Chat.loadUsers)
	bus.Sub(bus.Rocket_channels_load, Chat.loadChannels)
	bus.Sub(bus.Rocket_groups_load, Chat.loadGroups)
}
