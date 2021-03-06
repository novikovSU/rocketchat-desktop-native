package ui

/*
   Main send message textInput control handlers
*/

import (
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"github.com/novikovSU/rocketchat-desktop-native/rocket"
	"github.com/novikovSU/rocketchat-desktop-native/utils"
)

const (
	keyEnter = 65293
)

var (
	sendMsgPreventedKeyHeld int
	preventSendKeys         = [...]uint{gdk.KEY_Control_L, gdk.KEY_Control_R, gdk.KEY_Shift_L, gdk.KEY_Shift_R}
)

func InitSendMsgControl() {
	textInput := GetTextView("text_input")

	utils.Safe(textInput.Connect("key-press-event", onSendMsgKeyPress))
	utils.Safe(textInput.Connect("key-release-event", onSendMsgKeyUp))
}

/*
Handler for send message textInput keyPress event. Stores held keys for handle it in keyUp handler
*/
func onSendMsgKeyPress(tv *gtk.TextView, event *gdk.Event) {
	if isPreventedKeyHeld(&gdk.EventKey{event}) {
		sendMsgPreventedKeyHeld++
	}
}

/*
Handler for send message textInput keyUp event. Send message to chat. Supports special key-combinations for multi-line messages
*/
func onSendMsgKeyUp(tv *gtk.TextView, event *gdk.Event) {
	key := gdk.EventKey{event}
	if isPreventedKeyHeld(&key) {
		sendMsgPreventedKeyHeld--
	} else if sendMsgPreventedKeyHeld <= 0 && key.KeyVal() == keyEnter {
		buf := GetTextViewBuffer(tv)
		msgText := getText(buf)

		if utils.IsNotBlankString(msgText) {
			selectionText := GetTreeViewSelectionVal(contactList, 0)
			buf.SetText("")
			rocket.PostByNameRT(selectionText, msgText)
		}
	}
}

func isPreventedKeyHeld(key *gdk.EventKey) bool {
	for _, k := range preventSendKeys {
		if key.KeyVal() == k {
			return true
		}
	}
	return false
}

func getText(buf *gtk.TextBuffer) string {
	start, end := buf.GetBounds()

	inputText, err := buf.GetText(start, end, true)
	utils.AssertErrMsg(err, "Unable to get text: %s")

	return strings.TrimSuffix(inputText, "\n")
}
