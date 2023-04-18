package binder

import (
	"fmt"
	"net/http"
)

// BindForm binds the passed obj pointer to the request.
// It uses the application/x-www-form-urlencoded content type for binding.
// `obj` param should be a pointer to a struct with `form“ tags.
// Implements the binder.BinderFunc interface.
func BindForm(r *http.Request, obj interface{}) error {
	// Check if the request method is POST, PUT or PATCH
	if !isPostPutPatch(r) {
		return fmt.Errorf("%w: %s", ErrInvalidMethod, r.Method)
	}

	// Check if the request content type is form urlencoded
	if !isFormURLEncoded(r) {
		return fmt.Errorf("%w: %s", ErrInvalidContentType, r.Header.Get("Content-Type"))
	}

	// Check if the request body is empty
	if r.Body == nil {
		return ErrEmptyBody
	}

	// Parse the request body
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("%w: %s", ErrParseForm, err.Error())
	}

	// Decode the request body into the obj pointer
	if err := formDecoder.Decode(obj, r.PostForm); err != nil {
		return fmt.Errorf("%w: %s", ErrDecodeForm, err.Error())
	}

	return nil
}
