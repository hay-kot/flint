package builtins

import (
	"regexp"
	"strings"

	"github.com/hay-kot/flint/pkgs/frontmatter"
)

func (b BuiltIns) Match(fm frontmatter.FrontMatter, re, fields []string) error {
	compiled := make([]*regexp.Regexp, 0, len(re))

	for _, r := range re {
		compiled = append(compiled, regexp.MustCompile(r))
	}

	data := fm.Data()

	errors := ErrGroup{
		ID:          b.ID,
		Level:       b.Level,
		Description: b.Description,
	}

	for _, field := range fields {
		parts := strings.Split(field, ".")

		v, ok := extractValue(data, parts)
		if !ok {
			continue
		}

		switch v := v.(type) {
		case string:
			for _, re := range compiled {
				if !re.MatchString(v) {
					xy := fmtKeyCords(fm.KeyCords(field))

					errors.Errors = append(errors.Errors, ErrGroupValue{
						Line:  xy,
						Field: field,
					})
				}
			}
		}
	}

	if len(errors.Errors) > 0 {
		return errors
	}

	return nil
}
