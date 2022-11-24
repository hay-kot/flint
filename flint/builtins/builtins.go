package builtins

type BuiltIns struct {
	ID          string
	Level       string
	Description string
}

func New(id, level, desc string) BuiltIns {
	return BuiltIns{
		ID:          id,
		Level:       level,
		Description: desc,
	}
}
