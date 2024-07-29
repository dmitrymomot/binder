package binder_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dmitrymomot/binder"
)

func TestBindQuery(t *testing.T) {
	type User struct {
		ID     int    `query:"id"`
		Name   string `query:"name"`
		IsPaid bool   `query:"paid"`
	}

	t.Run("valid query", func(t *testing.T) {
		req, err := newQueryRequest(http.MethodGet, "/users", map[string]interface{}{
			"id":   42,
			"name": "john",
			"paid": true,
		}, nil)
		require.NoError(t, err)

		var user User
		err = binder.BindQuery(req, &user)
		require.NoError(t, err)
		require.Equal(t, 42, user.ID)
		require.Equal(t, "john", user.Name)
		require.Equal(t, true, user.IsPaid)
	})

	t.Run("invalid method", func(t *testing.T) {
		req, err := newQueryRequest(http.MethodPost, "/users", map[string]interface{}{}, nil)
		require.NoError(t, err)

		var user User
		err = binder.BindQuery(req, &user)
		require.Error(t, err)
		require.ErrorIs(t, err, binder.ErrInvalidMethod)
	})

	t.Run("empty query", func(t *testing.T) {
		req, err := newQueryRequest(http.MethodGet, "/users", map[string]interface{}{}, nil)
		require.NoError(t, err)

		var user User
		err = binder.BindQuery(req, &user)
		require.Error(t, err)
		require.ErrorIs(t, err, binder.ErrEmptyQuery)
	})

	t.Run("invalid input", func(t *testing.T) {
		req, err := newQueryRequest(http.MethodGet, "/users", map[string]interface{}{
			"id":   42,
			"name": "john",
			"paid": true,
		}, nil)
		require.NoError(t, err)

		var user User
		err = binder.BindQuery(req, user)
		require.Error(t, err)
		require.ErrorIs(t, err, binder.ErrInvalidInput)
	})

	t.Run("decode error", func(t *testing.T) {
		req, err := newQueryRequest(http.MethodGet, "/users", map[string]interface{}{
			"id":   42,
			"name": "john",
			"paid": "invalid value",
		}, nil)
		require.NoError(t, err)

		var user User
		err = binder.BindQuery(req, &user)
		require.Error(t, err)
		require.ErrorIs(t, err, binder.ErrDecodeQuery)
	})

	// Test the successful case where the one of the fields in the struct is a custom type.
	t.Run("custom type", func(t *testing.T) {
		// A custom type
		type CustomInt int
		type CustomString string

		// A test struct with query tags
		type RequestBody struct {
			FieldOne CustomString `query:"field_one"`
			FieldTwo CustomInt    `query:"field_two"`
		}

		// Initialize a dummy request payload
		payload := map[string]interface{}{
			"field_one": "value",
			"field_two": 123,
		}

		// Create a new query request with the test payload.
		req, err := newQueryRequest(http.MethodGet, "/", payload, nil)
		require.NoError(t, err)

		// Call BindQuery with the test request and a pointer to a RequestBody struct.
		var obj RequestBody
		err = binder.BindQuery(req, &obj)
		require.NoError(t, err)

		// Check that the struct was populated with the expected values.
		require.Equal(t, CustomString("value"), obj.FieldOne)
		require.Equal(t, CustomInt(123), obj.FieldTwo)
	})
}
