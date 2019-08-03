package model

import (
	"github.com/novikovSU/gorocket/api"
)

const (
	hashSign = "\u0023"     // Hash sign for channels
	lockSign = "\U0001F512" // Lock sign for private groups
)

type IContactModel interface {
	GetId() string
	GetName() string
	GetDisplayName() string
}

type UserModel struct {
	UnreadCount int
	User        api.User
	History     []api.Message
}

func (u UserModel) GetId() string {
	return u.User.ID
}

func (u UserModel) GetName() string {
	return u.User.Name
}

func (u UserModel) GetDisplayName() string {
	return u.GetName()
}

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

type GroupModel struct {
	UnreadCount int
	Group       api.Group
	History     []api.Message
}

func (g GroupModel) GetId() string {
	return g.Group.ID
}

func (g GroupModel) GetName() string {
	return g.Group.Name
}

func (g GroupModel) GetDisplayName() string {
	return lockSign + g.GetName()
}

type ChatModel struct {
	Users    map[string]UserModel
	Channels map[string]ChannelModel
	Groups   map[string]GroupModel
}

var (
	Chat ChatModel
)

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
