package binder

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// BindJSON binds the passed obj pointer to the request.
// It uses the JSON content type for binding.
// `obj` param should be a pointer to a struct with `json` tags.
// Implements the binder.BinderFunc interface.
func BindJSON(r *http.Request, obj interface{}) error {
	// Check if the request method is POST, PUT or PATCH
	if !isPostPutPatch(r) {
		return fmt.Errorf("%w: %s", ErrInvalidMethod, r.Method)
	}

	// Check if the request content type is JSON
	if !isJSON(r) {
		return fmt.Errorf("%w: %s", ErrInvalidContentType, r.Header.Get("Content-Type"))
	}

	// Check if the request body is empty
	if r.Body == nil {
		return ErrEmptyBody
	}

	// Decode the request body into the obj pointer
	if err := json.NewDecoder(r.Body).Decode(obj); err != nil {
		return err
	}

	return nil
}
