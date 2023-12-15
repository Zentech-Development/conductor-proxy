package adapters

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/redis/go-redis/v9"
)

type RedisAppRepo struct {
	Client *redis.Client
}

func newRedisAppRepo(client *redis.Client) RedisAppRepo {
	return RedisAppRepo{
		Client: client,
	}
}

func (r RedisAppRepo) GetByID(ctx context.Context, id string) (domain.App, error) {
	val, err := r.Client.Get(ctx, getRedisKey(appKey, id)).Result()
	if err != nil {
		if err == redis.Nil {
			return domain.App{}, errors.New("App not found")
		}

		return domain.App{}, err
	}

	app := domain.App{}

	if err = json.Unmarshal([]byte(val), &app); err != nil {
		return domain.App{}, err
	}

	return app, nil
}

func (r RedisAppRepo) Add(ctx context.Context, app domain.App) (domain.App, error) {
	valToSet, err := json.Marshal(&app)
	if err != nil {
		return domain.App{}, err
	}

	_, err = r.Client.Get(ctx, getRedisKey(appKey, app.ID)).Result()
	if err == nil {
		return domain.App{}, errors.New("App ID already exists")
	}

	_, err = r.Client.Set(ctx, getRedisKey(resourceKey, app.ID), valToSet, 0).Result()
	if err != nil {
		return domain.App{}, err
	}

	return app, nil
}
