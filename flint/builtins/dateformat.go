package builtins

import (
	"fmt"
	"time"

	"github.com/hay-kot/flint/pkgs/frontmatter"
)

func (b BuiltIns) DateFormatFunc(formats []string, fields []string) Checker {
	if len(formats) == 0 {
		formats = []string{
			time.Layout,
			time.ANSIC,
			time.UnixDate,
			time.RubyDate,
			time.RFC822,
			time.RFC822Z,
			time.RFC850,
			time.RFC1123,
			time.RFC1123Z,
			time.RFC3339,
			time.RFC3339Nano,
		}
	}

	return func(fm *frontmatter.FrontMatter) error {
		errGroup := newValueErrors(b.ID, b.Level, b.Description)

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
}
