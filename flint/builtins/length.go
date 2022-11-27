package builtins

import (
	"fmt"

	"github.com/hay-kot/flint/pkgs/frontmatter"
)

func (b *BuiltIns) LengthFunc(min, max int, fields []string) CheckerFunc {
	return func(fm frontmatter.FrontMatter) error {
		valueErrors := ValueErrors{
			ID:          b.ID,
			Level:       b.Level,
			Description: b.Description,
		}

		for _, field := range fields {
			v, ok := fm.Get(field)
			if !ok {
				continue
			}

			errAppender := func(v string) {
				xy := fmtKeyCords(fm.KeyCords(field))

				valueErrors.Errors = append(valueErrors.Errors, ValueError{
					Line:        xy,
					Description: fmt.Sprintf("length requirements min=%d, max=%d", min, max),
					Field:       field,
				})
			}

			switch v := v.(type) {
			case string:
				if len(v) < min || len(v) > max {
					errAppender(v)
				}

			case []any:
				if len(v) < min || len(v) > max {
					errAppender(fmt.Sprintf("%v", v))
				}
			}

		}

		if len(valueErrors.Errors) > 0 {
			return valueErrors
		}

		return nil
	}
}
