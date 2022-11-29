package builtins

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

func newValueErrors(id, level, description string) ValueErrors {
	return ValueErrors{
		ID:          id,
		Level:       level,
		Description: description,
	}
}

func (m ValueErrors) Error() string {
	return "match failed"
}

type FieldError = ValueError

type FieldErrors struct {
	ID          string
	Level       string
	Description string
	Fields      []FieldError
}

func (e FieldErrors) Error() string {
	return "required keys missing"
}
