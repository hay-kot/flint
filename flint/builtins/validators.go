package builtins

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/hay-kot/flint/pkgs/set"
)

func (b BuiltIns) AssetsFunc(fields []string, dirs []string) Checker {
	return b.valueCheckerFactory(fields, func(v string) bool {
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

	return b.valueCheckerFactory(fields, func(v string) bool {
		return valuesSet.Contains(v)
	})
}

func (b BuiltIns) MatchFunc(patterns []string, fields []string) Checker {
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, r := range patterns {
		compiled = append(compiled, regexp.MustCompile(r))
	}

	return b.valueCheckerFactory(fields, func(v string) bool {
		for _, r := range compiled {
			if r.MatchString(v) {
				return true
			}
		}
		return false
	})
}

func (b BuiltIns) LengthFunc(min, max int, fields []string) Checker {
	return b.valueCheckerFactory(fields, func(v string) bool {
		return len(v) >= min && len(v) <= max
	})
}
