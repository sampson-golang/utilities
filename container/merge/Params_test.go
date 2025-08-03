package merge_test

import (
	"testing"

	"github.com/sampson-golang/utilities/container/merge"
)

func TestParams_EmptyDestination(t *testing.T) {
	dest := make(map[string]string)
	src := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	merge.Params(dest, src)

	expected := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	if len(dest) != len(expected) {
		t.Errorf("Expected dest to have %d keys, got %d", len(expected), len(dest))
	}

	for key, expectedValue := range expected {
		if value, exists := dest[key]; !exists {
			t.Errorf("Expected key '%s' to exist in dest", key)
		} else if value != expectedValue {
			t.Errorf("Expected dest['%s'] to be '%s', got '%s'", key, expectedValue, value)
		}
	}
}

func TestParams_ExistingDestination(t *testing.T) {
	dest := map[string]string{
		"existing1": "original1",
		"existing2": "original2",
	}
	src := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	merge.Params(dest, src)

	expected := map[string]string{
		"existing1": "original1",
		"existing2": "original2",
		"key1":      "value1",
		"key2":      "value2",
	}

	if len(dest) != len(expected) {
		t.Errorf("Expected dest to have %d keys, got %d", len(expected), len(dest))
	}

	for key, expectedValue := range expected {
		if value, exists := dest[key]; !exists {
			t.Errorf("Expected key '%s' to exist in dest", key)
		} else if value != expectedValue {
			t.Errorf("Expected dest['%s'] to be '%s', got '%s'", key, expectedValue, value)
		}
	}
}

func TestParams_OverwriteExistingKeys(t *testing.T) {
	dest := map[string]string{
		"key1": "original1",
		"key2": "original2",
	}
	src := map[string]string{
		"key1": "new1",
		"key3": "value3",
	}

	merge.Params(dest, src)

	expected := map[string]string{
		"key1": "new1",      // overwritten
		"key2": "original2", // unchanged
		"key3": "value3",    // added
	}

	if len(dest) != len(expected) {
		t.Errorf("Expected dest to have %d keys, got %d", len(expected), len(dest))
	}

	for key, expectedValue := range expected {
		if value, exists := dest[key]; !exists {
			t.Errorf("Expected key '%s' to exist in dest", key)
		} else if value != expectedValue {
			t.Errorf("Expected dest['%s'] to be '%s', got '%s'", key, expectedValue, value)
		}
	}
}

func TestParams_MultipleSources(t *testing.T) {
	dest := make(map[string]string)
	src1 := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	src2 := map[string]string{
		"key3": "value3",
		"key4": "value4",
	}
	src3 := map[string]string{
		"key1": "overwritten1", // overwrites key1 from src1
		"key5": "value5",
	}

	merge.Params(dest, src1, src2, src3)

	expected := map[string]string{
		"key1": "overwritten1", // final value from src3
		"key2": "value2",       // from src1
		"key3": "value3",       // from src2
		"key4": "value4",       // from src2
		"key5": "value5",       // from src3
	}

	if len(dest) != len(expected) {
		t.Errorf("Expected dest to have %d keys, got %d", len(expected), len(dest))
	}

	for key, expectedValue := range expected {
		if value, exists := dest[key]; !exists {
			t.Errorf("Expected key '%s' to exist in dest", key)
		} else if value != expectedValue {
			t.Errorf("Expected dest['%s'] to be '%s', got '%s'", key, expectedValue, value)
		}
	}
}

func TestParams_EmptySources(t *testing.T) {
	dest := map[string]string{
		"existing": "value",
	}
	src := map[string]string{} // empty source

	merge.Params(dest, src)

	// dest should remain unchanged
	if len(dest) != 1 {
		t.Errorf("Expected dest to have 1 key, got %d", len(dest))
	}
	if dest["existing"] != "value" {
		t.Errorf("Expected dest['existing'] to be 'value', got '%s'", dest["existing"])
	}
}

func TestParams_NilSourceHandling(t *testing.T) {
	dest := map[string]string{
		"existing": "value",
	}

	// Test with no sources
	merge.Params(dest)

	// dest should remain unchanged
	if len(dest) != 1 {
		t.Errorf("Expected dest to have 1 key, got %d", len(dest))
	}
	if dest["existing"] != "value" {
		t.Errorf("Expected dest['existing'] to be 'value', got '%s'", dest["existing"])
	}
}

func TestParams_EmptyStringValues(t *testing.T) {
	dest := make(map[string]string)
	src := map[string]string{
		"key1": "",
		"key2": "value2",
		"key3": "",
	}

	merge.Params(dest, src)

	expected := map[string]string{
		"key1": "",
		"key2": "value2",
		"key3": "",
	}

	if len(dest) != len(expected) {
		t.Errorf("Expected dest to have %d keys, got %d", len(expected), len(dest))
	}

	for key, expectedValue := range expected {
		if value, exists := dest[key]; !exists {
			t.Errorf("Expected key '%s' to exist in dest", key)
		} else if value != expectedValue {
			t.Errorf("Expected dest['%s'] to be '%s', got '%s'", key, expectedValue, value)
		}
	}
}
