package auth

import (
	"os"
	"strings"
)

func LoadKeys() map[string]string {
	raw := os.Getenv("API_KEYS")
	keys := make(map[string]string)
	for _, pair := range strings.Split(raw, ",") {
		user, key, ok := strings.Cut(pair, ":")
		if !ok {
			continue
		}
		keys[user] = key
	}
	return keys
}
