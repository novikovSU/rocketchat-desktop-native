package bus

import (
	"github.com/asaskevich/EventBus"
	"github.com/novikovSU/rocketchat-desktop-native/config"
	"log"
)

const (
	/*
		Fires then application received new chat message
		args: api.Message
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
		Fires then application detects new user has been added to server
		args: api.User
	*/
	Contacts_users_added = "contacts.users.added"

	/*
		Fires then application detects existing user has been removed from server
		args: api.User
	*/
	Contacts_users_removed = "contacts.users.removed"

	/*
		Fires then application detects new channel has been added to server
		args: api.Channel
	*/
	Contacts_channels_added = "contacts.channels.added"

	/*
		Fires then application detects existing channel has been removed from server
		args: api.Channel
	*/
	Contacts_channels_removed = "contacts.channels.removed"

	/*
		Fires then application detects new group has been added to server
		args: api.Group
	*/
	Contacts_groups_added = "contacts.groups.added"

	/**
	Fires then application detects existing group has been removed from server
	args: api.Group
	*/
	Contacts_groups_removed = "contacts.groups.removed"

	/*
		Fires then user click on main window close button (Not implemented yet)
	*/
	Ui_mainwindow_closed = "ui.mainwindow.closed"
)

var (
	Bus = EventBus.New()
)

func Publish(topic string, args ...interface{}) {
	if config.Debug {
		log.Printf("Fire event: %s %s\n", topic, args)
	}
	Bus.Publish(topic, args...)
}

func SubscribeAsync(topic string, fn interface{}) {
	err := Bus.SubscribeAsync(topic, fn, false)
	if err != nil {
		log.Panicf("Invalid argument %s. It must be a function!", fn)
	}
}
