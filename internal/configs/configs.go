package configs

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug       bool `yaml:"is-debug" env:"IS_DEBUG" env-default:"false"`
	IsDevelopment bool `yaml:"is-development" env:"IS-DEV" env-default:"false"`
	IsProduction  bool `yaml:"is-production" env:"IS-PROD" env-default:"false"`
	Secret        string
	HTTP          *HTTP       `yaml:"http"`
	PostgreSQL    *PostgreSQL `yaml:"postgresql"`
}

type HTTP struct {
	IP           string        `yaml:"ip" env:"HTTP-IP"`
	Port         int           `yaml:"port" env:"HTTP-PORT"`
	ReadTimeout  time.Duration `yaml:"read-timeout" env:"HTTP-READ-TIMEOUT"`
	WriteTimeout time.Duration `yaml:"write-timeout" env:"HTTP-WRITE-TIMEOUT"`
}

type PostgreSQL struct {
	Username string `yaml:"username" env:"PSQL_USERNAME" env-required:"true"`
	Password string `yaml:"password" env:"PSQL_PASSWORD" env-required:"true"`
	Host     string `yaml:"host" env:"PSQL_HOST" env-required:"true"`
	Port     string `yaml:"port" env:"PSQL_PORT" env-required:"true"`
	Database string `yaml:"database" env:"PSQL_DATABASE" env-required:"true"`
}

const (
	EnvConfigPathName  = "CONFIG-PATH"
	FlagConfigPathName = "config"
)

var (
	configPath string
	instance   *Config
	once       sync.Once
)

func GetConfig() (*Config, error) {
	once.Do(func() {
		if configPath == "" {
			configPath = os.Getenv(EnvConfigPathName)
		}

		if configPath == "" {
			// Use the current working directory as the base path for the config file
			basePath, err := os.Getwd()
			if err != nil {
				log.Fatal("Failed to get current working directory")
			}
			configPath = filepath.Join(basePath, "configs", "env.dev.yaml")
		}

		instance = &Config{}

		if err := cleanenv.ReadConfig(configPath, instance); err != nil {
			helpText := "Flow-Gateway - service proxy-api-gateway"
			help, _ := cleanenv.GetDescription(instance, &helpText)
			log.Print(help)
			log.Fatal(err)
		}
	})

	return instance, nil
}
