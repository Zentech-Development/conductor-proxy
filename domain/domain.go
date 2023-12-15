package domain

type Repos struct {
	Resources interface {
		GetByID(id string) (Resource, error)
		Add(resource Resource) (Resource, error)
	}
	Apps interface {
		GetByID(id string) (App, error)
		Add(app App) (App, error)
	}
	Groups interface {
		Add(group Group) (Group, error)
	}
	Accounts interface {
		Add(account Account) (Account, error)
		Update(account Account) (Account, error)
	}
}

type Logger struct{}

type Handlers struct {
	Resources interface {
		GetByID(id string) (Resource, error)
		Add(resource ResourceInput) (Resource, error)
	}
	Apps interface {
		GetByID(id string) (App, error)
		Add(app AppInput) (App, error)
	}
	Groups interface {
		Add(group GroupInput) (Group, error)
	}
	Accounts interface {
		Add(account AccountInput) (Account, error)
		UpdateGroups(id string, groupsToAdd []string, groupsToRemove []string) error
	}
}

type Adapters struct {
	Repos  Repos
	Logger Logger
}
