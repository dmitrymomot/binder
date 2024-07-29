package binder_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dmitrymomot/binder"
)

func TestBindJSON(t *testing.T) {
	// A test struct with json tags
	type RequestBody struct {
		FieldOne string `json:"field_one"`
		FieldTwo int    `json:"field_two"`
	}

	// Initialize a dummy request payload
	payload := map[string]interface{}{
		"field_one": "value",
		"field_two": 123,
	}

	// Test the scenario where the request method is not POST, PUT or PATCH
	t.Run("invalid method", func(t *testing.T) {
		req, err := newQueryRequest(http.MethodGet, "/", payload, nil)
		require.NoError(t, err)

		err = binder.BindJSON(req, nil)
		require.Error(t, err)
		require.ErrorIs(t, err, binder.ErrInvalidMethod)
	})

	// Test the scenario where the request content type is not JSON
	t.Run("invalid content type", func(t *testing.T) {
		req, err := newFormRequest(http.MethodPost, "/", payload, map[string]string{
			"Content-Type": "application/xml",
		})
		require.NoError(t, err)

		err = binder.BindJSON(req, nil)
		require.Error(t, err)
		require.ErrorIs(t, err, binder.ErrInvalidContentType)
	})

	// Test the scenario where the request body is empty
	t.Run("invalid input", func(t *testing.T) {
		req, err := newJSONRequest(http.MethodPost, "/", nil, map[string]string{
			"Content-Type": "application/json",
		})
		require.NoError(t, err)

		var invalidValue interface{}
		err = binder.BindJSON(req, &invalidValue)
		require.Error(t, err)
		require.ErrorIs(t, err, binder.ErrInvalidInput)
	})

	// Test the scenario where the request body is empty
	t.Run("empty body", func(t *testing.T) {
		req, err := newJSONRequest(http.MethodPost, "/", nil, map[string]string{
			"Content-Type": "application/json",
		})
		require.NoError(t, err)

		err = binder.BindJSON(req, &RequestBody{})
		require.Error(t, err)
		require.ErrorIs(t, err, binder.ErrEmptyBody)
	})

	// Test the scenario where the decoding of the request body into the obj pointer fails
	// by passing a non-struct object as parameter
	t.Run("invalid obj", func(t *testing.T) {
		req, err := newJSONRequest(http.MethodPost, "/", payload, map[string]string{
			"Content-Type": "application/json",
		})
		require.NoError(t, err)

		err = binder.BindJSON(req, "invalid object")
		require.Error(t, err)
		require.ErrorIs(t, err, binder.ErrInvalidInput)
	})

	// Test the successful case where the request method is POST, PUT or PATCH and the content type is JSON.
	t.Run("successful case", func(t *testing.T) {
		// Create a new request with the JSON payload as the body
		req, err := newJSONRequest(http.MethodPost, "/", payload, map[string]string{
			"Content-Type": "application/json",
		})
		require.NoError(t, err)

		// Create an empty struct to be used as the obj parameter in the BindJSON call
		obj := &RequestBody{}

		// Call the BindJSON function to populate the obj struct with the request data
		err = binder.BindJSON(req, obj)
		require.NoError(t, err)

		// Check that the obj struct was populated correctly
		require.Equal(t, "value", obj.FieldOne)
		require.Equal(t, 123, obj.FieldTwo)
	})

	// Test the successful case where the one of the fields in the struct is a custom type.
	t.Run("custom type", func(t *testing.T) {
		// A custom type
		type CustomInt int
		type CustomString string

		// A test struct with json tags
		type RequestBody struct {
			FieldOne CustomString `json:"field_one"`
			FieldTwo CustomInt    `json:"field_two"`
		}

		// Create a new request with the JSON payload as the body
		req, err := newJSONRequest(http.MethodPost, "/", payload, map[string]string{
			"Content-Type": "application/json",
		})
		require.NoError(t, err)

		// Create an empty struct to be used as the obj parameter in the BindJSON call
		obj := &RequestBody{}

		// Call the BindJSON function to populate the obj struct with the request data
		err = binder.BindJSON(req, obj)
		require.NoError(t, err)

		// Check that the obj struct was populated correctly
		require.Equal(t, CustomString("value"), obj.FieldOne)
		require.Equal(t, CustomInt(123), obj.FieldTwo)
	})
}
