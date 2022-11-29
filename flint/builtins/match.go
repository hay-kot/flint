package builtins

import (
	"regexp"

	"github.com/hay-kot/flint/pkgs/frontmatter"
)

func (b BuiltIns) MatchFunc(patterns []string, fields []string) CheckerFunc {
	compiled := make([]*regexp.Regexp, 0, len(patterns))

	for _, r := range patterns {
		compiled = append(compiled, regexp.MustCompile(r))
	}

	return func(fm *frontmatter.FrontMatter) error {
		return b.Match(fm, compiled, fields)
	}
}

func (b BuiltIns) Match(fm *frontmatter.FrontMatter, rgx []*regexp.Regexp, fields []string) error {
	errors := ValueErrors{
		ID:          b.ID,
		Level:       b.Level,
		Description: b.Description,
	}

	for _, field := range fields {
		v, ok := fm.Get(field)
		if !ok {
			continue
		}

		check := func(r *regexp.Regexp, v string) {
			if r.MatchString(v) {
				return
			}

			errors.Errors = append(errors.Errors, ValueError{
				Line:  fmtKeyCords(fm.KeyCords(field)),
				Field: field,
			})
		}

		switch v := v.(type) {
		case string:
			for _, re := range rgx {
				check(re, v)
			}
		case []any:
			for _, vv := range v {
				if s, ok := vv.(string); ok {
					for _, re := range rgx {
						check(re, s)
					}
				}
			}
		}

	}

	if len(errors.Errors) > 0 {
		return errors
	}

	return nil
}
