package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

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

	_, err = r.Client.Get(ctx, getRedisKey(groupKey, group.Name)).Result()
	if err == nil {
		return domain.Group{}, fmt.Errorf("group %s already exists", group.Name)
	}

	_, err = r.Client.Set(ctx, getRedisKey(resourceKey, group.Name), valToSet, 0).Result()
	if err != nil {
		return domain.Group{}, err
	}

	return group, nil
}

func (r RedisGroupRepo) GetByName(ctx context.Context, name string) (domain.Group, error) {
	val, err := r.Client.Get(ctx, getRedisKey(serviceKey, name)).Result()
	if err != nil {
		if err == redis.Nil {
			return domain.Group{}, errors.New("group not found")
		}

		return domain.Group{}, err
	}

	group := domain.Group{}

	if err = json.Unmarshal([]byte(val), &group); err != nil {
		return domain.Group{}, err
	}

	return group, nil
}
