package adapters_test

import (
	"context"
	"testing"

	mockAdapters "github.com/Zentech-Development/conductor-proxy/adapters/mockdb"
	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/google/uuid"
)

func TestResourceGetByID(t *testing.T) {
	db := mockAdapters.NewMockDB(nil, "")

	id := uuid.NewString()
	name := "test-resource"
	friendlyName := "Test Resource"
	serviceId := uuid.NewString()

	service := domain.Service{
		ID:           serviceId,
		Name:         "test-service",
		FriendlyName: "Test Service",
		Host:         "localhost:8000",
		AdminGroups:  []string{"admins"},
		UserGroups:   []string{"test-service-users"},
		Type:         domain.ServiceTypeHTTP,
	}

	if _, err := db.Services.Add(context.Background(), service); err != nil {
		t.Fatal("Unexpected error occurred")
	}

	resource := domain.Resource{
		ID:           id,
		Name:         name,
		FriendlyName: friendlyName,
		ServiceID:    serviceId,
		Properties:   []domain.Property{},
		Endpoints:    []domain.Endpoint{},
	}

	if _, err := db.Resources.Add(context.Background(), resource); err != nil {
		t.Fatal("Unexpected error occurred")
	}

	savedResource, err := db.Resources.GetByID(context.Background(), id)
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	if savedResource.ID != id || savedResource.Name != name || savedResource.FriendlyName != friendlyName ||
		len(savedResource.Endpoints) != 0 || len(savedResource.Properties) != 0 {
		t.Fatal("Resource retrieved improperly")
	}

	if _, err := db.Resources.GetByID(context.Background(), "invalid"); err == nil {
		t.Fatal("Expected an error")
	}
}

func TestResourceAdd(t *testing.T) {
	db := mockAdapters.NewMockDB(nil, "")

	id := uuid.NewString()
	name := "test-resource"
	friendlyName := "Test Resource"
	serviceId := uuid.NewString()

	service := domain.Service{
		ID:           serviceId,
		Name:         "test-service",
		FriendlyName: "Test Service",
		Host:         "localhost:8000",
		AdminGroups:  []string{"admins"},
		UserGroups:   []string{"test-service-users"},
		Type:         domain.ServiceTypeHTTP,
	}

	if _, err := db.Services.Add(context.Background(), service); err != nil {
		t.Fatal("Unexpected error occurred")
	}

	resource := domain.Resource{
		ID:           id,
		Name:         name,
		FriendlyName: friendlyName,
		ServiceID:    serviceId,
		Properties:   []domain.Property{},
		Endpoints:    []domain.Endpoint{},
	}

	if _, err := db.Resources.Add(context.Background(), resource); err != nil {
		t.Fatal("Unexpected error occurred")
	}

	savedResource, err := db.Resources.GetByID(context.Background(), id)
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	if savedResource.ID != id || savedResource.Name != name || savedResource.FriendlyName != friendlyName ||
		len(savedResource.Endpoints) != 0 || len(savedResource.Properties) != 0 {
		t.Fatal("Resource saved improperly")
	}

	if _, err = db.Resources.Add(context.Background(), savedResource); err == nil {
		t.Fatal("Expected an error")
	}
}
