package builtins

import (
	"encoding/json"
	"reflect"
	"testing"
)

func prettyPrint(v interface{}) {
	marshalled, _ := json.MarshalIndent(v, "", "  ")

	println(string(marshalled))
}

func TestCreateStruct(t *testing.T) {
	fields := map[string]string{
		"name": `validate:"required,min=2,max=10"`,
		"age":  `validate:"required"`,
	}

	newStruct := createStruct(fields)

	newStruct2 := newStruct.Interface()

	// Ensure both fields are present
	for k := range fields {
		_, ok := reflect.TypeOf(newStruct2).Elem().FieldByName(ensureExported((k)))

		if !ok {
			t.Errorf("Field %s not found", k)
		}
	}

	v := fillStruct(newStruct, map[string]string{
		"name": "John",
		"age":  "30",
	})

	v2 := fillStruct(newStruct, map[string]string{
		"name": "John",
	})

	prettyPrint(v)
	prettyPrint(v2)

	println(checkStruct(v) == nil)
	println(checkStruct(v2) == nil)
}
