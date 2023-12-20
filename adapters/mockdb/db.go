package adapters

import (
	"github.com/Zentech-Development/conductor-proxy/domain"
)

type MockDBData struct {
	Accounts  []domain.Account
	Groups    []domain.Group
	Resources []domain.Resource
	Services  []domain.Service
}

type MockDB struct {
	Data *MockDBData
}

func NewMockDB() *domain.Repos {
	data := &MockDBData{
		Accounts:  make([]domain.Account, 0),
		Groups:    make([]domain.Group, 0),
		Resources: make([]domain.Resource, 0),
		Services:  make([]domain.Service, 0),
	}

	return &domain.Repos{
		Accounts:  newMockAccountRepo(data),
		Groups:    newMockGroupRepo(data),
		Resources: newMockResourceRepo(data),
		Services:  newMockServiceRepo(data),
	}
}
