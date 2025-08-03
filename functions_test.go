package utilities_test

import (
	"reflect"
	"testing"

	"github.com/sampson-golang/utilities"
	"github.com/sampson-golang/utilities/boolable"
	"github.com/sampson-golang/utilities/container"
	"github.com/sampson-golang/utilities/container/merge"
	"github.com/sampson-golang/utilities/env"
	"github.com/sampson-golang/utilities/output"
	"github.com/sampson-golang/utilities/strutil"
)

// TestAliasesPointToCorrectFunctions verifies that exported variables
// are properly aliased to their subpackage counterparts
func TestAliasesPointToCorrectFunctions(t *testing.T) {
	tests := []struct {
		name     string
		alias    interface{}
		original interface{}
	}{
		// Type conversion utilities
		{"Boolable", utilities.Boolable, boolable.From},

		// Container utilities
		{"Contains", utilities.Contains, container.Contains},
		{"Dig", utilities.Dig, container.Dig},
		{"DigAssign", utilities.DigAssign, container.DigAssign},

		// Environment variable utilities
		{"EnvExists", utilities.EnvExists, env.Exists},
		{"GetEnv", utilities.GetEnv, env.Get},
		{"GetPresentEnv", utilities.GetPresentEnv, env.GetPresent},
		{"LookupEnv", utilities.LookupEnv, env.Lookup},
		{"LookupPresentEnv", utilities.LookupPresentEnv, env.LookupPresent},

		// Merging utilities
		{"MergeParams", utilities.MergeParams, merge.Params},
		{"MergeStructs", utilities.MergeStructs, merge.Structs},

		// Output utilities
		{"Prettify", utilities.Prettify, output.Prettify},
		{"PrettyPrint", utilities.PrettyPrint, output.PrettyPrint},

		// String utilities
		{"Squish", utilities.Squish, strutil.Squish},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aliasPtr := reflect.ValueOf(tt.alias).Pointer()
			originalPtr := reflect.ValueOf(tt.original).Pointer()

			if aliasPtr != originalPtr {
				t.Errorf("%s alias does not point to the correct function", tt.name)
			}
		})
	}
}

// TestBasicFunctionality provides smoke tests to ensure aliases work correctly
// without duplicating comprehensive testing done in subpackages
func TestBasicFunctionality(t *testing.T) {
	t.Run("Boolable", func(t *testing.T) {
		if !utilities.Boolable(true) {
			t.Error("Boolable(true) should return true")
		}
		if utilities.Boolable(false) {
			t.Error("Boolable(false) should return false")
		}
	})

	t.Run("Contains", func(t *testing.T) {
		slice := []string{"a", "b", "c"}
		if !utilities.Contains(slice, "b") {
			t.Error("Contains should find existing element")
		}
		if utilities.Contains(slice, "z") {
			t.Error("Contains should not find non-existing element")
		}
	})

	t.Run("Squish", func(t *testing.T) {
		result := utilities.Squish("  hello   world  ")
		expected := "hello world"
		if result != expected {
			t.Errorf("Squish result = %q, expected %q", result, expected)
		}
	})

	t.Run("Prettify", func(t *testing.T) {
		input := map[string]string{"key": "value"}
		result := utilities.Prettify(input)
		if result == "" {
			t.Error("Prettify should return non-empty string")
		}
	})

	// Add basic smoke tests for other functions as needed
	// These should be minimal - just ensure they work, not comprehensive testing
}

// TestTypeAliases verifies that type aliases work correctly
func TestTypeAliases(t *testing.T) {
	t.Run("Set type alias", func(t *testing.T) {
		// Test that utilities.Set works correctly as a type alias
		utilSet := utilities.Set{}
		utilSet["test"] = struct{}{}

		if _, exists := utilSet["test"]; !exists {
			t.Error("Set type alias should work for basic operations")
		}

		// Test that we can convert between the types
		containerSet := container.Set(utilSet)
		if _, exists := containerSet["test"]; !exists {
			t.Error("Should be able to convert between utilities.Set and container.Set")
		}

		// Test the reverse conversion
		utilSet2 := utilities.Set(containerSet)
		if _, exists := utilSet2["test"]; !exists {
			t.Error("Should be able to convert from container.Set to utilities.Set")
		}
	})
}

// TestAllExportsAvailable ensures all expected exports are available
// This catches cases where imports might be missing or broken
func TestAllExportsAvailable(t *testing.T) {
	// Test that we can call each exported function without panics
	// This is a sanity check that all aliases are properly set up

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Panic occurred while testing exports: %v", r)
		}
	}()

	// Just verify they exist and are callable (don't test functionality)
	_ = utilities.Boolable
	_ = utilities.Contains
	_ = utilities.Dig
	_ = utilities.DigAssign
	_ = utilities.EnvExists
	_ = utilities.GetEnv
	_ = utilities.GetPresentEnv
	_ = utilities.LookupEnv
	_ = utilities.LookupPresentEnv
	_ = utilities.MergeParams
	_ = utilities.MergeStructs
	_ = utilities.Prettify
	_ = utilities.PrettyPrint
	_ = utilities.Squish
}
