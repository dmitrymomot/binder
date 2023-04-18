package binder_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"github.com/dmitrymomot/binder"
)

func TestGetFileData(t *testing.T) {
	// Create a test multipart form request.
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.txt")
	if err != nil {
		t.Fatalf("Failed to create form file: %s", err.Error())
	}
	_, err = part.Write([]byte("Test file data"))
	if err != nil {
		t.Fatalf("Failed to write form file data: %s", err.Error())
	}
	err = writer.Close()
	if err != nil {
		t.Fatalf("Failed to close multipart writer: %s", err.Error())
	}

	// Create a test HTTP request from the multipart form.
	req, err := http.NewRequest(http.MethodPost, "/upload", body)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %s", err.Error())
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Call GetFileData with the test request and field name.
	fileData, err := binder.GetFileData(req, "file")
	if err != nil {
		t.Fatalf("Unexpected error from GetFileData: %s", err.Error())
	}

	// Check that the returned file data matches the expected values.
	if fileData.Name != "test.txt" {
		t.Errorf("Name is %q, expected %q", fileData.Name, "test.txt")
	}
	if fileData.Size != int64(len("Test file data")) {
		t.Errorf("Size is %d, expected %d", fileData.Size, len("Test file data"))
	}
	if fileData.MimeType != "text/plain" {
		t.Errorf("MimeType is %q, expected %q", fileData.MimeType, "text/plain")
	}
	if string(fileData.Data) != "Test file data" {
		t.Errorf("Data is %q, expected %q", string(fileData.Data), "Test file data")
	}
}

// TestGetFileData_Image tests the GetFileData function with an image file.
// Unlike the TestGetFileData test, this test reads the file data from a file testdata/test.jpg
// and writes it to the multipart form.
// This test is useful for testing the MIME type detection.
func TestGetFileData_Image(t *testing.T) {
	// Read the test image file.
	testImage, err := os.ReadFile("testdata/test.jpg")
	if err != nil {
		t.Fatalf("Failed to read test image file: %s", err.Error())
	}

	// Create a test multipart form request.
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.jpg")
	if err != nil {
		t.Fatalf("Failed to create form file: %s", err.Error())
	}
	_, err = part.Write(testImage)
	if err != nil {
		t.Fatalf("Failed to write form file data: %s", err.Error())
	}
	err = writer.Close()
	if err != nil {
		t.Fatalf("Failed to close multipart writer: %s", err.Error())
	}

	// Create a test HTTP request from the multipart form.
	req, err := http.NewRequest(http.MethodPost, "/upload", body)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %s", err.Error())
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Call GetFileData with the test request and field name.
	fileData, err := binder.GetFileData(req, "file")
	if err != nil {
		t.Fatalf("Unexpected error from GetFileData: %s", err.Error())
	}

	// Check that the returned file data matches the expected values.
	if fileData.Name != "test.jpg" {
		t.Errorf("Name is %q, expected %q", fileData.Name, "test.jpg")
	}
	if fileData.Size != int64(len(testImage)) {
		t.Errorf("Size is %d, expected %d", fileData.Size, len(testImage))
	}
	if fileData.MimeType != "image/jpeg" {
		t.Errorf("MimeType is %q, expected %q", fileData.MimeType, "image/jpeg")
	}
	if !bytes.Equal(fileData.Data, testImage) {
		t.Errorf("Data is %q, expected %q", string(fileData.Data), "Test file data")
	}
}
