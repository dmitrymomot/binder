package binder

// Default tag names for binding
const (
	// JSON struct tag name for binding
	TagJSON = "json"
	// Form struct tag name for binding
	TagForm = "form"
	// Query struct tag name for binding
	TagQuery = "query"
	// Header struct tag name for binding
	TagHeader = "header"
)

// MultiPartFormMaxMemory is the maximum amount of memory to use when parsing a multipart form.
// It is passed to http.Request.ParseMultipartForm.
// Default value is 32 << 20 (32 MB).
const MultiPartFormMaxMemory int64 = 32 << 20
