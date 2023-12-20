package adapters

import (
	"context"

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
	return domain.Group{}, nil
}

func (r MockGroupRepo) GetByName(ctx context.Context, name string) (domain.Group, error) {
	return domain.Group{}, nil
}
