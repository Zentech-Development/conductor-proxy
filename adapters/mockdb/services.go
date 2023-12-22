package adapters

import (
	"context"

	"github.com/Zentech-Development/conductor-proxy/adapters"
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
	for _, service := range r.Data.Services {
		if service.ID == id {
			return service, nil
		}
	}

	return domain.Service{}, &adapters.NotFoundError{Name: "service"}
}

func (r MockServiceRepo) Add(ctx context.Context, service domain.Service) (domain.Service, error) {
	if _, err := r.GetByID(ctx, service.ID); err == nil {
		return domain.Service{}, &adapters.AlreadyExistsError{Name: "service"}
	}

	r.Data.Services = append(r.Data.Services, service)
	return service, nil
}
