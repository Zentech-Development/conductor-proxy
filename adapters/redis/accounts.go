package adapters

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/redis/go-redis/v9"
)

type RedisAccountRepo struct {
	Client *redis.Client
}

func newRedisAccountRepo(client *redis.Client) RedisAccountRepo {
	return RedisAccountRepo{
		Client: client,
	}
}

func (r RedisAccountRepo) GetByID(ctx context.Context, id string) (domain.Account, error) {
	return domain.Account{}, nil
}

func (r RedisAccountRepo) Add(ctx context.Context, account domain.Account) (domain.Account, error) {
	valToSet, err := json.Marshal(&account)
	if err != nil {
		return domain.Account{}, err
	}

	_, err = r.Client.Get(ctx, getRedisKey(accountKey, account.ID)).Result()
	if err == nil {
		return domain.Account{}, errors.New("Account ID already exists")
	}

	_, err = r.Client.Set(ctx, getRedisKey(resourceKey, account.ID), valToSet, 0).Result()
	if err != nil {
		return domain.Account{}, err
	}

	return account, nil
}

func (r RedisAccountRepo) Update(ctx context.Context, account domain.Account) (domain.Account, error) {
	valToSet, err := json.Marshal(&account)
	if err != nil {
		return domain.Account{}, err
	}

	_, err = r.Client.Get(ctx, getRedisKey(accountKey, account.ID)).Result()
	if err != nil {
		return domain.Account{}, errors.New("Account not found")
	}

	_, err = r.Client.Set(ctx, getRedisKey(resourceKey, account.ID), valToSet, 0).Result()
	if err != nil {
		return domain.Account{}, err
	}

	return account, nil
}
