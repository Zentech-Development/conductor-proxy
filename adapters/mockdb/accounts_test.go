package adapters_test

import (
	"context"
	"testing"

	mockAdapters "github.com/Zentech-Development/conductor-proxy/adapters/mockdb"
	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/google/uuid"
)

func TestAccountAdd(t *testing.T) {
	id := uuid.NewString()
	username := "test-account"

	db := mockAdapters.NewMockDB()

	account := domain.Account{
		ID:              id,
		Username:        username,
		Passkey:         "password123",
		Groups:          []string{"group1"},
		TokenExpiration: 60,
	}

	_, err := db.Accounts.Add(context.Background(), account)
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	savedAccount, err := db.Accounts.GetByUsername(context.Background(), username)
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	if savedAccount.ID != id || savedAccount.Username != username || len(savedAccount.Groups) != 1 ||
		savedAccount.Groups[0] != "group1" || savedAccount.TokenExpiration != 60 {
		t.Fatal("Account was saved improperly")
	}

	if _, err = db.Accounts.Add(context.Background(), account); err == nil {
		t.Fatal("Expected an already exists error")
	}
}

func TestAccountGetByUsername(t *testing.T) {
	id := uuid.NewString()
	username := "test-user"

	db := mockAdapters.NewMockDB()

	account := domain.Account{
		ID:              id,
		Username:        username,
		Passkey:         "password123",
		Groups:          []string{"group1"},
		TokenExpiration: 60,
	}

	_, err := db.Accounts.Add(context.Background(), account)
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	savedAccount, err := db.Accounts.GetByUsername(context.Background(), username)
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	if savedAccount.ID != id || savedAccount.Username != username || len(savedAccount.Groups) != 1 ||
		savedAccount.Groups[0] != "group1" || savedAccount.TokenExpiration != 60 {
		t.Fatal("Account was retrieved improperly")
	}

	if _, err := db.Accounts.GetByUsername(context.Background(), "bad"); err == nil {
		t.Fatal("Expected an error")
	}
}

func TestAccountUpdate(t *testing.T) {
	id := uuid.NewString()
	username := "test-user"

	db := mockAdapters.NewMockDB()

	account := domain.Account{
		ID:              id,
		Username:        username,
		Passkey:         "password123",
		Groups:          []string{"group1"},
		TokenExpiration: 60,
	}

	_, err := db.Accounts.Add(context.Background(), account)
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	account.Groups = append(account.Groups, "group2")

	if _, err := db.Accounts.Update(context.Background(), account); err != nil {
		t.Fatal("Unexpected error occurrred")
	}

	savedAccount, err := db.Accounts.GetByUsername(context.Background(), username)
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	if savedAccount.ID != id || savedAccount.Username != username || len(savedAccount.Groups) != 2 ||
		savedAccount.Groups[0] != "group1" || savedAccount.TokenExpiration != 60 || savedAccount.Groups[1] != "group2" {
		t.Fatal("Group was added improperly")
	}

	savedAccount.Groups = []string{}

	if _, err := db.Accounts.Update(context.Background(), savedAccount); err != nil {
		t.Fatal("Unexpected error occurrred")
	}

	savedAccount, err = db.Accounts.GetByUsername(context.Background(), username)
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	if savedAccount.ID != id || savedAccount.Username != username ||
		len(savedAccount.Groups) != 0 || savedAccount.TokenExpiration != 60 {
		t.Fatal("Group was removed improperly")
	}
}
