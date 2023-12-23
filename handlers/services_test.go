package handlers

import (
	"testing"

	"github.com/Zentech-Development/conductor-proxy/domain"
)

func TestServiceAdd(t *testing.T) {
	handlers := newHandlers()

	name := "test-service"
	friendlyName := "Test Service"
	host := "localhost:8000"

	service := domain.ServiceInput{
		Name:         name,
		FriendlyName: friendlyName,
		Host:         host,
		AdminGroups:  []string{"service-admins"},
		UserGroups:   []string{},
	}

	groupInput := domain.GroupInput{
		Name: "service-admins",
	}
	if _, err := handlers.Groups.Add(groupInput, []string{domain.GroupNameAdmin}); err != nil {
		t.Fatal("Unexpected error occurred")
	}

	savedService, err := handlers.Services.Add(service, []string{domain.GroupNameAdmin})
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	if savedService.Name != name || savedService.FriendlyName != friendlyName || savedService.Host != host {
		t.Fatal("Service saved improperly")
	}

	invalidService := domain.ServiceInput{
		Name:         name,
		FriendlyName: friendlyName,
		Host:         host,
		AdminGroups:  []string{"does-not-exist"},
		UserGroups:   []string{},
	}

	if _, err = handlers.Services.Add(invalidService, []string{domain.GroupNameAdmin}); err == nil {
		t.Fatal("Expected group does not exist error")
	}
}

func TestServiceNotAuthorized(t *testing.T) {
	handlers := newHandlers()

	name := "test-service"
	friendlyName := "Test Service"
	host := "localhost:8000"

	service := domain.ServiceInput{
		Name:         name,
		FriendlyName: friendlyName,
		Host:         host,
		AdminGroups:  []string{"service-admins"},
		UserGroups:   []string{},
	}

	groupInput := domain.GroupInput{
		Name: "service-admins",
	}
	if _, err := handlers.Groups.Add(groupInput, []string{domain.GroupNameAdmin}); err != nil {
		t.Fatal("Unexpected error occurred")
	}

	if _, err := handlers.Services.Add(service, []string{"not-admin"}); err == nil {
		t.Fatal("Expected not authorized error")
	}
}

func TestServiceGetByID(t *testing.T) {
	handlers := newHandlers()

	name := "test-service"
	friendlyName := "Test Service"
	host := "localhost:8000"

	service := domain.ServiceInput{
		Name:         name,
		FriendlyName: friendlyName,
		Host:         host,
		AdminGroups:  []string{"service-admins"},
		UserGroups:   []string{},
	}

	groupInput := domain.GroupInput{
		Name: "service-admins",
	}
	if _, err := handlers.Groups.Add(groupInput, []string{domain.GroupNameAdmin}); err != nil {
		t.Fatal("Unexpected error occurred")
	}

	savedService, err := handlers.Services.Add(service, []string{domain.GroupNameAdmin})
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	if _, err = handlers.Services.GetByID(savedService.ID, []string{"service-admins"}); err != nil {
		t.Fatal("Failed to get service")
	}
}

func TestServiceGetByIDAuthorized(t *testing.T) {
	handlers := newHandlers()

	name := "test-service"
	friendlyName := "Test Service"
	host := "localhost:8000"

	service := domain.ServiceInput{
		Name:         name,
		FriendlyName: friendlyName,
		Host:         host,
		AdminGroups:  []string{"service-admins"},
		UserGroups:   []string{},
	}

	groupInput := domain.GroupInput{
		Name: "service-admins",
	}
	if _, err := handlers.Groups.Add(groupInput, []string{domain.GroupNameAdmin}); err != nil {
		t.Fatal("Unexpected error occurred")
	}

	savedService, err := handlers.Services.Add(service, []string{domain.GroupNameAdmin})
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	if _, err = handlers.Services.GetByID(savedService.ID, []string{"not-admin"}); err == nil {
		t.Fatal("Expected not authorized error")
	}
}
