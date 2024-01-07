package binder_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"github.com/dmitrymomot/binder"
	"github.com/stretchr/testify/require"
)

func TestGetFileData(t *testing.T) {
	// Create a test multipart form request.
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.txt")
	require.NoError(t, err)

	_, err = part.Write([]byte("Test file data"))
	require.NoError(t, err)

	err = writer.Close()
	require.NoError(t, err)

	// Create a test HTTP request from the multipart form.
	req, err := http.NewRequest(http.MethodPost, "/upload", body)
	require.NoError(t, err)

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Call GetFileData with the test request and field name.
	fileData, err := binder.GetFileData(req, "file")
	require.NoError(t, err)

	// Check that the returned file data matches the expected values.
	require.Equal(t, "test.txt", fileData.Name)
	require.Equal(t, int64(len("Test file data")), fileData.Size)
	require.Equal(t, "text/plain", fileData.MimeType)
	require.Equal(t, []byte("Test file data"), fileData.Data)
}

// TestGetFileData_Image tests the GetFileData function with an image file.
// Unlike the TestGetFileData test, this test reads the file data from a file testdata/test.jpg
// and writes it to the multipart form.
// This test is useful for testing the MIME type detection.
func TestGetFileData_Image(t *testing.T) {
	// Read the test image file.
	testImage, err := os.ReadFile("testdata/test.jpg")
	require.NoError(t, err)

	// Create a test multipart form request.
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.jpg")
	require.NoError(t, err)

	_, err = part.Write(testImage)
	require.NoError(t, err)

	err = writer.Close()
	require.NoError(t, err)

	// Create a test HTTP request from the multipart form.
	req, err := http.NewRequest(http.MethodPost, "/upload", body)
	require.NoError(t, err)

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Call GetFileData with the test request and field name.
	fileData, err := binder.GetFileData(req, "file")
	require.NoError(t, err)

	// Check that the returned file data matches the expected values.
	require.Equal(t, "test.jpg", fileData.Name)
	require.Equal(t, int64(len(testImage)), fileData.Size)
	require.Equal(t, "image/jpeg", fileData.MimeType)
	require.Equal(t, testImage, fileData.Data)
}
