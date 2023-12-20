package adapters

import (
	"context"

	"github.com/Zentech-Development/conductor-proxy/adapters"
	"github.com/Zentech-Development/conductor-proxy/domain"
)

type MockAccountRepo struct {
	Data *MockDBData
}

func newMockAccountRepo(data *MockDBData) MockAccountRepo {
	return MockAccountRepo{
		Data: data,
	}
}

func (r MockAccountRepo) GetByUsername(ctx context.Context, id string) (domain.Account, error) {
	return domain.Account{}, nil
}

func (r MockAccountRepo) Add(ctx context.Context, account domain.Account) (domain.Account, error) {
	if _, err := r.GetByUsername(ctx, account.Username); err == nil {
		return domain.Account{}, &adapters.AlreadyExistsError{Name: "account"}
	}

	r.Data.Accounts = append(r.Data.Accounts, account)
	return domain.Account{}, nil
}

func (r MockAccountRepo) Update(ctx context.Context, account domain.Account) (domain.Account, error) {
	return domain.Account{}, nil
}
