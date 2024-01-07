package config

import (
	"encoding/json"
	"log"
	"slices"
	"sync"

	"github.com/spf13/viper"
)

var lock = &sync.Mutex{}

type ConductorConfig struct {
	Host         string `mapstructure:"CONDUCTOR_HOST"`
	DatabaseType string `mapstructure:"CONDUCTOR_DATABASE_TYPE"`
	SecureMode   bool   `mapstructure:"CONDUCTOR_SECURE_MODE"`

	AccessTokenSecret   string `mapstructure:"CONDUCTOR_ACCESS_TOKEN_SECRET" json:"-"`
	AccessTokenCost     int    `mapstructure:"CONDCUTOR_ACCESS_TOKEN_COST"`
	DefaultTokenTimeout int    `mapstructure:"CONDUCTOR_DEFAULT_TOKEN_TIMEOUT"`

	DefaultAdminUsername string `mapstructure:"CONDUCTOR_DEFAULT_ADMIN_USERNAME"`
	DefaultAdminPasskey  string `mapstructure:"CONDUCTOR_DEFAULT_ADMIN_PASSKEY" json:"-"`

	RedisHost     string `mapstructure:"CONDUCTOR_REDIS_HOST"`
	RedisPassword string `mapstructure:"CONDUCTOR_REDIS_PASSWORD" json:"-"`
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

	v.SetDefault("CONDUCTOR_HOST", "localhost:8000")
	v.SetDefault("CONDUCTOR_DATABASE_TYPE", "mock")
	v.SetDefault("CONDUCTOR_SECURE_MODE", true)
	v.SetDefault("CONDUCTOR_DEFAULT_TOKEN_TIMEOUT", 3600)
	v.SetDefault("CONDUCTOR_DEFAULT_ADMIN_USERNAME", "admin")
	v.SetDefault("CONDUCTOR_DEFAULT_ADMIN_PASSKEY", "password")
	v.SetDefault("CONDCUTOR_ACCESS_TOKEN_COST", 12)

	if envFilePath != "" {
		v.SetConfigFile(envFilePath)
		v.SetConfigType("env")
		if err := v.ReadInConfig(); err != nil {
			log.Fatal("Failed to load the configuration file: ", err)
		}
	}

	v.AutomaticEnv()

	if err := v.Unmarshal(conf); err != nil {
		log.Fatal("Invalid configuration: ", err)
	}

	validDatabaseTypes := []string{DatabaseTypeMongo, DatabaseTypePostgres, DatabaseTypeRedis, DatabaseTypeSQLite, DatabaseTypeMock}
	if !slices.Contains(validDatabaseTypes, conf.DatabaseType) {
		log.Fatalf("Invalid database type: %s", conf.DatabaseType)
	}

	if conf.SecureMode && conf.AccessTokenSecret == "" {
		log.Fatal("Must provide CONDUCTOR_ACCESS_TOKEN_SECRET")
	}

	log.Default().Println("Conductor configuration initialized")
	vals, _ := json.MarshalIndent(conf, "", "\t")
	log.Default().Println(string(vals))

	return conf
}
