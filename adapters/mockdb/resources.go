package adapters

import (
	"context"

	"github.com/Zentech-Development/conductor-proxy/domain"
)

type MockResourceRepo struct {
	Data *MockDBData
}

func newMockResourceRepo(data *MockDBData) MockResourceRepo {
	return MockResourceRepo{
		Data: data,
	}
}

func (r MockResourceRepo) GetByID(ctx context.Context, id string) (domain.Resource, error) {
	return domain.Resource{}, nil
}

func (r MockResourceRepo) Add(ctx context.Context, resource domain.Resource) (domain.Resource, error) {
	return domain.Resource{}, nil
}
