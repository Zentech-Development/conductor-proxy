package domain

type Resource struct {
	ID           string
	Name         string
	FriendlyName string
}

type ResourceInput struct {
	Name         string
	FriendlyName string
}
