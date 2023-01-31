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

// func toMapStr(v map[string]interface{}) (map[string]string, bool) {
// 	m := make(map[string]string, len(v))
// 	for k, v := range v {
// 		s, ok := v.(string)
// 		if !ok {
// 			return nil, false
// 		}
// 		m[k] = s
// 	}

// 	return m, true
// }

func (b BuiltIns) TypeCheck(fields []string, typeDef map[string]string) Checker {
	t := createStruct(typeDef)

	return func(fm *frontmatter.FrontMatter) error {
		valueErrors := b.valueError()

		stringMap := func(field string, m map[string]string, idx int) {
			err := checkStruct(fillStruct(t, m))

			if err == nil { // ! inverted guard
				return
			}

			validationErrors := err.(validator.ValidationErrors) // nolint:errorLint
			line, y := fm.KeyCords(field)

			for _, e := range validationErrors {
				localField := strings.ToLower(e.Field())

				fieldPath := fmt.Sprintf("%s.%s", field, localField)

				ve := newValueError(line, y, fieldPath, idx)
				if idx != -1 {
					ve.Field = fmt.Sprintf("%s[%d].%s", field, idx, localField)
				}

				ve.Description = fmt.Sprintf("failed on tag '%s'", e.ActualTag())
				valueErrors.Errors = append(valueErrors.Errors, ve)
			}
		}

		l := looper{stringMap: stringMap}

		l.Do(fields, fm)

		if len(valueErrors.Errors) > 0 {
			return valueErrors
		}

		return nil
	}
}
