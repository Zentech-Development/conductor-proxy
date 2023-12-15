package adapters

import (
	"fmt"

	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/redis/go-redis/v9"
)

const (
	resourceKey = "resource"
	appKey      = "app"
	accountKey  = "account"
	groupKey    = "group"
)

type RedisRepo struct {
	Apps      RedisAppRepo
	Accounts  RedisAccountRepo
	Resources RedisResourceRepo
	Groups    RedisGroupRepo
}

type RedisRepoConfig struct {
	Host     string
	Password string
}

func NewRedisRepo(config RedisRepoConfig) domain.Repos {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Host,
		Password: config.Password,
		DB:       0,
	})

	return domain.Repos{
		Apps:      newRedisAppRepo(client),
		Accounts:  newRedisAccountRepo(client),
		Resources: newRedisResourceRepo(client),
		Groups:    newRedisGroupRepo(client),
	}
}

func getRedisKey(model string, id string) string {
	return fmt.Sprintf("%s:%s", model, id)
}
