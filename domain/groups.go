package domain

import "context"

const (
	GroupNameAdmin = "admin"
)

type Group struct {
	ID   string
	Name string
}

type GroupInput struct {
	Name string `json:"name" binding:"required"`
}

type GroupRepo interface {
	Add(ctx context.Context, group Group) (Group, error)
	GetByName(ctx context.Context, name string) (Group, error)
}

type GroupHandlers interface {
	Add(group GroupInput, userGroups []string) (Group, error)
}
