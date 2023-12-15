package cmd

import (
	"log"

	redisAdapters "github.com/Zentech-Development/conductor-proxy/adapters/redis"
	"github.com/Zentech-Development/conductor-proxy/bindings"
	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/Zentech-Development/conductor-proxy/handlers"
)

func main() {
	redisRepo := redisAdapters.NewRedisRepo(redisAdapters.RedisRepoConfig{
		Host:     "localhost:7329",
		Password: "password123",
	})

	adapters := domain.Adapters{
		Repos: domain.Repos{
			Resources: redisRepo.Resources,
			Apps:      redisRepo.Apps,
			Groups:    redisRepo.Groups,
			Accounts:  redisRepo.Accounts,
		},
	}

	handlers := domain.Handlers{
		Apps:      handlers.NewAppHandler(&adapters),
		Accounts:  handlers.NewAccountHandler(&adapters),
		Groups:    handlers.NewGroupHandler(&adapters),
		Resources: handlers.NewResourceHandler(&adapters),
	}

	server := bindings.NewHTTPServerBinding(handlers, bindings.HTTPServerBindingConfig{
		Host:      ":8000",
		SecretKey: "asdf1234",
		GinMode:   "debug",
	})

	log.Fatal(server.Run())
}
