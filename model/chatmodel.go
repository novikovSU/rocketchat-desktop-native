package model

import (
	"github.com/novikovSU/gorocket/api"
)

var (
	Chat ChatModel
)

type ChatModel struct {
	Users    map[string]UserModel
	Channels map[string]ChannelModel
	Groups   map[string]GroupModel
}

func (chat *ChatModel) AddUser(usr api.User) bool {
	if _, exists := chat.Users[usr.ID]; exists {
		return false
	}

	chat.Users[usr.ID] = UserModel{User: usr}
	return true
}

func (chat *ChatModel) RemoveUser(usr api.User) bool {
	if _, exists := chat.Users[usr.ID]; exists {
		delete(chat.Users, usr.ID)
		return true
	}

	return false
}

func (chat *ChatModel) AddChannel(ch api.Channel) bool {
	if _, exists := chat.Channels[ch.ID]; exists {
		return false
	}

	chat.Channels[ch.ID] = ChannelModel{Channel: ch}
	return true
}

func (chat *ChatModel) RemoveChannel(ch api.Channel) bool {
	if _, exists := chat.Channels[ch.ID]; exists {
		delete(chat.Channels, ch.ID)
		return true
	}

	return false
}

func (chat *ChatModel) AddGroup(gr api.Group) bool {
	if _, exists := chat.Groups[gr.ID]; exists {
		return false
	}

	chat.Groups[gr.ID] = GroupModel{Group: gr}
	return true
}

func (chat *ChatModel) RemoveGroup(gr api.Group) bool {
	if _, exists := chat.Groups[gr.ID]; exists {
		delete(chat.Groups, gr.ID)
		return true
	}

	return false
}

func init() {
	Chat = ChatModel{Users: make(map[string]UserModel), Channels: make(map[string]ChannelModel), Groups: make(map[string]GroupModel)}
}
