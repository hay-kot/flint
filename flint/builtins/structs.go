package builtins

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func checkStruct(v interface{}) error {
	return validate.Struct(v)
}

// ensureExported ensures that the first letter of the key is capitalized.
func ensureExported(key string) string {
	if key == "" {
		return key
	}

	if key[0] >= 'a' && key[0] <= 'z' {
		key = string(key[0]+'A'-'a') + key[1:]
	}
	return key
}

// createStruct creates a new struct with the fields provided in the fields
// map. It automatically capitalizes the first letter of the field name in the
// map if it is not already capitalized.
//
// The fields map is a map of field names to struct tags.
func createStruct(fields map[string]string) reflect.Value {
	// Create a slice to hold the struct fields
	var structFields []reflect.StructField

	// Iterate through the fields and create a StructField for each
	for key, value := range fields {
		// ensure is exported
		key = ensureExported(key)

		structFields = append(structFields, reflect.StructField{
			Name: key,
			Type: reflect.TypeOf(""),
			Tag:  reflect.StructTag(value),
		})
	}

	structType := reflect.StructOf(structFields)
	structValue := reflect.New(structType).Elem()
	return structValue.Addr()
}

// fillStruct copies the reflect value of a struct and fills it with the values
// from the fields map. It automatically capitalizes the first letter of the
// field name in the map if it is not already capitalized.
//
// structValue should be a pointer to a struct, the easiest way to ensure this
// is to use the createStruct function to create the struct.
func fillStruct(structValue reflect.Value, fields map[string]string) interface{} {
	structValue = reflect.New(structValue.Elem().Type())

	// Iterate through the fields and set the value
	for key, value := range fields {
		// ensure is exported
		key = ensureExported(key)

		// check if field exists
		if _, ok := structValue.Elem().Type().FieldByName(key); !ok {
			continue
		}

		structValue.Elem().FieldByName(key).SetString(value)
	}

	// Unwrap Pointer
	return structValue.Elem().Interface()
}
