package binder

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

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

	targetVal := reflect.ValueOf(v)
	if targetVal.Kind() != reflect.Ptr {
		return errors.New("target must be a pointer to a struct")
	}

	targetElem := targetVal.Elem()
	if targetElem.Kind() != reflect.Struct {
		return errors.New("target must be a pointer to a struct")
	}

	targetType := targetElem.Type()

	// Decode the request body into the v pointer
	// if err := formDecoder.Decode(v, r.MultipartForm.Value); err != nil {
	// 	return fmt.Errorf("%w: %s", ErrDecodeForm, err.Error())
	// }

	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		tagStr := field.Tag.Get(TagForm)
		tag := strings.Split(tagStr, ",")[0]

		if tag == "" {
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
				case reflect.Array:
					arrSlice := strings.Split(formValue, ",")
					arr := reflect.New(fieldValue.Type()).Elem()
					fmt.Println(arrSlice)
					for i := 0; i < arr.Len(); i++ {
						arr.Index(i).Set(reflect.ValueOf(arrSlice[i]))
					}
					fieldValue.Set(reflect.ValueOf(arr))
				case reflect.Slice:
					fieldValue.Set(reflect.ValueOf(strings.Split(formValue, ",")))
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
			fileData, err := ioutil.ReadAll(formFile)
			if err != nil {
				return err
			}

			mime, err := GetFileMimeType(fileData)
			if err != nil {
				return err
			}

			fileStruct := &File{
				FileName:    header.Filename,
				FileSize:    header.Size,
				ContentType: mime,
				Data:        fileData,
			}

			jsonBytes, err := json.Marshal(fileStruct)
			if err != nil {
				return err
			}

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

type File struct {
	FileName    string
	FileSize    int64
	ContentType string
	Data        []byte
}
