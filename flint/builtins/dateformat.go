package builtins

import (
	"fmt"
	"strings"
	"time"

	"github.com/hay-kot/flint/pkgs/frontmatter"
)

func (b BuiltIns) DateFormat(fm frontmatter.FrontMatter, format []string, fields []string) error {
	errGroup := ErrGroup{
		ID:          b.ID,
		Level:       b.Level,
		Description: b.Description,
	}

	data := fm.Data()

outer:
	for _, field := range fields {
		value, ok := extractValue(data, strings.Split(field, "."))
		if !ok {
			continue
		}

		str, ok := value.(string)
		if !ok {
			continue
		}

		for _, f := range format {
			_, err := time.Parse(f, str)

			if err == nil {
				break outer
			}
		}

		errGroup.Errors = append(errGroup.Errors, ErrGroupValue{
			Line:        fmtKeyCords(fm.KeyCords(field)),
			Description: fmt.Sprintf("%q is not allowed format", str),
			Field:       field,
		})
	}

	if len(errGroup.Errors) > 0 {
		return errGroup
	}

	return nil
}
