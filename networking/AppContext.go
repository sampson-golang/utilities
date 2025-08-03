package networking

import (
	"context"
	"net/http"
)

type AppContext struct {
	token *ContextToken
}

func (c *AppContext) WithContext(request *http.Request, value interface{}) *http.Request {
	existing := request.Context().Value(c.token)
	if existing == nil {
		return request.WithContext(context.WithValue(request.Context(), c.token, value))
	}
	return request
}

func (c *AppContext) GetContext(request *http.Request) interface{} {
	return request.Context().Value(c.token)
}

func (c *AppContext) SetContext(request *http.Request, value interface{}) *http.Request {
	return request.WithContext(context.WithValue(request.Context(), c.token, value))
}

type ContextToken struct {
	name string
}

func (c *ContextToken) String() string {
	return c.name
}

func NewContextToken(name string) *ContextToken {
	return &ContextToken{name: name}
}

func NewAppContext(name string) *AppContext {
	return &AppContext{token: NewContextToken(name)}
}
