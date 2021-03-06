package model

import (
	"github.com/novikovSU/gorocket/api"
)

type ChannelModel struct {
	UnreadCount int
	Channel     api.Channel
	History     []api.Message
}

func (ch *ChannelModel) GetId() string {
	return ch.Channel.ID
}

func (ch *ChannelModel) GetName() string {
	return ch.Channel.Name
}

func (ch *ChannelModel) String() string {
	return hashSign + ch.GetName()
}

func (ch *ChannelModel) GetUnreadCount() int {
	return ch.UnreadCount
}

func (ch *ChannelModel) updateUnreadCount(change int) {
	ch.UnreadCount += change
}

func (ch *ChannelModel) clearUnreadCount() {
	ch.UnreadCount = 0
}

func ChannelsToModels(channels []*ChannelModel) []IContactModel {
	models := make([]IContactModel, 0, len(channels))
	for _, ch := range channels {
		models = append(models, ch)
	}

	return models
}

func ChannelsMapToModels(channels map[string]*ChannelModel) []IContactModel {
	models := make([]IContactModel, 0, len(channels))
	for _, ch := range channels {
		models = append(models, ch)
	}

	return models
}
