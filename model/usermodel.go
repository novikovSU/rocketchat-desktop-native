package model

import "github.com/novikovSU/gorocket/api"

type UserModel struct {
	UnreadCount int
	User        api.User
	History     []api.Message
}

func (u UserModel) GetId() string {
	return u.User.ID
}

func (u UserModel) GetName() string {
	return u.User.Name
}

func (u UserModel) GetDisplayName() string {
	return u.GetName()
}
