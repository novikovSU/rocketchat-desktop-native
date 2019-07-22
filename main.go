package main

import (
	"log"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/novikovSU/rocketchat-desktop-native/appconf"
)

const (
	iistItem = iota
	nColumns
	KEY_ENTER = 65293
)

var (
	CtrlPressed = false
)

func initList(list *gtk.TreeView) *gtk.ListStore {
	cellRenderer, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Fatal(err)
	}

	column, err := gtk.TreeViewColumnNewWithAttribute("List Items", cellRenderer, "text", 0)
	if err != nil {
		log.Fatal(err)
	}
	list.AppendColumn(column)

	store, err := gtk.ListStoreNew(glib.TYPE_STRING)
	if err != nil {
		log.Fatal(err)
	}
	list.SetModel(store)

	return store
}

func addToList(store *gtk.ListStore, text string) {
	iter := store.Append()
	store.SetValue(iter, 0, text)
}

func onChanged(selection *gtk.TreeSelection, label *gtk.Label) {
	var iter *gtk.TreeIter
	var model gtk.ITreeModel
	var ok bool
	model, iter, ok = selection.GetSelected()

	if ok {
		value, err := model.(*gtk.TreeModel).GetValue(iter, 0)
		if err != nil {
			log.Fatal(err)
		}

		text, err := value.GetString()
		if err != nil {
			log.Fatal(err)
		}

		label.SetText(text)
	}
}

func main() {
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal(err)
	}

	win.SetPosition(gtk.WIN_POS_CENTER)
	win.SetDefaultSize(800, 600)
	win.SetTitle("Rockaet.Chat Desktop native")

	mainBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		log.Fatal(err)
	}

	// -----------

	leftBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		log.Fatal(err)
	}
	//leftBox.Set

	contactList, err := gtk.TreeViewNew()
	if err != nil {
		log.Fatal(err)
	}
	contactList.SetHeadersVisible(false)
	leftBox.PackStart(contactList, true, true, 5)

	addAccBtn, err := gtk.ButtonNewWithLabel("Add account")
	if err != nil {
		log.Fatal(err)
	}
	leftBox.PackStart(addAccBtn, false, false, 0)

	addAccBtn.Connect("clicked", func() {
		appconf.GetConfig().Accounts = append(appconf.GetConfig().Accounts, appconf.Account{"http://example.com/chat", "Иванов Иван", "pswd"})
		appconf.StoreConfig()
	})
	
	connectButton, err := gtk.ButtonNewWithLabel("Connect")
	if err != nil {
		log.Fatal(err)
	}
	leftBox.PackStart(connectButton, false, false, 0)

	disconnectButton, err := gtk.ButtonNewWithLabel("Disconnect")
	if err != nil {
		log.Fatal(err)
	}
	leftBox.PackStart(disconnectButton, false, false, 0)

	// -------------

	chatList, err := gtk.TreeViewNew()
	if err != nil {
		log.Fatal(err)
	}
	chatList.SetHeadersVisible(false)

	rightBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		log.Fatal(err)
	}
	rightBox.PackStart(chatList, true, true, 5)

	textInput, err := gtk.TextViewNew()
	if err != nil {
		log.Fatal(err)
	}
	rightBox.PackStart(textInput, false, false, 5)

	// ------------------

	mainBox.PackStart(leftBox, true, true, 5)
	mainBox.PackStart(rightBox, true, true, 5)

	win.Add(mainBox)

	contactsStore := initList(contactList)
	addToList(contactsStore, "Иванов Иван")
	addToList(contactsStore, "Петров Петр")
	addToList(contactsStore, "Сидоров Сидор")
	addToList(contactsStore, "Иванов Иван")
	addToList(contactsStore, "Петров Петр")

	chatStore := initList(chatList)
	addToList(chatStore, "Вы вошли в чат.")
	addToList(chatStore, "11:12 Иванов Иван: строка 1\nстрока 2\nстрока 3")
	addToList(chatStore, "---")
	addToList(chatStore, "Иванов Иван")
	addToList(chatStore, "Петров Петр")

	/*	contactsSelection, err := contactList.GetSelection()
		if err != nil {
			log.Fatal(err)
		}

		contactsSelection.Connect("changed", func() {
			onChanged(contactsSelection, leftLabel)
		}) */

	/*chatSelection, err := chatList.GetSelection()
	if err != nil {
		log.Fatal(err)
	}

	chatSelection.Connect("changed", func() {
		onChanged(chatSelection, rightLabel)
	}) */

	textInput.Connect("key-press-event", func(tv *gtk.TextView, ev *gdk.Event) {
		keyEvent := &gdk.EventKey{ev}
		switch keyEvent.KeyVal() {
		case gdk.KEY_Control_L, gdk.KEY_Control_R:
			CtrlPressed = true
		}
	})

	textInput.Connect("key-release-event", func(tv *gtk.TextView, ev *gdk.Event) {
		keyEvent := &gdk.EventKey{ev}
		switch keyEvent.KeyVal() {
		case KEY_ENTER:
			if !CtrlPressed {
				buffer, err := tv.GetBuffer()
				if err != nil {
					log.Fatal("Unable to get buffer:", err)
				}
				start, end := buffer.GetBounds()

				text, err := buffer.GetText(start, end, true)
				if err != nil {
					log.Fatal("Unable to get text:", err)
				}
				text = strings.TrimSuffix(text, "\n")
				addToList(chatStore, text)

				buffer.SetText("")
			}
			break
		case gdk.KEY_Control_L, gdk.KEY_Control_R:
			CtrlPressed = false
		default:
			//log.Printf("Keycode: %d\n", keyEvent.KeyVal())
		}
	})

	win.Connect("destroy", gtk.MainQuit)

	win.ShowAll()
	gtk.Main()
}
