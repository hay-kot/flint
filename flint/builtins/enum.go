package builtins

import (
	"fmt"

	"github.com/hay-kot/flint/pkgs/frontmatter"
	"github.com/hay-kot/flint/pkgs/set"
)

func (b BuiltIns) EnumFunc(values []string, fields []string) Checker {
	return func(fm *frontmatter.FrontMatter) error {
		valueErrors := newValueErrors(b.ID, b.Level, b.Description)

		valuesSet := set.New(values...)

		for _, field := range fields {
			v, ok := fm.Get(field)
			if !ok {
				continue
			}

			addErr := func(v string) {
				xy := fmtKeyCords(fm.KeyCords(field))

				valueErrors.Errors = append(valueErrors.Errors, ValueError{
					Line:        xy,
					Description: fmt.Sprintf("%q is not an allowed values", v),
					Field:       field,
				})
			}

			switch v := v.(type) {
			case string:
				if !valuesSet.Contains(v) {
					addErr(v)
				}
			case []any:
				for _, vv := range v {
					if s, ok := vv.(string); ok {
						if !valuesSet.Contains(s) {
							addErr(s)
						}
					}
				}
			}

		}

		if len(valueErrors.Errors) > 0 {
			return valueErrors
		}

		return nil
	}
}
