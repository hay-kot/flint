package builtins

import "github.com/hay-kot/flint/pkgs/frontmatter"

type Checker func(fm *frontmatter.FrontMatter) error

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

// stringCheckFactory is a helper function that creates a Checker for a given field. It
// uses the given function to check the value of the field.
func (b BuiltIns) stringCheckFactory(fields []string, f func(v string) bool) Checker {
	return func(fm *frontmatter.FrontMatter) error {
		errors := newValueErrors(b.ID, b.Level, b.Description)

		for _, field := range fields {
			v, ok := fm.Get(field)
			if !ok {
				continue
			}

			check := func(v string, i int) {
				if f(v) {
					return
				}

				line, y := fm.KeyCords(field)
				errors.Errors = append(errors.Errors, newValueError(line, y, field, i))
			}

			switch v := v.(type) {
			case string:
				check(v, -1)
			case []any:
				mapAnyToStr(v, check)
			}
		}

		if len(errors.Errors) > 0 {
			return errors
		}

		return nil
	}
}
