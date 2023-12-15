package adapters

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/redis/go-redis/v9"
)

type RedisGroupRepo struct {
	Client *redis.Client
}

func newRedisGroupRepo(client *redis.Client) RedisGroupRepo {
	return RedisGroupRepo{
		Client: client,
	}
}

func (r RedisGroupRepo) Add(ctx context.Context, group domain.Group) (domain.Group, error) {
	valToSet, err := json.Marshal(&group)
	if err != nil {
		return domain.Group{}, err
	}

	_, err = r.Client.Get(ctx, getRedisKey(groupKey, group.ID)).Result()
	if err == nil {
		return domain.Group{}, errors.New("Group ID already exists")
	}

	_, err = r.Client.Set(ctx, getRedisKey(resourceKey, group.ID), valToSet, 0).Result()
	if err != nil {
		return domain.Group{}, err
	}

	return group, nil
}
