package ui

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/novikovSU/rocketchat-desktop-native/utils"
)

const (
	ChatListColumn = 0
)

var (
	chatStore *gtk.ListStore
)

func InitChatListControl() {
	chatList, store := createChatListTreeView()
	chatStore = store

	rightScrolledWindow := GetScrolledWindow("right_scrolled_window")
	utils.Safe(chatList.Connect("size-allocate", func() { chatListAutoScroll(rightScrolledWindow) }))
}

func createChatListTreeView() (*gtk.TreeView, *gtk.ListStore) {
	chatList := GetTreeView("chat_list")

	chatList.AppendColumn(createTextViewTextColumn("items", ChatListColumn))

	chatListStore := createListStore(glib.TYPE_STRING)
	chatList.SetModel(chatListStore)

	return chatList, chatListStore
}

/*
Autoscroll of chatList function
*/
func chatListAutoScroll(wnd *gtk.ScrolledWindow) {
	adj := wnd.GetVAdjustment()
	adj.SetValue(adj.GetUpper() - adj.GetPageSize())
}
