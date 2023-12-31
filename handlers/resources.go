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

	service, err := h.Adapters.Repos.Services.GetByID(ctx, resource.ServiceID)
	if err != nil {
		return domain.Resource{}, errors.New("failed to find service")
	}

	isServiceUser := checkForGroupMatch(userGroups, service.UserGroups)
	isServiceAdmin := checkForGroupMatch(userGroups, service.AdminGroups)

	if !isAdmin(userGroups) && !isServiceAdmin && !isServiceUser {
		return domain.Resource{}, errors.New("not authorized")
	}

	return resource, nil
}

func (h ResourceHandler) Add(resource domain.ResourceInput, userGroups []string) (domain.Resource, error) {
	ctx := context.Background()

	service, err := h.Adapters.Repos.Services.GetByID(ctx, resource.ServiceID)
	if err != nil {
		return domain.Resource{}, errors.New("failed to find service")
	}

	isServiceAdmin := checkForGroupMatch(userGroups, service.AdminGroups)

	if !isAdmin(userGroups) && !isServiceAdmin {
		return domain.Resource{}, errors.New("not authorized")
	}

	resourceToSave := domain.Resource{
		ID:           uuid.NewString(),
		Name:         resource.Name,
		FriendlyName: resource.FriendlyName,
		ServiceID:    resource.ServiceID,
		Properties:   resource.Properties,
		Endpoints:    resource.Endpoints,
	}

	savedResource, err := h.Adapters.Repos.Resources.Add(ctx, resourceToSave)
	if err != nil {
		return domain.Resource{}, nil
	}

	return savedResource, nil
}
