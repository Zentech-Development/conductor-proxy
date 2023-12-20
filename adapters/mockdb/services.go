package adapters

import (
	"context"

	"github.com/Zentech-Development/conductor-proxy/domain"
)

type MockServiceRepo struct {
	Data *MockDBData
}

func newMockServiceRepo(data *MockDBData) MockServiceRepo {
	return MockServiceRepo{
		Data: data,
	}
}

func (r MockServiceRepo) GetByID(ctx context.Context, id string) (domain.Service, error) {
	return domain.Service{}, nil
}

func (r MockServiceRepo) Add(ctx context.Context, resource domain.Service) (domain.Service, error) {
	return domain.Service{}, nil
}
