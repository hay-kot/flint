package builtins

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hay-kot/flint/pkgs/frontmatter"
)

type ErrorMatch struct {
	Line  string
	Match string
	Field string
	Value string
}

type MatchErrors struct {
	ID          string
	Level       string
	Description string
	Errors      []ErrorMatch
}

func (m MatchErrors) Error() string {
	return "match failed"
}

func (b BuiltIns) Match(fm frontmatter.FrontMatter, re, fields []string) error {
	compiled := make([]*regexp.Regexp, 0, len(re))

	for _, r := range re {
		compiled = append(compiled, regexp.MustCompile(r))
	}

	data := fm.Data()

	errors := MatchErrors{
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
					x, y := fm.KeyCords(field)

					xy := fmt.Sprintf("%d:%d", x, y)
					if x == -1 {
						xy = "0:0"
					}

					errors.Errors = append(errors.Errors, ErrorMatch{
						Line:  xy,
						Field: field,
						Value: v,
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
