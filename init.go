package binder

import "github.com/gorilla/schema"

// Default form decoder for binding form data
// It uses the gorilla/schema package.
// it caches meta-data about structs, and an instance can be shared safely.
var formDecoder = schema.NewDecoder()

// Default query decoder for binding query data
// It uses the gorilla/schema package.
// it caches meta-data about structs, and an instance can be shared safely.
var queryDecoder = schema.NewDecoder()

// Initialize the form & query decoders.
// It ignores unknown keys and sets zero values for empty fields.
func init() {
	formDecoder.IgnoreUnknownKeys(true)
	formDecoder.ZeroEmpty(true)
	formDecoder.SetAliasTag(TagForm)
	formDecoder.RegisterConverter(FileData{}, FileDataConverter)
	formDecoder.RegisterConverter(&FileData{}, FileDataConverterPtr)

	queryDecoder.IgnoreUnknownKeys(true)
	queryDecoder.ZeroEmpty(true)
	queryDecoder.SetAliasTag(TagQuery)
}
