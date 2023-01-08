package builtins

import (
	"errors"
	"fmt"
)

type ValueError struct {
	Line        string
	Description string
	Field       string
}

type ValueErrors struct {
	ID          string
	Level       string
	Description string
	Errors      []ValueError
}

func newValueErrors(id, level, description string) *ValueErrors {
	return &ValueErrors{
		ID:          id,
		Level:       level,
		Description: description,
	}
}

func (m *ValueErrors) Error() string {
	return "match failed"
}

// newValueError constructs a ValueError type, if a i >= 0 then the field
// is assumed to be an array and the index is appended to the field name.
//
// Example:
//
//	v := newValueError(1, 2, "foo", 3)
//	fmt.Println(v.Line)
//	// Output: 1:2
//	fmt.Println(v.Field)
//	// Output: foo[3]
func newValueError(line, y int, field string, i int) ValueError {
	fieldStr := field
	if i >= 0 {
		fieldStr = fmt.Sprintf("%s[%d]", fieldStr, i)
	}

	return ValueError{
		Line:  fmtKeyCords(y, line),
		Field: fieldStr,
	}
}

type FieldError = ValueError

type FieldErrors struct {
	ID          string
	Level       string
	Description string
	Fields      []FieldError
}

func (e *FieldErrors) Error() string {
	return "required keys missing"
}

func IsValueErrors(err error) bool {
	var e *ValueErrors
	return errors.As(err, &e)
}

func IsFieldErrors(err error) bool {
	var e *FieldErrors
	return errors.As(err, &e)
}
