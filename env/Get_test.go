package env_test

import (
	"os"
	"testing"

	"github.com/sampson-golang/utilities/env"
)

func TestGet(t *testing.T) {
	os.Setenv("TEST_GETENV_EXISTS", "test_env_value")
	defer os.Unsetenv("TEST_GETENV_EXISTS")
	os.Setenv("TEST_GETENV_EXISTS_2", "test_env_value_2")
	defer os.Unsetenv("TEST_GETENV_EXISTS_2")
	os.Setenv("TEST_GETENV_EXISTS_3", "test_env_value_3")
	defer os.Unsetenv("TEST_GETENV_EXISTS_3")
	os.Setenv("TEST_GETENV_EMPTY", "")
	defer os.Unsetenv("TEST_GETENV_EMPTY")

	t.Run("env var exists", func(t *testing.T) {
		got := env.Get("TEST_GETENV_EXISTS")

		if got != "test_env_value" {
			t.Errorf("Get() got = %v, want %v", got, "test_env_value")
		}

		t.Run("fallback string is not used", func(t *testing.T) {
			got := env.Get("TEST_GETENV_EXISTS", "fallback")

			if got != "test_env_value" {
				t.Errorf("Get() got = %v, want %v", got, "test_env_value")
			}
		})

		t.Run("fallback env var is not used", func(t *testing.T) {
			got := env.Get("TEST_GETENV_EXISTS", "TEST_GETENV_EXISTS_2", "TEST_GETENV_EXISTS_3")

			if got != "test_env_value" {
				t.Errorf("Get() got = %v, want %v", got, "test_env_value")
			}
		})
	})

	t.Run("env var does not exist", func(t *testing.T) {
		got := env.Get("TEST_GETENV_NOT_EXISTS")

		if got != "" {
			t.Errorf("Get() got = %v, want %v", got, "")
		}

		t.Run("with fallback string", func(t *testing.T) {
			got := env.Get("TEST_GETENV_NOT_EXISTS", "fallback")

			if got != "fallback" {
				t.Errorf("Get() got = %v, want %v", got, "fallback")
			}
		})

		t.Run("with fallback env", func(t *testing.T) {
			got := env.Get("TEST_GETENV_NOT_EXISTS", "TEST_GETENV_EXISTS", "")

			if got != "test_env_value" {
				t.Errorf("Get() got = %v, want %v", got, "test_env_value")
			}

			t.Run("does not exist", func(t *testing.T) {
				got := env.Get("TEST_GETENV_NOT_EXISTS", "TEST_GETENV_NOT_EXISTS_2", "fallback")

				if got != "fallback" {
					t.Errorf("Get() got = %v, want %v", got, "fallback")
				}
			})

			t.Run("is empty", func(t *testing.T) {
				got := env.Get("TEST_GETENV_NOT_EXISTS", "TEST_GETENV_EMPTY", "fallback")

				if got != "" {
					t.Errorf("Get() got = %v, want %v", got, "fallback")
				}
			})
		})
	})
}

// Benchmark tests to ensure the functions perform well
func BenchmarkGet(b *testing.B) {
	os.Setenv("BENCHMARK_TEST", "test_value")
	defer os.Unsetenv("BENCHMARK_TEST")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		env.Get("BENCHMARK_TEST")
	}
}

func BenchmarkGetWithFallbacks(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		env.Get("NON_EXISTENT_VAR", "fallback1", "fallback2", "fallback3")
	}
}
