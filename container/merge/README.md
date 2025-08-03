# Merge Subpackage

The `merge` subpackage provides utilities for merging maps and structs, allowing you to combine data from multiple sources.

## Installation

```bash
go get github.com/sampson-golang/utilities/container/merge
```

## Functions

### `Params`

Merge multiple string maps into a destination map. Values from source maps overwrite existing values in the destination.

```go
package main

import (
  "fmt"
  "github.com/sampson-golang/utilities/container/merge"
)

func main() {
  // Destination map
  config := map[string]string{
    "host": "localhost",
    "port": "8080",
  }

  // Source maps
  defaults := map[string]string{
    "port":    "3000",
    "timeout": "30s",
  }

  overrides := map[string]string{
    "host": "production.com",
    "ssl":  "true",
  }

  // Merge sources into destination
  merge.Params(config, defaults, overrides)

  fmt.Printf("%+v\n", config)
  // Output: map[host:production.com port:3000 ssl:true timeout:30s]
}
```

### `Structs`

Merge multiple structs into a destination struct. Only non-zero (non-empty) fields from source structs are copied to the destination.

```go
type Config struct {
  Host     string
  Port     int
  Database string
  Debug    bool
}

func exampleStructs() {
  // Destination struct
  config := &Config{
    Host: "localhost",
    Port: 8080,
  }

  // Source structs
  defaults := &Config{
    Port:     3000,    // Will not overwrite (config.Port is not zero)
    Database: "mydb",  // Will be copied (config.Database is zero)
    Debug:    true,    // Will be copied (config.Debug is zero)
  }

  overrides := &Config{
    Host:     "prod.com", // Will not overwrite (config.Host is not zero)
    Database: "proddb",   // Will overwrite defaults.Database
  }

  // Merge sources into destination
  merge.Structs(config, defaults, overrides)

  fmt.Printf("%+v\n", *config)
  // Output: {Host:localhost Port:8080 Database:proddb Debug:true}
}
```

## API Reference

### `Params(into map[string]string, from ...map[string]string)`

Merges multiple string maps into a destination map.

**Parameters:**
- `into` - The destination map that will receive merged values
- `from` - Variable number of source maps to merge from

**Behavior:**
- Sources are merged in order from left to right
- Later sources overwrite values from earlier sources
- All keys from all source maps are copied to the destination
- The destination map is modified in-place

**Example merge order:**
```go
dest := map[string]string{"a": "1"}
src1 := map[string]string{"a": "2", "b": "3"}
src2 := map[string]string{"a": "4", "c": "5"}

merge.Params(dest, src1, src2)
// Result: {"a": "4", "b": "3", "c": "5"}
```

### `Structs(into interface{}, from ...interface{})`

Merges multiple structs into a destination struct.

**Parameters:**
- `into` - Pointer to the destination struct
- `from` - Variable number of pointers to source structs

**Behavior:**
- Sources are merged in order from left to right
- Only non-zero (non-empty) fields from sources are copied
- values in the destination are overwritten by non-zero values from sources
- structs must have compatible type structures, only fields from the destination are considered

**Zero values by type:**
- `string`: `""`
- `int`, `int32`, `int64`, etc.: `0`
- `bool`: `false`
- `float32`, `float64`: `0.0`
- Pointers: `nil`
- Slices, maps, channels: `nil`

## Examples

### Configuration Management

```go
type ServerConfig struct {
  Host         string
  Port         int
  DatabaseURL  string
  Debug        bool
  Timeout      time.Duration
}

func loadConfiguration() *ServerConfig {
  // Start with base configuration
  config := &ServerConfig{
    Host: "localhost",
    Port: 8080,
  }

  // Apply defaults
  defaults := &ServerConfig{
    Port:        3000,
    DatabaseURL: "postgres://localhost/dev",
    Debug:       false,
    Timeout:     30 * time.Second,
  }

  // Apply environment-specific overrides
  prod := &ServerConfig{
    DatabaseURL: "postgres://prod-server/app",
    Debug:       false,
  }

  // Apply user overrides
  user := &ServerConfig{
    Debug: true,
  }

  merge.Structs(config, defaults, prod, user)
  return config
  // Result: localhost:8080, prod database, debug enabled
}
```

### HTTP Request Parameter Merging

```go
func buildRequestParams(base, query, body map[string]string) map[string]string {
  params := make(map[string]string)

  // Merge in order of precedence: base < query < body
  merge.Params(params, base, query, body)

  return params
}

func handleRequest() {
  baseParams := map[string]string{
    "version": "1.0",
    "format":  "json",
  }

  queryParams := map[string]string{
    "format": "xml",  // Overrides base
    "limit":  "10",
  }

  bodyParams := map[string]string{
    "limit": "20",    // Overrides query
    "sort":  "asc",
  }

  final := buildRequestParams(baseParams, queryParams, bodyParams)
  // Result: {"version": "1.0", "format": "xml", "limit": "20", "sort": "asc"}
}
```

### Feature Flag Merging

```go
type FeatureFlags struct {
  EnableNewUI     bool
  EnableBetaAPI   bool
  MaxConnections  int
  CacheEnabled    bool
}

func getFeatureFlags(userID string) *FeatureFlags {
  // Global defaults
  flags := &FeatureFlags{
    EnableNewUI:    false,
    EnableBetaAPI:  false,
    MaxConnections: 100,
    CacheEnabled:   true,
  }

  // Environment-specific flags
  if isProduction() {
    prodFlags := &FeatureFlags{
      EnableBetaAPI:  false, // Ensure beta is off in prod
      MaxConnections: 1000,  // Higher limits in prod
    }
    merge.Structs(flags, prodFlags)
  }

  // User-specific overrides
  if userFlags := getUserFlags(userID); userFlags != nil {
    merge.Structs(flags, userFlags)
  }

  return flags
}
```

### API Response Composition

```go
type APIResponse struct {
  Status  string
  Message string
  Data    map[string]string
  Meta    map[string]string
}

func buildResponse() *APIResponse {
  // Base response
  response := &APIResponse{
    Status: "success",
    Data:   make(map[string]string),
    Meta:   make(map[string]string),
  }

  // Add default metadata
  defaultMeta := map[string]string{
    "version":   "1.0",
    "timestamp": time.Now().Format(time.RFC3339),
  }
  merge.Params(response.Meta, defaultMeta)

  // Add request-specific data
  requestData := map[string]string{
    "user_id": "123",
    "action":  "create",
  }
  merge.Params(response.Data, requestData)

  // Add debug info if needed
  if isDebugMode() {
    debugMeta := map[string]string{
      "debug":        "true",
      "query_count":  "5",
      "response_time": "120ms",
    }
    merge.Params(response.Meta, debugMeta)
  }

  return response
}
```

## Testing

Run the tests with:

```bash
go test github.com/sampson-golang/utilities/container/merge
```

See the test files for comprehensive examples:
- [`Params_test.go`](./Params_test.go)
- [`Structs_test.go`](./Structs_test.go)
