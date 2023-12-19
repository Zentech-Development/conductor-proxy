package handlers

import (
	"context"
	"errors"

	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/google/uuid"
)

type GroupHandler struct {
	Adapters *domain.Adapters
}

func NewGroupHandler(adapters *domain.Adapters) GroupHandler {
	return GroupHandler{
		Adapters: adapters,
	}
}

func (h GroupHandler) Add(group domain.GroupInput, userGroups []string) (domain.Group, error) {
	ctx := context.Background()

	if !isAdmin(userGroups) {
		return domain.Group{}, errors.New("not authorized")
	}

	if _, err := h.Adapters.Repos.Groups.GetByName(ctx, group.Name); err == nil {
		return domain.Group{}, errors.New("group name already exists")
	}

	if group.Name == domain.GroupNameAdmin {
		return domain.Group{}, errors.New("group name is not allowed")
	}

	groupToSave := domain.Group{
		ID:   uuid.NewString(),
		Name: group.Name,
	}

	savedGroup, err := h.Adapters.Repos.Groups.Add(ctx, groupToSave)
	if err != nil {
		return domain.Group{}, nil
	}

	return savedGroup, nil
}
