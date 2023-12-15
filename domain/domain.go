package domain

type Repos struct {
	Resources ResourceRepo
	Apps      AppRepo
	Groups    GroupRepo
	Accounts  AccountRepo
}

type Logger struct{}

type Handlers struct {
	Resources ResourceHandlers
	Apps      AppHandlers
	Groups    GroupHandlers
	Accounts  AccountHandlers
	Proxy     ProxyHandlers
}

type Adapters struct {
	Repos  Repos
	Logger Logger
}
