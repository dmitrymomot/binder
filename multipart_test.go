package binder_test

import (
	"bytes"
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/dmitrymomot/binder"
)

func TestBindFormMultipart1(t *testing.T) {
	type objPayload struct {
		FieldString string  `form:"field_string,omitempty"`
		FieldInt    int     `form:"field_int,omitempty"`
		FieldUint   uint    `form:"field_uint,omitempty"`
		FieldBool   bool    `form:"field_bool,omitempty"`
		FieldFloat  float64 `form:"field_float,omitempty"`
		// FieldArray  [2]int            `form:"field_array,omitempty"`
		// FieldSlice  []int             `form:"field_slice,omitempty"`
		FieldMap map[string]string `form:"field_map,omitempty"`

		TestFile *binder.File `form:"test_file,omitempty"`
	}

	// Initialize a dummy request payload
	payload := map[string]interface{}{
		"field_string": "value",
		"field_int":    123,
		"field_uint":   123,
		"field_bool":   true,
		"field_float":  123.456,
		// "field_array":  [2]int{1, 2},
		// "field_slice":  []int{1, 2},
		"field_map": map[string]string{"key": "value"},
		"test_file": nil,
	}

	t.Run("invalid method", func(t *testing.T) {
		req, err := newQueryRequest(http.MethodGet, "/path", nil, nil)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		var obj objPayload
		if err := binder.BindFormMultipart(req, &obj); !errors.Is(err, binder.ErrInvalidMethod) {
			t.Errorf("expected error '%s', got '%s'", binder.ErrInvalidMethod, err)
		}
	})

	t.Run("invalid content type", func(t *testing.T) {
		req, err := newJSONRequest(http.MethodPost, "/path", payload, nil)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		var obj objPayload
		if err := binder.BindFormMultipart(req, &obj); !errors.Is(err, binder.ErrInvalidContentType) {
			t.Errorf("expected error '%s', got '%s'", binder.ErrInvalidContentType, err)
		}
	})

	t.Run("invalid input", func(t *testing.T) {
		req, err := newMultipartRequest(http.MethodPost, "/path", payload, nil)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		var obj objPayload
		if err := binder.BindFormMultipart(req, obj); !errors.Is(err, binder.ErrInvalidInput) {
			t.Errorf("expected error '%s', got '%s'", binder.ErrInvalidInput, err)
		}
	})

	t.Run("success: file", func(t *testing.T) {
		// Read the test image file.
		testImage, err := os.ReadFile("testdata/test.jpg")
		if err != nil {
			t.Fatalf("Failed to read test image file: %s", err.Error())
		}
		payload["test_file"] = bytes.NewReader(testImage)
		req, err := newMultipartRequest(http.MethodPost, "/path", payload, nil)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		var obj objPayload
		if err := binder.BindFormMultipart(req, &obj); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		if obj.FieldString != "value" {
			t.Errorf("expected 'value', got '%s'", obj.FieldString)
		}
		if obj.FieldInt != 123 {
			t.Errorf("expected 123, got %d", obj.FieldInt)
		}
		if obj.FieldUint != 123 {
			t.Errorf("expected 123, got %d", obj.FieldUint)
		}
		if obj.FieldBool != true {
			t.Errorf("expected true, got %t", obj.FieldBool)
		}
		if obj.FieldFloat != 123.456 {
			t.Errorf("expected 123.456, got %f", obj.FieldFloat)
		}
		// if obj.FieldArray != [2]int{1, 2} {
		// 	t.Errorf("expected [1, 2], got %v", obj.FieldArray)
		// }
		// if obj.FieldSlice == nil || len(obj.FieldSlice) != 2 || obj.FieldSlice[0] != 1 || obj.FieldSlice[1] != 2 {
		// 	t.Errorf("expected [1, 2], got %v", obj.FieldSlice)
		// }
		if obj.FieldMap == nil || len(obj.FieldMap) != 1 || obj.FieldMap["key"] != "value" {
			t.Errorf("expected {'key': 'value'}, got %v", obj.FieldMap)
		}

		if obj.TestFile == nil {
			t.Error("expected non-nil file data")
		} else {
			if obj.TestFile.FileName != "test_file" {
				t.Errorf("expected 'test_file', got '%s'", obj.TestFile.FileName)
			}
			if obj.TestFile.FileSize != int64(len(testImage)) {
				t.Errorf("expected %d, got %d", len(testImage), obj.TestFile.FileSize)
			}
			if obj.TestFile.ContentType != "image/jpeg" {
				t.Errorf("expected 'image/jpeg', got '%s'", obj.TestFile.ContentType)
			}
		}
	})
}
