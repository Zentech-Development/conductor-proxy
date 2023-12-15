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

	isAdmin := checkForGroupMatch(userGroups, make([]string, 0))

	if !isAdmin {
		return domain.Group{}, errors.New("Not authorized")
	}

	// TODO: check if group name already exists
	if group.Name == "admin" {
		return domain.Group{}, errors.New("Group name is not allowed")
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
