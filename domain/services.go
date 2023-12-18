package domain

import "context"

type Service struct {
	ID           string
	Name         string
	FriendlyName string
	Host         string
	AdminGroups  []string
	UserGroups   []string
	Type         string
}

type ServiceInput struct {
	Name         string   `json:"name" binding:"required"`
	FriendlyName string   `json:"friendlyName" binding:"required"`
	Host         string   `json:"host" binding:"required"`
	AdminGroups  []string `json:"adminGroups" binding:"required"`
	UserGroups   []string `json:"userGroups" binding:"required"`
}

type ServiceRepo interface {
	GetByID(ctx context.Context, id string) (Service, error)
	Add(ctx context.Context, service Service) (Service, error)
}

type ServiceHandlers interface {
	GetByID(id string, userGroups []string) (Service, error)
	Add(service ServiceInput, userGroups []string) (Service, error)
}

const (
	ServiceTypeHTTP  = "http"
	ServiceTypeHTTPS = "https"
)
