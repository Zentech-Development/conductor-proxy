package domain

type Repos struct {
	Resources ResourceRepo
	Services  ServiceRepo
	Groups    GroupRepo
	Accounts  AccountRepo
}

type Logger struct{}

type Handlers struct {
	Resources ResourceHandlers
	Services  ServiceHandlers
	Groups    GroupHandlers
	Accounts  AccountHandlers
	Proxy     ProxyHandlers
}

type Adapters struct {
	Repos  Repos
	Logger Logger
}
