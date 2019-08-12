package ui

import (
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"

	"github.com/novikovSU/rocketchat-desktop-native/utils"
)

var (
	GtkBuilder     gtk.Builder
	gtkApplication *gtk.Application
)

func InitUI(app *gtk.Application) {
	gtkApplication = app
	GtkBuilder = *createGtkBuilder()
}

func createGtkBuilder() *gtk.Builder {
	// Get the GtkBuilder UI definition in the glade file.
	builder, err := gtk.BuilderNewFromFile("main.glade")
	if err != nil {
		log.Panic(err)
	}

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
	if err != nil {
		log.Fatal(err)
	}
	cellRenderer.SetProperty("wrap-mode", pango.WRAP_WORD_CHAR)
	cellRenderer.SetProperty("wrap-width", 530)
	cellRenderer.SetProperty("ypad", 5)
	cellRenderer.SetProperty("xpad", 3)

	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "markup", colNum)
	if err != nil {
		log.Fatal(err)
	}
	column.SetSizing(gtk.TREE_VIEW_COLUMN_AUTOSIZE)

	return column
}

func createListStore(types ...glib.Type) *gtk.ListStore {
	store, err := gtk.ListStoreNew(types...)
	if err != nil {
		log.Fatal(err)
	}

	return store
}

// onMainWindowDestory is the callback that is linked to the
// on_main_window_destroy handler. It is not required to map this,
// and is here to simply demo how to hook-up custom callbacks.
func onMainWindowDestroy() {
	log.Println("onMainWindowDestroy")
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
	if err != nil {
		log.Panicf("Can't add close action to window %s. Cause: %s\n", id, err)
	}
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
