package binder_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/dmitrymomot/binder"
)

func TestBindForm(t *testing.T) {
	type objPayload struct {
		FieldOne string `form:"field_one"`
		FieldTwo int    `form:"field_two"`
	}

	// Initialize a dummy request payload
	payload := map[string]interface{}{
		"field_one": "value",
		"field_two": 123,
	}

	t.Run("invalid method", func(t *testing.T) {
		req, err := newQueryRequest(http.MethodGet, "/path", nil, nil)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		var obj objPayload
		err = binder.BindForm(req, &obj)

		if !errors.Is(err, binder.ErrInvalidMethod) {
			t.Errorf("expected error '%s', got '%s'", binder.ErrInvalidMethod, err)
		}
	})

	t.Run("invalid content type", func(t *testing.T) {
		req, err := newJSONRequest(http.MethodPost, "/path", payload, nil)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		var obj objPayload
		err = binder.BindForm(req, &obj)

		if !errors.Is(err, binder.ErrInvalidContentType) {
			t.Errorf("expected error '%s', got '%s'", binder.ErrInvalidContentType, err)
		}
	})

	t.Run("empty body", func(t *testing.T) {
		req, err := newFormRequest(http.MethodPost, "/path", nil, nil)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		var obj objPayload
		err = binder.BindForm(req, &obj)
		if !errors.Is(err, binder.ErrEmptyBody) {
			t.Errorf("expected error '%s', got '%s'", binder.ErrEmptyBody, err)
		}
	})

	t.Run("decode form error", func(t *testing.T) {
		req, err := newFormRequest(http.MethodPost, "/path", payload, nil)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		var obj objPayload
		err = binder.BindForm(req, obj)
		if !errors.Is(err, binder.ErrDecodeForm) {
			t.Errorf("expected error '%s', got '%s'", binder.ErrDecodeForm, err)
		}
	})

	t.Run("success", func(t *testing.T) {
		req, err := newFormRequest(http.MethodPost, "/path", payload, nil)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		var obj objPayload
		err = binder.BindForm(req, &obj)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		if obj.FieldOne != "value" {
			t.Errorf("expected 'value', got '%s'", obj.FieldOne)
		}
		if obj.FieldTwo != 123 {
			t.Errorf("expected 123, got %d", obj.FieldTwo)
		}
	})
}
