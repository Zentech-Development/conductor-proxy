package handlers

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/google/uuid"
)

type ServiceHandler struct {
	Adapters *domain.Adapters
}

func NewServiceHandler(adapters *domain.Adapters) ServiceHandler {
	return ServiceHandler{
		Adapters: adapters,
	}
}

func (h ServiceHandler) GetByID(id string, userGroups []string) (domain.Service, error) {
	ctx := context.Background()

	service, err := h.Adapters.Repos.Services.GetByID(ctx, id)
	if err != nil {
		return domain.Service{}, err
	}

	isServiceUser := checkForGroupMatch(userGroups, service.UserGroups)
	isServiceAdmin := checkForGroupMatch(userGroups, service.AdminGroups)

	if !isAdmin(userGroups) && !isServiceAdmin && !isServiceUser {
		return domain.Service{}, errors.New("not authorized")
	}

	return service, nil
}

func (h ServiceHandler) Add(service domain.ServiceInput, userGroups []string) (domain.Service, error) {
	ctx := context.Background()

	if !isAdmin(userGroups) {
		return domain.Service{}, errors.New("not authorized")
	}

	service.AdminGroups = slices.Compact(append(service.AdminGroups, domain.GroupNameAdmin))

	for _, group := range service.AdminGroups {
		if _, err := h.Adapters.Repos.Groups.GetByName(ctx, group); err != nil {
			return domain.Service{}, fmt.Errorf("group name %s not found", group)
		}
	}

	serviceToSave := domain.Service{
		ID:           uuid.NewString(),
		Name:         service.Name,
		FriendlyName: service.FriendlyName,
		Host:         service.Host,
		AdminGroups:  service.AdminGroups,
		UserGroups:   service.UserGroups,
		Type:         "http",
	}

	savedService, err := h.Adapters.Repos.Services.Add(ctx, serviceToSave)
	if err != nil {
		return domain.Service{}, nil
	}

	return savedService, nil
}
