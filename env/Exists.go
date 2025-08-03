package env

func Exists(key string) bool {
	_, exists := Lookup(key)
	return exists
}
