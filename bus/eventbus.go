package bus

import (
	"log"

	"github.com/asaskevich/EventBus"
	"github.com/novikovSU/rocketchat-desktop-native/config"
)

const (
	/*
		Fires then new chat message is received by rocket package
		args: api.Message
	*/
	Rocket_messages_new = "rocket.messages.new"

	/*
		Fires then users is loaded by rocket package
		args: []api.User
	*/
	Rocket_users_load = "rocket.users.load"

	/*
		Fires then channels is loaded by rocket package
		args: []api.Channel
	*/
	Rocket_channels_load = "rocket.channels.load"

	/*
		Fires then groups is loaded by rocket package
		args: []api.Group
	*/
	Rocket_groups_load = "rocket.groups.load"

	/*
		Fires then application received new chat message
		args: api.Message

		deprecated
	*/
	Messages_new = "messages.new"

	/*
		Fires then user read the chat message (Not implemented yet)
		args: api.Message
	*/
	Messages_read = "messages.read"

	/*
		Fires then application starts to load/update contact list
	*/
	Contacts_update_started = "contacts.update.started"

	/*
		Fires then application finish to load/update contact list
	*/
	Contacts_update_finished = "contacts.update.finished"

	/*
		Fires then user adds to model
		args: model.ChatModel, model.UserModel
	*/
	Model_user_added = "model.user.added"

	/*
		Fires then user removes from model
		args: model.ChatModel, model.UserModel
	*/
	Model_user_removed = "model.user.removed"

	/*
		Fires then channel adds to model
		args: model.ChatModel, model.ChannelModel
	*/
	Model_channel_added = "model.channel.added"

	/*
		Fires then channel removes from model
		args: model.ChatModel, model.ChannelModel
	*/
	Model_channel_removed = "model.channel.removed"

	/*
		Fires then group adds to model
		args: model.ChatModel, model.GroupModel
	*/
	Model_group_added = "model.group.added"

	/*
		Fires then group removes from model
		args: model.ChatModel, model.GroupModel
	*/
	Model_group_removed = "model.group.removed"

	/*
		Fires then user click on main window close button (Not implemented yet)
	*/
	Ui_mainwindow_closed = "ui.mainwindow.closed"
)

var (
	b = EventBus.New()
)

// Pub AAA
func Pub(topic string, args ...interface{}) {
	if config.Debug {
		log.Printf("Fire event: %s %s\n", topic, args)
	}
	b.Publish(topic, args...)
}

// Sub AAA
func Sub(topic string, fn interface{}) {
	err := b.SubscribeAsync(topic, fn, false)
	if err != nil {
		log.Panicf("Invalid argument %s. It must be a function!: %s\n", fn, err)
	}
}
