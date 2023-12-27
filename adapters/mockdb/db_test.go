package adapters_test

import (
	"context"
	"testing"

	mockAdapters "github.com/Zentech-Development/conductor-proxy/adapters/mockdb"
	"github.com/Zentech-Development/conductor-proxy/pkg/config"
)

func TestNewMockDB(t *testing.T) {
	_ = mockAdapters.NewMockDB()

	db := mockAdapters.NewMockDB()

	conf := config.GetConfig()

	if _, err := db.Accounts.GetByUsername(context.Background(), conf.DefaultAdminUsername); err != nil {
		t.Fatal("Failed to create initial admin account")
	}
}
