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
	for _, account := range r.Data.Accounts {
		if account.Username == id {
			return account, nil
		}
	}

	return domain.Account{}, &adapters.NotFoundError{Name: "account"}
}

func (r MockAccountRepo) Add(ctx context.Context, account domain.Account) (domain.Account, error) {
	if _, err := r.GetByUsername(ctx, account.Username); err == nil {
		return domain.Account{}, &adapters.AlreadyExistsError{Name: "account"}
	}

	r.Data.Accounts = append(r.Data.Accounts, account)
	return account, nil
}

func (r MockAccountRepo) Update(ctx context.Context, account domain.Account) (domain.Account, error) {
	for i, savedAccount := range r.Data.Accounts {
		if savedAccount.ID == account.ID {
			r.Data.Accounts[i].Groups = account.Groups
			return r.Data.Accounts[i], nil
		}
	}

	return domain.Account{}, &adapters.NotFoundError{Name: "account"}
}
