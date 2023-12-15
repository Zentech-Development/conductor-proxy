package handlers

import (
	"context"
	"errors"

	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

type AccountHandler struct {
	Adapters *domain.Adapters
}

func NewAccountHandler(adapters *domain.Adapters) AccountHandler {
	return AccountHandler{
		Adapters: adapters,
	}
}

func (h AccountHandler) Add(account domain.AccountInput, userGroups []string) (domain.Account, error) {
	ctx := context.Background()

	isAdmin := checkForGroupMatch(userGroups, make([]string, 0))

	if !isAdmin {
		return domain.Account{}, errors.New("Not authorized")
	}

	if account.Username == "admin" {
		return domain.Account{}, errors.New("Account name is not allowed")
	}

	accountToSave := domain.Account{
		ID:              uuid.NewString(),
		Username:        account.Username,
		Passkey:         account.Passkey, // TODO: HASH THIS
		Groups:          account.Groups,
		TokenExpiration: account.TokenExpiration,
	}

	savedAccount, err := h.Adapters.Repos.Accounts.Add(ctx, accountToSave)
	if err != nil {
		return domain.Account{}, err
	}

	return savedAccount, nil
}

func (h AccountHandler) UpdateGroups(id string, groupsToAdd []string, groupsToRemove []string, userGroups []string) error {
	ctx := context.Background()

	isAdmin := checkForGroupMatch(userGroups, make([]string, 0))

	if !isAdmin {
		return errors.New("Not authorized")
	}

	account, err := h.Adapters.Repos.Accounts.GetByID(ctx, id)
	if err != nil {
		return err
	}

	account.Groups = slices.Compact(append(account.Groups, groupsToAdd...))

	validGroups := make([]string, 0)

	for _, group := range account.Groups {
		if !slices.Contains(groupsToRemove, group) {
			validGroups = append(validGroups, group)
		}
	}

	account.Groups = validGroups

	_, err = h.Adapters.Repos.Accounts.Update(ctx, account)
	if err != nil {
		return err
	}

	return nil
}
