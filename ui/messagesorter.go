package ui

import "github.com/novikovSU/gorocket/api"

/*
Sorts messages by it timestamp
*/
type MessageSorter []api.Message

func (a MessageSorter) Len() int {
	return len(a)
}

func (a MessageSorter) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a MessageSorter) Less(i, j int) bool {
	return a[i].Timestamp.Before(*a[j].Timestamp)
}
