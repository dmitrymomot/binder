package binder

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

// FileData represents the data contained in a file,
// including metadata and the file contents.
type FileData struct {
	// Name is the filename.
	Name string
	// Size is the size of the file in bytes.
	Size int64
	// MimeType is the MIME type of the file, as determined by its contents or extension.
	MimeType string
	// Data contains the raw bytes of the file.
	Data []byte
}

// GetFileData extracts file data from a multipart request.
func GetFileData(req *http.Request, fieldName string) (*FileData, error) {
	// Parse the multipart form data.
	err := req.ParseMultipartForm(32 << 20) // MaxMemory is 32MB
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrParseForm, err.Error())
	}

	// Get the file from the request.
	file, fileHeader, err := req.FormFile(fieldName)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrGetFile, err.Error())
	}
	defer file.Close()

	// Read the file data into memory.
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrReadFile, err.Error())
	}

	// Get the MIME type of the file.
	mime, err := GetFileMimeType(data)
	if err != nil {
		return nil, err
	}

	// Return the file data.
	return &FileData{
		Name:     fileHeader.Filename,
		Size:     fileHeader.Size,
		MimeType: mime,
		Data:     data,
	}, nil
}

// GetFileMimeType returns the MIME type of a file using the
// github.com/gabriel-vasile/mimetype package.
func GetFileMimeType(input []byte) (string, error) {
	if input == nil {
		return "", fmt.Errorf("%w: input is nil", ErrGetFileMimeType)
	}

	mtype := mimetype.Detect(input)
	if mtype == nil {
		return "", fmt.Errorf("%w: unknown MIME type", ErrGetFileMimeType)
	}

	parts := strings.Split(mtype.String(), ";")
	return parts[0], nil
}

// GetFileDataFromMultipartFileHeader extracts file data from a multipart file header.
func GetFileDataFromMultipartFileHeader(header *multipart.FileHeader) (*FileData, error) {
	// Open the file.
	file, err := header.Open()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrGetFile, err.Error())
	}
	defer file.Close()

	// Read the file data into memory.
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrReadFile, err.Error())
	}

	// Get the MIME type of the file.
	mime, err := GetFileMimeType(data)
	if err != nil {
		return nil, err
	}

	// Return the file data.
	return &FileData{
		Name:     header.Filename,
		Size:     header.Size,
		MimeType: mime,
		Data:     data,
	}, nil
}
