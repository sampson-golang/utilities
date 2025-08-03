package httputil_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sampson-golang/utilities/httputil"
)

func TestNewContextToken(t *testing.T) {
	name := "test-token"
	token := httputil.NewContextToken(name)

	if token == nil {
		t.Error("Expected non-nil token")
	}

	if token.String() != name {
		t.Errorf("Expected token name %s, got %s", name, token.String())
	}
}

func TestNewAppContext(t *testing.T) {
	name := "test-context"
	appCtx := httputil.NewAppContext(name)

	if appCtx == nil {
		t.Error("Expected non-nil AppContext")
	}

	// Test that the internal token was created correctly
	req, _ := http.NewRequest("GET", "/", nil)
	testValue := "test-value"

	updatedReq := appCtx.SetContext(req, testValue)
	retrievedValue := appCtx.GetContext(updatedReq)

	if retrievedValue != testValue {
		t.Errorf("Expected value %v, got %v", testValue, retrievedValue)
	}
}

func TestAppContext_SetContext(t *testing.T) {
	appCtx := httputil.NewAppContext("test")
	req, _ := http.NewRequest("GET", "/", nil)

	testValue := "test-value"
	updatedReq := appCtx.SetContext(req, testValue)

	if updatedReq == req {
		t.Error("Expected new request object, got same reference")
	}

	retrievedValue := appCtx.GetContext(updatedReq)
	if retrievedValue != testValue {
		t.Errorf("Expected value %v, got %v", testValue, retrievedValue)
	}
}

func TestAppContext_GetContext(t *testing.T) {
	appCtx := httputil.NewAppContext("test")
	req, _ := http.NewRequest("GET", "/", nil)

	// Test getting from request without context value
	value := appCtx.GetContext(req)
	if value != nil {
		t.Errorf("Expected nil value for empty context, got %v", value)
	}

	// Test getting after setting value
	testValue := "test-value"
	updatedReq := appCtx.SetContext(req, testValue)
	retrievedValue := appCtx.GetContext(updatedReq)

	if retrievedValue != testValue {
		t.Errorf("Expected value %v, got %v", testValue, retrievedValue)
	}
}

func TestAppContext_WithContext(t *testing.T) {
	appCtx := httputil.NewAppContext("test")
	req, _ := http.NewRequest("GET", "/", nil)

	testValue := "test-value"
	updatedReq := appCtx.WithContext(req, testValue)

	if updatedReq == req {
		t.Error("Expected new request object, got same reference")
	}

	retrievedValue := appCtx.GetContext(updatedReq)
	if retrievedValue != testValue {
		t.Errorf("Expected value %v, got %v", testValue, retrievedValue)
	}
}

func TestAppContext_WithContext_ExistingValue(t *testing.T) {
	appCtx := httputil.NewAppContext("test")
	req, _ := http.NewRequest("GET", "/", nil)

	// Set initial value
	initialValue := "initial-value"
	reqWithValue := appCtx.SetContext(req, initialValue)

	// Try to set another value using WithContext (should not override)
	newValue := "new-value"
	finalReq := appCtx.WithContext(reqWithValue, newValue)

	// Should return the same request since context already has a value
	if finalReq != reqWithValue {
		t.Error("Expected same request object when context already has value")
	}

	// Value should still be the initial value
	retrievedValue := appCtx.GetContext(finalReq)
	if retrievedValue != initialValue {
		t.Errorf("Expected initial value %v, got %v", initialValue, retrievedValue)
	}
}

func TestAppContext_DifferentTypes(t *testing.T) {
	appCtx := httputil.NewAppContext("test")
	req, _ := http.NewRequest("GET", "/", nil)

	// Test string
	stringVal := "test-string"
	req = appCtx.SetContext(req, stringVal)
	if appCtx.GetContext(req) != stringVal {
		t.Error("String value test failed")
	}

	// Test int
	req, _ = http.NewRequest("GET", "/", nil)
	intVal := 42
	req = appCtx.SetContext(req, intVal)
	if appCtx.GetContext(req) != intVal {
		t.Error("Int value test failed")
	}

	// Test struct
	type TestStruct struct {
		Name string
		ID   int
	}
	req, _ = http.NewRequest("GET", "/", nil)
	structVal := TestStruct{Name: "test", ID: 123}
	req = appCtx.SetContext(req, structVal)
	if appCtx.GetContext(req) != structVal {
		t.Error("Struct value test failed")
	}

	// Test map
	req, _ = http.NewRequest("GET", "/", nil)
	mapVal := map[string]string{"key": "value"}
	req = appCtx.SetContext(req, mapVal)
	retrieved := appCtx.GetContext(req).(map[string]string)
	if retrieved["key"] != "value" {
		t.Error("Map value test failed")
	}
}

func TestAppContext_MultipleContexts(t *testing.T) {
	// Test that different AppContext instances don't interfere with each other
	ctx1 := httputil.NewAppContext("context1")
	ctx2 := httputil.NewAppContext("context2")

	req, _ := http.NewRequest("GET", "/", nil)

	value1 := "value-for-context1"
	value2 := "value-for-context2"

	// Set values in both contexts
	req = ctx1.SetContext(req, value1)
	req = ctx2.SetContext(req, value2)

	// Both values should be retrievable
	if ctx1.GetContext(req) != value1 {
		t.Errorf("Expected value1 %v, got %v", value1, ctx1.GetContext(req))
	}

	if ctx2.GetContext(req) != value2 {
		t.Errorf("Expected value2 %v, got %v", value2, ctx2.GetContext(req))
	}
}

func TestContextMiddleware(t *testing.T) {
	appCtx := httputil.NewAppContext("test-middleware")
	testValue := "middleware-value"

	// Create a simple handler that checks the context
	var receivedValue interface{}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedValue = appCtx.GetContext(r)
		w.WriteHeader(http.StatusOK)
	})

	// Wrap with middleware
	middleware := httputil.ContextMiddleware(appCtx, testValue)
	wrappedHandler := middleware(handler)

	// Create test request and response
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Execute the wrapped handler
	wrappedHandler.ServeHTTP(w, req)

	// Check that the value was set correctly
	if receivedValue != testValue {
		t.Errorf("Expected middleware value %v, got %v", testValue, receivedValue)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestContextMiddleware_ChainedMiddleware(t *testing.T) {
	ctx1 := httputil.NewAppContext("middleware1")
	ctx2 := httputil.NewAppContext("middleware2")

	value1 := "value1"
	value2 := "value2"

	var receivedValue1, receivedValue2 interface{}

	// Final handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedValue1 = ctx1.GetContext(r)
		receivedValue2 = ctx2.GetContext(r)
		w.WriteHeader(http.StatusOK)
	})

	// Chain middlewares
	middleware1 := httputil.ContextMiddleware(ctx1, value1)
	middleware2 := httputil.ContextMiddleware(ctx2, value2)

	wrappedHandler := middleware1(middleware2(handler))

	// Execute
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w, req)

	// Both values should be present
	if receivedValue1 != value1 {
		t.Errorf("Expected value1 %v, got %v", value1, receivedValue1)
	}

	if receivedValue2 != value2 {
		t.Errorf("Expected value2 %v, got %v", value2, receivedValue2)
	}
}

func TestContextToken_String(t *testing.T) {
	names := []string{"test", "user-id", "session", "auth-token", ""}

	for _, name := range names {
		token := httputil.NewContextToken(name)
		if token.String() != name {
			t.Errorf("Expected token name %s, got %s", name, token.String())
		}
	}
}

// Benchmarks
func BenchmarkAppContext_SetContext(b *testing.B) {
	appCtx := httputil.NewAppContext("benchmark")
	req, _ := http.NewRequest("GET", "/", nil)
	testValue := "benchmark-value"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = appCtx.SetContext(req, testValue)
	}
}

func BenchmarkAppContext_GetContext(b *testing.B) {
	appCtx := httputil.NewAppContext("benchmark")
	req, _ := http.NewRequest("GET", "/", nil)
	testValue := "benchmark-value"
	req = appCtx.SetContext(req, testValue)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = appCtx.GetContext(req)
	}
}

func BenchmarkAppContext_WithContext(b *testing.B) {
	appCtx := httputil.NewAppContext("benchmark")
	req, _ := http.NewRequest("GET", "/", nil)
	testValue := "benchmark-value"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = appCtx.WithContext(req, testValue)
	}
}

func BenchmarkAppContext_WithContext_ExistingValue(b *testing.B) {
	appCtx := httputil.NewAppContext("benchmark")
	req, _ := http.NewRequest("GET", "/", nil)
	req = appCtx.SetContext(req, "initial-value")
	newValue := "new-value"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = appCtx.WithContext(req, newValue)
	}
}

func BenchmarkContextMiddleware(b *testing.B) {
	appCtx := httputil.NewAppContext("benchmark")
	testValue := "benchmark-value"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = appCtx.GetContext(r)
		w.WriteHeader(http.StatusOK)
	})

	middleware := httputil.ContextMiddleware(appCtx, testValue)
	wrappedHandler := middleware(handler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(w, req)
	}
}

func BenchmarkMultipleContexts(b *testing.B) {
	contexts := make([]*httputil.AppContext, 10)
	for i := 0; i < 10; i++ {
		contexts[i] = httputil.NewAppContext("context" + string(rune(i+48)))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/", nil)

		// Set values in all contexts
		for j, ctx := range contexts {
			req = ctx.SetContext(req, "value"+string(rune(j+48)))
		}

		// Retrieve values from all contexts
		for _, ctx := range contexts {
			_ = ctx.GetContext(req)
		}
	}
}

func BenchmarkNewAppContext(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = httputil.NewAppContext("benchmark-context")
	}
}

func BenchmarkNewContextToken(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = httputil.NewContextToken("benchmark-token")
	}
}
