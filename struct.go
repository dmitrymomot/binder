package binder

import (
	"reflect"
)

// Converts a struct to a map[string]interface{}
func structToMap(v interface{}) (map[string]interface{}, error) {
	structValue := reflect.ValueOf(v).Elem()
	structType := structValue.Type()

	structMap := make(map[string]interface{})

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		value := structValue.Field(i).Addr().Interface()

		structMap[field.Name] = value
	}

	return structMap, nil
}

// Converts a map[string]interface{} to a struct
func mapToStruct(m map[string]interface{}, v interface{}) error {
	structValue := reflect.ValueOf(v).Elem()

	for key, value := range m {
		field := structValue.FieldByName(key)
		if !field.IsValid() {
			continue
		}

		fieldValue := reflect.ValueOf(value)
		if fieldValue.Type().AssignableTo(field.Type()) {
			field.Set(fieldValue)
		}
	}

	return nil
}
