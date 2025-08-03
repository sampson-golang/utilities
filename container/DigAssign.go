package container

import (
	"reflect"
)

// DigAssign assigns the result of Dig(data, path...) to the field of result with the given key.
// If Dig returns nil, the field is not assigned.
// If the value is not assignable or not convertible to result[key], type it is not assigned.
func DigAssign(result interface{}, key string, data interface{}, path ...any) {
	if value := Dig(data, path...); value != nil {
		val := reflect.ValueOf(value)

		if !val.IsValid() || val.IsNil() {
			return
		}

		for val.Kind() == reflect.Interface || (val.Kind() == reflect.Ptr && !val.IsNil()) {
			val = val.Elem()
		}

		field := reflect.ValueOf(result).Elem().FieldByName(key)
		if !field.IsValid() {
			return
		}

		if val.Type().AssignableTo(field.Type()) {
			field.Set(val)
		} else if val.Type().ConvertibleTo(field.Type()) {
			field.Set(val.Convert(field.Type()))
		}
	}
}
