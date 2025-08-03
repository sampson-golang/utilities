package merge

// Params merges the src map into the dest map.
// Params assigns from src to dest for all keys in src.
// the destination/into param is first, and sources are merged in order from left to right.
func Params(into map[string]string, from ...map[string]string) {
	for _, src := range from {
		for key, value := range src {
			into[key] = value
		}
	}
}
