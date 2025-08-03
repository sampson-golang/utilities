package env

func LookupPresent(key string, fallbacks ...string) (string, bool) {
	value, exists := Lookup(key)

	if exists && value != "" {
		return value, true
	}

	if len(fallbacks) == 0 {
		return "", false
	}

	if len(fallbacks) > 1 {
		return LookupPresent(fallbacks[0], fallbacks[1:]...)
	}

	return fallbacks[0], false
}
