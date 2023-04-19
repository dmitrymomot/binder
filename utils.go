package binder

import (
	"net/http"
	"reflect"
	"strings"
)

// check if the request method is POST, PUT or PATCH
func isPostPutPatch(r *http.Request) bool {
	return r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch
}

// check if the request method is GET, HEAD, OPTIONS or DELETE
func isGetHeadOptionDelete(r *http.Request) bool {
	return r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions || r.Method == http.MethodDelete
}

// check if the request content type is form urlencoded
func isFormURLEncoded(r *http.Request) bool {
	return strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded")
}

// check if the request content type is multipart/form-data
func isMultipartFormData(r *http.Request) bool {
	return strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data")
}

// check if the request content type is json
func isJSON(r *http.Request) bool {
	return strings.HasPrefix(r.Header.Get("Content-Type"), "application/json")
}

// check if the passed value is a pointer
func isPointer(v interface{}) bool {
	if v == nil {
		return false
	}
	t := reflect.TypeOf(v)
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}
