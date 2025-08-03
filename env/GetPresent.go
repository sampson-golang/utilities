package env

func GetPresent(key string, fallbacks ...string) string {
	value, _ := LookupPresent(key, fallbacks...)
	return value
}
