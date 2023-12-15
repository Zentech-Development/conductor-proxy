package handlers

import (
	"context"
	"errors"

	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/google/uuid"
)

type ResourceHandler struct {
	Adapters *domain.Adapters
}

func NewResourceHandler(adapters *domain.Adapters) ResourceHandler {
	return ResourceHandler{
		Adapters: adapters,
	}
}

func (h ResourceHandler) GetByID(id string, userGroups []string) (domain.Resource, error) {
	ctx := context.Background()

	resource, err := h.Adapters.Repos.Resources.GetByID(ctx, id)
	if err != nil {
		return domain.Resource{}, err
	}

	app, err := h.Adapters.Repos.Apps.GetByID(ctx, resource.AppID)
	if err != nil {
		return domain.Resource{}, errors.New("Failed to find app")
	}

	appHasGroups := (len(app.AdminGroups) + len(app.UserGroups)) > 0
	isAppUser := checkForGroupMatch(userGroups, app.UserGroups)
	isAppAdmin := checkForGroupMatch(userGroups, app.AdminGroups)

	if appHasGroups && !isAppAdmin && !isAppUser {
		return domain.Resource{}, errors.New("Not authorized")
	}

	return resource, nil
}

func (h ResourceHandler) Add(resource domain.ResourceInput, userGroups []string) (domain.Resource, error) {
	ctx := context.Background()

	app, err := h.Adapters.Repos.Apps.GetByID(ctx, resource.AppID)
	if err != nil {
		return domain.Resource{}, errors.New("Failed to find app")
	}

	isAppAdmin := checkForGroupMatch(userGroups, app.AdminGroups)

	if isAppAdmin {
		return domain.Resource{}, errors.New("Not authorized")
	}

	resourceToSave := domain.Resource{
		ID:           uuid.NewString(),
		Name:         resource.Name,
		FriendlyName: app.FriendlyName,
		AppID:        resource.AppID,
	}

	savedResource, err := h.Adapters.Repos.Resources.Add(ctx, resourceToSave)
	if err != nil {
		return domain.Resource{}, nil
	}

	return savedResource, nil
}
