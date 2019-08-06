package ui

import "github.com/novikovSU/rocketchat-desktop-native/model"

/*
Sorts contactModels by it names
*/
type ContactsSorter []model.IContactModel

func (c ContactsSorter) Len() int {
	return len(c)
}

func (c ContactsSorter) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c ContactsSorter) Less(i, j int) bool {
	return c[i].GetName() < c[j].GetName()
}
