package builtins

import (
	"fmt"
	"time"

	"github.com/hay-kot/flint/pkgs/frontmatter"
)

func (b BuiltIns) DateFormatFunc(formats []string, fields []string) CheckerFunc {
	return func(fm frontmatter.FrontMatter) error {
		return b.DateFormat(fm, formats, fields)
	}
}

func (b BuiltIns) DateFormat(fm frontmatter.FrontMatter, formats []string, fields []string) error {
	errGroup := ValueErrors{
		ID:          b.ID,
		Level:       b.Level,
		Description: b.Description,
	}

	for _, field := range fields {
		value, ok := fm.Get(field)
		if !ok {
			continue
		}

		str, ok := value.(string)
		if !ok {
			continue
		}

		match := false

	inner:
		for _, f := range formats {
			_, err := time.Parse(f, str)

			if err == nil {
				match = true
				break inner
			}
		}

		if !match {
			errGroup.Errors = append(errGroup.Errors, ValueError{
				Line:        fmtKeyCords(fm.KeyCords(field)),
				Description: fmt.Sprintf("%q is not allowed format", str),
				Field:       field,
			})
		}
	}

	if len(errGroup.Errors) > 0 {
		return errGroup
	}

	return nil
}
