package binder

import (
	"fmt"
	"net/http"
	"reflect"
)

// BindQuery binds the passed obj pointer to the request.
// It uses the query string for binding.
// `obj` param should be a pointer to a struct with `query“ tags.
// Implements the binder.BinderFunc interface.
func BindQuery(r *http.Request, obj interface{}) error {
	// Check if the request method is GET, HEAD or DELETE
	switch r.Method {
	case http.MethodGet, http.MethodHead, http.MethodDelete:
	default:
		return fmt.Errorf("%w: %s", ErrInvalidMethod, r.Method)
	}

	// Check if the request query is empty
	if r.URL.RawQuery == "" {
		return ErrEmptyQuery
	}

	// Validate obj pointer before decoding query into it
	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("%w: obj must be a pointer to a struct", ErrInvalidInput)
	}

	// Decode the request query into the obj pointer and handle decoding errors
	if err := queryDecoder.Decode(obj, r.URL.Query()); err != nil {
		return fmt.Errorf("%w: %s", ErrDecodeQuery, err.Error())
	}

	return nil
}
