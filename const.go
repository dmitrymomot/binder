package binder

// Default tag names for binding
const (
	// TagForm Form struct tag name for binding
	TagForm = "form"
	// TagQuery Query struct tag name for binding
	TagQuery = "query"
)

// MultiPartFormMaxMemory is the maximum amount of memory to use when parsing a multipart form.
// It is passed to http.Request.ParseMultipartForm.
// Default value is 32 << 20 (32 MB).
var MultiPartFormMaxMemory int64 = 32 << 20
