package env_test

import (
	"os"
	"testing"

	"github.com/sampson-golang/utilities/env"
)

func TestLookupPresent(t *testing.T) {
	os.Setenv("TEST_GETPRESENT_EXISTS", "test_env_value")
	defer os.Unsetenv("TEST_GETPRESENT_EXISTS")
	os.Setenv("TEST_GETPRESENT_EXISTS_2", "test_env_value_2")
	defer os.Unsetenv("TEST_GETPRESENT_EXISTS_2")
	os.Setenv("TEST_GETPRESENT_EXISTS_3", "test_env_value_3")
	defer os.Unsetenv("TEST_GETPRESENT_EXISTS_3")
	os.Setenv("TEST_GETPRESENT_EMPTY", "")
	defer os.Unsetenv("TEST_GETPRESENT_EMPTY")

	t.Run("env var exists", func(t *testing.T) {
		got, gotExists := env.LookupPresent("TEST_GETPRESENT_EXISTS")
		if got != "test_env_value" {
			t.Errorf("LookupPresent() got = %v, want %v", got, "test_env_value")
		}
		if gotExists != true {
			t.Errorf("LookupPresent() gotExists = %v, want %v", gotExists, true)
		}

		t.Run("fallback string is not used", func(t *testing.T) {
			got, gotExists := env.LookupPresent("TEST_GETPRESENT_EXISTS", "fallback")
			if got != "test_env_value" {
				t.Errorf("LookupPresent() got = %v, want %v", got, "test_env_value")
			}
			if gotExists != true {
				t.Errorf("LookupPresent() gotExists = %v, want %v", gotExists, true)
			}
		})

		t.Run("fallback env var is not used", func(t *testing.T) {
			got, gotExists := env.LookupPresent("TEST_GETPRESENT_EXISTS", "TEST_GETPRESENT_EXISTS_2", "TEST_GETPRESENT_EXISTS_3")
			if got != "test_env_value" {
				t.Errorf("LookupPresent() got = %v, want %v", got, "test_env_value")
			}
			if gotExists != true {
				t.Errorf("LookupPresent() gotExists = %v, want %v", gotExists, true)
			}
		})
	})

	t.Run("env var does not exist", func(t *testing.T) {
		got, gotExists := env.LookupPresent("TEST_GETPRESENT_NOT_EXISTS")
		if got != "" {
			t.Errorf("LookupPresent() got = %v, want %v", got, "")
		}
		if gotExists != false {
			t.Errorf("LookupPresent() gotExists = %v, want %v", gotExists, false)
		}

		t.Run("with fallback string", func(t *testing.T) {
			got, gotExists := env.LookupPresent("TEST_GETPRESENT_NOT_EXISTS", "fallback")
			if got != "fallback" {
				t.Errorf("LookupPresent() got = %v, want %v", got, "fallback")
			}
			if gotExists != false {
				t.Errorf("LookupPresent() gotExists = %v, want %v", gotExists, true)
			}
		})

		t.Run("with fallback env", func(t *testing.T) {
			got, gotExists := env.LookupPresent("TEST_GETPRESENT_NOT_EXISTS", "TEST_GETPRESENT_EXISTS", "")

			if got != "test_env_value" {
				t.Errorf("LookupPresent() got = %v, want %v", got, "test_env_value")
			}
			if gotExists != true {
				t.Errorf("LookupPresent() gotExists = %v, want %v", gotExists, true)
			}

			t.Run("does not exist", func(t *testing.T) {
				got, gotExists := env.LookupPresent("TEST_GETPRESENT_NOT_EXISTS", "TEST_GETPRESENT_NOT_EXISTS_2", "fallback")
				if got != "fallback" {
					t.Errorf("LookupPresent() got = %v, want %v", got, "fallback")
				}
				if gotExists != false {
					t.Errorf("LookupPresent() gotExists = %v, want %v", gotExists, false)
				}
			})

			t.Run("is empty", func(t *testing.T) {
				got, gotExists := env.LookupPresent("TEST_GETPRESENT_NOT_EXISTS", "TEST_GETPRESENT_EMPTY", "fallback")
				if got != "fallback" {
					t.Errorf("LookupPresent() got = %v, want %v", got, "fallback")
				}
				if gotExists != false {
					t.Errorf("LookupPresent() gotExists = %v, want %v", gotExists, false)
				}
			})
		})
	})
}

func BenchmarkLookupPresent(b *testing.B) {
	os.Setenv("BENCHMARK_TEST", "test_value")
	defer os.Unsetenv("BENCHMARK_TEST")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		env.LookupPresent("BENCHMARK_TEST")
	}
}

func BenchmarkLookupPresentWithFallbacks(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		env.LookupPresent("NON_EXISTENT_VAR", "fallback1", "fallback2", "fallback3")
	}
}
