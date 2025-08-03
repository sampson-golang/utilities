package env

import (
	"os"
)

func Lookup(key string, fallbacks ...string) (string, bool) {
	value, exists := os.LookupEnv(key)

	if exists {
		return value, true
	}

	if len(fallbacks) == 0 {
		return "", false
	}

	if len(fallbacks) > 1 {
		return Lookup(fallbacks[0], fallbacks[1:]...)
	}

	return fallbacks[0], false
}
