# HTTPUtil Package

The `httputil` package provides utilities for HTTP server development including context management, port checking, request parsing, and response handling.

## Installation

```bash
go get github.com/sampson-golang/utilities/httputil
```

## Functions and Types

### `AppContext` and Context Management

Type-safe context management for HTTP requests with middleware support.

```go
package main

import (
  "fmt"
  "net/http"
  "github.com/sampson-golang/utilities/httputil"
)

type User struct {
  ID   int    `json:"id"`
  Name string `json:"name"`
}

func main() {
  // Create a new app context for user data
  userCtx := httputil.NewAppContext("user")

  // Middleware to add user to context
  userMiddleware := httputil.ContextMiddleware(userCtx, &User{ID: 123, Name: "John"})

  // Handler that uses context
  handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if user := userCtx.GetContext(r); user != nil {
      userObj := user.(*User)
      fmt.Fprintf(w, "Hello, %s (ID: %d)", userObj.Name, userObj.ID)
    } else {
      fmt.Fprint(w, "No user in context")
    }
  })

  // Chain middleware and handler
  http.Handle("/", userMiddleware(handler))
}
```

### `PortInUse`

Check if a network port is currently in use on both IPv4 and IPv6.

```go
func examplePortCheck() {
  port := 8080

  if httputil.PortInUse(port) {
    fmt.Printf("Port %d is already in use\n", port)
    // Try a different port
    for i := 8081; i <= 8090; i++ {
      if !httputil.PortInUse(i) {
        port = i
        fmt.Printf("Using port %d instead\n", port)
        break
      }
    }
  }

  // Start server on available port
  server := &http.Server{Addr: fmt.Sprintf(":%d", port)}
  fmt.Printf("Server starting on port %d\n", port)
  // server.ListenAndServe()
}
```

### `QueryOrBody`

Extract specified parameters from either URL query parameters or request body (JSON/form), with query parameters taking precedence.

```go
func handleAPIRequest(w http.ResponseWriter, r *http.Request) {
  // Extract specific parameters from query or body
  params, err := httputil.QueryOrBody(r, "username", "email", "token")
  if err != nil {
    http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
    return
  }

  username := params["username"]
  email := params["email"]
  token := params["token"]

  if username == "" {
    http.Error(w, "username is required", http.StatusBadRequest)
    return
  }

  fmt.Printf("Processing request for user: %s, email: %s\n", username, email)

  // Process the request...
}

// Example requests this handles:
// GET /api/user?username=john&email=john@example.com
// POST /api/user with JSON: {"username": "john", "email": "john@example.com"}
// POST /api/user with form data: username=john&email=john@example.com
```

### `UnmarshalResponse`

Parse JSON response bodies with better error messages.

```go
type APIResponse struct {
  Success bool   `json:"success"`
  Message string `json:"message"`
  Data    any    `json:"data"`
}

func handleAPIResponse(responseBody []byte) {
  var response APIResponse

  err := httputil.UnmarshalResponse(responseBody, &response)
  if err != nil {
    fmt.Printf("Failed to parse API response: %v\n", err)
    return
  }

  if response.Success {
    fmt.Printf("API call successful: %s\n", response.Message)
    fmt.Printf("Data: %+v\n", response.Data)
  } else {
    fmt.Printf("API call failed: %s\n", response.Message)
  }
}
```

## API Reference

### `AppContext` Type

A type-safe context manager for HTTP requests.

#### `NewAppContext(name string) *AppContext`

Creates a new AppContext with the given name.

**Parameters:**
- `name` - A unique identifier for this context type

**Returns:**
- `*AppContext` - A new context manager

#### `WithContext(request *http.Request, value interface{}) *http.Request`

Adds a value to the request context only if no value already exists for this context.

**Parameters:**
- `request` - The HTTP request
- `value` - The value to store in context

**Returns:**
- `*http.Request` - Request with updated context

#### `SetContext(request *http.Request, value interface{}) *http.Request`

Sets a value in the request context, overwriting any existing value.

**Parameters:**
- `request` - The HTTP request
- `value` - The value to store in context

**Returns:**
- `*http.Request` - Request with updated context

#### `GetContext(request *http.Request) interface{}`

Retrieves the value from the request context.

**Parameters:**
- `request` - The HTTP request

**Returns:**
- `interface{}` - The stored value, or `nil` if not found

### `ContextMiddleware(ctx *AppContext, value interface{}) func(next http.Handler) http.Handler`

Creates HTTP middleware that automatically adds a value to the context for all requests.

**Parameters:**
- `ctx` - The AppContext to use
- `value` - The value to add to all requests

**Returns:**
- `func(next http.Handler) http.Handler` - Standard HTTP middleware function

### `PortInUse(port int) bool`

Checks if a port is in use on both IPv4 and IPv6.

**Parameters:**
- `port` - The port number to check

**Returns:**
- `bool` - `true` if the port is in use, `false` if available

**Note:** This function uses a mutex to ensure thread-safe port checking.

### `QueryOrBody(request *http.Request, keys ...string) (map[string]string, error)`

Extracts specified parameters from either URL query or request body.

**Parameters:**
- `request` - The HTTP request
- `keys` - Parameter names to extract

**Returns:**
- `map[string]string` - Map of found parameter values
- `error` - Error if JSON parsing fails

**Behavior:**
1. First checks URL query parameters for each key
2. For GET requests, stops after query parameters
3. For other methods, also checks request body:
   - `application/json`: Parses JSON body
   - Other content types: Parses as form data
4. Query parameters take precedence over body parameters

### `UnmarshalResponse(body []byte, data interface{}) error`

Unmarshals JSON response with enhanced error messages.

**Parameters:**
- `body` - The JSON byte data to unmarshal
- `data` - Pointer to the struct to unmarshal into

**Returns:**
- `error` - Enhanced error message if unmarshaling fails, `nil` on success

## Examples

### Authentication Middleware

```go
type AuthContext struct {
  UserID   int
  Username string
  Roles    []string
}

func AuthenticationMiddleware() func(http.Handler) http.Handler {
  authCtx := httputil.NewAppContext("auth")

  return func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      // Extract token from query or body
      params, err := httputil.QueryOrBody(r, "token", "api_key")
      if err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
      }

      token := params["token"]
      if token == "" {
        token = params["api_key"]
      }

      if token == "" {
        http.Error(w, "Authentication required", http.StatusUnauthorized)
        return
      }

      // Validate token and get user info
      user := validateToken(token) // Your validation logic
      if user == nil {
        http.Error(w, "Invalid token", http.StatusUnauthorized)
        return
      }

      // Add user to context
      r = authCtx.SetContext(r, user)
      next.ServeHTTP(w, r)
    })
  }
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
  authCtx := httputil.NewAppContext("auth")
  user := authCtx.GetContext(r).(*AuthContext)

  fmt.Fprintf(w, "Welcome, %s! (User ID: %d)", user.Username, user.UserID)
}
```

### Dynamic Server Startup

```go
func StartServer(preferredPort int) error {
  port := preferredPort

  // Find an available port starting from preferred
  for port <= preferredPort+10 {
    if !httputil.PortInUse(port) {
        break
    }
    port++
  }

  if port > preferredPort+10 {
    return fmt.Errorf("no available ports in range %d-%d", preferredPort, preferredPort+10)
  }

  if port != preferredPort {
    fmt.Printf("Preferred port %d unavailable, using %d\n", preferredPort, port)
  }

  server := &http.Server{
    Addr: fmt.Sprintf(":%d", port),
    // ... other config
  }

  fmt.Printf("Server starting on http://localhost:%d\n", port)
  return server.ListenAndServe()
}
```

### API Client with Response Parsing

```go
type APIClient struct {
  BaseURL string
  Client  *http.Client
}

type APIError struct {
  Code    int    `json:"code"`
  Message string `json:"message"`
}

func (c *APIClient) makeRequest(endpoint string, result interface{}) error {
  resp, err := c.Client.Get(c.BaseURL + endpoint)
  if err != nil {
    return err
  }
  defer resp.Body.Close()

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return err
  }

  if resp.StatusCode >= 400 {
    var apiErr APIError
    if err := httputil.UnmarshalResponse(body, &apiErr); err != nil {
      return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
    }
    return fmt.Errorf("API Error %d: %s", apiErr.Code, apiErr.Message)
  }

  return httputil.UnmarshalResponse(body, result)
}
```

### Multi-Source Parameter Handler

```go
func createUserHandler(w http.ResponseWriter, r *http.Request) {
  // Extract user data from query or body
  params, err := httputil.QueryOrBody(r, "name", "email", "role", "department")
  if err != nil {
    http.Error(w, "Invalid request format: "+err.Error(), http.StatusBadRequest)
    return
  }

  // Validate required fields
  required := []string{"name", "email"}
  for _, field := range required {
    if params[field] == "" {
      http.Error(w, fmt.Sprintf("%s is required", field), http.StatusBadRequest)
      return
    }
  }

  // Set defaults
  if params["role"] == "" {
    params["role"] = "user"
  }
  if params["department"] == "" {
    params["department"] = "general"
  }

  // Create user
  user := User{
    Name:       params["name"],
    Email:      params["email"],
    Role:       params["role"],
    Department: params["department"],
  }

  // Save user...
  fmt.Printf("Created user: %+v\n", user)

  // Return success response
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(map[string]interface{}{
    "success": true,
    "user":    user,
  })
}
```

### Context-Aware Logging

```go
type RequestContext struct {
  RequestID string
  UserID    string
  StartTime time.Time
}

func RequestLoggingMiddleware() func(http.Handler) http.Handler {
  reqCtx := httputil.NewAppContext("request")

  return func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      // Create request context
      ctx := &RequestContext{
        RequestID: generateRequestID(),
        StartTime: time.Now(),
      }

      // Try to get user ID from auth context
      authCtx := httputil.NewAppContext("auth")
      if auth := authCtx.GetContext(r); auth != nil {
        ctx.UserID = auth.(*AuthContext).Username
      }

      r = reqCtx.SetContext(r, ctx)

      // Log request start
      logRequest(ctx, r)

      next.ServeHTTP(w, r)

      // Log request completion
      logRequestComplete(ctx, time.Since(ctx.StartTime))
    })
  }
}
```

## Testing

Run the tests with:

```bash
go test github.com/sampson-golang/utilities/httputil
```

See the test files for comprehensive examples:
- [`AppContext_test.go`](./AppContext_test.go)
- [`PortInUse_test.go`](./PortInUse_test.go)
- [`QueryOrBody_test.go`](./QueryOrBody_test.go)
- [`UnmarshalResponse_test.go`](./UnmarshalResponse_test.go)

## Best Practices

1. **Use AppContext for type safety**: Create specific context types rather than storing arbitrary data
2. **Check ports before binding**: Always use `PortInUse` before starting servers in production
3. **Handle both query and body**: Use `QueryOrBody` for flexible API endpoints
4. **Validate extracted parameters**: Always validate required parameters after extraction
5. **Use meaningful context names**: Choose descriptive names for your AppContext instances
6. **Chain middleware properly**: Order middleware to ensure dependencies are available
