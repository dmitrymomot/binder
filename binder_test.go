package binder_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dmitrymomot/binder"
)

func TestBindFunc(t *testing.T) {
	// Test GET request
	t.Run("GET", func(t *testing.T) {
		getReq, err := http.NewRequest(http.MethodGet, "/data?id=123", nil)
		require.NoError(t, err)
		var getParams struct {
			ID string `json:"id"`
		}
		err = binder.BindFunc(getReq, &getParams)
		require.NoError(t, err)
		require.Equal(t, "123", getParams.ID)
	})

	// Test POST request with JSON content type
	t.Run("POST JSON", func(t *testing.T) {
		jsonReq, err := http.NewRequest(http.MethodPost, "/data", strings.NewReader(`{"name":"John"}`))
		jsonReq.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)
		var jsonBody struct {
			Name string `json:"name"`
		}
		err = binder.BindFunc(jsonReq, &jsonBody)
		require.NoError(t, err)
		require.Equal(t, "John", jsonBody.Name)
	})

	// Test POST request with form-urlencoded content type
	t.Run("POST form-urlencoded", func(t *testing.T) {
		formReq, err := http.NewRequest(http.MethodPost, "/data", strings.NewReader("name=John"))
		formReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		require.NoError(t, err)
		var formBody struct {
			Name string `form:"name"`
		}
		err = binder.BindFunc(formReq, &formBody)
		require.NoError(t, err)
		require.Equal(t, "John", formBody.Name)
	})

	// Test POST request with multipart/form-data content type
	t.Run("POST multipart/form-data", func(t *testing.T) {
		// Create a test multipart form request.
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "test.txt")
		require.NoError(t, err)

		_, err = part.Write([]byte("Test file data"))
		require.NoError(t, err)

		err = writer.Close()
		require.NoError(t, err)

		multipartReq, err := http.NewRequest(http.MethodPost, "/data", body)
		require.NoError(t, err)
		multipartReq.Header.Set("Content-Type", writer.FormDataContentType())

		var multipartBody struct {
			File binder.FileData `form:"file"`
		}
		err = binder.BindFunc(multipartReq, &multipartBody)
		require.NoError(t, err)
		// Add assertions for multipart form data binding
	})

	// Test invalid content type
	t.Run("Invalid content type", func(t *testing.T) {
		invalidContentTypeReq, err := http.NewRequest(http.MethodPost, "/data", nil)
		require.NoError(t, err)
		invalidContentTypeReq.Header.Set("Content-Type", "text/plain")
		var invalidContentTypeBody struct{}
		err = binder.BindFunc(invalidContentTypeReq, &invalidContentTypeBody)
		require.ErrorIs(t, err, binder.ErrInvalidContentType)
	})
}
