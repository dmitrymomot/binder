package binder

import "errors"

// Predefined errors
var (
	ErrInvalidMethod      = errors.New("invalid http method for binding")
	ErrInvalidContentType = errors.New("invalid content type for binding")
	ErrEmptyBody          = errors.New("empty request body")
	ErrParseForm          = errors.New("failed to parse form")
	ErrDecodeForm         = errors.New("failed to decode form")
)
