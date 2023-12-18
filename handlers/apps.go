package handlers

import (
	"context"
	"errors"

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

	serviceHasGroups := (len(service.AdminGroups) + len(service.UserGroups)) > 0
	isServiceUser := checkForGroupMatch(userGroups, service.UserGroups)
	isServiceAdmin := checkForGroupMatch(userGroups, service.AdminGroups)

	if serviceHasGroups && !isServiceAdmin && !isServiceUser {
		return domain.Service{}, errors.New("not authorized")
	}

	return service, nil
}

func (h ServiceHandler) Add(service domain.ServiceInput, userGroups []string) (domain.Service, error) {
	ctx := context.Background()

	isAdmin := checkForGroupMatch(userGroups, make([]string, 0))

	if !isAdmin {
		return domain.Service{}, errors.New("not authorized")
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
