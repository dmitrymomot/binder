package binder

import "errors"

// Predefined errors
var (
	ErrInvalidMethod        = errors.New("invalid http method for binding")
	ErrInvalidContentType   = errors.New("invalid content type for binding")
	ErrEmptyBody            = errors.New("empty request body")
	ErrParseForm            = errors.New("failed to parse form")
	ErrDecodeForm           = errors.New("failed to decode form")
	ErrGetFile              = errors.New("failed to get file from request")
	ErrReadFile             = errors.New("failed to read file from request")
	ErrGetFileMimeType      = errors.New("failed to get file mime type")
	ErrEmptyQuery           = errors.New("empty query string")
	ErrDecodeQuery          = errors.New("failed to decode request query")
	ErrInvalidInput         = errors.New("invalid input")
	ErrUnsupportedType      = errors.New("unsupported type")
	ErrTargetMustBeAPointer = errors.New("target must be a pointer")
	ErrTargetMustBeAStruct  = errors.New("target must be a struct")
	ErrInputIsNil           = errors.New("input is nil")
	ErrDecodeJSON           = errors.New("failed to decode json")
)
