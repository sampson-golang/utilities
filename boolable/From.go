package boolable

import (
	"reflect"
	"strings"
)

var falseValues = map[string]bool{
	"":      true,
	"0":     true,
	"f":     true,
	"false": true,
	"off":   true,
	"n":     true,
	"no":    true,
}

func fromPointer(value *reflect.Value) bool {
	if value.IsNil() {
		return false
	}

	elem := value.Elem()
	switch elem.Kind() {
	case reflect.Bool:
		return elem.Bool()
	case reflect.Int:
		return elem.Int() != 0
	case reflect.String:
		return !falseValues[strings.ToLower(elem.String())]
	case reflect.Interface:
		return fromPointer(&elem)
	default:
		return true
	}
}

func From(value interface{}, dereference ...bool) bool {
	reflectedValue := reflect.ValueOf(value)

	if reflectedValue.Kind() == reflect.Ptr && (len(dereference) == 0 || dereference[0]) {
		return fromPointer(&reflectedValue)
	}

	switch v := value.(type) {
	case bool:
		return v
	case int:
		return v != 0
	case string:
		return !falseValues[strings.ToLower(v)]
	case nil:
		return false
	default:
		return true
	}
}
