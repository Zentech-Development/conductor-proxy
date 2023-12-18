package adapters

import (
	"context"
	"fmt"

	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/redis/go-redis/v9"
)

const (
	resourceKey = "resource"
	serviceKey  = "service"
	accountKey  = "account"
	groupKey    = "group"
)

type RedisRepo struct {
	Services  RedisServiceRepo
	Accounts  RedisAccountRepo
	Resources RedisResourceRepo
	Groups    RedisGroupRepo
}

type RedisRepoConfig struct {
	Host     string
	Password string
	Mock     string
}

func NewRedisRepo(config RedisRepoConfig) domain.Repos {
	client := redis.NewClient(&redis.Options{
		Addr: config.Host,
		// Password: config.Password,
		DB: 0,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		fmt.Print(err.Error())
		panic("Connection to Redis failed")
	}

	return domain.Repos{
		Services:  newRedisServiceRepo(client),
		Accounts:  newRedisAccountRepo(client),
		Resources: newRedisResourceRepo(client),
		Groups:    newRedisGroupRepo(client),
	}
}

func getRedisKey(model string, id string) string {
	return fmt.Sprintf("%s:%s", model, id)
}
