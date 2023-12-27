package main

import (
	"log"

	mockAdapters "github.com/Zentech-Development/conductor-proxy/adapters/mockdb"
	redisAdapters "github.com/Zentech-Development/conductor-proxy/adapters/redis"
	bindings "github.com/Zentech-Development/conductor-proxy/bindings/gin"
	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/Zentech-Development/conductor-proxy/handlers"
	"github.com/Zentech-Development/conductor-proxy/pkg/config"
	"github.com/gin-gonic/gin"
)

func setupApp(conf *config.ConductorConfig) *gin.Engine {
	var repos domain.Repos

	switch conf.Database {
	case config.DatabaseTypeRedis:
		repos = redisAdapters.NewRedisRepo(redisAdapters.RedisRepoConfig{
			Host:     conf.RedisHost,
			Password: conf.RedisPassword,
		})
	case config.DatabaseTypeMock:
		repos = mockAdapters.NewMockDB()
	default:
		panic("Bad database type")
	}

	adapters := domain.Adapters{
		Repos: repos,
	}

	handlers := domain.Handlers{
		Services:  handlers.NewServiceHandler(&adapters),
		Accounts:  handlers.NewAccountHandler(&adapters),
		Groups:    handlers.NewGroupHandler(&adapters),
		Resources: handlers.NewResourceHandler(&adapters),
		Proxy:     handlers.NewProxyHandler(&adapters),
	}

	return bindings.NewHTTPServerBinding(handlers, bindings.HTTPServerBindingConfig{
		SecretKey: conf.SecretKey,
		GinMode:   conf.GinMode,
	})
}

func main() {
	config := config.GetConfig()
	server := setupApp(config)
	log.Fatal(server.Run(config.Host))
}
