package model

import (
	"strings"

	"github.com/novikovSU/gorocket/api"
)

var (
	Chat ChatModel
)

type ChatModel struct {
	ActiveContactId string
	Users           map[string]*UserModel
	Channels        map[string]*ChannelModel
	Groups          map[string]*GroupModel
}

func (chat *ChatModel) AddUser(usr api.User) bool {
	if _, exists := chat.Users[usr.ID]; exists {
		return false
	}

	chat.Users[usr.ID] = &UserModel{User: usr}
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

	chat.Channels[ch.ID] = &ChannelModel{Channel: ch}
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

	chat.Groups[gr.ID] = &GroupModel{Group: gr}
	return true
}

func (chat *ChatModel) RemoveGroup(gr api.Group) bool {
	if _, exists := chat.Groups[gr.ID]; exists {
		delete(chat.Groups, gr.ID)
		return true
	}

	return false
}

func (chat *ChatModel) GetModelById(id string) IContactModel {
	if model, exists := chat.Users[id]; exists {
		return model
	}

	if model, exists := chat.Groups[id]; exists {
		return model
	}

	if model, exists := chat.Channels[id]; exists {
		return model
	}

	return nil
}

//TODO do we need to store total count as variable?
func (chat *ChatModel) GetTotalUnreadCount() int {
	count := 0
	for _, user := range chat.Users {
		count += user.UnreadCount
	}

	for _, user := range chat.Channels {
		count += user.UnreadCount
	}

	for _, user := range chat.Groups {
		count += user.UnreadCount
	}

	return count
}

func (chat *ChatModel) GetUnreadCount(id string) int {
	model := chat.GetModelById(id)
	if model == nil {
		return 0
	}
	return model.GetUnreadCount()
}

func (chat *ChatModel) ClearUnreadCount(id string) {
	model := chat.GetModelById(id)
	if model == nil {
		return
	}
	model.ClearUnreadCount()
}

func (chat *ChatModel) AddMessage(msg api.Message, me *api.User) {
	model := chat.GetModelById(strings.Replace(msg.ChannelID, me.ID, "", 1))
	if model != nil {
		//TODO handle activeContactId
		// add to model variable: type. It should means strategy of unread count: differs when chat window visible or hide
		model.UpdateUnreadCount(1)
	}
}

func init() {
	Chat = ChatModel{Users: make(map[string]*UserModel), Channels: make(map[string]*ChannelModel), Groups: make(map[string]*GroupModel)}
}
