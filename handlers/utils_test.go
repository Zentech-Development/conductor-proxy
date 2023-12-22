package handlers

import (
	adapters "github.com/Zentech-Development/conductor-proxy/adapters/mockdb"
	"github.com/Zentech-Development/conductor-proxy/domain"
)

func newHandlers() domain.Handlers {
	mockDB := adapters.NewMockDB(nil, domain.GroupNameAdmin)

	adapts := domain.Adapters{
		Repos:  mockDB,
		Logger: domain.Logger{},
	}

	handlers := domain.Handlers{
		Services:  NewServiceHandler(&adapts),
		Accounts:  NewAccountHandler(&adapts),
		Groups:    NewGroupHandler(&adapts),
		Resources: NewResourceHandler(&adapts),
		Proxy:     NewProxyHandler(&adapts),
	}

	return handlers
}
