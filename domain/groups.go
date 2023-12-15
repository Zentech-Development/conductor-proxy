package domain

import "context"

type Group struct {
	ID   string
	Name string
}

type GroupInput struct {
	Name string
}

type GroupRepo interface {
	Add(ctx context.Context, group Group) (Group, error)
}

type GroupHandlers interface {
	Add(group GroupInput, userGroups []string) (Group, error)
}
