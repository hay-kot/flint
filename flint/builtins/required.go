package builtins

import (
	"github.com/hay-kot/flint/pkgs/frontmatter"
	"github.com/hay-kot/flint/pkgs/set"
)

type ErrorKeysRequired struct {
	ID          string
	Level       string
	Description string
	Fields      []string
}

func (e ErrorKeysRequired) Error() string {
	return "required keys missing"
}

// Required is a builtin rule for flint the allows the user/caller to specify a
// set of Required keys in the frontmatter. These respect dot seperated keys.
// so you can require nested keys by providing "author.name" as a key.
func (b BuiltIns) Required(fm frontmatter.FrontMatter, keys *set.Set[string]) error {
	fmKeys := extractKeys(fm.Data())
	missing := fmKeys.Missing(keys)

	if missing.Len() > 0 {
		var errors []string
		for _, key := range missing.Slice() {
			errors = append(errors, key)
		}

		return ErrorKeysRequired{
			ID:          b.ID,
			Level:       b.Level,
			Description: b.Description,
			Fields:      errors,
		}
	}

	return nil
}
