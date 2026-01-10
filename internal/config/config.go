package config

import (
	"errors"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type HTTPServer struct {
	Host        string        `yaml:"host" env:"HTTP_HOST" env-default:"0.0.0.0"`
	Port        string        `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
	Timeout     time.Duration `yaml:"timeout" env:"HTTP_TIMEOUT" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idleTimeout" env:"HTTP_IDLE_TIMEOUT" env-default:"60s"`
}

type Config struct {
	AppEnv     string `yaml:"appEnv" env:"APP_ENV" env-default:"local"`
	StorageDSN string `yaml:"storageDsn" env:"STORAGE_DSN" env-required:"true"`
	HTTPServer `yaml:"http"`
}

func Get() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	configPath := os.Getenv("CONFIG_PATH")

	if configPath != "" {
		if err := ReadFromFile(configPath, &cfg); err != nil {
			return nil, err
		}
	}

	return &cfg, nil
}

func ReadFromFile[T any](configPath string, cfg *T) error {
	if _, err := os.Stat(configPath); configPath != "" && os.IsNotExist(err) {
		return errors.New("Can't find config file at path: " + configPath)
	}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		return err
	}

	fmt.Println("Read config file:", configPath, cfg)
	return nil
}
