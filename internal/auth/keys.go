package auth

import (
	"github.com/nafine/url-shortener/internal/config"
)

type APIKey struct {
	User string `yaml:"user"`
	Key  string `yaml:"key"`
}

type ApiKeys struct {
	Keys []APIKey `yaml:"keys"`
}

func LoadKeys() (map[string]string, error) {
	configPath := "/etc/shortener/apiKeys.yaml"

	var apiKeys ApiKeys

	if err := config.ReadFromFile(configPath, &apiKeys); err != nil {
		return nil, err
	}

	valid := make(map[string]string)
	for _, k := range apiKeys.Keys {
		valid[k.Key] = k.User
	}

	return valid, nil
}
