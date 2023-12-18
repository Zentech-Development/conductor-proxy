package main

import (
	"log"

	redisAdapters "github.com/Zentech-Development/conductor-proxy/adapters/redis"
	bindings "github.com/Zentech-Development/conductor-proxy/bindings/gin"
	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/Zentech-Development/conductor-proxy/handlers"
	"github.com/Zentech-Development/conductor-proxy/pkg/config"
)

func main() {
	config := config.GetConfig()

	redisRepo := redisAdapters.NewRedisRepo(redisAdapters.RedisRepoConfig{
		Host:     config.RedisHost,
		Password: config.RedisPassword,
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
		Proxy:     handlers.NewProxyHandler(&adapters),
	}

	server := bindings.NewHTTPServerBinding(handlers, bindings.HTTPServerBindingConfig{
		SecretKey: config.SecretKey,
		GinMode:   config.GinMode,
	})

	log.Fatal(server.Run(config.Host))
}
