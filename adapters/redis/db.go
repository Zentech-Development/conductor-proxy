package adapters

import (
	"context"
	"fmt"

	"github.com/Zentech-Development/conductor-proxy/domain"
	conf "github.com/Zentech-Development/conductor-proxy/pkg/config"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
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

	repos := domain.Repos{
		Services:  newRedisServiceRepo(client),
		Accounts:  newRedisAccountRepo(client),
		Resources: newRedisResourceRepo(client),
		Groups:    newRedisGroupRepo(client),
	}

	if _, err := repos.Groups.GetByName(context.Background(), domain.GroupNameAdmin); err != nil {
		fmt.Println("admin group not found, creating")

		adminGroup := domain.Group{
			ID:   uuid.NewString(),
			Name: domain.GroupNameAdmin,
		}

		if _, err := repos.Groups.Add(context.Background(), adminGroup); err != nil {
			panic("Failed to automatically create admin group")
		}
	}

	conductorConfig := conf.GetConfig()

	addFirstAdminUserIfRequired(client, &repos, conductorConfig.DefaultAdminUsername, conductorConfig.DefaultAdminPasskey)

	return repos
}

func getRedisKey(model string, id string) string {
	return fmt.Sprintf("%s:%s", model, id)
}

func addFirstAdminUserIfRequired(r *redis.Client, repos *domain.Repos, username string, passkey string) {
	if _, err := repos.Accounts.GetByUsername(context.Background(), username); err == nil {
		return
	}

	if _, err := r.Get(context.Background(), firstAdminFlagName).Result(); err != nil {
		adminUser := domain.Account{
			ID:              uuid.NewString(),
			Username:        username,
			Passkey:         passkey,
			Groups:          []string{domain.GroupNameAdmin},
			TokenExpiration: 3600,
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(adminUser.Passkey), 12)
		if err != nil {
			panic(err)
		}

		adminUser.Passkey = string(hash)

		if _, err := repos.Accounts.Add(context.Background(), adminUser); err != nil {
			panic("Failed to automically create admin user")
		}
	}
}
