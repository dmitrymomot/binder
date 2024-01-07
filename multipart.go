package binder

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// File represents a file that was uploaded via a multipart form.
type File struct {
	// FileName stores the name of the file
	FileName string
	// FileSize stores the size of the file in bytes
	FileSize int64
	// ContentType stores the MIME type of the file
	ContentType string
	// Data is a byte slice that holds the contents of the file
	Data []byte
}

// BindFormMultipart binds the passed v pointer to the request.
// It uses the multipart/form-data content type for binding.
// `v` param should be a pointer to a struct with `form“ tags.
// Implements the binder.BinderFunc interface.
func BindFormMultipart(r *http.Request, v interface{}) error {
	// Check if the request method is POST, PUT or PATCH
	if !isPostPutPatch(r) {
		return fmt.Errorf("%w: %s", ErrInvalidMethod, r.Method)
	}

	// Check if request is multipart
	if !isMultipartFormData(r) {
		return ErrInvalidContentType
	}

	// Parse the request body
	if err := r.ParseMultipartForm(MultiPartFormMaxMemory); err != nil {
		return errors.Join(ErrParseForm, err)
	}

	// Check if the request body is empty
	if r.Body == nil {
		return ErrEmptyBody
	}

	// Get the target value
	targetVal := reflect.ValueOf(v)
	if targetVal.Kind() != reflect.Ptr {
		return ErrTargetMustBeAPointer
	}

	// Get the target element
	targetElem := targetVal.Elem()
	if targetElem.Kind() != reflect.Struct {
		return ErrTargetMustBeAStruct
	}

	// Get the target type
	targetType := targetElem.Type()

	// Iterate over the target fields
	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		tagStr := field.Tag.Get(TagForm)
		tag := strings.Split(tagStr, ",")[0]

		// Skip if tag is empty or "-"
		if tag == "" || tag == "-" {
			continue
		}

		// Bind form values
		if formValue := r.FormValue(tag); formValue != "" {
			fieldValue := targetElem.Field(i)
			if fieldValue.CanSet() {
				switch fieldValue.Kind() {
				case reflect.String:
					fieldValue.SetString(formValue)
				case reflect.Complex128, reflect.Complex64:
					complexValue, err := strconv.ParseComplex(formValue, 64)
					if err != nil {
						return err
					}
					fieldValue.SetComplex(complexValue)
				// FIXME: fix mapping of arrays and slices to a struct
				// case reflect.Array:
				// 	arrSlice := strings.Split(formValue, ",")
				// 	arr := reflect.New(fieldValue.Type()).Elem()
				// 	fmt.Println(arrSlice)
				// 	for i := 0; i < arr.Len(); i++ {
				// 		arr.Index(i).Set(reflect.ValueOf(arrSlice[i]))
				// 	}
				// 	fieldValue.Set(reflect.ValueOf(arr))
				// case reflect.Slice:
				// 	fieldValue.Set(reflect.ValueOf(strings.Split(formValue, ",")))
				case reflect.Map:
					var mapValue map[string]interface{}
					if err := json.Unmarshal([]byte(formValue), &mapValue); err != nil {
						return err
					}
					fieldValue.Set(reflect.ValueOf(mapValue))
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					intValue, err := strconv.ParseInt(formValue, 10, 64)
					if err != nil {
						return err
					}
					fieldValue.SetInt(intValue)
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					uintValue, err := strconv.ParseUint(formValue, 10, 64)
					if err != nil {
						return err
					}
					fieldValue.SetUint(uintValue)
				case reflect.Float32, reflect.Float64:
					floatValue, err := strconv.ParseFloat(formValue, 64)
					if err != nil {
						return err
					}
					fieldValue.SetFloat(floatValue)
				case reflect.Bool:
					boolValue, err := strconv.ParseBool(formValue)
					if err != nil {
						return err
					}
					fieldValue.SetBool(boolValue)
				case reflect.Ptr, reflect.Struct:
					err := json.Unmarshal([]byte(formValue), fieldValue.Addr().Interface())
					if err != nil {
						return err
					}
				default:
					return fmt.Errorf("%w: %s", ErrUnsupportedType, fieldValue.Kind())
				}
			}
		}

		// skip if the field is not a binder.File or *binder.File
		if field.Type.Kind() != reflect.ValueOf(File{}).Kind() &&
			field.Type.Kind() != reflect.ValueOf(&File{}).Kind() {
			continue
		}

		// Bind file data
		if formFile, header, err := r.FormFile(tag); err == nil {
			defer formFile.Close()
			fileData, err := io.ReadAll(formFile)
			if err != nil {
				return err
			}

			// Get the file mime type
			mime, err := GetFileMimeType(fileData)
			if err != nil {
				return err
			}

			// Create a new File struct
			fileStruct := &File{
				FileName:    header.Filename,
				FileSize:    header.Size,
				ContentType: mime,
				Data:        fileData,
			}

			// Marshal the File struct to JSON
			jsonBytes, err := json.Marshal(fileStruct)
			if err != nil {
				return err
			}

			// Unmarshal the JSON to the target field
			fieldValue := targetElem.Field(i)
			if fieldValue.CanSet() {
				switch fieldValue.Kind() {
				case reflect.Ptr, reflect.Struct:
					err = json.Unmarshal(jsonBytes, fieldValue.Addr().Interface())
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}
