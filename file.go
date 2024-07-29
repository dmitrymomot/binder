package binder

import (
	"errors"
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
		return nil, errors.Join(ErrParseForm, err)
	}

	// Get the file from the request.
	file, fileHeader, err := req.FormFile(fieldName)
	if err != nil {
		return nil, errors.Join(ErrGetFile, err)
	}
	defer func(file multipart.File) {
		_ = file.Close()
	}(file)

	// Read the file data into memory.
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.Join(ErrReadFile, err)
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

// GetFileDataFromMultipartFileHeader extracts file data from a multipart file header.
//
//goland:noinspection GoUnusedExportedFunction
func GetFileDataFromMultipartFileHeader(header *multipart.FileHeader) (*FileData, error) {
	// Open the file.
	file, err := header.Open()
	if err != nil {
		return nil, errors.Join(ErrGetFile, err)
	}
	defer func(file multipart.File) {
		_ = file.Close()
	}(file)

	// Read the file data into memory.
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.Join(ErrReadFile, err)
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

// GetFileMimeType returns the MIME type of a file using the
// github.com/gabriel-vasile/mimetype package.
func GetFileMimeType(input []byte) (string, error) {
	if input == nil {
		return "", errors.Join(ErrGetFileMimeType, ErrInputIsNil)
	}

	mtype := mimetype.Detect(input)
	if mtype == nil {
		return "", errors.Join(ErrGetFileMimeType, ErrUnsupportedType)
	}

	parts := strings.Split(mtype.String(), ";")
	return parts[0], nil
}
