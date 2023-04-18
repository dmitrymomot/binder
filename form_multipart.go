package binder

import (
	"fmt"
	"net/http"
	"strings"
)

// BindFormMultipart binds the passed obj pointer to the request.
// It uses the multipart/form-data content type for binding.
// `obj` param should be a pointer to a struct with `form“ tags.
// Implements the binder.BinderFunc interface.
func BindFormMultipart(r *http.Request, obj interface{}) error {
	// Check if the request method is POST, PUT or PATCH
	if r.Method != http.MethodPost && r.Method != http.MethodPut && r.Method != http.MethodPatch {
		return fmt.Errorf("%w: %s", ErrInvalidMethod, r.Method)
	}

	// Check if the request content type is form urlencoded
	if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		return fmt.Errorf("%w: %s", ErrInvalidContentType, r.Header.Get("Content-Type"))
	}

	// Check if the request body is empty
	if r.Body == nil {
		return ErrEmptyBody
	}

	// Parse the request body
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		return fmt.Errorf("%w: %s", ErrParseForm, err.Error())
	}

	// Decode the request body into the obj pointer
	if err := formDecoder.Decode(obj, r.PostForm); err != nil {
		return fmt.Errorf("%w: %s", ErrDecodeForm, err.Error())
	}

	return nil
}
