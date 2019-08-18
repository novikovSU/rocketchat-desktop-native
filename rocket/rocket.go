package rocket

import (
	"errors"
	log "github.com/chaykin/log4go"
	"github.com/novikovSU/rocketchat-desktop-native/settings"
	"github.com/novikovSU/rocketchat-desktop-native/utils"
	"time"

	bolt "go.etcd.io/bbolt"

	"github.com/novikovSU/gorocket/api"
	"github.com/novikovSU/gorocket/realtime"
	"github.com/novikovSU/gorocket/rest"

	"github.com/novikovSU/rocketchat-desktop-native/bus"
	"github.com/novikovSU/rocketchat-desktop-native/model"
)

//TODO refactor it
const (
	hashSign = "\u0023"     // Hash sign for channels
	lockSign = "\U0001F512" // Lock sign for private groups
)

const (
	contactListUpdateInterval = 45 * time.Minute
)

var (
	client   *rest.Client
	clientRT *realtime.Client
	msgChan  []api.Message

	allHistory map[string]chatHistory
	pullChan   chan api.Message

	logger *log.Filter
)

type chatHistory struct {
	lastTime string
	msgs     []api.Message
}

func InitRocket() {
	client = initRESTConnection()
	clientRT = initRTConnection()
	loadContactListAsync()
	subscribeToMessages()
}

func init() {
	logger = utils.CreateLogger("rocket")
}

func initRESTConnection() *rest.Client {
	client := rest.NewClient(settings.Conf.Server, settings.Conf.Port, settings.Conf.UseTLS, settings.Conf.Debug)
	err := client.Login(api.UserCredentials{Email: settings.Conf.Email, Name: settings.Conf.User, Password: settings.Conf.Password})
	utils.AssertErrMsg(err, "login err: %s")

	return client
}

func initRTConnection() *realtime.Client {
	proto := "ws"
	if settings.Conf.UseTLS {
		proto = "wss"
	}
	client, _ := realtime.NewClient(proto, settings.Conf.Server, settings.Conf.Port, settings.Conf.Debug)
	client.Login(&api.UserCredentials{Email: settings.Conf.Email, Name: settings.Conf.User, Password: settings.Conf.Password})

	return client
}

//deprecated
func GetConnectionSafe(config *settings.Config) error {
	client = rest.NewClient(config.Server, config.Port, config.UseTLS, config.Debug)
	return client.Login(api.UserCredentials{Email: config.Email, Name: config.User, Password: config.Password})
}

// getChannelByName AAA
func getChannelByName(name string) (*api.Channel, error) {
	channels, err := client.Channel().List()
	if err != nil {
		logger.Error("Can't get channels list from server: %s", err)
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
		logger.Error("Can't get groups list from server: %s", err)
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
		logger.Error("Can't get users list from server: %s", err)
		return nil, err
	}

	for _, user := range users {
		if user.Name == name {
			return &user, nil
		}
	}

	return nil, errors.New("can't find user by name")
}

func getUserByUsername(username string) (*api.User, error) {
	users, err := client.Users().List()
	if err != nil {
		logger.Error("Can't get users list from server: %s", err)
		return nil, err
	}

	for _, user := range users {
		if user.UserName == username {
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
			logger.Error("Get channel id for name %s err: %s", name, err)
			return "", err
		}
		//log.Printf("Channel ID: %s\n", channel.ID)
		return channel.ID, nil
		//		break
	case lockSign:
		group, err := getGroupByName(string([]rune(name)[1:]))
		if err != nil {
			logger.Error("Get group id for name %s err: %s", name, err)
			return "", err
		}
		return group.ID, nil
		//		break
	default:
		user, err := getUserByName(name)
		if err != nil {
			logger.Error("Get user id for name %s err: %s", name, err)
			return "", err
		}
		return user.ID, nil
	}

	//	return "", nil
}

func GetRIDByName(name string) (string, error) {
	firstSymbol := string([]rune(name)[0])

	switch firstSymbol {
	case hashSign:
		channel, err := getChannelByName(string([]rune(name)[1:]))
		if err != nil {
			logger.Error("Get channel id for name %s err: %s", name, err)
			return "", err
		}
		//log.Printf("Channel ID: %s\n", channel.ID)
		return channel.ID, nil
		//		break
	case lockSign:
		group, err := getGroupByName(string([]rune(name)[1:]))
		if err != nil {
			logger.Error("Get group id for name %s err: %s", name, err)
			return "", err
		}
		return group.ID, nil
		//		break
	default:
		user, err := getUserByName(name)
		if err != nil {
			logger.Error("Get user id for name %s err: %s", name, err)
			return "", err
		}
		return model.Chat.GetMe().User.ID + user.ID, nil
	}
	//	return "", nil
}

func getHistoryByName(name string) ([]api.Message, error) {
	var msgs []api.Message
	rID, _ := getIDByName(name)
	msgs = allHistory[rID].msgs
	return msgs, nil
}

//deprecated TODO: delete?
func postByNameREST(name string, text string) {
	_, _ = client.Chat().Post(&rest.ChatPostOptions{Channel: model.Chat.ActiveContactId, Text: text})
}

func PostByNameRT(name string, text string) {
	roomID, err := GetRIDByName(name)
	if err != nil {
		logger.Error("Can't get room by name %s: %s", name, err)
		return
	}
	room := api.Channel{ID: model.Chat.ActiveContactId}
	_, err = clientRT.SendMessage(&room, text)
	if err != nil {
		logger.Error("Send message (to: %s[%s], text: %s) err: %s", name, roomID, text, err)
	}
}

func OwnMessage(msg api.Message) bool {
	return model.Chat.GetMe().User.ID == msg.User.ID
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
			logger.Error("Get messages from channel %s err: %s", channel.Name, err)
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
			logger.Error("Get messages from group %s err: %s", group.Name, err)
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
			logger.Error("Get messages from IMs err: %s", err)
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

func GetHistoryByID(id string) ([]api.Message, error) {
	msgs, err := clientRT.LoadHistory(&realtime.HistoryOptions{RoomID: id})
	if err != nil {
		logger.Error("Get messages for room with id (%s) err: %s", id, err)
		return nil, err
	}

	return msgs, nil
}

func getNewMessagesRT(c *realtime.Client) []api.Message {
	var result []api.Message

	for _, channel := range model.Chat.Channels {
		/*var lastTime string
		if hist, ok := allHistory[channel.ID]; ok {
			lastTime = hist.lastTime
		}*/
		msgs, err := c.LoadHistory(&realtime.HistoryOptions{RoomID: channel.Channel.ID})
		if err != nil {
			logger.Error("Get messages from channel %s err: %s", channel.Channel.Name, err)
		} else {
			if len(msgs) != 0 {
				chat, ok := allHistory[channel.Channel.ID]
				if ok {
					chat.lastTime = msgs[0].Timestamp.String()
					chat.msgs = append(chat.msgs, msgs...)
				} else {
					chat = chatHistory{lastTime: msgs[0].Timestamp.String(), msgs: msgs}
				}
				allHistory[channel.Channel.ID] = chat
				result = append(result, msgs...)
			}
		}
	}

	for _, group := range model.Chat.Groups {
		/*var lastTime string
		if hist, ok := allHistory[group.ID]; ok {
			lastTime = hist.lastTime
		}*/
		msgs, err := c.LoadHistory(&realtime.HistoryOptions{RoomID: group.Group.ID})
		if err != nil {
			logger.Error("Get messages from group %s err: %s", group.Group.Name, err)
		} else {
			if len(msgs) != 0 {
				chat, ok := allHistory[group.Group.ID]
				if ok {
					chat.lastTime = msgs[0].Timestamp.String()
					chat.msgs = append(chat.msgs, msgs...)
				} else {
					chat = chatHistory{lastTime: msgs[0].Timestamp.String(), msgs: msgs}
				}
				allHistory[group.Group.ID] = chat
				result = append(result, msgs...)
			}
		}
	}

	for _, user := range model.Chat.Users {
		/*var lastTime string
		if hist, ok := allHistory[user.ID]; ok {
			lastTime = hist.lastTime
		}*/
		msgs, err := c.LoadHistory(&realtime.HistoryOptions{RoomID: user.User.ID})
		if err != nil {
			logger.Error("Get messages from IMs err: %s", err)
		} else {
			if len(msgs) != 0 {
				chat, ok := allHistory[user.User.ID]
				if ok {
					chat.lastTime = msgs[0].Timestamp.String()
					chat.msgs = append(chat.msgs, msgs...)
				} else {
					chat = chatHistory{lastTime: msgs[0].Timestamp.String(), msgs: msgs}
				}
				allHistory[user.User.ID] = chat
				result = append(result, msgs...)
			}
		}
	}

	return result
}

/**
Loads async from server: channels, groups, users.
Stay active for changes. Use subscribers to get them
*/
func loadContactListAsync() {
	db, err := bolt.Open("data.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	utils.AssertErr(err)

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("cache"))
		if err != nil {
			return err
		}
		return nil
	})
	/*
		err = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("cache"))
			channels = b.Get([]byte("channels"))
			return nil
		})
	*/
	go func() {
		for {
			bus.Pub(bus.Contacts_update_started)
			loadUsers()
			loadChannels()
			loadGroups()
			bus.Pub(bus.Contacts_update_finished)

			time.Sleep(contactListUpdateInterval)
		}
	}()
}

func subscribeToMessages() {
	msgChan := make(chan api.Message, 1024)

	// Subscribe to message stream
	allMessages := api.Channel{ID: "__my_messages__"}

	//TODO handle errors
	msgChan, _ = clientRT.SubscribeToMessageStream(&allMessages)

	go func() {
		for msg := range msgChan {
			logger.Debug("CurrentChatID: %s", model.Chat.ActiveContactId)
			logger.Debug("Incoming message: %+v", msg)

			model.Chat.AddMessage(msg)
			bus.Pub(bus.Messages_new, msg)
		}
	}()
}

func loadUsers() {
	restUsers, err := client.Users().List()
	if err == nil {
		bus.Pub(bus.Rocket_users_load, restUsers)
	} else {
		logger.Error("Can't get users: %s", err)
	}
}

func loadChannels() {
	restChannels, err := client.Channel().List()
	if err == nil {
		bus.Pub(bus.Rocket_channels_load, restChannels)
	} else {
		logger.Error("Can't get channels: %s", err)
	}
}

func loadGroups() {
	restGroups, err := client.Groups().ListGroups()
	if err == nil {
		bus.Pub(bus.Rocket_groups_load, restGroups)
	} else {
		logger.Error("Can't get groups: %s", err)
	}
}
