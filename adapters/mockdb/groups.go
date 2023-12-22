package adapters

import (
	"context"

	"github.com/Zentech-Development/conductor-proxy/adapters"
	"github.com/Zentech-Development/conductor-proxy/domain"
)

type MockGroupRepo struct {
	Data *MockDBData
}

func newMockGroupRepo(data *MockDBData) MockGroupRepo {
	return MockGroupRepo{
		Data: data,
	}
}

func (r MockGroupRepo) Add(ctx context.Context, group domain.Group) (domain.Group, error) {
	if _, err := r.GetByName(ctx, group.Name); err == nil {
		return domain.Group{}, &adapters.AlreadyExistsError{Name: "group"}
	}

	r.Data.Groups = append(r.Data.Groups, group)
	return group, nil
}

func (r MockGroupRepo) GetByName(ctx context.Context, name string) (domain.Group, error) {
	for _, group := range r.Data.Groups {
		if group.Name == name {
			return group, nil
		}
	}

	return domain.Group{}, &adapters.NotFoundError{Name: "group"}
}
