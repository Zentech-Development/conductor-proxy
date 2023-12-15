package domain

import "context"

type Resource struct {
	ID           string
	Name         string
	FriendlyName string
	AppID        string
}

type ResourceInput struct {
	Name         string
	FriendlyName string
	AppID        string
}

type ResourceRepo interface {
	GetByID(ctx context.Context, id string) (Resource, error)
	Add(ctx context.Context, resource Resource) (Resource, error)
}

type ResourceHandlers interface {
	GetByID(id string, userGroups []string) (Resource, error)
	Add(resource ResourceInput, userGroups []string) (Resource, error)
}
