package builtins

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/hay-kot/flint/pkgs/frontmatter"
	"github.com/hay-kot/flint/pkgs/set"
)

func (b BuiltIns) AssetsFunc(fields []string, dirs []string) Checker {
	return b.stringCheckFactory(fields, func(v string) bool {
		for _, dir := range dirs {
			path := filepath.Join(dir, v)
			if _, err := os.Stat(path); err == nil {
				return true
			}
		}

		return false
	})
}

func (b BuiltIns) EnumFunc(values []string, fields []string) Checker {
	valuesSet := set.New(values...)

	return b.stringCheckFactory(fields, func(v string) bool {
		return valuesSet.Contains(v)
	})
}

func (b BuiltIns) MatchFunc(patterns []string, fields []string) Checker {
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, r := range patterns {
		compiled = append(compiled, regexp.MustCompile(r))
	}

	return b.stringCheckFactory(fields, func(v string) bool {
		for _, r := range compiled {
			if r.MatchString(v) {
				return true
			}
		}
		return false
	})
}

func (b BuiltIns) LengthFunc(min, max int, fields []string) Checker {
	return b.stringCheckFactory(fields, func(v string) bool {
		return len(v) >= min && len(v) <= max
	})
}

func toMapStr(v map[string]interface{}) (map[string]string, bool) {
	m := make(map[string]string, len(v))
	for k, v := range v {
		s, ok := v.(string)
		if !ok {
			return nil, false
		}
		m[k] = s
	}

	return m, true
}

func (b BuiltIns) TypeCheck(fields []string, typeDef map[string]string) Checker {
	t := createStruct(typeDef)

	return func(fm *frontmatter.FrontMatter) error {
		errors := newValueErrors(b.ID, b.Level, b.Description)

		for _, field := range fields {
			v, ok := fm.Get(field)
			if !ok {
				continue
			}

			check := func(m map[string]string, idx int) {
				err := checkStruct(fillStruct(t, m))
				if err != nil {
					validationErrors := err.(validator.ValidationErrors)
					line, y := fm.KeyCords(field)

					for _, e := range validationErrors {
						localField := strings.ToLower(e.Field())

						fieldPath := fmt.Sprintf("%s.%s", field, localField)

						ve := newValueError(line, y, fieldPath, idx)
						if idx != -1 {
							ve.Field = fmt.Sprintf("%s[%d].%s", field, idx, localField)
						}

						ve.Description = fmt.Sprintf("failed on tag '%s'", e.ActualTag())
						errors.Errors = append(errors.Errors, ve)
					}
				}
			}

			switch v := v.(type) {
			case map[string]interface{}:
				m, ok := toMapStr(v)
				if !ok {
					continue
				}

				check(m, -1)
			case []interface{}:
				var mapSlice []map[string]string

				for _, v := range v {
					m, ok := toMapStr(v.(map[string]interface{}))
					if !ok {
						continue
					}
					mapSlice = append(mapSlice, m)
				}

				for i, m := range mapSlice {
					check(m, i)
				}
			default:
				continue
			}
		}

		if len(errors.Errors) > 0 {
			return errors
		}

		return nil
	}
}
