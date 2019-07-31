package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/novikovSU/rocketchat-desktop-native/bus"

	"github.com/novikovSU/gorocket/api"
	"github.com/novikovSU/gorocket/realtime"
	"github.com/novikovSU/gorocket/rest"
)

const (
	contactListUpdateInterval = 45 * time.Second
)

var (
	client   *rest.Client
	clientRT *realtime.Client
	msgChan  []api.Message

	me            *api.User
	allHistory    map[string]chatHistory
	pullChan      chan api.Message
	currentChatID string

	users    = make(map[string]*api.User)
	channels = make(map[string]*api.Channel)
	groups   = make(map[string]*api.Group)
)

type chatHistory struct {
	lastTime string
	msgs     []api.Message
}

func initRocket() {
	client = initRESTConnection()
	clientRT = initRTConnection()
	loadContactListAsync()
}

func initRESTConnection() *rest.Client {
	client := rest.NewClient(config.RestServer, config.RestPort, config.UseTLS, config.Debug)
	err := client.Login(api.UserCredentials{Email: config.Email, Name: config.User, Password: config.Password})
	if err != nil {
		log.Fatalf("login err: %s\n", err)
	}

	return client
}

func initRTConnection() *realtime.Client {
	client, _ := realtime.NewClient("ws", config.RTServer, config.RTPort, config.Debug)
	client.Login(&api.UserCredentials{Email: config.Email, Name: config.User, Password: config.Password})

	return client
}

//deprecated
func getConnection() (err error) {
	err = getConnectionSafe(config)
	if err != nil {
		log.Fatalf("login err: %s\n", err)
	}
	return
}

//deprecated
func getConnectionSafe(config *Config) error {
	client = rest.NewClient(config.RestServer, config.RestPort, config.UseTLS, config.Debug)
	return client.Login(api.UserCredentials{Email: config.Email, Name: config.User, Password: config.Password})
}

// getChannelByName AAA
func getChannelByName(name string) (*api.Channel, error) {
	channels, err := client.Channel().List()
	if err != nil {
		log.Printf("ERROR: can't get channels list from server: %s\n", err)
		return nil, err
	}

	for _, channel := range channels {
		if channel.Name == name {
			return &channel, nil
		}
	}

	return nil, errors.New("can't find channel by name")
}

// getGroupByName AAA
func getGroupByName(name string) (*api.Group, error) {
	groups, err := client.Groups().ListGroups()
	if err != nil {
		log.Printf("ERROR: can't get groups list from server: %s\n", err)
		return nil, err
	}

	for _, group := range groups {
		if group.Name == name {
			return &group, nil
		}
	}

	return nil, errors.New("can't find group by name")
}

// getUserByName AAA
func getUserByName(name string) (*api.User, error) {
	users, err := client.Users().List()
	if err != nil {
		log.Printf("ERROR: can't get users list from server: %s\n", err)
		return nil, err
	}

	for _, user := range users {
		if user.Name == name {
			return &user, nil
		}
	}

	return nil, errors.New("can't find user by name")
}

func getIDByName(name string) (string, error) {
	firstSymbol := string([]rune(name)[0])

	switch firstSymbol {
	case hashSign:
		channel, err := getChannelByName(string([]rune(name)[1:]))
		if err != nil {
			log.Printf("ERROR: get channel id for name %s err: %s\n", name, err)
			return "", err
		}
		//log.Printf("Channel ID: %s\n", channel.ID)
		return channel.ID, nil
		//		break
	case lockSign:
		group, err := getGroupByName(string([]rune(name)[1:]))
		if err != nil {
			log.Printf("ERROR: get group id for name %s err: %s\n", name, err)
			return "", err
		}
		return group.ID, nil
		//		break
	default:
		user, err := getUserByName(name)
		if err != nil {
			log.Printf("ERROR: get user id for name %s err: %s\n", name, err)
			return "", err
		}
		return user.ID, nil
	}

	//	return "", nil
}

func getHistoryByName(name string) ([]api.Message, error) {
	var msgs []api.Message
	rID, _ := getIDByName(name)
	msgs = allHistory[rID].msgs
	return msgs, nil
}

func postByName(name string, text string) {
	roomID, err := getIDByName(name)
	if err != nil {
		log.Printf("can't get room by name %s: %s\n", name, err)
		return
	}
	room := api.Channel{ID: roomID}
	_, err = clientRT.SendMessage(&room, text)
	if err != nil {
		if config.Debug {
			log.Printf("send message err: %s\n", err)
		}
	}
}

func ownMessage(c *Config, msg api.Message) bool {
	return c.User == msg.User.UserName
}

func getNewMessages(c *rest.Client) []api.Message {

	var result []api.Message

	channels, _ := c.Channel().List()
	for _, channel := range channels {
		var lastTime string
		if hist, ok := allHistory[channel.ID]; ok {
			lastTime = hist.lastTime
		}
		msgs, err := c.Channel().History(&rest.HistoryOptions{RoomID: channel.ID, Oldest: lastTime})
		if err != nil {
			log.Printf("ERROR: get messages from channel %s err: %s\n", channel.Name, err)
		} else {
			if len(msgs) != 0 {
				chat, ok := allHistory[channel.ID]
				if ok {
					chat.lastTime = msgs[0].Timestamp.String()
					chat.msgs = append(chat.msgs, msgs...)
				} else {
					chat = chatHistory{lastTime: msgs[0].Timestamp.String(), msgs: msgs}
				}
				allHistory[channel.ID] = chat
				result = append(result, msgs...)
			}
		}
	}

	groups, _ := c.Groups().ListGroups()
	for _, group := range groups {
		var lastTime string
		if hist, ok := allHistory[group.ID]; ok {
			lastTime = hist.lastTime
		}
		msgs, err := c.Groups().History(&rest.HistoryOptions{RoomID: group.ID, Oldest: lastTime})
		if err != nil {
			log.Printf("ERROR: get messages from group %s err: %s\n", group.Name, err)
		} else {
			if len(msgs) != 0 {
				chat, ok := allHistory[group.ID]
				if ok {
					chat.lastTime = msgs[0].Timestamp.String()
					chat.msgs = append(chat.msgs, msgs...)
				} else {
					chat = chatHistory{lastTime: msgs[0].Timestamp.String(), msgs: msgs}
				}
				allHistory[group.ID] = chat
				result = append(result, msgs...)
			}
		}
	}

	users, _ := c.Users().List()
	for _, user := range users {
		var lastTime string
		if hist, ok := allHistory[user.ID]; ok {
			lastTime = hist.lastTime
		}
		msgs, err := c.Im().History(&rest.HistoryOptions{RoomID: user.ID, Oldest: lastTime})
		if err != nil {
			log.Printf("ERROR: get messages from IMs err: %s\n", err)
		} else {
			if len(msgs) != 0 {
				chat, ok := allHistory[user.ID]
				if ok {
					chat.lastTime = msgs[0].Timestamp.String()
					chat.msgs = append(chat.msgs, msgs...)
				} else {
					chat = chatHistory{lastTime: msgs[0].Timestamp.String(), msgs: msgs}
				}
				allHistory[user.ID] = chat
				result = append(result, msgs...)
			}
		}
	}

	return result
}

func subscribeToUpdates(c *rest.Client, freq time.Duration) chan api.Message {
	msgChan := make(chan api.Message, 1024)

	// Subscribe to message stream
	allMessages := api.Channel{ID: "__my_messages__"}
	msgChan, _ = clientRT.SubscribeToMessageStream(&allMessages)

	go func() {
		var msg api.Message

		for {
			msg = <-msgChan
			//log.Printf("CurrentChatID: %s\n", currentChatID)
			//log.Printf("Incoming message: %+v\n", msg)

			chat, ok := allHistory[msg.ChannelID]
			if ok {
				chat.lastTime = msg.Timestamp.String()
				chat.msgs = append(chat.msgs, msg)
			} else {
				msgs := make([]api.Message, 1)
				msgs = append(msgs, msg)
				chat = chatHistory{lastTime: msg.Timestamp.String(), msgs: msgs}
			}
			allHistory[msg.ChannelID] = chat

			if msg.ChannelID == currentChatID || msg.ChannelID == currentChatID+currentChatID {
				text := strings.Replace(msg.Text, "&nbsp;", "", -1)
				text = strings.Replace(text, "<", "", -1)
				text = strings.Replace(text, ">", "", -1)
				//log.Printf("Text: %s\n", text)
				text = fmt.Sprintf("<b>%s</b> <i>%s</i>\n%s", msg.User.Name, msg.Timestamp.Format("2006-01-02 15:04:05"), text)
				addToList(chatStore, text)
			}

			bus.Publish("messages.new", msg)
		}
	}()

	return msgChan
}

/**
Loads async from server: channels, groups, users.
Stay active for changes. Use subscribers to get them
*/
func loadContactListAsync() {
	go func() {
		for {
			loadUsers()
			loadChannels()
			loadGroups()
			time.Sleep(contactListUpdateInterval)
		}
	}()
}

func loadUsers() {
	restUsers, err := client.Users().List()
	if err != nil {
		log.Printf("Can't get users: %s\n", err)
	}

	for _, existsUser := range users {
		if !containsUser(&restUsers, existsUser) {
			removeUser(existsUser)
		}
	}

	for _, restUser := range restUsers {
		if users[restUser.ID] == nil {
			addUser(&restUser)
		}
	}
}

func loadChannels() {
	restChannels, err := client.Channel().List()
	if err != nil {
		log.Printf("Can't get channels: %s\n", err)
	}

	for _, existsChannel := range channels {
		if !containsChannel(&restChannels, existsChannel) {
			removeChannel(existsChannel)
		}
	}

	for _, restChannel := range restChannels {
		if channels[restChannel.ID] == nil {
			addChannel(&restChannel)
		}
	}
}

func loadGroups() {
	restGroups, err := client.Groups().ListGroups()
	if err != nil {
		log.Printf("Can't get groups: %s\n", err)
	}

	for _, existsGroup := range groups {
		if !containsGroup(&restGroups, existsGroup) {
			removeGroup(existsGroup)
		}
	}

	for _, restGroup := range restGroups {
		if groups[restGroup.ID] == nil {
			addGroup(&restGroup)
		}
	}
}

func addUser(user *api.User) {
	users[user.ID] = user
	bus.Publish("contacts.users.added", user)
}

func removeUser(user *api.User) {
	delete(users, user.ID)
	bus.Publish("contacts.users.removed", user)
}

func addChannel(channel *api.Channel) {
	channels[channel.ID] = channel
	bus.Publish("contacts.channels.added", channel)
}

func removeChannel(channel *api.Channel) {
	delete(channels, channel.ID)
	bus.Publish("contacts.channels.removed", channel)
}

func addGroup(group *api.Group) {
	groups[group.ID] = group
	bus.Publish("contacts.groups.added", group)
}

func removeGroup(group *api.Group) {
	delete(groups, group.ID)
	bus.Publish("contacts.groups.removed", group)
}

/*---------------------------------------------------------------------------
Very common and dummy functions
TODO codgen?
---------------------------------------------------------------------------*/

func containsUser(users *[]api.User, cmpUser *api.User) bool {
	for _, user := range *users {
		if strings.Compare(user.ID, cmpUser.ID) == 0 {
			return true
		}

	}
	return false
}

func containsChannel(channels *[]api.Channel, cmpChannel *api.Channel) bool {
	for _, channel := range *channels {
		if strings.Compare(channel.ID, cmpChannel.ID) == 0 {
			return true
		}

	}
	return false
}

func containsGroup(groups *[]api.Group, cmpGroup *api.Group) bool {
	for _, group := range *groups {
		if strings.Compare(group.ID, cmpGroup.ID) == 0 {
			return true
		}

	}
	return false
}
