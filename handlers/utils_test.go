package handlers

import (
	"os"
	"testing"

	adapters "github.com/Zentech-Development/conductor-proxy/adapters/mockdb"
	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/Zentech-Development/conductor-proxy/pkg/config"
)

func newHandlers() domain.Handlers {
	mockDB := adapters.NewMockDB()

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

func TestMain(m *testing.M) {
	config.SetAndGetConfig("")
	m.Run()
	os.Exit(0)
}
