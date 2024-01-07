package binder

import (
	"errors"
	"fmt"
	"net/http"
)

// BindQuery binds the passed v pointer to the request.
// It uses the query string for binding.
// `v` param should be a pointer to a struct with `query“ tags.
// Implements the binder.BinderFunc interface.
func BindQuery(r *http.Request, v interface{}) error {
	// Check if the request method is GET, HEAD or DELETE
	if !isGetHeadOptionDelete(r) {
		return fmt.Errorf("%w: %s", ErrInvalidMethod, r.Method)
	}

	// Validate v pointer before decoding query into it
	if !isPointer(v) {
		return errors.Join(ErrInvalidInput, ErrTargetMustBeAPointer)
	}

	// Check if the request query is empty
	if r.URL.RawQuery == "" {
		return ErrEmptyQuery
	}

	// Decode the request query into the v pointer and handle decoding errors
	if err := queryDecoder.Decode(v, r.URL.Query()); err != nil {
		return errors.Join(ErrDecodeQuery, err)
	}

	return nil
}
