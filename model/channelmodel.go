package model

import "github.com/novikovSU/gorocket/api"

type ChannelModel struct {
	UnreadCount int
	Channel     api.Channel
	History     []api.Message
}

func (ch ChannelModel) GetId() string {
	return ch.Channel.ID
}

func (ch ChannelModel) GetName() string {
	return ch.Channel.Name
}

func (ch ChannelModel) GetDisplayName() string {
	return hashSign + ch.GetName()
}
