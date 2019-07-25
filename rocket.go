package main

import (
	"errors"
	"log"
	"time"

	bolt "go.etcd.io/bbolt"

	"github.com/novikovSU/gorocket/api"
	"github.com/novikovSU/gorocket/rest"
)

var (
	client  *rest.Client
	msgChan []api.Message
)

func getConnection() (err error) {
	client = rest.NewClient(config.Server, config.Port, config.UseTLS, config.Debug)
	err = client.Login(api.UserCredentials{Email: config.Email, Name: config.User, Password: config.Password})
	if err != nil {
		log.Fatalf("login err: %s\n", err)
	}
	return
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

func getHistoryByName(name string) ([]api.Message, error) {
	firstSymbol := string([]rune(name)[0])

	var msgs []api.Message

	switch firstSymbol {
	case hashSign:
		channel, err := getChannelByName(string([]rune(name)[1:]))
		if err != nil {
			log.Printf("ERROR: get channel id for name %s err: %s\n", name, err)
			return nil, err
		}
		msgs, err = client.Channel().History(&rest.HistoryOptions{RoomID: channel.ID})
		if err != nil {
			log.Printf("ERROR: get messages from channel %s err: %s\n", channel.Name, err)
			return nil, err
		}
		break
	case lockSign:
		group, err := getGroupByName(string([]rune(name)[1:]))
		if err != nil {
			log.Printf("ERROR: get group id for name %s err: %s\n", name, err)
			return nil, err
		}
		msgs, err = client.Groups().History(&rest.HistoryOptions{RoomID: group.ID})
		if err != nil {
			log.Printf("ERROR: get messages from group %s err: %s\n", group.Name, err)
			return nil, err
		}
		break
	default:
		user, err := getUserByName(name)
		if err != nil {
			log.Printf("ERROR: get user id for name %s err: %s\n", name, err)
			return nil, err
		}
		msgs, err = client.Im().History(&rest.HistoryOptions{RoomID: user.ID})
		if err != nil {
			log.Printf("ERROR: get messages from im %s err: %s\n", user.Name, err)
			return nil, err
		}
	}

	return msgs, nil
}

func ownMessage(c *Config, msg api.Message) bool {
	return c.User == msg.User.UserName
}

// Rewrite github.com/pyinx/gorocket/rest.GetAllMessages
// CHANGES:
//     Add dirty persistent storage for request last time (need for dedup responses after restart)
//     Add private groups and direct chats for polling (TODO)
func getAllMessages(c *rest.Client) chan []api.Message {

	msgChan := make(chan []api.Message, 1024)

	go func() {
		db, err := bolt.Open("rocket.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte("cache"))
			if err != nil {
				return err
			}
			return nil
		})
		//msgMap := make(map[string]string)

		for {
			channels, _ := c.Channel().ListJoined()
			for _, channel := range channels {
				var lastTime string
				err := db.View(func(tx *bolt.Tx) error {
					b := tx.Bucket([]byte("cache"))
					lastTime = string(b.Get([]byte(channel.Name)))
					return nil
				})
				if err != nil {
					log.Fatal(err)
				}
				msgs, err := c.Channel().History(&rest.HistoryOptions{RoomID: channel.ID, Oldest: lastTime})
				if err != nil {
					log.Printf("ERROR: get messages from channel %s err: %s\n", channel.Name, err)
				} else {
					if len(msgs) != 0 {
						err := db.Update(func(tx *bolt.Tx) error {
							b := tx.Bucket([]byte("cache"))
							err := b.Put([]byte(channel.Name), []byte(msgs[0].Timestamp.String()))
							return err
						})
						if err != nil {
							log.Fatal(err)
						}
						msgChan <- msgs
					}
				}
			}

			groups, _ := c.Groups().ListGroups()
			for _, group := range groups {
				var lastTime string
				err := db.View(func(tx *bolt.Tx) error {
					b := tx.Bucket([]byte("cache"))
					lastTime = string(b.Get([]byte(group.Name)))
					return nil
				})
				if err != nil {
					log.Fatal(err)
				}
				msgs, err := c.Groups().History(&rest.HistoryOptions{RoomID: group.ID, Oldest: lastTime})
				if err != nil {
					log.Printf("ERROR: get messages from group %s err: %s\n", group.Name, err)
				} else {
					if len(msgs) != 0 {
						err := db.Update(func(tx *bolt.Tx) error {
							b := tx.Bucket([]byte("cache"))
							err := b.Put([]byte(group.Name), []byte(msgs[0].Timestamp.String()))
							return err
						})
						if err != nil {
							log.Fatal(err)
						}
						msgChan <- msgs
					}
				}
			}

			ims, _ := c.Im().List()
			for _, im := range ims {
				var lastTime string
				err := db.View(func(tx *bolt.Tx) error {
					b := tx.Bucket([]byte("cache"))
					lastTime = string(b.Get([]byte(im.ID)))
					return nil
				})
				if err != nil {
					log.Fatal(err)
				}
				msgs, err := c.Im().History(&rest.HistoryOptions{RoomID: im.ID, Oldest: lastTime})
				if err != nil {
					log.Printf("ERROR: get messages from IMs err: %s\n", err)
				} else {
					if len(msgs) != 0 {
						err := db.Update(func(tx *bolt.Tx) error {
							b := tx.Bucket([]byte("cache"))
							err := b.Put([]byte(im.ID), []byte(msgs[0].Timestamp.String()))
							return err
						})
						if err != nil {
							log.Fatal(err)
						}
						msgChan <- msgs
					}
				}
			}
			time.Sleep(200 * time.Microsecond)
		}
	}()
	return msgChan
}
