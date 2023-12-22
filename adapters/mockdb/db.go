package adapters

import (
	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/google/uuid"
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

func NewMockDB(initialAdminAccount *domain.Account, initialAdminGroup string) domain.Repos {
	data := &MockDBData{
		Accounts:  make([]domain.Account, 0),
		Groups:    make([]domain.Group, 0),
		Resources: make([]domain.Resource, 0),
		Services:  make([]domain.Service, 0),
	}

	if initialAdminAccount != nil {
		data.Accounts = append(data.Accounts, *initialAdminAccount)
	}

	if initialAdminGroup != "" {
		data.Groups = append(data.Groups, domain.Group{
			ID:   uuid.NewString(),
			Name: initialAdminGroup,
		})
	}

	return domain.Repos{
		Accounts:  newMockAccountRepo(data),
		Groups:    newMockGroupRepo(data),
		Resources: newMockResourceRepo(data),
		Services:  newMockServiceRepo(data),
	}
}
