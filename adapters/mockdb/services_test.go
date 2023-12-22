package adapters_test

import (
	"context"
	"testing"

	mockAdapters "github.com/Zentech-Development/conductor-proxy/adapters/mockdb"
	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/google/uuid"
)

func TestServiceGetByID(t *testing.T) {
	db := mockAdapters.NewMockDB(nil)

	id := uuid.NewString()
	name := "test-service"
	friendlyName := "Test Service"
	host := "localhost:8000"
	serviceType := domain.ServiceTypeHTTP

	service := domain.Service{
		ID:           id,
		Name:         name,
		FriendlyName: friendlyName,
		Host:         host,
		AdminGroups:  []string{"admins"},
		UserGroups:   []string{"test-service-users"},
		Type:         serviceType,
	}

	if _, err := db.Services.Add(context.Background(), service); err != nil {
		t.Fatal("Unexpected error occurred")
	}

	savedService, err := db.Services.GetByID(context.Background(), id)
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	if savedService.ID != id || savedService.Name != name || savedService.FriendlyName != friendlyName ||
		savedService.Host != host || savedService.Type != serviceType {
		t.Fatal("Service saved improperly")
	}

	if _, err = db.Services.GetByID(context.Background(), "invalid"); err == nil {
		t.Fatal("Expected an error")
	}
}

func TestServiceAdd(t *testing.T) {
	db := mockAdapters.NewMockDB(nil)

	id := uuid.NewString()
	name := "test-service"
	friendlyName := "Test Service"
	host := "localhost:8000"
	serviceType := domain.ServiceTypeHTTP

	service := domain.Service{
		ID:           id,
		Name:         name,
		FriendlyName: friendlyName,
		Host:         host,
		AdminGroups:  []string{"admins"},
		UserGroups:   []string{"test-service-users"},
		Type:         serviceType,
	}

	if _, err := db.Services.Add(context.Background(), service); err != nil {
		t.Fatal("Unexpected error occurred")
	}

	savedService, err := db.Services.GetByID(context.Background(), id)
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	if savedService.ID != id || savedService.Name != name || savedService.FriendlyName != friendlyName ||
		savedService.Host != host || savedService.Type != serviceType {
		t.Fatal("Service saved improperly")
	}

	if _, err = db.Services.Add(context.Background(), savedService); err == nil {
		t.Fatal("Expected an error")
	}
}
