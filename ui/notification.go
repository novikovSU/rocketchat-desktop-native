package ui

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/novikovSU/rocketchat-desktop-native/config"
)

func SendNotification(title, text string) {
	notif := glib.NotificationNew(title)
	notif.SetBody(text)
	gtkApplication.SendNotification(config.AppID, notif)
}
