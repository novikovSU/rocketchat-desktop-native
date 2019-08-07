package model

import (
	"github.com/novikovSU/gorocket/api"
)

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

func (u UserModel) GetUnreadCount() int {
	return u.UnreadCount
}

func UsersToModels(users []UserModel) []IContactModel {
	models := make([]IContactModel, len(users))
	for _, u := range users {
		models = append(models, u)
	}

	return models
}

func UsersMapToModels(users map[string]UserModel) []IContactModel {
	models := make([]IContactModel, len(users))
	for _, u := range users {
		models = append(models, u)
	}

	return models
}
