# Output Package

The `output` package provides utilities for pretty-printing and formatting of data.

## Installation

```bash
go get github.com/sampson-golang/utilities/output
```

## Functions

### `Prettify`

Convert any Go value to a pretty-formatted JSON string with proper indentation.

```go
  package main

import (
  "fmt"
  "github.com/sampson-golang/utilities/output"
)

func main() {
  // Simple data structures
  data := map[string]interface{}{
    "name":  "John Doe",
    "age":   30,
    "active": true,
    "scores": []int{95, 87, 92},
    "address": map[string]string{
      "street": "123 Main St",
      "city":   "New York",
      "zip":    "10001",
    },
  }

  // Get pretty JSON string
  jsonStr := output.Prettify(data)
  fmt.Println(jsonStr)
}
```

Output:
```json
{
  "active": true,
  "address": {
    "city": "New York",
    "street": "123 Main St",
    "zip": "10001"
  },
  "age": 30,
  "name": "John Doe",
  "scores": [
    95,
    87,
    92
  ]
}
```

### `PrettyPrint`

Directly print any Go value as pretty-formatted JSON to stdout.

```go
func examplePrettyPrint() {
  user := struct {
    ID       int      `json:"id"`
    Name     string   `json:"name"`
    Email    string   `json:"email"`
    Tags     []string `json:"tags"`
    Settings map[string]bool `json:"settings"`
  }{
    ID:    123,
    Name:  "Alice Smith",
    Email: "alice@example.com",
    Tags:  []string{"admin", "developer", "team-lead"},
    Settings: map[string]bool{
      "notifications": true,
      "dark_mode":     false,
      "beta_features": true,
    },
  }

  // Print directly to stdout
  output.PrettyPrint(user)
}
```

Output:
```json
{
  "id": 123,
  "name": "Alice Smith",
  "email": "alice@example.com",
  "tags": [
    "admin",
    "developer",
    "team-lead"
  ],
  "settings": {
    "beta_features": true,
    "dark_mode": false,
    "notifications": true
  }
}
```

## API Reference

### `Prettify(i interface{}, indent ...string) string`

Converts any Go value to a pretty-formatted JSON string.

**Parameters:**
- `i` - Any Go value to convert to JSON
- `indent` - Optional indent string, defaults to 2 spaces
- ...rest - ignored

**Returns:**
- `string` - Pretty-formatted JSON string using provided or 2-space indentation

**Behavior:**
- Uses `json.MarshalIndent` with `""` prefix
  - if `indent` is passed, it is used
  - if `indent` is not passed, 2-spaces are used
- Handles all JSON-serializable Go types
- Returns empty string if JSON marshaling fails (error is ignored)
- Does not panic on invalid input

### `PrettyPrint(i interface{}, indent ...string)`

Prints any Go value as pretty-formatted JSON directly to stdout.

**Parameters:**
- `i` - Any Go value to print

**Behavior:**
- Internally calls `Prettify` and prints the result with `fmt.Println`
- Adds a newline after the JSON output
- Does not return anything

## Examples

### Debug Logging

```go
import (
    "log"
    "github.com/sampson-golang/utilities/output"
)

func debugLog(label string, data interface{}) {
    if isDebugMode() {
        log.Printf("DEBUG [%s]:\n%s", label, output.Prettify(data))
    }
}

func processAPIRequest(request map[string]interface{}) {
    debugLog("Incoming Request", request)

    // Process request...
    result := map[string]interface{}{
        "status": "success",
        "processed_at": time.Now(),
        "items": []string{"item1", "item2"},
    }

    debugLog("Processing Result", result)
}
```

### Configuration Display

```go
type AppConfig struct {
    Database DatabaseConfig `json:"database"`
    Server   ServerConfig   `json:"server"`
    Features FeatureFlags   `json:"features"`
}

type DatabaseConfig struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Name     string `json:"name"`
    SSL      bool   `json:"ssl"`
}

type ServerConfig struct {
    Port         int           `json:"port"`
    ReadTimeout  time.Duration `json:"read_timeout"`
    WriteTimeout time.Duration `json:"write_timeout"`
}

type FeatureFlags struct {
    EnableNewUI    bool `json:"enable_new_ui"`
    EnableBetaAPI  bool `json:"enable_beta_api"`
    DebugMode      bool `json:"debug_mode"`
}

func displayConfig(config *AppConfig) {
    fmt.Println("Current Application Configuration:")
    fmt.Println("=" * 40)
    output.PrettyPrint(config)
    fmt.Println("=" * 40)
}
```

### API Response Debugging

```go
func makeAPICall(endpoint string) (*APIResponse, error) {
    resp, err := http.Get(endpoint)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    // Debug: Show raw response
    if isDebugMode() {
        var rawData interface{}
        json.Unmarshal(body, &rawData)
        fmt.Printf("API Response from %s:\n", endpoint)
        output.PrettyPrint(rawData)
    }

    var response APIResponse
    err = json.Unmarshal(body, &response)
    return &response, err
}
```

### Test Data Comparison

```go
func TestUserCreation(t *testing.T) {
    expected := User{
        ID:    123,
        Name:  "John Doe",
        Email: "john@example.com",
        Roles: []string{"user", "admin"},
    }

    actual := createUser("John Doe", "john@example.com")

    if !reflect.DeepEqual(expected, actual) {
        t.Errorf("User creation failed.\nExpected:\n%s\nActual:\n%s",
            output.Prettify(expected),
            output.Prettify(actual))
    }
}
```

### Data Structure Exploration

```go
func exploreDataStructure(data interface{}) {
    fmt.Printf("Data type: %T\n", data)
    fmt.Println("Data structure:")
    output.PrettyPrint(data)

    // Also show as one-line for comparison
    if jsonBytes, err := json.Marshal(data); err == nil {
        fmt.Printf("Compact JSON: %s\n", string(jsonBytes))
    }
}

func main() {
    complexData := map[string]interface{}{
        "users": []map[string]interface{}{
            {"id": 1, "name": "Alice", "active": true},
            {"id": 2, "name": "Bob", "active": false},
        },
        "metadata": map[string]interface{}{
            "total": 2,
            "page":  1,
            "filters": map[string]string{
                "status": "all",
                "sort":   "name",
            },
        },
    }

    exploreDataStructure(complexData)
}
```

### Environment Configuration Viewer

```go
func showEnvironmentConfig() {
    config := map[string]interface{}{
        "environment": os.Getenv("ENVIRONMENT"),
        "database": map[string]string{
            "host": os.Getenv("DB_HOST"),
            "port": os.Getenv("DB_PORT"),
            "name": os.Getenv("DB_NAME"),
        },
        "features": map[string]bool{
            "debug":      os.Getenv("DEBUG") == "true",
            "enable_api": os.Getenv("ENABLE_API") != "false",
        },
        "timestamps": map[string]string{
            "startup": time.Now().Format(time.RFC3339),
            "build":   getBuildTime(), // Your build time function
        },
    }

    fmt.Println("Environment Configuration:")
    output.PrettyPrint(config)
}
```

### HTTP Request/Response Logging

```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Log request
        requestData := map[string]interface{}{
            "method":     r.Method,
            "url":        r.URL.String(),
            "headers":    r.Header,
            "remote_ip":  r.RemoteAddr,
            "user_agent": r.UserAgent(),
            "timestamp":  time.Now().Format(time.RFC3339),
        }

        log.Printf("Incoming Request:\n%s", output.Prettify(requestData))

        // Capture response
        recorder := httptest.NewRecorder()
        next.ServeHTTP(recorder, r)

        // Log response
        responseData := map[string]interface{}{
            "status_code": recorder.Code,
            "headers":     recorder.Header(),
            "body_size":   len(recorder.Body.Bytes()),
        }

        log.Printf("Outgoing Response:\n%s", output.Prettify(responseData))

        // Copy response to actual writer
        for k, v := range recorder.Header() {
            w.Header()[k] = v
        }
        w.WriteHeader(recorder.Code)
        w.Write(recorder.Body.Bytes())
    })
}
```

## Use Cases

1. **Debugging**: Pretty-print complex data structures for easier inspection
2. **Logging**: Format structured logs with readable JSON output
3. **Testing**: Display expected vs actual data in test failures
4. **Configuration**: Show application configuration in a readable format
5. **API Development**: Debug API responses and request data
6. **Data Exploration**: Understand the structure of unknown data

## Testing

Run the tests with:

```bash
go test github.com/sampson-golang/utilities/output
```

See the test files for comprehensive examples:
- [`Prettify_test.go`](./Prettify_test.go)
- [`PrettyPrint_test.go`](./PrettyPrint_test.go)

## Notes

- Both functions handle all JSON-serializable Go types
- Circular references will cause infinite recursion (same as standard `json` package)
- Time values are formatted according to their JSON marshaling implementation
- Private struct fields (lowercase) are not included in output
- Functions use struct tags like `json:"field_name"` for field naming
