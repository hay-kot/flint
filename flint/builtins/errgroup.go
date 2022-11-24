package builtins

type ErrGroupValue struct {
	Line        string
	Description string
	Field       string
}

type ErrGroup struct {
	ID          string
	Level       string
	Description string
	Errors      []ErrGroupValue
}

func (m ErrGroup) Error() string {
	return "match failed"
}
