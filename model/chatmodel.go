package model

import (
	"strings"
	"sync"

	"github.com/novikovSU/gorocket/api"
	"github.com/novikovSU/rocketchat-desktop-native/bus"
)

var (
	Chat ChatModel

	meMutex sync.Mutex
	meCond  = sync.NewCond(&meMutex)
)

type ChatModel struct {
	ActiveContactId string
	Users           map[string]*UserModel
	Channels        map[string]*ChannelModel
	Groups          map[string]*GroupModel
	meUserName      string
	me              *UserModel
}

func (chat *ChatModel) GetMe() *UserModel {
	if chat.me == nil {
		meMutex.Lock()
		if chat.me == nil {
			meCond.Wait()
		}
		meMutex.Unlock()
	}

	return chat.me
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

func (chat *ChatModel) GetSenderId(msg *api.Message) string {
	meId := chat.GetMe().User.ID
	return strings.Replace(msg.ChannelID, meId, "", 1)
}

func (chat *ChatModel) addMessage(msg api.Message) {
	model := chat.GetModelById(chat.GetSenderId(&msg))
	if model != nil {
		//TODO handle activeContactId
		// add to model variable: type. It should means strategy of unread count: differs when chat window visible or hide
		model.updateUnreadCount(1)

		bus.Pub(bus.Model_unreadCounters_updated, chat, model.GetId())
		bus.Pub(bus.Model_messages_received, chat, model.GetId(), msg)
	}
}

func (chat *ChatModel) loadUsers(users []api.User) {
	for _, existsUser := range chat.Users {
		if !containsUser(&users, existsUser) {
			chat.removeUser(existsUser.User)
			bus.Pub(bus.Model_user_removed, chat, existsUser)
		}
	}

	for _, restUser := range users {
		addedUser := chat.addUser(restUser)
		if addedUser != nil {
			bus.Pub(bus.Model_user_added, chat, addedUser)
		}
	}
}

func (chat *ChatModel) loadChannels(channels []api.Channel) {
	for _, existsChannel := range chat.Channels {
		if !containsChannel(&channels, existsChannel) {
			chat.removeChannel(existsChannel.Channel)
			bus.Pub(bus.Model_channel_removed, chat, existsChannel)
		}
	}

	for _, restChannel := range channels {
		addedChannel := chat.addChannel(restChannel)
		if addedChannel != nil {
			bus.Pub(bus.Model_channel_added, chat, addedChannel)
		}
	}
}

func (chat *ChatModel) loadGroups(groups []api.Group) {
	for _, existsGroup := range chat.Groups {
		if !containsGroup(&groups, existsGroup) {
			chat.removeGroup(existsGroup.Group)
			bus.Pub(bus.Model_group_removed, chat, existsGroup)
		}
	}

	for _, restGroup := range groups {
		addedGroup := chat.addGroup(restGroup)
		if addedGroup != nil {
			bus.Pub(bus.Model_group_added, chat, addedGroup)
		}
	}
}

func (chat *ChatModel) addUser(usr api.User) *UserModel {
	if _, exists := chat.Users[usr.ID]; exists {
		return nil
	}

	model := &UserModel{User: usr}
	chat.Users[usr.ID] = model

	if strings.Compare(usr.UserName, chat.meUserName) == 0 {
		meMutex.Lock()

		chat.me = model

		meMutex.Unlock()
		meCond.Signal()
	}

	return model
}

func (chat *ChatModel) removeUser(usr api.User) bool {
	if _, exists := chat.Users[usr.ID]; exists {
		delete(chat.Users, usr.ID)
		return true
	}

	return false
}

func (chat *ChatModel) addChannel(ch api.Channel) *ChannelModel {
	if _, exists := chat.Channels[ch.ID]; exists {
		return nil
	}

	model := &ChannelModel{Channel: ch}
	chat.Channels[ch.ID] = model
	return model
}

func (chat *ChatModel) removeChannel(ch api.Channel) bool {
	if _, exists := chat.Channels[ch.ID]; exists {
		delete(chat.Channels, ch.ID)
		return true
	}

	return false
}

func (chat *ChatModel) addGroup(gr api.Group) *GroupModel {
	if _, exists := chat.Groups[gr.ID]; exists {
		return nil
	}

	model := &GroupModel{Group: gr}
	chat.Groups[gr.ID] = model
	return model
}

func (chat *ChatModel) removeGroup(gr api.Group) bool {
	if _, exists := chat.Groups[gr.ID]; exists {
		delete(chat.Groups, gr.ID)
		return true
	}

	return false
}

func init() {
	Chat = ChatModel{Users: make(map[string]*UserModel), Channels: make(map[string]*ChannelModel), Groups: make(map[string]*GroupModel)}
}

/*---------------------------------------------------------------------------
Very common and dummy functions
TODO codgen?
---------------------------------------------------------------------------*/

func containsUser(users *[]api.User, cmpUser *UserModel) bool {
	for _, user := range *users {
		if strings.Compare(user.ID, cmpUser.User.ID) == 0 {
			return true
		}

	}
	return false
}

func containsChannel(channels *[]api.Channel, cmpChannel *ChannelModel) bool {
	for _, channel := range *channels {
		if strings.Compare(channel.ID, cmpChannel.Channel.ID) == 0 {
			return true
		}

	}
	return false
}

func containsGroup(groups *[]api.Group, cmpGroup *GroupModel) bool {
	for _, group := range *groups {
		if strings.Compare(group.ID, cmpGroup.Group.ID) == 0 {
			return true
		}

	}
	return false
}
