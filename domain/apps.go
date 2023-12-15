package domain

type App struct {
	ID           string
	Name         string
	FriendlyName string
	Host         string
	AdminGroups  []string
	UserGroups   []string
}

type AppInput struct {
	Name         string
	FriendlyName string
	Host         string
	AdminGroups  []string
	UserGroups   []string
}
