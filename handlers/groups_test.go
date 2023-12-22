package handlers

import (
	"testing"

	"github.com/Zentech-Development/conductor-proxy/domain"
)

func TestGroupAdd(t *testing.T) {
	handlers := newHandlers()

	groupName := "test-group"
	group := domain.GroupInput{
		Name: groupName,
	}

	savedGroup, err := handlers.Groups.Add(group, []string{domain.GroupNameAdmin})
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	if savedGroup.Name != groupName {
		t.Fatal("Group saved improperly")
	}

	if _, err = handlers.Groups.Add(group, []string{domain.GroupNameAdmin}); err == nil {
		t.Fatal("Expected already exists error")
	}

	invalidGroup := domain.GroupInput{
		Name: domain.GroupNameAdmin,
	}
	if _, err = handlers.Groups.Add(invalidGroup, []string{domain.GroupNameAdmin}); err == nil {
		t.Fatal("Expected group name not allowed error")
	}
}

func TestGroupAddNotAuthorized(t *testing.T) {
	handlers := newHandlers()

	groupName := "test-group"
	group := domain.GroupInput{
		Name: groupName,
	}

	_, err := handlers.Groups.Add(group, []string{"not-admin"})
	if err == nil {
		t.Fatal("Expected not authorized error")
	}
}
