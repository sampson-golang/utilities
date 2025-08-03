# Golang Utilities

A collection of useful Go utility functions and packages for common programming tasks.

## Version

Current version: **0.0.1**

## Installation

```bash
go get github.com/sampson-golang/utilities
```

## Quick Start

```go
package main

import (
  "fmt"
  "github.com/sampson-golang/utilities"
)

func main() {
  // Convert various types to boolean
  fmt.Println(utilities.Boolable("true"))  // true
  fmt.Println(utilities.Boolable(0))       // false

  // Check if slice contains item
  slice := []string{"apple", "banana", "orange"}
  fmt.Println(utilities.Contains(slice, "banana")) // true

  // Get environment variable with fallback
  value, exists := utilities.LookupEnv("DATABASE_URL", "localhost:5432")
  fmt.Println(value, exists)

  // Pretty print JSON
  data := map[string]interface{}{"name": "John", "age": 30}
  utilities.PrettyPrint(data)

  // Clean up whitespace
  clean := utilities.Squish("  hello    world  ")
  fmt.Println(clean) // "hello world"
}
```

## Ease of Use API

The utilities package makes a select list of functions and types available directly from the root package
NOTE: Not all functions and types are available through the root package

| Convenience Function | Actual Location |
|---------------------|-----------------|
| `utilities.Boolable` | `boolable.From` |
| `utilities.Contains` | `container.Contains` |
| `utilities.Dig` | `container.Dig` |
| `utilities.DigAssign` | `container.DigAssign` |
| `utilities.EnvExists` | `env.Exists` |
| `utilities.GetEnv` | `env.Get` |
| `utilities.GetPresentEnv` | `env.GetPresent` |
| `utilities.LookupEnv` | `env.Lookup` |
| `utilities.LookupPresentEnv` | `env.LookupPresent` |
| `utilities.MergeParams` | `merge.Params` |
| `utilities.MergeStructs` | `merge.Structs` |
| `utilities.Prettify` | `output.Prettify` |
| `utilities.PrettyPrint` | `output.PrettyPrint` |
| `utilities.Set` | `container.Set` |
| `utilities.Squish` | `strutil.Squish` |

## Packages

Users are encouraged to import subpackages directly instead of using the root-level ease-of-use access

### [`boolable`](./boolable/README.md)
Convert various types (strings, numbers, pointers) to boolean values with intelligent defaults.

### [`container`](./container/README.md)
Utilities for working with slices, maps, and data structures including:
- `Contains` - Check if slice contains an item
- `Dig` - Safely traverse nested maps/slices
- `Set` - Set data structure implementation
- [`merge`](./container/merge/) - Merge maps and structs

### [`env`](./env/README.md)
Environment variable utilities with fallback support:
- `Lookup` - Get env var with fallbacks
- `LookupPresent` - Get non-empty env var with fallbacks
- `Exists` - Check if env var exists

### [`httputil`](./httputil/README.md)
HTTP utilities for web applications including context management, port checking, and request parsing.

### [`output`](./output/README.md)
Pretty printing and JSON formatting utilities.

### [`strutil`](./strutil/README.md)
String manipulation utilities including whitespace normalization.

## Development

### Running Tests

```bash
# Run all tests
./bin/test

# Or use go test
go test ./...
```

### Code Formatting

```bash
./bin/hooks/format
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Run tests (`./bin/test`)
4. Commit your changes (`git commit -m 'Add amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
