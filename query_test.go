package binder_test

import (
	"errors"
	"net/http"
	"testing"

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
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		var user User
		if err := binder.BindQuery(req, &user); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if user.ID != 42 || user.Name != "john" || !user.IsPaid {
			t.Errorf("unexpected result: %+v", user)
		}
	})

	t.Run("invalid method", func(t *testing.T) {
		req, err := newQueryRequest(http.MethodPost, "/users", map[string]interface{}{}, nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		var user User
		if err := binder.BindQuery(req, &user); !errors.Is(err, binder.ErrInvalidMethod) {
			t.Errorf("expected %v, got %v", binder.ErrInvalidMethod, err)
		}
	})

	t.Run("empty query", func(t *testing.T) {
		req, err := newQueryRequest(http.MethodGet, "/users", map[string]interface{}{}, nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		var user User
		if err := binder.BindQuery(req, &user); !errors.Is(err, binder.ErrEmptyQuery) {
			t.Errorf("expected %v, got %v", binder.ErrEmptyQuery, err)
		}
	})

	t.Run("invalid input", func(t *testing.T) {
		req, err := newQueryRequest(http.MethodGet, "/users", map[string]interface{}{
			"id":   42,
			"name": "john",
			"paid": true,
		}, nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		var user User
		if err := binder.BindQuery(req, user); !errors.Is(err, binder.ErrInvalidInput) {
			t.Errorf("expected %v, got %v", binder.ErrInvalidInput, err)
		}
	})

	t.Run("decode error", func(t *testing.T) {
		req, err := newQueryRequest(http.MethodGet, "/users", map[string]interface{}{
			"id":   42,
			"name": "john",
			"paid": "invalid value",
		}, nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		var user User
		if err := binder.BindQuery(req, &user); !errors.Is(err, binder.ErrDecodeQuery) {
			t.Errorf("expected %v, got %v", binder.ErrDecodeQuery, err)
		}
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
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		// Call BindQuery with the test request and a pointer to a RequestBody struct.
		var obj RequestBody
		err = binder.BindQuery(req, &obj)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		// Check that the struct was populated with the expected values.
		if obj.FieldOne != CustomString("value") {
			t.Errorf("FieldOne is %q, expected %q", obj.FieldOne, "value")
		}
		if obj.FieldTwo != CustomInt(123) {
			t.Errorf("FieldTwo is %d, expected %d", obj.FieldTwo, 123)
		}
	})
}
