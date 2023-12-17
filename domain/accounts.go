package domain

import "context"

type Account struct {
	ID              string   `json:"id"`
	Username        string   `json:"username"`
	Passkey         string   `json:"passkey,omitempty"`
	Groups          []string `json:"groups"`
	TokenExpiration int      `json:"tokenExpiration"`
}

type AccountInput struct {
	Username        string   `json:"username"`
	Passkey         string   `json:"passkey"`
	Groups          []string `json:"groups"`
	TokenExpiration int      `json:"tokenExpiration"`
}

type LoginInput struct {
	Username string
	Passkey  string
}

type AccountRepo interface {
	GetByID(ctx context.Context, id string) (Account, error)
	Add(ctx context.Context, account Account) (Account, error)
	Update(ctx context.Context, account Account) (Account, error)
}

type AccountHandlers interface {
	Add(account AccountInput, userGroups []string) (Account, error)
	UpdateGroups(id string, groupsToAdd []string, groupsToRemove []string, userGroups []string) error
}
