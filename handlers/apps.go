package handlers

import (
	"context"
	"errors"

	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/google/uuid"
)

type AppHandler struct {
	Adapters *domain.Adapters
}

func NewAppHandler(adapters *domain.Adapters) AppHandler {
	return AppHandler{
		Adapters: adapters,
	}
}

func (h AppHandler) GetByID(id string, userGroups []string) (domain.App, error) {
	ctx := context.Background()

	app, err := h.Adapters.Repos.Apps.GetByID(ctx, id)
	if err != nil {
		return domain.App{}, err
	}

	appHasGroups := (len(app.AdminGroups) + len(app.UserGroups)) > 0
	isAppUser := checkForGroupMatch(userGroups, app.UserGroups)
	isAppAdmin := checkForGroupMatch(userGroups, app.AdminGroups)

	if appHasGroups && !isAppAdmin && !isAppUser {
		return domain.App{}, errors.New("Not authorized")
	}

	return app, nil
}

func (h AppHandler) Add(app domain.AppInput, userGroups []string) (domain.App, error) {
	ctx := context.Background()

	isAdmin := checkForGroupMatch(userGroups, make([]string, 0))

	if !isAdmin {
		return domain.App{}, errors.New("Not authorized")
	}

	appToSave := domain.App{
		ID:           uuid.NewString(),
		Name:         app.Name,
		FriendlyName: app.FriendlyName,
		Host:         app.Host,
		AdminGroups:  app.AdminGroups,
		UserGroups:   app.UserGroups,
	}

	savedApp, err := h.Adapters.Repos.Apps.Add(ctx, appToSave)
	if err != nil {
		return domain.App{}, nil
	}

	return savedApp, nil
}
