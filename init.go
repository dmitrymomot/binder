package binder

import "github.com/gorilla/schema"

// Default form decoder for binding form data
// It uses the gorilla/schema package.
// it caches meta-data about structs, and an instance can be shared safely.
var formDecoder = schema.NewDecoder()

// Initialize the form decoder.
// It ignores unknown keys and sets zero values for empty fields.
func init() {
	formDecoder.IgnoreUnknownKeys(true)
	formDecoder.ZeroEmpty(true)
	formDecoder.SetAliasTag(TagForm)
}
