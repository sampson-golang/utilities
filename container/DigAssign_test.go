package container_test

import (
	"testing"

	"github.com/sampson-golang/utilities/container"
)

// Test structures for DigAssign
type digAssignTestStruct struct {
	StringField string
	IntField    int
	BoolField   bool
	FloatField  float64
}

type DigAssignTestCase struct {
	name        string
	data        interface{}
	path        []any
	key         string
	expectField interface{}
	shouldSet   bool
}

var digAssignTestCases = []DigAssignTestCase{
	// Successful assignments
	{
		"assign string",
		map[string]interface{}{"source": "test value"},
		[]any{"source"},
		"StringField",
		"test value",
		true,
	},
	{
		"assign int",
		map[string]interface{}{"number": 42},
		[]any{"number"},
		"IntField",
		42,
		true,
	},
	{
		"assign bool",
		map[string]interface{}{"flag": true},
		[]any{"flag"},
		"BoolField",
		true,
		true,
	},
	{
		"assign with type conversion int to float",
		map[string]interface{}{"number": 42},
		[]any{"number"},
		"FloatField",
		float64(42),
		true,
	},
	{
		"nested path assignment",
		map[string]interface{}{
			"nested": map[string]interface{}{
				"value": "deep assignment",
			},
		},
		[]any{"nested", "value"},
		"StringField",
		"deep assignment",
		true,
	},

	// Failed assignments
	{
		"non-existent field",
		map[string]interface{}{"source": "test"},
		[]any{"source"},
		"NonExistentField",
		"",
		false,
	},
	{
		"non-existent path",
		map[string]interface{}{"other": "test"},
		[]any{"nonexistent"},
		"StringField",
		"",
		false,
	},
	{
		"incompatible type assignment",
		map[string]interface{}{"complex": map[string]interface{}{"inner": "value"}},
		[]any{"complex"},
		"StringField",
		"",
		false,
	},
}

func TestDigAssign(t *testing.T) {
	for _, tc := range digAssignTestCases {
		t.Run(tc.name, func(t *testing.T) {
			result := &digAssignTestStruct{}
			container.DigAssign(result, tc.key, tc.data, tc.path...)

			switch tc.key {
			case "StringField":
				if tc.shouldSet {
					if result.StringField != tc.expectField {
						t.Errorf("DigAssign failed: StringField = %v; expected %v", result.StringField, tc.expectField)
					}
				} else {
					if result.StringField != "" {
						t.Errorf("DigAssign should not have set StringField, but got %v", result.StringField)
					}
				}
			case "IntField":
				if tc.shouldSet {
					if result.IntField != tc.expectField {
						t.Errorf("DigAssign failed: IntField = %v; expected %v", result.IntField, tc.expectField)
					}
				} else {
					if result.IntField != 0 {
						t.Errorf("DigAssign should not have set IntField, but got %v", result.IntField)
					}
				}
			case "BoolField":
				if tc.shouldSet {
					if result.BoolField != tc.expectField {
						t.Errorf("DigAssign failed: BoolField = %v; expected %v", result.BoolField, tc.expectField)
					}
				} else {
					if result.BoolField != false {
						t.Errorf("DigAssign should not have set BoolField, but got %v", result.BoolField)
					}
				}
			case "FloatField":
				if tc.shouldSet {
					if result.FloatField != tc.expectField {
						t.Errorf("DigAssign failed: FloatField = %v; expected %v", result.FloatField, tc.expectField)
					}
				} else {
					if result.FloatField != 0.0 {
						t.Errorf("DigAssign should not have set FloatField, but got %v", result.FloatField)
					}
				}
			}
		})
	}
}

// Test DigAssign with nil result - should panic
func TestDigAssignNilResult(t *testing.T) {
	data := map[string]interface{}{"key": "value"}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("DigAssign with nil result should panic, but it didn't")
		}
	}()

	container.DigAssign(nil, "StringField", data, "key")
}

// Test DigAssign with non-struct result - should panic
func TestDigAssignNonStruct(t *testing.T) {
	data := map[string]interface{}{"key": "value"}
	var result string

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("DigAssign with non-struct should panic, but it didn't")
		}
	}()

	container.DigAssign(&result, "StringField", data, "key")
}

func BenchmarkDigAssign(b *testing.B) {
	data := map[string]interface{}{
		"nested": map[string]interface{}{
			"value": "test string",
		},
	}
	path := []any{"nested", "value"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := &digAssignTestStruct{}
		container.DigAssign(result, "StringField", data, path...)
	}
}
