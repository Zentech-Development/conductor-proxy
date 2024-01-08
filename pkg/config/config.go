package config

import (
	"encoding/json"
	"log"
	"os"
	"slices"
	"sync"

	"github.com/spf13/viper"
)

var lock = &sync.Mutex{}

type ConductorConfig struct {
	Host         string `mapstructure:"HOST"`
	DatabaseType string `mapstructure:"DATABASE_TYPE"`
	SecureMode   bool   `mapstructure:"SECURE_MODE"`

	AccessTokenSecret   string `mapstructure:"ACCESS_TOKEN_SECRET_KEY" json:"-"`
	AccessTokenCost     int    `mapstructure:"ACCESS_TOKEN_COST"`
	DefaultTokenTimeout int    `mapstructure:"DEFAULT_TOKEN_TIMEOUT_SECONDS"`

	DefaultAdminUsername string `mapstructure:"DEFAULT_ADMIN_USERNAME"`
	DefaultAdminPasskey  string `mapstructure:"DEFAULT_ADMIN_PASSKEY" json:"-"`

	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD" json:"-"`
}

const (
	DatabaseTypeRedis    = "redis"
	DatabaseTypeMongo    = "mongo"
	DatabaseTypeSQLite   = "sqlite"
	DatabaseTypePostgres = "postgres"
	DatabaseTypeMock     = "mock"
)

var conductorConfig *ConductorConfig

func SetAndGetConfig(configFilePath string) *ConductorConfig {
	// extra check here to avoid using the (very expensive) lock whenever possible
	if conductorConfig == nil {
		lock.Lock()
		defer lock.Unlock()

		if conductorConfig == nil {
			conductorConfig = NewConfig(configFilePath)
		}
	}

	return conductorConfig
}

func GetConfig() *ConductorConfig {
	return conductorConfig
}

func NewConfig(envFilePath string) *ConductorConfig {
	conf := &ConductorConfig{}

	v := viper.New()

	v.SetDefault("HOST", "localhost:8000")
	v.SetDefault("DATABASE_TYPE", "mock")
	v.SetDefault("SECURE_MODE", true)
	v.SetDefault("DEFAULT_TOKEN_TIMEOUT_SECONDS", 3600)
	v.SetDefault("DEFAULT_ADMIN_USERNAME", "admin")
	v.SetDefault("DEFAULT_ADMIN_PASSKEY", "password")
	v.SetDefault("ACCESS_TOKEN_COST", 12)

	if envFilePath != "" {
		v.SetConfigFile(envFilePath)
		v.SetConfigType("env")
		if err := v.ReadInConfig(); err != nil {
			log.Fatal("Failed to load the configuration file: ", err)
		}
	}

	secretKey := os.Getenv("CONDUCTOR_SECRET")
	if secretKey == "" && v.GetString("ACCESS_TOKEN_SECRET_KEY") == "" {
		log.Fatal("Must supply CONDUCTOR_SECRET env variable or ACCESS_TOKEN_SECRET_KEY in config file")
	}

	if secretKey != "" {
		v.Set("ACCESS_TOKEN_SECRET_KEY", secretKey)
	}

	if err := v.Unmarshal(conf); err != nil {
		log.Fatal("Invalid configuration: ", err)
	}

	validDatabaseTypes := []string{DatabaseTypeMongo, DatabaseTypePostgres, DatabaseTypeRedis, DatabaseTypeSQLite, DatabaseTypeMock}
	if !slices.Contains(validDatabaseTypes, conf.DatabaseType) {
		log.Fatalf("Invalid database type: %s", conf.DatabaseType)
	}

	log.Default().Println("Conductor configuration initialized")
	vals, _ := json.MarshalIndent(conf, "", "\t")
	log.Default().Println(string(vals))

	return conf
}
