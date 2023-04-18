package binder

import (
	"net/http"
)

// Binder is the interface that wraps the Bind method.
//
// Bind should bind the passed obj pointer to the request.
// For example, the implementation could bind the request body to the obj pointer.
type Binder interface {
	Bind(*http.Request, interface{}) error
}

// BinderFunc is the function type that implements the Binder interface.
type BinderFunc func(*http.Request, interface{}) error
