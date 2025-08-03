package httputil_test

import (
	"strings"
	"testing"

	"github.com/sampson-golang/utilities/httputil"
)

type TestStruct struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

func TestUnmarshalResponse_ValidJSON(t *testing.T) {
	jsonData := `{"name":"John Doe","age":30,"email":"john@example.com"}`
	var result TestStruct

	err := httputil.UnmarshalResponse([]byte(jsonData), &result)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if result.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got: %s", result.Name)
	}
	if result.Age != 30 {
		t.Errorf("Expected age 30, got: %d", result.Age)
	}
	if result.Email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got: %s", result.Email)
	}
}

func TestUnmarshalResponse_InvalidJSON(t *testing.T) {
	invalidJSON := `{"name":"John Doe","age":30,"email":` // Missing closing quote and bracket
	var result TestStruct

	err := httputil.UnmarshalResponse([]byte(invalidJSON), &result)

	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}

	if !strings.Contains(err.Error(), "failed to parse response body") {
		t.Errorf("Expected error message to contain 'failed to parse response body', got: %v", err)
	}
}

func TestUnmarshalResponse_EmptyJSON(t *testing.T) {
	var result TestStruct

	err := httputil.UnmarshalResponse([]byte("{}"), &result)

	if err != nil {
		t.Errorf("Expected no error for empty JSON, got: %v", err)
	}

	// Should have zero values
	if result.Name != "" || result.Age != 0 || result.Email != "" {
		t.Error("Expected zero values for empty JSON")
	}
}

func TestUnmarshalResponse_NilBytes(t *testing.T) {
	var result TestStruct

	err := httputil.UnmarshalResponse(nil, &result)

	// Nil bytes should actually cause an error since it's invalid JSON
	if err == nil {
		t.Error("Expected error for nil bytes (invalid JSON), got nil")
	}

	if !strings.Contains(err.Error(), "failed to parse response body") {
		t.Errorf("Expected error message to contain 'failed to parse response body', got: %v", err)
	}
}

func TestUnmarshalResponse_Array(t *testing.T) {
	jsonData := `[{"name":"John"},{"name":"Jane"}]`
	var result []TestStruct

	err := httputil.UnmarshalResponse([]byte(jsonData), &result)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 items, got: %d", len(result))
	}

	if result[0].Name != "John" || result[1].Name != "Jane" {
		t.Error("Array unmarshaling failed")
	}
}

func TestUnmarshalResponse_TypeMismatch(t *testing.T) {
	jsonData := `{"name":"John Doe","age":"thirty","email":"john@example.com"}` // age as string instead of int
	var result TestStruct

	err := httputil.UnmarshalResponse([]byte(jsonData), &result)

	if err == nil {
		t.Error("Expected error for type mismatch, got nil")
	}
}

func TestUnmarshalResponse_Map(t *testing.T) {
	jsonData := `{"key1":"value1","key2":"value2"}`
	var result map[string]string

	err := httputil.UnmarshalResponse([]byte(jsonData), &result)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if result["key1"] != "value1" || result["key2"] != "value2" {
		t.Error("Map unmarshaling failed")
	}
}

// Benchmarks
func BenchmarkUnmarshalResponse_SmallStruct(b *testing.B) {
	jsonData := []byte(`{"name":"John Doe","age":30,"email":"john@example.com"}`)
	var result TestStruct

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = httputil.UnmarshalResponse(jsonData, &result)
	}
}

func BenchmarkUnmarshalResponse_LargeArray(b *testing.B) {
	// Create a large JSON array
	var items []string
	for i := 0; i < 1000; i++ {
		items = append(items, `{"name":"User `+string(rune(i))+`","age":25,"email":"user@example.com"}`)
	}
	jsonData := []byte(`[` + strings.Join(items, ",") + `]`)
	var result []TestStruct

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = httputil.UnmarshalResponse(jsonData, &result)
	}
}

func BenchmarkUnmarshalResponse_Map(b *testing.B) {
	jsonData := []byte(`{"key1":"value1","key2":"value2","key3":"value3","key4":"value4","key5":"value5"}`)
	var result map[string]string

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = httputil.UnmarshalResponse(jsonData, &result)
	}
}

func BenchmarkUnmarshalResponse_ErrorCase(b *testing.B) {
	invalidJSON := []byte(`{"name":"John Doe","age":30,"email":`) // Invalid JSON
	var result TestStruct

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = httputil.UnmarshalResponse(invalidJSON, &result)
	}
}
