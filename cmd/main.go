package main

import (
	"flag"
	"log"

	mockAdapters "github.com/Zentech-Development/conductor-proxy/adapters/mockdb"
	redisAdapters "github.com/Zentech-Development/conductor-proxy/adapters/redis"
	bindings "github.com/Zentech-Development/conductor-proxy/bindings/gin"
	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/Zentech-Development/conductor-proxy/handlers"
	"github.com/Zentech-Development/conductor-proxy/pkg/config"
	conf "github.com/Zentech-Development/conductor-proxy/pkg/config"
	"github.com/gin-gonic/gin"
)

func setupApp(config *conf.ConductorConfig) *gin.Engine {
	var repos domain.Repos

	switch config.DatabaseType {
	case conf.DatabaseTypeRedis:
		repos = redisAdapters.NewRedisRepo(redisAdapters.RedisRepoConfig{
			Host:     config.RedisHost,
			Password: config.RedisPassword,
		})
	case conf.DatabaseTypeMock:
		repos = mockAdapters.NewMockDB()
	default:
		log.Fatalf("Bad database type: %s", config.DatabaseType)
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

	ginMode := gin.ReleaseMode
	if !config.SecureMode {
		ginMode = gin.DebugMode
	}

	log.Default().Printf("Starting app in %s mode\n", ginMode)

	return bindings.NewHTTPServerBinding(handlers, bindings.HTTPServerBindingConfig{
		SecretKey: config.AccessTokenSecret,
		GinMode:   ginMode,
	})
}

func main() {
	var configFilePath string
	flag.StringVar(&configFilePath, "config", "", "optional, path to config file")
	flag.Parse()

	config := config.SetAndGetConfig(configFilePath)
	server := setupApp(config)
	log.Fatal(server.Run(config.Host))
}
