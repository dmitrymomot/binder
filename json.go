package binder

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// BindJSON binds the passed v pointer to the request.
// It uses the JSON content type for binding.
// `v` param should be a pointer to a struct with `json` tags.
// Implements the binder.BinderFunc interface.
func BindJSON(r *http.Request, v interface{}) error {
	// Check if the request method is POST, PUT or PATCH
	if !isPostPutPatch(r) {
		return fmt.Errorf("%w: %s", ErrInvalidMethod, r.Method)
	}

	// Check if the request content type is JSON
	if !isJSON(r) {
		return fmt.Errorf("%w: %s", ErrInvalidContentType, r.Header.Get("Content-Type"))
	}

	// Validate v pointer before decoding query into it
	if !isPointer(v) {
		return fmt.Errorf("%w: v must be a pointer to a struct", ErrInvalidInput)
	}

	// Check if the request body is empty
	if r.Body == nil {
		return ErrEmptyBody
	}

	// Decode the request body into the v pointer
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}

	return nil
}
