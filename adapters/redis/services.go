package adapters

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/redis/go-redis/v9"
)

type RedisServiceRepo struct {
	Client *redis.Client
}

func newRedisServiceRepo(client *redis.Client) RedisServiceRepo {
	return RedisServiceRepo{
		Client: client,
	}
}

func (r RedisServiceRepo) GetByID(ctx context.Context, id string) (domain.Service, error) {
	val, err := r.Client.Get(ctx, getRedisKey(serviceKey, id)).Result()
	if err != nil {
		if err == redis.Nil {
			return domain.Service{}, errors.New("service not found")
		}

		return domain.Service{}, err
	}

	service := domain.Service{}

	if err = json.Unmarshal([]byte(val), &service); err != nil {
		return domain.Service{}, err
	}

	return service, nil
}

func (r RedisServiceRepo) Add(ctx context.Context, service domain.Service) (domain.Service, error) {
	valToSet, err := json.Marshal(&service)
	if err != nil {
		return domain.Service{}, err
	}

	_, err = r.Client.Get(ctx, getRedisKey(serviceKey, service.ID)).Result()
	if err == nil {
		return domain.Service{}, errors.New("service ID already exists")
	}

	_, err = r.Client.Set(ctx, getRedisKey(serviceKey, service.ID), valToSet, 0).Result()
	if err != nil {
		return domain.Service{}, err
	}

	return service, nil
}
