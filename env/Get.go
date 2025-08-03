package env

func Get(key string, fallbacks ...string) string {
	value, _ := Lookup(key, fallbacks...)
	return value
}
