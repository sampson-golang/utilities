# StrUtil Package

The `strutil` package provides string manipulation utilities for common text processing tasks.

## Installation

```bash
go get github.com/sampson-golang/utilities/strutil
```

## Functions

### `Squish`

Remove leading and trailing whitespace from a string and replace any sequence of internal whitespace characters with a single space.

```go
package main

import (
    "fmt"
    "github.com/sampson-golang/utilities/strutil"
)

func main() {
  // Basic whitespace cleanup
  text := "  hello    world  "
  cleaned := strutil.Squish(text)
  fmt.Printf("'%s' -> '%s'\n", text, cleaned) // '  hello    world  ' -> 'hello world'

  // Multiple types of whitespace
  text = "\t\n  hello  \t\n  world  \r\n  "
  cleaned = strutil.Squish(text)
  fmt.Printf("'%s' -> '%s'\n", text, cleaned) // (various whitespace) -> 'hello world'

  // No whitespace
  text = "hello"
  cleaned = strutil.Squish(text)
  fmt.Printf("'%s' -> '%s'\n", text, cleaned) // 'hello' -> 'hello'

  // Only whitespace
  text = "   \t\n  "
  cleaned = strutil.Squish(text)
  fmt.Printf("'%s' -> '%s'\n", text, cleaned) // (whitespace) -> ''

  // Empty string
  text = ""
  cleaned = strutil.Squish(text)
  fmt.Printf("'%s' -> '%s'\n", text, cleaned) // '' -> ''
}
```

## API Reference

### `Squish(s string) string`

Normalizes whitespace in a string by trimming leading/trailing whitespace and collapsing internal whitespace sequences.

**Parameters:**
- `s` - The input string to process

**Returns:**
- `string` - The processed string with normalized whitespace

**Behavior:**
1. Removes all leading and trailing whitespace using `strings.TrimSpace`
2. Replaces any sequence of one or more whitespace characters with a single space using regex `\s+`
3. Handles all Unicode whitespace characters (spaces, tabs, newlines, etc.)

**Whitespace characters handled:**
- Space (` `)
- Tab (`\t`)
- Newline (`\n`)
- Carriage return (`\r`)
- Form feed (`\f`)
- Vertical tab (`\v`)
- Other Unicode whitespace characters

## Examples

### User Input Processing

```go
func processUserInput(input string) string {
  // Clean up user input from forms, text areas, etc.
  cleaned := strutil.Squish(input)

  if cleaned == "" {
    return "No input provided"
  }

  return cleaned
}

func main() {
  inputs := []string{
    "  John   Doe  ",           // Form input with extra spaces
    "\t\nHello\n\tWorld\r\n",   // Text area with mixed whitespace
    "   ",                      // Only whitespace
    "",                         // Empty string
    "SingleWord",               // No whitespace
    "Normal text here",         // Already clean
  }

  for _, input := range inputs {
    result := processUserInput(input)
    fmt.Printf("Input: '%s' -> Output: '%s'\n", input, result)
  }
}
```

### Search Query Normalization

```go
func normalizeSearchQuery(query string) string {
  // Clean up search queries for better matching
  return strutil.Squish(strings.ToLower(query))
}

func searchProducts(query string) []Product {
  normalizedQuery := normalizeSearchQuery(query)

  if normalizedQuery == "" {
    return getAllProducts()
  }

  var results []Product
  for _, product := range products {
    if strings.Contains(normalizeSearchQuery(product.Name), normalizedQuery) ||
      strings.Contains(normalizeSearchQuery(product.Description), normalizedQuery) {
      results = append(results, product)
    }
  }

  return results
}

// Example usage:
// searchProducts("  LAPTOP   GAMING  ") -> searches for "laptop gaming"
// searchProducts("\toffice\n chair\t") -> searches for "office chair"
```

### Log Message Cleaning

```go
import (
  "log"
  "github.com/sampson-golang/utilities/strutil"
)

func logMessage(level, message string) {
  // Clean up log messages to ensure proper formatting
  cleanMessage := strutil.Squish(message)
  if cleanMessage == "" {
    cleanMessage = "Empty log message"
  }

  timestamp := time.Now().Format("2006-01-02 15:04:05")
  log.Printf("[%s] %s: %s", timestamp, level, cleanMessage)
}

func main() {
  // These all produce clean, well-formatted log entries
  logMessage("INFO", "  Server   started   successfully  ")
  logMessage("ERROR", "\n\tDatabase\t\tconnection\n\tfailed\r\n")
  logMessage("DEBUG", "Processing user ID: 123")
  logMessage("WARN", "   ")  // Handles empty/whitespace-only messages
}
```

### CSV/Data Processing

```go
func processCSVData(csvData [][]string) [][]string {
  var cleaned [][]string

  for _, row := range csvData {
    var cleanRow []string
    for _, cell := range row {
      // Clean each cell value
      cleanCell := strutil.Squish(cell)
      cleanRow = append(cleanRow, cleanCell)
    }
    cleaned = append(cleaned, cleanRow)
  }

  return cleaned
}

// Example: Clean up messy CSV data
func main() {
  messyData := [][]string{
    {"  John  ", "  Doe  ", "  john@email.com  "},
    {"\tJane\t", "\nSmith\n", " jane@email.com "},
    {"", "   Bob   ", "bob@email.com"},
  }

  cleanData := processCSVData(messyData)

  for i, row := range cleanData {
    fmt.Printf("Row %d: %v\n", i+1, row)
  }
  // Output:
  // Row 1: [John Doe john@email.com]
  // Row 2: [Jane Smith jane@email.com]
  // Row 3: [ Bob bob@email.com]
}
```

### API Parameter Cleaning

```go
func sanitizeAPIParams(params map[string]string) map[string]string {
  cleaned := make(map[string]string)

  for key, value := range params {
    cleanKey := strutil.Squish(key)
    cleanValue := strutil.Squish(value)

    // Only include non-empty parameters
    if cleanKey != "" && cleanValue != "" {
      cleaned[cleanKey] = cleanValue
    }
  }

  return cleaned
}

func handleAPIRequest(w http.ResponseWriter, r *http.Request) {
  // Extract and clean parameters
  rawParams := map[string]string{
    "  name  ":    "  John Doe  ",
    "\temail\n":   " john@example.com ",
    "":            "ignored",  // Empty key
    "department":  "   ",      // Empty value
    "role":        "admin",
  }

  cleanParams := sanitizeAPIParams(rawParams)
  // Result: {"name": "John Doe", "email": "john@example.com", "role": "admin"}

  // Process clean parameters...
}
```

### Text Content Processing

```go
func cleanTextContent(content string) string {
  // Clean up text content from various sources
  lines := strings.Split(content, "\n")
  var cleanLines []string

  for _, line := range lines {
    cleanLine := strutil.Squish(line)
    if cleanLine != "" {  // Skip empty lines
      cleanLines = append(cleanLines, cleanLine)
    }
  }

  return strings.Join(cleanLines, "\n")
}

func main() {
  messyText := `
    This is a paragraph with
        inconsistent    spacing


    Another paragraph here
        with   more   spacing   issues

  `

  clean := cleanTextContent(messyText)
  fmt.Println(clean)
  // Output:
  // This is a paragraph with inconsistent spacing
  // Another paragraph here with more spacing issues
}
```

### Configuration Value Processing

```go
type Config struct {
  DatabaseURL string
  AppName     string
  Environment string
}

func loadConfig() *Config {
  config := &Config{
    DatabaseURL: os.Getenv("DATABASE_URL"),
    AppName:     os.Getenv("APP_NAME"),
    Environment: os.Getenv("ENVIRONMENT"),
  }

  // Clean all configuration values
  config.DatabaseURL = strutil.Squish(config.DatabaseURL)
  config.AppName = strutil.Squish(config.AppName)
  config.Environment = strutil.Squish(config.Environment)

  // Set defaults for empty values
  if config.AppName == "" {
    config.AppName = "MyApp"
  }
  if config.Environment == "" {
    config.Environment = "development"
  }

  return config
}
```

### Form Data Validation

```go
func validateFormData(form map[string]string) (map[string]string, []string) {
  var errors []string
  cleaned := make(map[string]string)

  required := []string{"first_name", "last_name", "email"}

  for key, value := range form {
    cleanValue := strutil.Squish(value)
    cleaned[key] = cleanValue

    // Check required fields
    for _, req := range required {
      if key == req && cleanValue == "" {
        errors = append(errors, fmt.Sprintf("%s is required", req))
      }
    }
  }

  return cleaned, errors
}
```

## Use Cases

1. **User Input Sanitization**: Clean form data and user input
2. **Search Query Processing**: Normalize search terms for better matching
3. **Log Message Formatting**: Ensure consistent log message formatting
4. **CSV/Data Processing**: Clean imported data from external sources
5. **API Parameter Cleaning**: Sanitize API request parameters
6. **Configuration Processing**: Clean environment variables and config values
7. **Text Content Processing**: Normalize text content from various sources

## Testing

Run the tests with:

```bash
go test github.com/sampson-golang/utilities/strutil
```

See [`Squish_test.go`](./Squish_test.go) for comprehensive test cases covering edge cases and various whitespace scenarios.

## Performance Notes

- Uses `strings.TrimSpace` for efficient leading/trailing whitespace removal
- Uses compiled regex pattern `\s+` for internal whitespace replacement
- Suitable for most text processing scenarios
- For high-performance applications processing large amounts of text, consider caching the regex pattern

## Related Functions

This package works well with other string utilities:
- Use with `env.Get()` to clean environment variables
- Combine with validation libraries for input processing
- Use in logging pipelines with `output.Prettify()` for clean formatted output
