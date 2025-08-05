package output_test

import (
	"testing"

	"github.com/sampson-golang/utilities/output"
)

func TestPrettify_SimpleMap(t *testing.T) {
	input := map[string]interface{}{
		"name": "John",
		"age":  30,
	}

	expected := "{\n  \"age\": 30,\n  \"name\": \"John\"\n}"

	raw := output.PrettifyBytes(input)
	result := output.Prettify(input)

	if string(raw) != result {
		t.Errorf("Bytes Don't Match: %v | %v", string(raw), result)
	}

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	t.Run("with custom indent", func(t *testing.T) {
		raw := output.PrettifyBytes(input, "\t")
		result := output.Prettify(input, "\t")

		if string(raw) != result {
			t.Errorf("Bytes Don't Match: %v | %v", string(raw), result)
		}

		expected = "{\n\t\"age\": 30,\n\t\"name\": \"John\"\n}"
		if result != expected {
			t.Errorf("Expected %s, got %s", expected, result)
		}
	})
}

func TestPrettify_NestedStructure(t *testing.T) {
	input := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "Jane",
			"details": map[string]interface{}{
				"email":  "jane@example.com",
				"active": true,
			},
		},
		"count": 42,
	}

	expected := "{\n  \"count\": 42,\n  \"user\": {\n    \"details\": {\n      \"active\": true,\n      \"email\": \"jane@example.com\"\n    },\n    \"name\": \"Jane\"\n  }\n}"

	raw := output.PrettifyBytes(input)
	result := output.Prettify(input)

	if string(raw) != result {
		t.Errorf("Bytes Don't Match: %v | %v", string(raw), result)
	}

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	t.Run("with custom indent", func(t *testing.T) {
		raw := output.PrettifyBytes(input, "\t")
		result := output.Prettify(input, "\t")

		if string(raw) != result {
			t.Errorf("Bytes Don't Match: %v | %v", string(raw), result)
		}

		expected = "{\n\t\"count\": 42,\n\t\"user\": {\n\t\t\"details\": {\n\t\t\t\"active\": true,\n\t\t\t\"email\": \"jane@example.com\"\n\t\t},\n\t\t\"name\": \"Jane\"\n\t}\n}"
		if result != expected {
			t.Errorf("Expected %s, got %s", expected, result)
		}
	})
}

func TestPrettify_Array(t *testing.T) {
	input := []interface{}{"apple", "banana", "cherry"}
	expected := "[\n  \"apple\",\n  \"banana\",\n  \"cherry\"\n]"

	raw := output.PrettifyBytes(input)
	result := output.Prettify(input)

	if string(raw) != result {
		t.Errorf("Bytes Don't Match: %v | %v", string(raw), result)
	}

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	t.Run("with custom indent", func(t *testing.T) {
		raw := output.PrettifyBytes(input, "	")
		result := output.Prettify(input, "	")

		if string(raw) != result {
			t.Errorf("Bytes Don't Match: %v | %v", string(raw), result)
		}

		expected = "[\n\t\"apple\",\n\t\"banana\",\n\t\"cherry\"\n]"
		if result != expected {
			t.Errorf("Expected %s, got %s", expected, result)
		}
	})
}

func TestPrettify_SimpleTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"string", "hello", `"hello"`},
		{"number", 42, "42"},
		{"boolean", true, "true"},
		{"null", nil, "null"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			raw := output.PrettifyBytes(test.input)
			result := output.Prettify(test.input)

			if string(raw) != result {
				t.Errorf("Bytes Don't Match: %v | %v", string(raw), result)
			}

			if result != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, result)
			}
		})
	}
}

func TestPrettify_EmptyStructures(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"empty map", map[string]interface{}{}, "{}"},
		{"empty array", []interface{}{}, "[]"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			raw := output.PrettifyBytes(test.input)
			result := output.Prettify(test.input)

			if string(raw) != result {
				t.Errorf("Bytes Don't Match: %v | %v", string(raw), result)
			}

			if result != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, result)
			}
		})
	}

	t.Run("with custom indent is still empty", func(t *testing.T) {
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				raw := output.PrettifyBytes(test.input, "\t")
				result := output.Prettify(test.input, "\t")

				if string(raw) != result {
					t.Errorf("Bytes Don't Match: %v | %v", string(raw), result)
				}

				if result != test.expected {
					t.Errorf("Expected %s, got %s", test.expected, result)
				}
			})
		}
	})
}
