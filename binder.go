package binder

import (
	"net/http"

	"github.com/gorilla/schema"
)

var (
	// JSON struct tag name for binding
	TagJSON = "json"
	// Form struct tag name for binding
	TagForm = "form"
	// Query struct tag name for binding
	TagQuery = "query"
	// Header struct tag name for binding
	TagHeader = "header"
)

// Default form decoder for binding form data
// It uses the gorilla/schema package.
// it caches meta-data about structs, and an instance can be shared safely.
var formDecoder = schema.NewDecoder()

// Binder is the interface that wraps the Bind method.
//
// Bind should bind the passed obj pointer to the request.
// For example, the implementation could bind the request body to the obj pointer.
type Binder interface {
	Bind(*http.Request, interface{}) error
}

// BinderFunc is the function type that implements the Binder interface.
type BinderFunc func(*http.Request, interface{}) error

// Initialize the form decoder.
// It ignores unknown keys and sets zero values for empty fields.
func init() {
	formDecoder.IgnoreUnknownKeys(true)
	formDecoder.ZeroEmpty(true)
	formDecoder.SetAliasTag(TagForm)
}
