package adapters

import (
	"context"
	"encoding/json"

	"github.com/Zentech-Development/conductor-proxy/adapters"
	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/redis/go-redis/v9"
)

type RedisResourceRepo struct {
	Client *redis.Client
}

func newRedisResourceRepo(client *redis.Client) RedisResourceRepo {
	return RedisResourceRepo{
		Client: client,
	}
}

func (r RedisResourceRepo) GetByID(ctx context.Context, id string) (domain.Resource, error) {
	val, err := r.Client.Get(ctx, getRedisKey(resourceKey, id)).Result()
	if err != nil {
		if err == redis.Nil {
			return domain.Resource{}, &adapters.NotFoundError{Name: "resource"}
		}

		return domain.Resource{}, err
	}

	resource := domain.Resource{}

	if err = json.Unmarshal([]byte(val), &resource); err != nil {
		return domain.Resource{}, err
	}

	return resource, nil
}

func (r RedisResourceRepo) Add(ctx context.Context, resource domain.Resource) (domain.Resource, error) {
	valToSet, err := json.Marshal(&resource)
	if err != nil {
		return domain.Resource{}, err
	}

	_, err = r.Client.Get(ctx, getRedisKey(resourceKey, resource.ID)).Result()
	if err == nil {
		return domain.Resource{}, &adapters.AlreadyExistsError{Name: "resource"}
	}

	_, err = r.Client.Set(ctx, getRedisKey(resourceKey, resource.ID), valToSet, 0).Result()
	if err != nil {
		return domain.Resource{}, err
	}

	return resource, nil
}
