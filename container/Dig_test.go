package container_test

import (
	"testing"

	"github.com/sampson-golang/utilities/container"
)

// Test structures for Dig
type digTestCase struct {
	name     string
	data     interface{}
	path     []any
	expected interface{}
	isNil    bool
}

var digTestCases = []digTestCase{
	// Basic map navigation
	{
		"simple map access",
		map[string]interface{}{"key": "value"},
		[]any{"key"},
		"value",
		false,
	},
	{
		"nested map access",
		map[string]interface{}{
			"level1": map[string]interface{}{
				"level2": "deep value",
			},
		},
		[]any{"level1", "level2"},
		"deep value",
		false,
	},
	{
		"non-existent key",
		map[string]interface{}{"key": "value"},
		[]any{"nonexistent"},
		nil,
		true,
	},

	// Array/slice navigation
	{
		"simple array access",
		[]interface{}{"first", "second", "third"},
		[]any{1},
		"second",
		false,
	},
	{
		"array access with string index",
		[]interface{}{"first", "second", "third"},
		[]any{"2"},
		"third",
		false,
	},
	{
		"out of bounds index",
		[]interface{}{"first", "second"},
		[]any{5},
		nil,
		true,
	},
	{
		"negative index",
		[]interface{}{"first", "second"},
		[]any{-1},
		nil,
		true,
	},
	{
		"invalid string index",
		[]interface{}{"first", "second"},
		[]any{"invalid"},
		nil,
		true,
	},

	// Mixed navigation
	{
		"map to array to map",
		map[string]interface{}{
			"items": []interface{}{
				map[string]interface{}{"name": "item1"},
				map[string]interface{}{"name": "item2"},
			},
		},
		[]any{"items", 1, "name"},
		"item2",
		false,
	},
	{
		"complex nested structure",
		map[string]interface{}{
			"users": []interface{}{
				map[string]interface{}{
					"profile": map[string]interface{}{
						"details": map[string]interface{}{
							"email": "user@example.com",
						},
					},
				},
			},
		},
		[]any{"users", 0, "profile", "details", "email"},
		"user@example.com",
		false,
	},

	// Edge cases
	{
		"empty path with string",
		"simple string",
		[]any{},
		"simple string",
		false,
	},
	{
		"nil data",
		nil,
		[]any{"key"},
		nil,
		true,
	},
	{
		"invalid key type for map",
		map[string]interface{}{"key": "value"},
		[]any{123},
		nil,
		true,
	},
	{
		"invalid key type for array",
		[]interface{}{"first", "second"},
		[]any{map[string]interface{}{}},
		nil,
		true,
	},
	{
		"path through non-container type",
		"simple string",
		[]any{"key"},
		nil,
		true,
	},
}

func TestDig(t *testing.T) {
	for _, tc := range digTestCases {
		t.Run(tc.name, func(t *testing.T) {
			result := container.Dig(tc.data, tc.path...)

			if tc.isNil {
				if result != nil {
					t.Errorf("Dig(%v, %v) = %v; expected nil", tc.data, tc.path, result)
				}
			} else {
				if result == nil {
					t.Errorf("Dig(%v, %v) = nil; expected %v", tc.data, tc.path, tc.expected)
					return
				}

				// Dig returns a pointer, so we need to dereference it
				actualValue := *result.(*interface{})
				if actualValue != tc.expected {
					t.Errorf("Dig(%v, %v) = %v; expected %v", tc.data, tc.path, actualValue, tc.expected)
				}
			}
		})
	}

	// Special test for empty path with map (can't compare maps directly)
	t.Run("empty path with map", func(t *testing.T) {
		data := map[string]interface{}{"key": "value"}
		result := container.Dig(data, []any{}...)

		if result == nil {
			t.Errorf("Dig with empty path should return pointer to original data, got nil")
			return
		}

		// Verify it returns a pointer to the same data
		actualValue := *result.(*interface{})
		actualMap, ok := actualValue.(map[string]interface{})
		if !ok {
			t.Errorf("Expected map[string]interface{}, got %T", actualValue)
			return
		}

		if actualMap["key"] != "value" {
			t.Errorf("Expected map to contain original data")
		}
	})
}

// Benchmark tests
func BenchmarkDig(b *testing.B) {
	data := map[string]interface{}{
		"level1": map[string]interface{}{
			"level2": []interface{}{
				map[string]interface{}{
					"target": "found",
				},
			},
		},
	}
	path := []any{"level1", "level2", 0, "target"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		container.Dig(data, path...)
	}
}

func BenchmarkDigDeepNesting(b *testing.B) {
	// Create deeply nested structure
	data := make(map[string]interface{})
	current := data
	for i := 0; i < 10; i++ {
		next := make(map[string]interface{})
		current["next"] = next
		current = next
	}
	current["target"] = "deep value"

	path := make([]any, 11)
	for i := 0; i < 10; i++ {
		path[i] = "next"
	}
	path[10] = "target"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		container.Dig(data, path...)
	}
}
