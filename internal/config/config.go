package config

import (
	"errors"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type HTTPServer struct {
	Address     string        `yaml:"address" env:"HTTP_ADDRESS" env-default:":8080"`
	Timeout     time.Duration `yaml:"timeout" env:"HTTP_TIMEOUT" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idleTimeout" env:"HTTP_IDLE_TIMEOUT" env-default:"60s"`
}

type Config struct {
	AppEnv     string `yaml:"appEnv" env:"APP_ENV" env-default:"local"`
	StorageDSN string `yaml:"storageDsn" env:"STORAGE_DSN" env-required:"true"`
	HTTPServer
}

func Get() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	configPath := os.Getenv("CONFIG_PATH")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, errors.New("can't find config file at path: " + configPath)
	}

	if path := os.Getenv("CONFIG_PATH"); path != "" {
		if err := cleanenv.ReadConfig(path, &cfg); err != nil {
			return nil, err
		}
	}

	fmt.Println("Read config file:", configPath, cfg)
	return &cfg, nil
}
