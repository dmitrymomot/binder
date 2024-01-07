package binder

import (
	"net/http"
	"strings"
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

// DefaultBinder is the default implementation of the Binder interface.
type DefaultBinder struct{}

// Bind binds the passed v pointer to the request.
// For example, the implementation could bind the request body to the v pointer.
// Bind implements the Binder interface.
// It returns an error if the binding fails.
// Binding depends on the request method and the content type.
// If the request method is GET, HEAD, DELETE, or OPTIONS, then the binding is done from the query.
// If the request method is POST, PUT, or PATCH, then the binding is done from the request body.
// If the content type is JSON, then the binding is done from the request body.
// If the content type is form, then the binding is done from the request body.
func (b *DefaultBinder) Bind(r *http.Request, v interface{}) error {
	return BindFunc(r, v)
}

// BindFunc is the function type that implements the BinderFunc interface.
// It returns an error if the binding fails.
// Binding depends on the request method and the content type.
// If the request method is GET, HEAD, DELETE, or OPTIONS, then the binding is done from the query.
// If the request method is POST, PUT, or PATCH, then the binding is done from the request body.
// If the content type is JSON, then the binding is done from the request body.
// If the content type is form, then the binding is done from the request body.
func BindFunc(r *http.Request, v interface{}) error {
	switch r.Method {
	case http.MethodGet, http.MethodHead, http.MethodDelete, http.MethodOptions:
		return BindQuery(r, v)
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		contentType := strings.ToLower(r.Header.Get("Content-Type"))
		switch {
		case strings.HasPrefix(contentType, "application/json"):
			return BindJSON(r, v)
		case strings.HasPrefix(contentType, "application/x-www-form-urlencoded"):
			return BindForm(r, v)
		case strings.HasPrefix(contentType, "multipart/form-data"):
			return BindFormMultipart(r, v)
		default:
			return ErrInvalidContentType
		}
	default:
		return ErrInvalidMethod
	}
}
