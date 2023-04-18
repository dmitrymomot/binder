package binder

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
)

type FormData struct {
	Fields url.Values
	Files  map[string][]*FileData
}

// BindFormMultipart binds the passed v pointer to the request.
// It uses the multipart/form-data content type for binding.
// `v` param should be a pointer to a struct with `form“ tags.
// Implements the binder.BinderFunc interface.
func BindFormMultipart1(r *http.Request, v interface{}) error {
	// Check if the request method is POST, PUT or PATCH
	if !isPostPutPatch(r) {
		return fmt.Errorf("%w: %s", ErrInvalidMethod, r.Method)
	}

	// Check if request is multipart
	if !isMultipartFormData(r) {
		return ErrInvalidContentType
	}

	// Validate v pointer before decoding query into it
	if !isPointer(v) {
		return fmt.Errorf("%w: v must be a pointer to a struct", ErrInvalidInput)
	}

	// Parse the request body
	if err := r.ParseMultipartForm(MultiPartFormMaxMemory); err != nil {
		return fmt.Errorf("%w: %s", ErrParseForm, err.Error())
	}

	// Check if the request body is empty
	if r.Body == nil {
		return ErrEmptyBody
	}

	// Create form data struct
	form := FormData{
		Fields: make(url.Values),
		Files:  make(map[string][]*FileData),
	}

	// Decode the request body into the v pointer
	if err := formDecoder.Decode(v, r.MultipartForm.Value); err != nil {
		return fmt.Errorf("%w: %s", ErrDecodeForm, err.Error())
	}

	// Convert struct to map
	structMap, err := structToMap(v)
	if err != nil {
		return err
	}

	// Loop through struct fields and add to form data
	for key, value := range structMap {
		// Check if field is a file
		if isFileData(value) {
			// Get file data from request
			fileData, err := GetFileData(r, key)
			if err != nil {
				return err
			}
			if _, ok := form.Files[key]; !ok {
				form.Files[key] = make([]*FileData, 0, 1)
			}
			form.Files[key] = append(form.Files[key], fileData)
			continue
		}

		// Check if field is a slice of files
		if isSliceOfFileData(value) {
			// slice of files is not supported
			continue
		}

		// Add field to form data
		// form.Fields.Add(key, fmt.Sprintf("%v", value))
	}

	// Convert struct map back to struct
	if err := mapToStruct(structMap, v); err != nil {
		return err
	}

	return nil
}

// Populates a struct with form data
func populate(v interface{}, form FormData) error {
	// Convert struct to map
	structMap, err := structToMap(v)
	if err != nil {
		return err
	}

	if len(form.Fields) == 0 && len(form.Files) == 0 {
		return ErrEmptyRequestData
	}

	if len(form.Fields) > 0 {
		// Decode the request body into the v pointer
		if err := formDecoder.Decode(v, form.Fields); err != nil {
			return fmt.Errorf("%w: %s", ErrDecodeForm, err.Error())
		}
	}

	// Loop through files and set values in struct map
	for key, value := range form.Files {
		if f, ok := structMap[key]; ok {
			field := reflect.ValueOf(f)
			if field.Kind() == reflect.Slice {
				if structMap[key] == nil {
					structMap[key] = make([]FileData, 0, len(value))
				}
				for _, file := range value {
					structMap[key] = append(structMap[key].([]FileData), *file)
				}
			} else {
				structMap[key] = *value[0]
			}
		}
	}

	// Convert struct map back to struct
	if err := mapToStruct(structMap, v); err != nil {
		return err
	}

	return nil
}
