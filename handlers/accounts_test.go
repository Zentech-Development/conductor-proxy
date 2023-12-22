package handlers

import (
	"os"
	"strings"
	"testing"

	"github.com/Zentech-Development/conductor-proxy/domain"
)

func TestHashFunctions(t *testing.T) {
	testPassword := "test123"
	hash, err := hashPassword(testPassword, 12)
	if err != nil {
		t.Fatal("Failed to hash password")
	}

	if !strings.HasPrefix(hash, "$2a") {
		t.Fatal("Token is bad")
	}

	if !checkPassword(testPassword, hash) {
		t.Fatal("Password should have been correct")
	}

	if checkPassword("wrong", hash) {
		t.Fatal("Password should have been incorrect")
	}
}

func TestNewAccountAdd(t *testing.T) {
	handlers := newHandlers()

	username := "test-account"
	account := domain.AccountInput{
		Username:        username,
		Passkey:         "password123",
		Groups:          []string{},
		TokenExpiration: 60,
	}

	os.Setenv("CONDUCTOR_SECRET_KEY", "asdfasdfasdfasdfsdf")

	savedAccount, err := handlers.Accounts.Add(account, []string{domain.GroupNameAdmin})
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	if savedAccount.ID == "" || savedAccount.Username != username || savedAccount.Passkey != "" {
		t.Fatal("Account saved improperly")
	}

	invalidAccount := domain.AccountInput{
		Username:        "admin",
		Passkey:         "password123",
		Groups:          []string{},
		TokenExpiration: 60,
	}

	if _, err := handlers.Accounts.Add(invalidAccount, []string{domain.GroupNameAdmin}); err == nil {
		t.Fatal("Expected username not allowed error")
	}

	groupNameNotFoundAccount := domain.AccountInput{
		Username:        "admin",
		Passkey:         "password123",
		Groups:          []string{"i-dont-exist"},
		TokenExpiration: 60,
	}

	if _, err := handlers.Accounts.Add(groupNameNotFoundAccount, []string{domain.GroupNameAdmin}); err == nil {
		t.Fatal("Expected group name not found error")
	}
}

func TestNewAccountNotAuthorized(t *testing.T) {
	handlers := newHandlers()

	username := "test-account"
	account := domain.AccountInput{
		Username:        username,
		Passkey:         "password123",
		Groups:          []string{},
		TokenExpiration: 60,
	}

	os.Setenv("CONDUCTOR_SECRET_KEY", "asdfasdfasdfasdfsdf")

	if _, err := handlers.Accounts.Add(account, []string{"nonadmin"}); err == nil {
		t.Fatal("Expected not authorized error")
	}
}

func TestLogin(t *testing.T) {
	handlers := newHandlers()

	username := "test-account"
	passkey := "password123"
	account := domain.AccountInput{
		Username:        username,
		Passkey:         passkey,
		Groups:          []string{},
		TokenExpiration: 60,
	}

	os.Setenv("CONDUCTOR_SECRET_KEY", "asdfasdfasdfasdfsdf")

	_, err := handlers.Accounts.Add(account, []string{domain.GroupNameAdmin})
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	validCreds := domain.LoginInput{
		Username: username,
		Passkey:  passkey,
	}

	result, err := handlers.Accounts.Login(validCreds)
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	if result.Username != username {
		t.Fatal("Login returned wrong account")
	}

	badUsername := domain.LoginInput{
		Username: "bad",
		Passkey:  passkey,
	}
	if _, err = handlers.Accounts.Login(badUsername); err == nil {
		t.Fatal("Expected bad credentials error")
	}

	badPasskey := domain.LoginInput{
		Username: username,
		Passkey:  "bad",
	}
	if _, err = handlers.Accounts.Login(badPasskey); err == nil {
		t.Fatal("Expected bad credentials error")
	}

	badUsernameAndPasskey := domain.LoginInput{
		Username: "bad",
		Passkey:  "also-bad",
	}
	if _, err = handlers.Accounts.Login(badUsernameAndPasskey); err == nil {
		t.Fatal("Expected bad credentials error")
	}
}

func TestUpdateGroups(t *testing.T) {
	handlers := newHandlers()

	username := "test-account"
	account := domain.AccountInput{
		Username:        username,
		Passkey:         "password123",
		Groups:          []string{},
		TokenExpiration: 60,
	}

	os.Setenv("CONDUCTOR_SECRET_KEY", "asdfasdfasdfasdfsdf")

	savedAccount, err := handlers.Accounts.Add(account, []string{domain.GroupNameAdmin})
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	savedAccount.Groups = append(savedAccount.Groups, "test-group")

	if err := handlers.Accounts.UpdateGroups(savedAccount.ID, []string{}, []string{}, []string{domain.GroupNameAdmin}); err == nil {
		t.Fatal("Unexpected error occurred")
	}
}

func TestUpdateGroupsNotAuthorized(t *testing.T) {
	handlers := newHandlers()

	username := "test-account"
	account := domain.AccountInput{
		Username:        username,
		Passkey:         "password123",
		Groups:          []string{},
		TokenExpiration: 60,
	}

	os.Setenv("CONDUCTOR_SECRET_KEY", "asdfasdfasdfasdfsdf")

	savedAccount, err := handlers.Accounts.Add(account, []string{domain.GroupNameAdmin})
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	if _, err := handlers.Groups.Add(domain.GroupInput{Name: "test-group"}, []string{domain.GroupNameAdmin}); err != nil {
		t.Fatal("Unexpected error occurred")
	}

	account.Groups = append(account.Groups, "test-group")

	if err := handlers.Accounts.UpdateGroups(savedAccount.ID, []string{}, []string{}, []string{"not-admin"}); err == nil {
		t.Fatal("Unexpected error occurred")
	}
}
