package env_test

import (
	"os"
	"testing"

	"github.com/sampson-golang/utilities/env"
)

func TestLookup(t *testing.T) {
	os.Setenv("TEST_GETENV_EXISTS", "test_env_value")
	defer os.Unsetenv("TEST_GETENV_EXISTS")
	os.Setenv("TEST_GETENV_EXISTS_2", "test_env_value_2")
	defer os.Unsetenv("TEST_GETENV_EXISTS_2")
	os.Setenv("TEST_GETENV_EXISTS_3", "test_env_value_3")
	defer os.Unsetenv("TEST_GETENV_EXISTS_3")
	os.Setenv("TEST_GETENV_EMPTY", "")
	defer os.Unsetenv("TEST_GETENV_EMPTY")

	t.Run("env var exists", func(t *testing.T) {
		got, gotExists := env.Lookup("TEST_GETENV_EXISTS")
		if got != "test_env_value" {
			t.Errorf("Lookup() got = %v, want %v", got, "test_env_value")
		}
		if gotExists != true {
			t.Errorf("Lookup() gotExists = %v, want %v", gotExists, true)
		}

		t.Run("fallback string is not used", func(t *testing.T) {
			got, gotExists := env.Lookup("TEST_GETENV_EXISTS", "fallback")
			if got != "test_env_value" {
				t.Errorf("Lookup() got = %v, want %v", got, "test_env_value")
			}
			if gotExists != true {
				t.Errorf("Lookup() gotExists = %v, want %v", gotExists, true)
			}
		})

		t.Run("fallback env var is not used", func(t *testing.T) {
			got, gotExists := env.Lookup("TEST_GETENV_EXISTS", "TEST_GETENV_EXISTS_2", "TEST_GETENV_EXISTS_3")
			if got != "test_env_value" {
				t.Errorf("Lookup() got = %v, want %v", got, "test_env_value")
			}
			if gotExists != true {
				t.Errorf("Lookup() gotExists = %v, want %v", gotExists, true)
			}
		})
	})

	t.Run("env var does not exist", func(t *testing.T) {
		got, gotExists := env.Lookup("TEST_GETENV_NOT_EXISTS")
		if got != "" {
			t.Errorf("Lookup() got = %v, want %v", got, "")
		}
		if gotExists != false {
			t.Errorf("Lookup() gotExists = %v, want %v", gotExists, false)
		}

		t.Run("with fallback string", func(t *testing.T) {
			got, gotExists := env.Lookup("TEST_GETENV_NOT_EXISTS", "fallback")
			if got != "fallback" {
				t.Errorf("Lookup() got = %v, want %v", got, "fallback")
			}
			if gotExists != false {
				t.Errorf("Lookup() gotExists = %v, want %v", gotExists, true)
			}
		})

		t.Run("with fallback env", func(t *testing.T) {
			got, gotExists := env.Lookup("TEST_GETENV_NOT_EXISTS", "TEST_GETENV_EXISTS", "")

			if got != "test_env_value" {
				t.Errorf("Lookup() got = %v, want %v", got, "test_env_value")
			}
			if gotExists != true {
				t.Errorf("Lookup() gotExists = %v, want %v", gotExists, true)
			}

			t.Run("does not exist", func(t *testing.T) {
				got, gotExists := env.Lookup("TEST_GETENV_NOT_EXISTS", "TEST_GETENV_NOT_EXISTS_2", "fallback")
				if got != "fallback" {
					t.Errorf("Lookup() got = %v, want %v", got, "fallback")
				}
				if gotExists != false {
					t.Errorf("Lookup() gotExists = %v, want %v", gotExists, false)
				}
			})

			t.Run("is empty", func(t *testing.T) {
				got, gotExists := env.Lookup("TEST_GETENV_NOT_EXISTS", "TEST_GETENV_EMPTY", "fallback")
				if got != "" {
					t.Errorf("Lookup() got = %v, want %v", got, "fallback")
				}
				if gotExists != true {
					t.Errorf("Lookup() gotExists = %v, want %v", gotExists, true)
				}
			})
		})
	})
}

// Benchmark tests to ensure the functions perform well
func BenchmarkLookup(b *testing.B) {
	os.Setenv("BENCHMARK_TEST", "test_value")
	defer os.Unsetenv("BENCHMARK_TEST")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		env.Lookup("BENCHMARK_TEST")
	}
}

func BenchmarkLookupWithFallbacks(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		env.Lookup("NON_EXISTENT_VAR", "fallback1", "fallback2", "fallback3")
	}
}
