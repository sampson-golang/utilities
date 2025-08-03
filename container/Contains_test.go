package container_test

import (
	"testing"

	"github.com/sampson-golang/utilities/container"
)

type containsTestCase struct {
	name             string
	slice            []string
	nonExistingItems []string
}

var containsTestCases = []containsTestCase{
	// Basic functionality
	{"matches on strings", []string{"apple", "banana", "cherry"}, []string{"orange"}},
	{"is case sensitive", []string{"Apple", "banana"}, []string{"apple", "baNana"}},
	{"matches on special characters", []string{"!@#$%^&*()"}, []string{}},
	{"matches on unicode characters", []string{"café", "naïve", "résumé"}, []string{"cafe", "naive", "resume"}},
	{"matches on numbers", []string{"123", "abc456", "789"}, []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}},
	{
		"whitespace is not ignored",
		[]string{"hello world", " leading", "trailing ", " l and t ", ""},
		[]string{"helloworld", "leading", "trailing", "l and t", " "},
	},
}

func TestContains(t *testing.T) {
	for _, tc := range containsTestCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, item := range tc.slice {
				result := container.Contains(tc.slice, item)
				if result != true {
					t.Errorf("Contains(%v, %q) = %v; expected true", tc.slice, item, result)
				}
			}

			for _, item := range tc.nonExistingItems {
				result := container.Contains(tc.slice, item)
				if result != false {
					t.Errorf("Contains(%v, %q) = %v; expected false", tc.slice, item, result)
				}
			}
		})
	}

	t.Run("duplicate items are matched", func(t *testing.T) {
		result := container.Contains([]string{"apple", "banana", "apple"}, "apple")
		if result != true {
			t.Errorf("Contains(apple, banana, apple, apple) = %v; expected true", result)
		}

		result = container.Contains([]string{"apple", "banana", "apple"}, "banana")
		if result != true {
			t.Errorf("Contains(apple, banana, apple, banana) = %v; expected true", result)
		}

		result = container.Contains([]string{"apple", "banana", "apple"}, "cherry")
		if result != false {
			t.Errorf("Contains(apple, banana, apple, cherry) = %v; expected false", result)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		result := container.Contains([]string{}, "item")
		if result != false {
			t.Errorf("Contains(empty, \"item\") = %v; expected false", result)
		}
	})

	t.Run("empty item", func(t *testing.T) {
		result := container.Contains([]string{"apple", "banana"}, "")
		if result != false {
			t.Errorf("Contains(apple, banana, \"\") = %v; expected false", result)
		}

		result = container.Contains([]string{"", "banana"}, "apple")
		if result != false {
			t.Errorf("Contains(empty, banana, apple) = %v; expected false", result)
		}

		result = container.Contains([]string{"", "banana"}, "")
		if result != true {
			t.Errorf("Contains(empty, banana, empty) = %v; expected true", result)
		}
	})

	t.Run("nil slice", func(t *testing.T) {
		var nilSlice []string
		result := container.Contains(nilSlice, "item")
		if result != false {
			t.Errorf("Contains(nil, \"item\") = %v; expected false", result)
		}
	})
}

// Benchmark test to measure performance
func BenchmarkContains(b *testing.B) {
	slice := []string{"apple", "banana", "cherry", "date", "elderberry", "fig", "grape"}
	item := "elderberry"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		container.Contains(slice, item)
	}
}

// Benchmark with large slice
func BenchmarkContainsLargeSlice(b *testing.B) {
	slice := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		slice[i] = string(rune('a' + i%26))
	}
	item := "z"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		container.Contains(slice, item)
	}
}
