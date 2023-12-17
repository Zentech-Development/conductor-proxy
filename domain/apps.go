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
	Name         string   `json:"name" binding:"required"`
	FriendlyName string   `json:"friendlyName" binding:"required"`
	Host         string   `json:"host" binding:"required"`
	AdminGroups  []string `json:"adminGroups" binding:"required"`
	UserGroups   []string `json:"userGroups" binding:"required"`
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
