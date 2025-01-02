package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var instance *Config

type Config struct {
	HTTP_PORT             string `envconfig:"HTTP_PORT"  required:"true"`
	EXECUTABLES_JSON_PATH string `envconfig:"EXECUTABLES_JSON_PATH" required:"true"`
	AUTOSETRUN            bool   `envconfig:"AUTOSETRUN" required:"true"`
}

func load() (*Config, error) {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return nil, fmt.Errorf("failed to load configuration from the environment: %w", err)
	}
	return &c, nil
}

func mustLoad() {
	_ = godotenv.Load()

	config, err := load()
	if err != nil {
		panic(err)
	}
	instance = config
}

func GetConfig() *Config {
	mustLoad()

	return instance
}
