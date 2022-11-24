package builtins

import (
	"fmt"
	"strings"

	"github.com/hay-kot/flint/pkgs/frontmatter"
	"github.com/hay-kot/flint/pkgs/set"
)

func (b BuiltIns) Enum(fm frontmatter.FrontMatter, values []string, fields []string) error {
	errGroup := ErrGroup{
		ID:          b.ID,
		Level:       b.Level,
		Description: b.Description,
	}

	valuesSet := set.New(values...)
	data := fm.Data()

	for _, field := range fields {
		v, ok := extractValue(data, strings.Split(field, "."))
		if !ok {
			continue
		}

		errAppender := func(v string) {
			xy := fmtKeyCords(fm.KeyCords(field))

			errGroup.Errors = append(errGroup.Errors, ErrGroupValue{
				Line:        xy,
				Description: fmt.Sprintf("%q is not an allowed values", v),
				Field:       field,
			})
		}

		switch v := v.(type) {
		case string:
			if !valuesSet.Contains(v) {
				errAppender(v)
			}
		case []any:
			for _, vv := range v {
				if s, ok := vv.(string); ok {
					if !valuesSet.Contains(s) {
						errAppender(s)
					}
				}
			}
		}

	}

	if len(errGroup.Errors) > 0 {
		return errGroup
	}

	return nil
}
