package domain

type Account struct {
	ID              string
	Username        string
	Passkey         string
	Groups          []string
	TokenExpiration int
}

type AccountInput struct {
	Username        string
	Passkey         string
	Groups          []string
	TokenExpiration int
}
