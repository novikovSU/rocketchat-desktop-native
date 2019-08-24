package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"

	log "github.com/chaykin/log4go"

	"github.com/novikovSU/gorocket/api"
	"github.com/novikovSU/rocketchat-desktop-native/bus"
	"github.com/novikovSU/rocketchat-desktop-native/model"
	"github.com/novikovSU/rocketchat-desktop-native/rocket"
	"github.com/novikovSU/rocketchat-desktop-native/utils"
)

var (
	GtkBuilder     gtk.Builder
	gtkApplication *gtk.Application

	logger *log.Filter
)

func InitUI(app *gtk.Application) {
	gtkApplication = app
	GtkBuilder = *createGtkBuilder()

	SendNotification("Rocket.Chat Desktop native", "application activated")
}

func init() {
	logger = utils.CreateLogger("ui")
}
func createGtkBuilder() *gtk.Builder {
	// Get the GtkBuilder UI definition in the glade file.
	builder, err := gtk.BuilderNewFromFile("main.glade")
	utils.AssertErr(err)

	// Map the handlers to callback functions, and connect the signals to the Builder.
	signals := map[string]interface{}{
		//TODO do we need event for this?
		"on_main_window_destroy": onMainWindowDestroy,
	}
	builder.ConnectSignals(signals)

	return builder
}

func createTextViewTextColumn(title string, colNum int) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererTextNew()
	utils.AssertErr(err)

	utils.AssertErr(cellRenderer.SetProperty("wrap-mode", pango.WRAP_WORD_CHAR))
	utils.AssertErr(cellRenderer.SetProperty("wrap-width", 530))
	utils.AssertErr(cellRenderer.SetProperty("ypad", 5))
	utils.AssertErr(cellRenderer.SetProperty("xpad", 3))

	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "markup", colNum)
	utils.AssertErr(err)

	column.SetSizing(gtk.TREE_VIEW_COLUMN_AUTOSIZE)

	return column
}

func createListStore(types ...glib.Type) *gtk.ListStore {
	store, err := gtk.ListStoreNew(types...)
	utils.AssertErr(err)

	return store
}

// onMainWindowDestory is the callback that is linked to the
// on_main_window_destroy handler. It is not required to map this,
// and is here to simply demo how to hook-up custom callbacks.
func onMainWindowDestroy() {
	logger.Trace("onMainWindowDestroy")
}

/*
Creates GTK window with default close action
*/
func CreateWindow(id string) *gtk.Window {
	wnd := GetWindow(id)

	wndCloseAction := glib.SimpleActionNew("close", nil)
	_, err := wndCloseAction.Connect("activate", func() {
		wnd.Close()
	})
	utils.AssertErrMsg(err, fmt.Sprintf("Can't add close action to window %s. Cause: ", id)+"%s")

	gtkApplication.AddAction(wndCloseAction)

	return wnd
}

/**
Returns the string value of specified column of treeview row selected
*/
func GetTreeViewSelectionVal(tv *gtk.TreeView, column int) string {
	selection, err := tv.GetSelection()
	utils.AssertErr(err)

	model, iter, ok := selection.GetSelected()
	if ok {
		value, err := model.(*gtk.TreeModel).GetValue(iter, column)
		utils.AssertErr(err)

		val, err := value.GetString()
		utils.AssertErr(err)

		return val
	}

	return ""
}

func InitSubscribers() {
	bus.Sub(bus.Rocket_contacts_update_finished, func() {
		//TODO react to remove users/groups/channels: if it happens, what should we do with selection, for example?
		cs := contactsStore
		cs.Clear()

		addContactsToList(cs, getSortedChannels())
		addContactsToList(cs, getSortedGroups())
		addContactsToList(cs, getSortedUsers())
	})

	bus.Sub(bus.Model_messages_received, onMessageReceived)
	bus.Sub(bus.Model_messages_received, sendMessageReceivedNotification)
	bus.Sub(bus.Model_unreadCounters_updated, onUnreadCounterUpdated)
}

func getSortedChannels() []model.IContactModel {
	contacts := model.ChannelsMapToModels(model.Chat.Channels)
	sort.Sort(ContactsSorter(contacts))

	return contacts
}

func getSortedGroups() []model.IContactModel {
	contacts := model.GroupsMapToModels(model.Chat.Groups)
	sort.Sort(ContactsSorter(contacts))

	return contacts
}

func getSortedUsers() []model.IContactModel {
	contacts := model.UsersMapToModels(model.Chat.Users)
	sort.Sort(ContactsSorter(contacts))

	return contacts
}

func onMessageReceived(chat *model.ChatModel, modelId string, msg api.Message) {
	senderId := chat.GetSenderId(&msg)
	if strings.Compare(senderId, chat.ActiveContactId) == 0 {
		addTextMessageToActiveChat(&msg)
	}
}

func sendMessageReceivedNotification(chat *model.ChatModel, modelId string, msg api.Message) {
	logger.Debug("Prepare to notificate")
	if rocket.OwnMessage(msg) {
		return
	}
	if mainWindowIsFocused && msg.ChannelID == chat.ActiveContactId {
		return
	}

	notifTitle := fmt.Sprintf("%s (%s)", msg.User.Name, msg.User.UserName)
	SendNotification(notifTitle, msg.Text)
}

func addTextMessageToActiveChat(msg *api.Message) {
	text := strings.Replace(msg.Text, "&nbsp;", "", -1)
	text = strings.Replace(text, "<", "", -1)
	text = strings.Replace(text, ">", "", -1)
	text = fmt.Sprintf("<b>%s</b> <i>%s</i>\n%s", msg.User.Name, msg.Timestamp.Format("2006-01-02 15:04:05"), text)
	addToList(chatStore, text)
}

func onUnreadCounterUpdated(chat *model.ChatModel, modelId string) {
	contactModel := chat.GetModelById(modelId)
	if contactModel != nil {
		for iter, exists := contactsStore.GetIterFirst(); exists; exists = contactsStore.IterNext(iter) {
			val, err := contactsStore.GetValue(iter, ContactListNameColumn)
			utils.AssertErr(err)

			strVal, err := val.GetString()
			utils.AssertErr(err)

			if strings.Compare(strVal, contactModel.String()) == 0 {
				utils.AssertErr(contactsStore.SetValue(iter, ContactListUnreadCountColumn, getUnreadCount(&contactModel)))
				break
			}
		}
	}
}

func addContactsToList(cs *gtk.ListStore, contacts []model.IContactModel) {
	for _, contact := range contacts {
		if contact.GetName() != "" {
			iter := cs.Append()
			utils.AssertErr(cs.SetValue(iter, ContactListNameColumn, contact.String()))
			utils.AssertErr(cs.SetValue(iter, ContactListUnreadCountColumn, getUnreadCount(&contact)))
		}
	}
}

func addToList(store *gtk.ListStore, text string) {
	utils.AssertErr(store.SetValue(store.Append(), 0, text))
}
