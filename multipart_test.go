package binder_test

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/dmitrymomot/binder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type FormData struct {
	Name   string       `form:"name,omitempty"`
	Email  string       `form:"email,omitempty"`
	Age    int          `form:"age,omitempty"`
	Avatar *binder.File `form:"avatar,omitempty"`
}

func TestBindFormMultipart(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		setupForm   func(*multipart.Writer) error
		expectErr   bool
	}{
		{
			name:        "ValidFormData",
			contentType: "multipart/form-data",
			setupForm: func(w *multipart.Writer) error {
				if err := w.WriteField("name", "John Doe"); err != nil {
					return err
				}
				if err := w.WriteField("email", "john.doe@example.com"); err != nil {
					return err
				}
				if err := w.WriteField("age", "30"); err != nil {
					return err
				}

				filePath := "testdata/test.jpg"
				file, err := os.Open(filePath)
				if err != nil {
					return err
				}
				defer file.Close()

				part, err := w.CreateFormFile("avatar", filepath.Base(filePath))
				if err != nil {
					return err
				}
				if _, err := io.Copy(part, file); err != nil {
					return err
				}

				return nil
			},
			expectErr: false,
		},
		{
			name:        "InvalidContentType",
			contentType: "application/x-www-form-urlencoded",
			setupForm:   func(w *multipart.Writer) error { return nil },
			expectErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			multipartWriter := multipart.NewWriter(&buf)

			err := test.setupForm(multipartWriter)
			require.NoError(t, err)

			err = multipartWriter.Close()
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/upload", &buf)
			req.Header.Set("Content-Type", test.contentType+"; boundary="+multipartWriter.Boundary())

			var data FormData
			err = binder.BindFormMultipart(req, &data)

			if test.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "John Doe", data.Name)
				assert.Equal(t, "john.doe@example.com", data.Email)
				assert.Equal(t, 30, data.Age)
				if assert.NotNil(t, data.Avatar) {
					assert.Equal(t, "test.jpg", data.Avatar.FileName)
					assert.NotEmpty(t, data.Avatar.Data)
					assert.Equal(t, "image/jpeg", data.Avatar.ContentType)
					assert.Greater(t, data.Avatar.FileSize, int64(0))
				}
			}
		})
	}
}
