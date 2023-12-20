package handlers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/Zentech-Development/conductor-proxy/pkg/config"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

	if !isAdmin(userGroups) {
		return domain.Account{}, errors.New("not authorized")
	}

	if account.Username == "admin" {
		return domain.Account{}, errors.New("account name is not allowed")
	}

	for _, group := range userGroups {
		if _, err := h.Adapters.Repos.Groups.GetByName(ctx, group); err != nil {
			return domain.Account{}, fmt.Errorf("group name %s not found", group)
		}
	}

	hashedPasskey, err := hashPassword(account.Passkey, config.GetConfig().JwtHashCost)
	if err != nil {
		return domain.Account{}, errors.New("failed to generate password hash")
	}

	accountToSave := domain.Account{
		ID:              uuid.NewString(),
		Username:        account.Username,
		Passkey:         hashedPasskey,
		Groups:          account.Groups,
		TokenExpiration: account.TokenExpiration,
	}

	savedAccount, err := h.Adapters.Repos.Accounts.Add(ctx, accountToSave)
	if err != nil {
		return domain.Account{}, err
	}

	savedAccount.Passkey = ""

	return savedAccount, nil
}

func (h AccountHandler) Login(credentials domain.LoginInput) (domain.Account, error) {
	ctx := context.Background()

	account, err := h.Adapters.Repos.Accounts.GetByUsername(ctx, credentials.Username)
	if err != nil {
		time.Sleep(time.Second)
		return domain.Account{}, errors.New("invalid credentials")
	}

	if !checkPassword(credentials.Passkey, account.Passkey) {
		time.Sleep(time.Second)
		return domain.Account{}, errors.New("invalid credentials")
	}

	return account, nil
}

func (h AccountHandler) UpdateGroups(id string, groupsToAdd []string, groupsToRemove []string, userGroups []string) error {
	ctx := context.Background()

	if !isAdmin(userGroups) {
		return errors.New("not authorized")
	}

	for _, group := range groupsToAdd {
		if _, err := h.Adapters.Repos.Groups.GetByName(ctx, group); err != nil {
			return fmt.Errorf("group name %s not found", group)
		}
	}

	account, err := h.Adapters.Repos.Accounts.GetByUsername(ctx, id)
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

func hashPassword(password string, cost int) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func checkPassword(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}
