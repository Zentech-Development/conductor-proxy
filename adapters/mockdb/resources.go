package adapters

import (
	"context"

	"github.com/Zentech-Development/conductor-proxy/adapters"
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
	for _, resource := range r.Data.Resources {
		if resource.ID == id {
			return resource, nil
		}
	}

	return domain.Resource{}, &adapters.NotFoundError{Name: "resource"}
}

func (r MockResourceRepo) Add(ctx context.Context, resource domain.Resource) (domain.Resource, error) {
	if _, err := r.GetByID(ctx, resource.ID); err == nil {
		return domain.Resource{}, &adapters.AlreadyExistsError{Name: "resource"}
	}

	foundService := false
	for _, service := range r.Data.Services {
		if service.ID == resource.ServiceID {
			foundService = true
		}
	}

	if !foundService {
		return domain.Resource{}, &adapters.NotFoundError{Name: "service"}
	}

	r.Data.Resources = append(r.Data.Resources, resource)
	return resource, nil
}
