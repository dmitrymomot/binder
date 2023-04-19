package binder

import (
	"net/http"
)

// Binder is the interface that wraps the Bind method.
//
// Bind should bind the passed v pointer to the request.
// For example, the implementation could bind the request body to the v pointer.
type Binder interface {
	Bind(r *http.Request, v interface{}) error
}

// BinderFunc is the function type that implements the Binder interface.
type BinderFunc func(*http.Request, interface{}) error
