package model

const (
	hashSign = "\u0023"     // Hash sign for channels
	lockSign = "\U0001F512" // Lock sign for private groups
)

type IContactModel interface {
	GetId() string
	GetName() string
	String() string
	GetUnreadCount() int
	UpdateUnreadCount(change int)
	ClearUnreadCount()
}
