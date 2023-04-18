package binder

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

type FormData struct {
	Fields map[string]string
	Files  map[string][]byte
}

func Bind(r *http.Request, v interface{}) error {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		return err
	}

	form := FormData{
		Fields: make(map[string]string),
		Files:  make(map[string][]byte),
	}

	// Loop through form values and add to form data
	for key, values := range r.MultipartForm.Value {
		if len(values) > 0 {
			form.Fields[key] = values[0]
		}
	}

	// Loop through files and add to form data
	for key, headers := range r.MultipartForm.File {
		for _, header := range headers {
			file, err := header.Open()
			if err != nil {
				return err
			}
			defer file.Close()

			data, err := ReadAll(file)
			if err != nil {
				return err
			}

			form.Files[key] = data
		}
	}

	// Populate struct with form data
	if err := populate(v, form); err != nil {
		return err
	}

	return nil
}

// Reads all bytes from a reader and returns them as a byte slice
func ReadAll(r io.Reader) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, r); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Populates a struct with form data
func populate(v interface{}, form FormData) error {
	// Convert struct to map
	structMap, err := structToMap(v)
	if err != nil {
		return err
	}

	// Loop through fields and set values in struct map
	for key, value := range form.Fields {
		if field, ok := structMap[key]; ok {
			if err := setValue(field, value); err != nil {
				return err
			}
		}
	}

	// Loop through files and set values in struct map
	for key, value := range form.Files {
		if field, ok := structMap[key]; ok {
			if err := setFileValue(field, key, value); err != nil {
				return err
			}
		}
	}

	// Convert struct map back to struct
	if err := mapToStruct(structMap, v); err != nil {
		return err
	}

	return nil
}

// Sets the value of a field in a struct
func setValue(field interface{}, value string) error {
	switch f := field.(type) {
	case *string:
		*f = value
	case *int:
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		*f = intValue
	// Add more types as needed
	default:
		return fmt.Errorf("unsupported field type: %T", field)
	}

	return nil
}

// Sets the value of a file field in a struct
func setFileValue(field interface{}, key string, value []byte) error {
	switch f := field.(type) {
	case *multipart.FileHeader:
		// Create temporary file to store file data
		file, err := ioutil.TempFile("", "")
		if err != nil {
			return err
		}
		defer file.Close()

		// Write file data to temporary file
		if _, err := file.Write(value); err != nil {
			return err
		}

		// Set file header fields
		f.Filename = key
		f.Size = int64(len(value))
		f.Header.Set("Content-Type", http.DetectContentType(value))
		f.Header.Set("Content-Disposition", fmt.Sprintf("form-data; name=\"%s\"; filename=\"%s\"", key, key))

		// Set file pointer to temporary file
		file.Seek(0, 0)
		// f.File = &fileWrapper{file}

	default:
		return fmt.Errorf("unsupported file field type: %T", field)
	}

	return nil
}

// Wrapper around a *os.File that implements io.ReadCloser interface
type fileWrapper struct {
	*os.File
}

func (f *fileWrapper) Close() error {
	return nil
}
