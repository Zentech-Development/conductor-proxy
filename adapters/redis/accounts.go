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

func (r RedisAccountRepo) GetByUsername(ctx context.Context, id string) (domain.Account, error) {
	val, err := r.Client.Get(ctx, getRedisKey(accountKey, id)).Result()
	if err != nil {
		if err == redis.Nil {
			return domain.Account{}, errors.New("account not found")
		}

		return domain.Account{}, err
	}

	account := domain.Account{}

	if err = json.Unmarshal([]byte(val), &account); err != nil {
		return domain.Account{}, err
	}

	return account, nil
}

func (r RedisAccountRepo) Add(ctx context.Context, account domain.Account) (domain.Account, error) {
	valToSet, err := json.Marshal(&account)
	if err != nil {
		return domain.Account{}, err
	}

	_, err = r.Client.Get(ctx, getRedisKey(accountKey, account.Username)).Result()
	if err == nil {
		return domain.Account{}, errors.New("account username already exists")
	}

	_, err = r.Client.HSet(ctx, getRedisKey(accountKey, account.Username), valToSet, 0).Result()
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

	_, err = r.Client.Get(ctx, getRedisKey(accountKey, account.Username)).Result()
	if err != nil {
		return domain.Account{}, errors.New("Account not found")
	}

	_, err = r.Client.Set(ctx, getRedisKey(accountKey, account.Username), valToSet, 0).Result()
	if err != nil {
		return domain.Account{}, err
	}

	return account, nil
}
