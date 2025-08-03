# Container Package

The `container` package provides utilities for working with data structures including slices, maps, sets, and nested data traversal.

## Installation

```bash
go get github.com/sampson-golang/utilities/container
```

## Functions

### `Contains`

Check if a string slice contains a specific item.

```go
package main

import (
  "fmt"
  "github.com/sampson-golang/utilities/container"
)

func main() {
  fruits := []string{"apple", "banana", "orange"}

  fmt.Println(container.Contains(fruits, "banana"))  // true
  fmt.Println(container.Contains(fruits, "grape"))   // false
}
```

### `Dig`

Safely traverse nested map/slice structures using a path of keys and indexes.

```go
func exampleDig() {
  data := map[string]interface{}{
    "users": []interface{}{
      map[string]interface{}{
        "name": "John",
        "profile": map[string]interface{}{
          "age": 30,
          "city": "New York",
        },
      },
      map[string]interface{}{
        "name": "Jane",
        "profile": map[string]interface{}{
          "age": 25,
          "city": "San Francisco",
        },
      },
    },
  }

  // Get the first user's name
  name := container.Dig(data, "users", 0, "name")
  if name != nil {
    fmt.Println(*name.(*interface{})) // "John"
  }

  // Get the second user's city
  city := container.Dig(data, "users", 1, "profile", "city")
  if city != nil {
    fmt.Println(*city.(*interface{})) // "San Francisco"
  }

  // Safe access - returns nil if path doesn't exist
  missing := container.Dig(data, "users", 5, "name")
  fmt.Println(missing) // nil
}
```

### `DigAssign`

Assign the result of `Dig` to a struct field with type conversion.

```go
type User struct {
  Name string
  Age  int
  City string
}

func exampleDigAssign() {
  data := map[string]interface{}{
    "user": map[string]interface{}{
      "name": "Alice",
      "age":  28,
      "location": map[string]interface{}{
        "city": "Boston",
      },
    },
  }

  var user User

  // Assign values from nested data to struct fields
  container.DigAssign(&user, "Name", data, "user", "name")
  container.DigAssign(&user, "Age", data, "user", "age")
  container.DigAssign(&user, "City", data, "user", "location", "city")

  fmt.Printf("%+v\n", user) // {Name:Alice Age:28 City:Boston}
}
```

### `Set`

A deduplicating data structure implementation using Go maps.

```go
func exampleSet() {
  // Create a new set
  s := container.Set{}

  // Add elements
  s.Add("apple")
  s.Add("banana")
  s.Add("apple") // Duplicates are ignored

  // Check membership
  fmt.Println(s.Has("apple"))   // true
  fmt.Println(s.Has("grape"))   // false

  // Get all values
  values := s.Values()
  fmt.Println(values) // [apple banana] (order may vary)

  // Remove elements
  s.Remove("banana")
  fmt.Println(s.Has("banana"))  // false
}
```

## API Reference

### `Contains(slice []string, item string) bool`

**Parameters:**
- `slice` - The string slice to search in
- `item` - The string to search for

**Returns:**
- `bool` - `true` if the item is found, `false` otherwise

### `Dig(data interface{}, path ...any) interface{}`

**Parameters:**
- `data` - The nested map/slice structure to traverse
- `path` - Variable number of keys (strings) and indexes (integers) defining the path

**Returns:**
- `interface{}` - Pointer to the value if found, `nil` if path doesn't exist

**Supported path elements:**
- `string` - Key for `map[string]interface{}`
- `int` - Index for `[]interface{}`
- `string` representing a number - Converted to int for slice indexing

### `DigAssign(result interface{}, key string, data interface{}, path ...any)`

**Parameters:**
- `result` - Pointer to struct where the value should be assigned
- `key` - Field name in the struct
- `data` - The nested data structure (same as `Dig`)
- `path` - Path to the value (same as `Dig`)

**Behavior:**
- Only assigns if the path exists and value is not nil
- Performs type conversion if possible
- Does nothing if field doesn't exist or types are incompatible

### `Set` Type

A set implementation with the following methods:

#### `Add(value any)`
Add an element to the set. Duplicates are ignored.

#### `Remove(value any)`
Remove an element from the set.

#### `Has(value any) bool`
Check if the set contains an element.

#### `Values() []any`
Get all elements in the set as a slice. Order is not guaranteed.

## Merge Subpackage

The `merge` subpackage provides utilities for merging maps and structs. See [`merge/README.md`](./merge/README.md) for detailed documentation.

```go
import "github.com/sampson-golang/utilities/container/merge"

// Quick examples:
params := map[string]string{"key1": "value1"}
merge.Params(params, map[string]string{"key2": "value2"})

var dest MyStruct
merge.Structs(&dest, &source1, &source2)
```

## Examples

### JSON Processing

```go
func processJSONData(jsonData map[string]interface{}) {
  // Extract nested API response data
  if userId := container.Dig(jsonData, "response", "data", "user", "id"); userId != nil {
    fmt.Println("User ID:", *userId.(*interface{}))
  }

  // Process array of items
  if items := container.Dig(jsonData, "response", "items"); items != nil {
    if itemsSlice, ok := (*items.(*interface{})).([]interface{}); ok {
      for i := range itemsSlice {
        if name := container.Dig(jsonData, "response", "items", i, "name"); name != nil {
          fmt.Println("Item:", *name.(*interface{}))
        }
      }
    }
  }
}
```

### Configuration Management

```go
type Config struct {
  DatabaseURL string
  Port       int
  Debug      bool
}

func loadConfig(data map[string]interface{}) Config {
  var config Config

  container.DigAssign(&config, "DatabaseURL", data, "database", "url")
  container.DigAssign(&config, "Port", data, "server", "port")
  container.DigAssign(&config, "Debug", data, "debug", "enabled")

  return config
}
```

### Unique Value Tracking

```go
func removeDuplicates(items []string) []string {
  seen := make(container.Set)
  var result []string

  for _, item := range items {
    if !seen.Has(item) {
      seen.Add(item)
      result = append(result, item)
    }
  }

  return result
}
```

## Testing

Run the tests with:

```bash
go test github.com/sampson-golang/utilities/container
go test github.com/sampson-golang/utilities/container/merge
```

See the test files for comprehensive examples:
- [`Contains_test.go`](./Contains_test.go)
- [`Dig_test.go`](./Dig_test.go)
- [`merge/Params_test.go`](./merge/Params_test.go)
- [`merge/Structs_test.go`](./merge/Structs_test.go)
