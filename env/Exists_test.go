package env_test

import (
	"os"
	"testing"

	"github.com/sampson-golang/utilities/env"
)

var Exists = env.Exists

func TestExists(t *testing.T) {
	t.Run("env var exists", func(t *testing.T) {
		os.Setenv("TEST_EXISTS_SET", "some_value")
		defer os.Unsetenv("TEST_EXISTS_SET")

		if got := Exists("TEST_EXISTS_SET"); got != true {
			t.Errorf("Exists() = %v, want %v", got, true)
		}

		t.Run("is empty", func(t *testing.T) {
			os.Setenv("TEST_EXISTS_EMPTY", "")
			defer os.Unsetenv("TEST_EXISTS_EMPTY")

			if got := Exists("TEST_EXISTS_EMPTY"); got != true {
				t.Errorf("Exists() = %v, want %v", got, true)
			}
		})
	})

	t.Run("env var does not exist", func(t *testing.T) {
		if got := Exists("TEST_EXISTS_NOT_SET"); got != false {
			t.Errorf("Exists() = %v, want %v", got, false)
		}
	})
}

func BenchmarkExists(b *testing.B) {
	os.Setenv("BENCHMARK_TEST", "test_value")
	defer os.Unsetenv("BENCHMARK_TEST")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Exists("BENCHMARK_TEST")
		Exists("BENCHMARK_TEST_EMPTY")
	}
}
