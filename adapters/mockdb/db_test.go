package adapters_test

import (
	"context"
	"testing"

	mockAdapters "github.com/Zentech-Development/conductor-proxy/adapters/mockdb"
	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/google/uuid"
)

func TestNewMockDB(t *testing.T) {
	_ = mockAdapters.NewMockDB(nil, "")

	username := "test-admin"
	adminAccount := &domain.Account{
		ID:              uuid.NewString(),
		Username:        username,
		Passkey:         "password123",
		Groups:          []string{domain.GroupNameAdmin},
		TokenExpiration: 0,
	}

	db := mockAdapters.NewMockDB(adminAccount, "")

	if _, err := db.Accounts.GetByUsername(context.Background(), username); err != nil {
		t.Fatal("Failed to create initial admin account")
	}
}
