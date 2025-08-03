package merge

import "reflect"

// Structs assigns non-empty fields from src to dest if they are not empty in src.
// the destination/into param is first, and sources are merged in order from left to right.
func Structs(into interface{}, from ...interface{}) {
	for _, src := range from {
		destVal := reflect.ValueOf(into).Elem()
		srcVal := reflect.ValueOf(src).Elem()

		for i := 0; i < destVal.NumField(); i++ {
			destField := destVal.Field(i)
			srcField := srcVal.Field(i)

			// Check if the field is zero in src
			if !isEmpty(srcField) {
				destField.Set(srcField)
			}
		}
	}
}

// isEmpty checks if a reflect.Value is considered empty (zero value).
func isEmpty(v reflect.Value) bool {
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}
