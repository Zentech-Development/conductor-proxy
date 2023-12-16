package domain

import "context"

type App struct {
	ID           string
	Name         string
	FriendlyName string
	Host         string
	AdminGroups  []string
	UserGroups   []string
	Type         string
}

type AppInput struct {
	Name         string
	FriendlyName string
	Host         string
	AdminGroups  []string
	UserGroups   []string
}

type AppRepo interface {
	GetByID(ctx context.Context, id string) (App, error)
	Add(ctx context.Context, app App) (App, error)
}

type AppHandlers interface {
	GetByID(id string, userGroups []string) (App, error)
	Add(app AppInput, userGroups []string) (App, error)
}

const (
	AppTypeHTTP  = "http"
	AppTypeHTTPS = "https"
)
