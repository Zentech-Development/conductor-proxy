package config

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

var lock = &sync.Mutex{}

type ConductorConfig struct {
	Host                 string
	Database             string
	Secure               bool
	SecretKey            string
	GinMode              string
	DefaultTokenTimeout  int
	RedisHost            string
	RedisPassword        string
	DefaultAdminUsername string
	DefaultAdminPasskey  string
}

const (
	DatabaseTypeRedis    = "redis"
	DatabaseTypeMongo    = "mongo"
	DatabaseTypeSQLite   = "sqlite"
	DatabaseTypePostgres = "postgres"
)

var config *ConductorConfig

func GetConfig() *ConductorConfig {
	// extra check here to avoid using the (very expensive) lock whenever possible
	if config == nil {
		lock.Lock()
		defer lock.Unlock()

		if config == nil {
			config = NewConfig("./pkg/config/.env")
		}
	}

	return config
}

func NewConfig(envFilePath string) *ConductorConfig {
	if err := godotenv.Load(envFilePath); err != nil {
		fmt.Println("Failed to load .env file, using system environment variables")
	}

	conductorConfig := &ConductorConfig{
		Host:                 "localhost:8080",
		Database:             DatabaseTypeRedis,
		Secure:               true,
		DefaultTokenTimeout:  60 * 60,
		GinMode:              "release",
		SecretKey:            "",
		DefaultAdminUsername: "admin",
		DefaultAdminPasskey:  "admin",
	}

	conductorConfig.SecretKey = os.Getenv("CONDUCTOR_SECRET_KEY")
	if conductorConfig.SecretKey == "" {
		panic("Must provide a CONDUCTOR_SECRET_KEY environment variable")
	}

	if defaultTimeout := os.Getenv("CONDUCTOR_DEFAULT_TOKEN_TIMEOUT"); defaultTimeout != "" {
		parsedVal, err := strconv.Atoi(defaultTimeout)
		if err != nil {
			panic("Bad value supplied for CONDUCTOR_DEFAULT_TOKEN_TIMEOUT")
		}
		conductorConfig.DefaultTokenTimeout = parsedVal
	}

	if ginMode := os.Getenv("CONDUCTOR_GIN_MODE"); ginMode == "debug" {
		conductorConfig.GinMode = "debug"
	}

	if secure := os.Getenv("CONDUCTOR_SECURE"); secure == "false" {
		conductorConfig.Secure = false
	}

	if host := os.Getenv("CONDUCTOR_HOST"); host != "" {
		conductorConfig.Host = host
	}

	if database := os.Getenv("CONDUCTOR_DATABASE"); database != "" {
		validDatabases := []string{DatabaseTypeMongo, DatabaseTypePostgres, DatabaseTypeRedis, DatabaseTypeSQLite}
		if !slices.Contains(validDatabases, database) {
			panic("Invalid value for CONDUCTOR_DATABASE")
		}
		conductorConfig.Database = database
	}

	if conductorConfig.Database == DatabaseTypeRedis {
		conductorConfig.RedisHost = os.Getenv("CONDUCTOR_REDIS_HOST")
		conductorConfig.RedisPassword = os.Getenv("CONDUCTOR_REDIS_PASSWORD")
	}

	if defaultAdminUsername := os.Getenv("CONDUCTOR_DEFAULT_ADMIN_USERNAME"); defaultAdminUsername != "" {
		conductorConfig.DefaultAdminUsername = defaultAdminUsername
	}

	if defaultAdminPasskey := os.Getenv("CONDUCTOR_DEFAULT_ADMIN_PASSKEY"); defaultAdminPasskey != "" {
		conductorConfig.DefaultAdminPasskey = defaultAdminPasskey
	}

	return conductorConfig
}
