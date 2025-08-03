# Env Package

The `env` package provides utilities for working with environment variables, including fallback support and presence checking.

## Installation

```bash
go get github.com/sampson-golang/utilities/env
```

## Functions

### `Exists`

Check if an environment variable exists (regardless of its value).

```go
func exampleExists() {
  os.Setenv("EXISTING_VAR", "value")
  os.Setenv("EMPTY_VAR", "")

  fmt.Println(env.Exists("EXISTING_VAR"))  // true
  fmt.Println(env.Exists("EMPTY_VAR"))     // true
  fmt.Println(env.Exists("MISSING_VAR"))   // false
}
```

### `Lookup`

Get an environment variable with fallback support.

the first argument is required, and an optional string to use as a fallback

you can also pass an arbitrary number of env vars to fallback to. if more than one falback string is given, all of them are treated as ENV keys, except for the final static string value.

```go
package main

import (
  "fmt"
  "os"
  "github.com/sampson-golang/utilities/env"
)

func main() {
  // Lookup with only string fallback
  url, exists := env.Lookup("DATABASE_URL", "sqlite://default.db")
  fmt.Println(url, exists)

  // Lookup with both ENVs and string fallback
  dbUrl, exists := env.Lookup("PRIMARY_DB", "BACKUP_DB", "OTHER_BACKUP" "sqlite://fallback.db")
  fmt.Println(dbUrl, exists)

  // Lookup with only ENV fallback
  dbUrl, exists := env.Lookup("PRIMARY_DB", "BACKUP_DB", "")
  fmt.Println(dbUrl, exists)

  // Get without fallbacks
  missing, exists := env.Lookup("PRIMARY_DB")
  fmt.Println(missing, exists)
}
```

### Get

Same as `Lookup` except only returns the string value, not the existance bool
Useful for inline function calls

### `LookupPresent`

Function signature is exactly the same as `Lookup`, except treats empty string ENVs as non-existent

```go
func exampleLookupPresent() {
  // Set variables including an empty one
  os.Setenv("API_KEY", "secret123")
  os.Setenv("EMPTY_VAR", "")
  os.Setenv("FALLBACK_KEY", "fallback123")

  // LookupPresent ignores empty values
  key, exists := env.LookupPresent("API_KEY")
  fmt.Println(key, exists) // "secret123", true

  // Empty variable falls through to fallback
  key, exists = env.LookupPresent("EMPTY_VAR", "FALLBACK_KEY", "fallback_static_value")
  fmt.Println(key, exists) // "fallback123", true

  // Compare with regular Lookup (which would return exists true)
  empty, exists := env.Lookup("EMPTY_VAR")
  fmt.Println(empty, exists) // "", false
}
```

### GetPresent

Same as `LookupPresent` except only returns the string value, not the existance bool
Useful for inline function calls

## API Reference

### `Lookup(key string, fallbacks ...string) (string, bool)`

Retrieves an environment variable with fallback support.

**Parameters:**
- `key` - The primary environment variable name to look up
- `fallbacks` - Variable number of fallback values.

**Returns:**
- `string` - The first existing key value found or the final fallback value
- `bool` - `true` if any environment variable was found, `false` if using a literal fallback

**Behavior:**
1. Checks if `key` exists in environment; if found, returns its value and `true`
2. If not found and no fallbacks provided, returns `""` and `false`
3. If fallbacks provided, recursively calls `Lookup` with the first fallback as the key and remaining fallbacks
4. If only one fallback remains, returns that literal value and `false`

### `Get(key string, fallbacks ...string) string`

Retrieves an environment variable with fallback support.

**Parameters:**
- `key` - The primary environment variable name to look up
- `fallbacks` - Variable number of fallback values.

**Returns:**
- `string` - The first existing key value found or the final fallback value

**Behavior:**
- Same as `Lookup`, but only returns the string value

### `LookupPresent(key string, fallbacks ...string) (string, bool)`

Similar to `Lookup` but only considers non-empty values.

**Parameters:**
- `key` - The primary environment variable name to look up
- `fallbacks` - Variable number of fallback variable names or literal value

**Returns:**
- `string` - The non-empty value found or the final fallback value
- `bool` - `true` if any non-empty environment variable was found, `false` if using a literal fallback

**Behavior:**
- Same as `Lookup`, but treats empty string values as if the variable doesn't exist

### `GetPresent(key string, fallbacks ...string) string`

Similar to `Get` but only considers non-empty values.

**Parameters:**
- `key` - The primary environment variable name to look up
- `fallbacks` - Variable number of fallback variable names or literal value

**Returns:**
- `string` - The non-empty value found or the final fallback value

**Behavior:**
- Same as `LookupPresent`, but treats empty string values as if the variable doesn't exist

### `Exists(key string) bool`

Checks if an environment variable exists.

**Parameters:**
- `key` - The environment variable name to check

**Returns:**
- `bool` - `true` if the variable exists (even if empty), `false` otherwise

## Examples

### Configuration Loading

```go
type Config struct {
  DatabaseURL string
  RedisURL    string
  Port        string
  Debug       bool
}

func LoadConfig() *Config {
  config := &Config{}

  // Database with fallbacks
  config.DatabaseURL = env.Get("DATABASE_URL", "DB_URL", "sqlite://app.db")

  // Redis with environment-specific fallbacks
  config.RedisURL = env.Get("REDIS_URL", "CACHE_URL", "redis://localhost:6379")

  // Port with numeric fallback
  config.Port = env.GetPresent("PORT", "HTTP_PORT", "8080")

  // Debug mode - only if explicitly set and non-empty
  debugStr, exists := env.LookupPresent("DEBUG", "ENABLE_DEBUG", "")
  config.Debug = exists && (debugStr == "true" || debugStr == "1")

  return config
}
```

### Environment Detection

```go
func GetEnvironment() string {
  // Check for explicit environment setting
  if env, exists := env.LookupPresent("ENVIRONMENT", "ENV", ""); exists {
    return env
  }

  // Detect based on other variables
  if env.Exists("HEROKU_APP_NAME") {
    return "heroku"
  }

  if env.Exists("AWS_EXECUTION_ENV") {
    return "aws"
  }

  if env.Exists("KUBERNETES_SERVICE_HOST") {
    return "kubernetes"
  }

  return "development"
}
```

### Feature Flags

```go
func IsFeatureEnabled(featureName string) bool {
  // Check specific feature flag
  flagName := "FEATURE_" + strings.ToUpper(featureName)
  if value, exists := env.LookupPresent(flagName); exists {
    return value == "true" || value == "1" || value == "enabled"
  }

  // Check global feature flags
  if value, exists := env.LookupPresent("ENABLE_ALL_FEATURES", "BETA_MODE"); exists {
    return value == "true" || value == "1"
  }

  return false
}

func main() {
  os.Setenv("FEATURE_NEW_UI", "true")
  os.Setenv("ENABLE_ALL_FEATURES", "1")

  fmt.Println(IsFeatureEnabled("new_ui"))    // true
  fmt.Println(IsFeatureEnabled("beta_api"))  // true (global flag)
  fmt.Println(IsFeatureEnabled("unknown"))   // false
}
```

### Database Connection

```go
func ConnectDatabase() (*sql.DB, error) {
  // Try multiple database URL formats and variables
  dbURL, exists := env.Lookup(
    "DATABASE_URL",           // Primary
    "DB_CONNECTION_STRING",   // Alternative name
    "POSTGRES_URL",          // Postgres-specific
    "MYSQL_URL",             // MySQL-specific
  )

  if !exists {
    return nil, fmt.Errorf("no database URL configured")
  }

  // Determine driver from URL or separate variable
  driver, _ := env.Lookup("DB_DRIVER", "postgres") // Default to postgres

  return sql.Open(driver, dbURL)
}
```

### Multi-Environment Setup

```go
func getConfigByEnvironment() map[string]string {
  environment, _ := env.Lookup("ENVIRONMENT", "development")

  config := make(map[string]string)

  switch environment {
  case "production":
    config["db"], _ = env.LookupPresent("PROD_DATABASE_URL")
    config["cache"], _ = env.LookupPresent("PROD_REDIS_URL")
    config["api"], _ = env.LookupPresent("PROD_API_URL")
  case "staging":
    config["db"], _ = env.LookupPresent("STAGING_DATABASE_URL", "PROD_DATABASE_URL")
    config["cache"], _ = env.LookupPresent("STAGING_REDIS_URL", "PROD_REDIS_URL")
    config["api"], _ = env.LookupPresent("STAGING_API_URL", "PROD_API_URL")
  default: // development
    config["db"], _ = env.Lookup("DEV_DATABASE_URL", "sqlite://dev.db")
    config["cache"], _ = env.Lookup("DEV_REDIS_URL", "redis://localhost:6379")
    config["api"], _ = env.Lookup("DEV_API_URL", "http://localhost:3000")
  }

  return config
}
```

### Validation and Reporting

```go
func validateEnvironment() []string {
  var missing []string

  required := []string{
    "DATABASE_URL",
    "API_SECRET_KEY",
    "JWT_SECRET",
  }

  optional := map[string]string{
    "REDIS_URL":     "redis://localhost:6379",
    "LOG_LEVEL":     "info",
    "MAX_WORKERS":   "10",
  }

  // Check required variables
  for _, key := range required {
    if !env.Exists(key) {
      missing = append(missing, key)
    }
  }

  // Set defaults for optional variables
  for key, defaultValue := range optional {
    if !env.Exists(key) {
      os.Setenv(key, defaultValue)
    }
  }

  return missing
}
```

## Testing

Run the tests with:

```bash
go test github.com/sampson-golang/utilities/env
```

See the test files for comprehensive examples:
- [`Get_test.go`](./Get_test.go)
- [`LookupPresent_test.go`](./LookupPresent_test.go)
- [`Exists_test.go`](./Exists_test.go)

## Best Practices

1. **Use fallback chains**: Order fallbacks from most specific to most general
2. **Use `LookupPresent` for optional configs**: When empty strings should be treated as unset
3. **Provide sensible defaults**: Always include a final literal fallback for non-critical settings
4. **Validate critical variables**: Use `Exists` to ensure required variables are set
5. **Document your variables**: Clearly document which environment variables your application uses
