package binder

import "reflect"

// FileDataConverter is the function type that converts a string to a reflect.Value.
// The string is the file data.
func FileDataConverter(_ string) reflect.Value {
	return reflect.ValueOf(FileData{})
}

// FileDataConverterPtr is the function type that converts a string to a reflect.Value.
// The string is the file data.
func FileDataConverterPtr(_ string) reflect.Value {
	return reflect.ValueOf(&FileData{})
}
