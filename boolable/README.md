# Boolable Package

The `boolable` package provides intelligent boolean conversion from various types including strings, integers, pointers, and interfaces.

## Installation

```bash
go get github.com/sampson-golang/utilities/boolable
```

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/sampson-golang/utilities/boolable"
)

func main() {
    // String conversions
    fmt.Println(boolable.From("true"))    // true
    fmt.Println(boolable.From("false"))   // false
    fmt.Println(boolable.From("yes"))     // true
    fmt.Println(boolable.From("no"))      // false
    fmt.Println(boolable.From("1"))       // true
    fmt.Println(boolable.From("0"))       // false
    fmt.Println(boolable.From(""))        // false

    // Integer conversions
    fmt.Println(boolable.From(1))         // true
    fmt.Println(boolable.From(0))         // false
    fmt.Println(boolable.From(-1))        // true

    // Boolean values
    fmt.Println(boolable.From(true))      // true
    fmt.Println(boolable.From(false))     // false

    // Nil values
    fmt.Println(boolable.From(nil))       // false
}
```

### Working with Pointers

```go
func examplePointers() {
    str := "true"
    num := 42
    zero := 0

    // Automatic dereferencing (default behavior)
    fmt.Println(boolable.From(&str))      // true
    fmt.Println(boolable.From(&num))      // true
    fmt.Println(boolable.From(&zero))     // false

    var nilPtr *string
    fmt.Println(boolable.From(nilPtr))    // false

    // Disable dereferencing
    fmt.Println(boolable.From(&str, false))  // true (pointer exists)
    fmt.Println(boolable.From(nilPtr, false)) // false (nil pointer)
}
```

## API Reference

### `From(value interface{}, dereference ...bool) bool`

Converts any value to a boolean using intelligent rules.

**Parameters:**
- `value` - The value to convert to boolean
- `dereference` - Optional boolean flag. If true (default), pointers are dereferenced before conversion. If false, only checks if pointer is nil.

**Returns:**
- `bool` - The boolean representation of the input value

**Conversion Rules:**

| Input Type | Conversion Logic |
|------------|------------------|
| `bool` | Returns the boolean value as-is |
| `int` | `false` if zero, `true` otherwise |
| `string` | `false` for empty string or "false", "f", "0", "off", "n", "no" (case-insensitive), `true` otherwise |
| `nil` | Always returns `false` |
| `*T` (pointer) | If `dereference` is true (default), dereferences and applies rules to the pointed value. If false, returns `false` only for nil pointers |
| Other types | Always returns `true` |

**False Values for Strings:**
- `""` (empty string)
- `"0"`
- `"f"` / `"F"`
- `"false"` / `"FALSE"` / `"False"`
- `"off"` / `"OFF"` / `"Off"`
- `"n"` / `"N"`
- `"no"` / `"NO"` / `"No"`

All string comparisons are case-insensitive.

## Examples

### Environment Variable Parsing

```go
import (
  "os"
  "github.com/sampson-golang/utilities/boolable"
)

func isFeatureEnabled() bool {
  return boolable.From(os.Getenv("ENABLE_FEATURE"))
}

// Usage:
// ENABLE_FEATURE=true  -> true
// ENABLE_FEATURE=1     -> true
// ENABLE_FEATURE=yes   -> true
// ENABLE_FEATURE=false -> false
// ENABLE_FEATURE=0     -> false
// ENABLE_FEATURE=      -> false (empty)
```

### Configuration Struct

```go
type Config struct {
  Debug   *bool   `json:"debug"`
  Verbose *string `json:"verbose"`
}

func (c *Config) IsDebugEnabled() bool {
  return boolable.From(c.Debug)
}

func (c *Config) IsVerboseEnabled() bool {
  return boolable.From(c.Verbose)
}
```

### JSON API Parsing

```go
func handleAPIRequest(data map[string]interface{}) {
  if boolable.From(data["force_update"]) {
    // Handle forced update
  }

  if boolable.From(data["send_notifications"]) {
    // Send notifications
  }
}
```

## Testing

Run the tests with:

```bash
go test github.com/sampson-golang/utilities/boolable
```

See [`From_test.go`](./From_test.go) for comprehensive test cases covering all conversion scenarios.
