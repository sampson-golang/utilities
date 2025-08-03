package networking_test

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/sampson-golang/utilities/networking"
)

func TestQueryOrBody_GET_Request(t *testing.T) {
	// Create a GET request with query parameters
	req, _ := http.NewRequest("GET", "http://example.com?name=John&age=30&email=john@example.com", nil)

	values, err := networking.QueryOrBody(req, "name", "age", "email")

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	expectedValues := map[string]string{
		"name":  "John",
		"age":   "30",
		"email": "john@example.com",
	}

	for key, expected := range expectedValues {
		if values[key] != expected {
			t.Errorf("Expected %s=%s, got %s=%s", key, expected, key, values[key])
		}
	}
}

func TestQueryOrBody_GET_PartialKeys(t *testing.T) {
	// Create a GET request with only some of the requested keys
	req, _ := http.NewRequest("GET", "http://example.com?name=John&city=NYC", nil)

	values, err := networking.QueryOrBody(req, "name", "age", "email")

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if values["name"] != "John" {
		t.Errorf("Expected name=John, got name=%s", values["name"])
	}

	// age and email should not be present
	if values["age"] != "" {
		t.Errorf("Expected empty age, got: %s", values["age"])
	}
	if values["email"] != "" {
		t.Errorf("Expected empty email, got: %s", values["email"])
	}
}

func TestQueryOrBody_POST_JSON(t *testing.T) {
	jsonBody := `{"name":"Jane","age":"25","email":"jane@example.com"}`
	req, _ := http.NewRequest("POST", "http://example.com", bytes.NewBufferString(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	values, err := networking.QueryOrBody(req, "name", "age", "email")

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	expectedValues := map[string]string{
		"name":  "Jane",
		"age":   "25",
		"email": "jane@example.com",
	}

	for key, expected := range expectedValues {
		if values[key] != expected {
			t.Errorf("Expected %s=%s, got %s=%s", key, expected, key, values[key])
		}
	}
}

func TestQueryOrBody_POST_JSON_InvalidJSON(t *testing.T) {
	invalidJSON := `{"name":"Jane","age":"25","email":` // Invalid JSON
	req, _ := http.NewRequest("POST", "http://example.com", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	values, err := networking.QueryOrBody(req, "name", "age", "email")

	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}

	if err.Error() != "invalid JSON" {
		t.Errorf("Expected 'invalid JSON' error, got: %v", err)
	}

	// Should still return any query parameters
	if len(values) != 0 {
		t.Errorf("Expected empty values for invalid JSON, got: %v", values)
	}
}

func TestQueryOrBody_POST_Form(t *testing.T) {
	formData := url.Values{}
	formData.Set("name", "Bob")
	formData.Set("age", "35")
	formData.Set("email", "bob@example.com")

	req, _ := http.NewRequest("POST", "http://example.com", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.ParseForm() // Important: must parse form data

	values, err := networking.QueryOrBody(req, "name", "age", "email")

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	expectedValues := map[string]string{
		"name":  "Bob",
		"age":   "35",
		"email": "bob@example.com",
	}

	for key, expected := range expectedValues {
		if values[key] != expected {
			t.Errorf("Expected %s=%s, got %s=%s", key, expected, key, values[key])
		}
	}
}

func TestQueryOrBody_POST_JSON_ContentTypeWithCharset(t *testing.T) {
	jsonBody := `{"name":"Alice","age":"28"}`
	req, _ := http.NewRequest("POST", "http://example.com", bytes.NewBufferString(jsonBody))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	values, err := networking.QueryOrBody(req, "name", "age")

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if values["name"] != "Alice" {
		t.Errorf("Expected name=Alice, got name=%s", values["name"])
	}
	if values["age"] != "28" {
		t.Errorf("Expected age=28, got age=%s", values["age"])
	}
}

func TestQueryOrBody_QueryAndBody_Combined(t *testing.T) {
	// Test case where we have both query parameters and body data
	jsonBody := `{"age":"30","email":"john@example.com"}`
	req, _ := http.NewRequest("POST", "http://example.com?name=John", bytes.NewBufferString(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	values, err := networking.QueryOrBody(req, "name", "age", "email")

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Query parameter should be included
	if values["name"] != "John" {
		t.Errorf("Expected name=John from query, got name=%s", values["name"])
	}

	// Body parameters should be included
	if values["age"] != "30" {
		t.Errorf("Expected age=30 from body, got age=%s", values["age"])
	}
	if values["email"] != "john@example.com" {
		t.Errorf("Expected email=john@example.com from body, got email=%s", values["email"])
	}
}

func TestQueryOrBody_EmptyKeys(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com?name=John&age=30", nil)

	values, err := networking.QueryOrBody(req)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(values) != 0 {
		t.Errorf("Expected empty values for no keys, got: %v", values)
	}
}

func TestQueryOrBody_EmptyValues(t *testing.T) {
	// Test with keys that have empty values
	req, _ := http.NewRequest("GET", "http://example.com?name=&age=30", nil)

	values, err := networking.QueryOrBody(req, "name", "age")

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// name should not be included (empty value)
	if values["name"] != "" {
		t.Errorf("Expected empty name to be excluded, got: %s", values["name"])
	}

	// age should be included
	if values["age"] != "30" {
		t.Errorf("Expected age=30, got age=%s", values["age"])
	}
}

// Benchmarks
func BenchmarkQueryOrBody_GET_Request(b *testing.B) {
	req, _ := http.NewRequest("GET", "http://example.com?name=John&age=30&email=john@example.com&city=NYC&country=USA", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = networking.QueryOrBody(req, "name", "age", "email", "city", "country")
	}
}

func BenchmarkQueryOrBody_POST_JSON(b *testing.B) {
	jsonBody := `{"name":"Jane","age":"25","email":"jane@example.com","city":"LA","country":"USA"}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("POST", "http://example.com", bytes.NewBufferString(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		_, _ = networking.QueryOrBody(req, "name", "age", "email", "city", "country")
	}
}

func BenchmarkQueryOrBody_POST_Form(b *testing.B) {
	formData := url.Values{}
	formData.Set("name", "Bob")
	formData.Set("age", "35")
	formData.Set("email", "bob@example.com")
	formData.Set("city", "Chicago")
	formData.Set("country", "USA")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("POST", "http://example.com", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.ParseForm()
		_, _ = networking.QueryOrBody(req, "name", "age", "email", "city", "country")
	}
}

func BenchmarkQueryOrBody_Large_JSON(b *testing.B) {
	// Create a large JSON with many fields
	var jsonParts []string
	keys := make([]string, 100)
	for i := 0; i < 100; i++ {
		key := "field" + string(rune(i+48))
		keys[i] = key
		jsonParts = append(jsonParts, `"`+key+`":"value`+string(rune(i+48))+`"`)
	}
	jsonBody := `{` + strings.Join(jsonParts, ",") + `}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("POST", "http://example.com", bytes.NewBufferString(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		_, _ = networking.QueryOrBody(req, keys...)
	}
}
