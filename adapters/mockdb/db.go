package adapters

import (
	"slices"

	"github.com/Zentech-Development/conductor-proxy/domain"
	conf "github.com/Zentech-Development/conductor-proxy/pkg/config"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

func NewMockDB() domain.Repos {
	data := &MockDBData{
		Accounts:  make([]domain.Account, 0),
		Groups:    make([]domain.Group, 0),
		Resources: make([]domain.Resource, 0),
		Services:  make([]domain.Service, 0),
	}

	addFirstAdminUserIfRequired(data)

	return domain.Repos{
		Accounts:  newMockAccountRepo(data),
		Groups:    newMockGroupRepo(data),
		Resources: newMockResourceRepo(data),
		Services:  newMockServiceRepo(data),
	}
}

func addFirstAdminUserIfRequired(data *MockDBData) {
	config := conf.GetConfig()
	foundAdminGroup := false

	for _, group := range data.Groups {
		if group.Name == domain.GroupNameAdmin {
			foundAdminGroup = true
		}
	}

	if !foundAdminGroup {
		data.Groups = append(data.Groups, domain.Group{
			ID:   uuid.NewString(),
			Name: domain.GroupNameAdmin,
		})
	}

	for _, account := range data.Accounts {
		if slices.Contains(account.Groups, domain.GroupNameAdmin) {
			return
		}
	}

	adminUser := domain.Account{
		ID:              uuid.NewString(),
		Username:        config.DefaultAdminUsername,
		Passkey:         config.DefaultAdminPasskey,
		Groups:          []string{domain.GroupNameAdmin},
		TokenExpiration: 3600,
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(adminUser.Passkey), config.AccessTokenCost)
	if err != nil {
		panic(err)
	}

	adminUser.Passkey = string(hash)

	data.Accounts = append(data.Accounts, adminUser)
}
