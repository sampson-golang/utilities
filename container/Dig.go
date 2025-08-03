package container

import (
	"strconv"
)

// Dig traverses a nested map/slice structure using a path of keys/indexes.
// Example: Dig(m, "foo", 0, "bar") will return a pointer to m["foo"][0]["bar"] if it exists.
func Dig(data interface{}, path ...any) interface{} {
	current := data
	for _, key := range path {
		switch c := current.(type) {
		case map[string]interface{}:
			ks, ok := key.(string)
			if !ok {
				return nil
			}
			val, exists := c[ks]
			if !exists {
				return nil
			}
			current = val
		case []interface{}:
			var idx int
			switch k := key.(type) {
			case int:
				idx = k
			case string:
				parsed, err := strconv.Atoi(k)
				if err != nil {
					return nil
				}
				idx = parsed
			default:
				return nil
			}
			if idx < 0 || idx >= len(c) {
				return nil
			}
			current = c[idx]
		default:
			return nil
		}
	}
	return &current
}
