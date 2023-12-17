package handlers

import (
	"context"
	"errors"
	"time"

	"github.com/Zentech-Development/conductor-proxy/domain"
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

	isAdmin := checkForGroupMatch(userGroups, make([]string, 0))

	if !isAdmin {
		return domain.Account{}, errors.New("Not authorized")
	}

	if account.Username == "admin" {
		return domain.Account{}, errors.New("Account name is not allowed")
	}

	hashedPasskey, err := hashPassword(account.Passkey)
	if err != nil {
		return domain.Account{}, errors.New("Failed to generate password hash")
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

	account, err := h.Adapters.Repos.Accounts.GetByID(ctx, credentials.Username)
	if err != nil {
		time.Sleep(time.Second)
		return domain.Account{}, errors.New("Invalid credentials")
	}

	if !checkPassword(credentials.Passkey, account.Passkey) {
		time.Sleep(time.Second)
		return domain.Account{}, errors.New("Invalid credentials")
	}

	return account, nil
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

const hashCost = 12

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func checkPassword(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}
