package adapters_test

import (
	"context"
	"testing"

	mockAdapters "github.com/Zentech-Development/conductor-proxy/adapters/mockdb"
	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/google/uuid"
)

func TestGroupGetByName(t *testing.T) {
	db := mockAdapters.NewMockDB()

	id := uuid.NewString()
	name := "test-group"

	group := domain.Group{
		ID:   id,
		Name: name,
	}

	if _, err := db.Groups.Add(context.Background(), group); err != nil {
		t.Fatal("Unexpected error occurred")
	}

	savedGroup, err := db.Groups.GetByName(context.Background(), name)
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	if savedGroup.ID != id || savedGroup.Name != name {
		t.Fatal("Group retrieved improperly")
	}

	if _, err := db.Groups.GetByName(context.Background(), "invalid"); err == nil {
		t.Fatal("Expected an error")
	}
}

func TestGroupAdd(t *testing.T) {
	db := mockAdapters.NewMockDB()

	id := uuid.NewString()
	name := "test-group"

	group := domain.Group{
		ID:   id,
		Name: name,
	}

	if _, err := db.Groups.Add(context.Background(), group); err != nil {
		t.Fatal("Unexpected error occurred")
	}

	savedGroup, err := db.Groups.GetByName(context.Background(), name)
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	if savedGroup.ID != id || savedGroup.Name != name {
		t.Fatal("Group saved improperly")
	}

	if _, err := db.Groups.Add(context.Background(), savedGroup); err == nil {
		t.Fatal("Expected an error")
	}
}
