package binder_test

import (
	"errors"
	"net/http"
	"testing"

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
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		if err := binder.BindJSON(req, nil); err == nil || !errors.Is(err, binder.ErrInvalidMethod) {
			t.Errorf("BindJSON() error = %v, wantErr %v", err, binder.ErrInvalidMethod)
		}
	})

	// Test the scenario where the request content type is not JSON
	t.Run("invalid content type", func(t *testing.T) {
		req, err := newFormRequest(http.MethodPost, "/", payload, map[string]string{
			"Content-Type": "application/xml",
		})
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		if err := binder.BindJSON(req, nil); err == nil || !errors.Is(err, binder.ErrInvalidContentType) {
			t.Errorf("BindJSON() error = %v, wantErr %v", err, binder.ErrInvalidContentType)
		}
	})

	// Test the scenario where the request body is empty
	t.Run("invalid input", func(t *testing.T) {
		req, err := newJSONRequest(http.MethodPost, "/", nil, map[string]string{
			"Content-Type": "application/json",
		})
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		var invalidValue interface{}
		if err := binder.BindJSON(req, &invalidValue); err == nil || !errors.Is(err, binder.ErrInvalidInput) {
			t.Errorf("BindJSON() error = %v, wantErr %v", err, binder.ErrInvalidInput)
		}
	})

	// Test the scenario where the request body is empty
	t.Run("empty body", func(t *testing.T) {
		req, err := newJSONRequest(http.MethodPost, "/", nil, map[string]string{
			"Content-Type": "application/json",
		})
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		if err := binder.BindJSON(req, &RequestBody{}); err == nil || !errors.Is(err, binder.ErrEmptyBody) {
			t.Errorf("BindJSON() error = %v, wantErr %v", err, binder.ErrEmptyBody)
		}
	})

	// Test the scenario where the decoding of the request body into the obj pointer fails
	// by passing a non-struct object as parameter
	t.Run("invalid obj", func(t *testing.T) {
		req, err := newJSONRequest(http.MethodPost, "/", payload, map[string]string{
			"Content-Type": "application/json",
		})
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		if err := binder.BindJSON(req, "invalid object"); err == nil {
			t.Error("BindJSON() error = nil, wantErr non-nil")
		}
	})

	// Test the successful case where the request method is POST, PUT or PATCH and the content type is JSON.
	t.Run("successful case", func(t *testing.T) {
		// Create a new request with the JSON payload as the body
		req, err := newJSONRequest(http.MethodPost, "/", payload, map[string]string{
			"Content-Type": "application/json",
		})
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		// Create an empty struct to be used as the obj parameter in the BindJSON call
		obj := &RequestBody{}

		// Call the BindJSON function to populate the obj struct with the request data
		if err := binder.BindJSON(req, obj); err != nil {
			t.Errorf("BindJSON() error = %v, wantErr nil", err)
		}

		// Check that the obj struct was populated correctly
		if obj.FieldOne != "value" || obj.FieldTwo != 123 {
			t.Errorf("BindJSON() obj = %v, want %v", obj, payload)
		}
	})
}
