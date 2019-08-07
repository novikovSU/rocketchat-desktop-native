package model

import (
	"github.com/novikovSU/gorocket/api"
)

type GroupModel struct {
	UnreadCount int
	Group       api.Group
	History     []api.Message
}

func (g GroupModel) GetId() string {
	return g.Group.ID
}

func (g GroupModel) GetName() string {
	return g.Group.Name
}

func (g GroupModel) GetDisplayName() string {
	return lockSign + g.GetName()
}

func (g GroupModel) GetUnreadCount() int {
	return g.UnreadCount
}

func GroupsToModels(groups []GroupModel) []IContactModel {
	models := make([]IContactModel, len(groups))
	for _, g := range groups {
		models = append(models, g)
	}

	return models
}

func GroupsMapToModels(groups map[string]GroupModel) []IContactModel {
	models := make([]IContactModel, len(groups))
	for _, g := range groups {
		models = append(models, g)
	}

	return models
}
