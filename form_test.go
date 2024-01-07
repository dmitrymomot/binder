package binder_test

import (
	"net/http"
	"testing"

	"github.com/dmitrymomot/binder"
	"github.com/stretchr/testify/require"
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
		require.NoError(t, err)

		var obj objPayload
		err = binder.BindForm(req, &obj)
		require.Error(t, err)
		require.ErrorIs(t, err, binder.ErrInvalidMethod)
	})

	t.Run("invalid content type", func(t *testing.T) {
		req, err := newJSONRequest(http.MethodPost, "/path", payload, nil)
		require.NoError(t, err)

		var obj objPayload
		err = binder.BindForm(req, &obj)
		require.Error(t, err)
		require.ErrorIs(t, err, binder.ErrInvalidContentType)
	})

	t.Run("empty body", func(t *testing.T) {
		req, err := newFormRequest(http.MethodPost, "/path", nil, nil)
		require.NoError(t, err)

		var obj objPayload
		err = binder.BindForm(req, &obj)
		require.Error(t, err)
		require.ErrorIs(t, err, binder.ErrEmptyBody)
	})

	t.Run("invalid input", func(t *testing.T) {
		req, err := newFormRequest(http.MethodPost, "/path", payload, nil)
		require.NoError(t, err)

		var obj objPayload
		err = binder.BindForm(req, obj)
		require.Error(t, err)
		require.ErrorIs(t, err, binder.ErrInvalidInput)
	})

	t.Run("success", func(t *testing.T) {
		req, err := newFormRequest(http.MethodPost, "/path", payload, nil)
		require.NoError(t, err)

		var obj objPayload
		err = binder.BindForm(req, &obj)
		require.NoError(t, err)
		require.Equal(t, "value", obj.FieldOne)
		require.Equal(t, 123, obj.FieldTwo)
	})

	t.Run("success with empty payload", func(t *testing.T) {
		req, err := newFormRequest(http.MethodPost, "/path", map[string]interface{}{}, nil)
		require.NoError(t, err)

		var obj objPayload
		err = binder.BindForm(req, &obj)
		require.NoError(t, err)
		require.Equal(t, "", obj.FieldOne)
		require.Equal(t, 0, obj.FieldTwo)
	})
}
