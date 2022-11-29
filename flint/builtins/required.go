package builtins

import (
	"github.com/hay-kot/flint/pkgs/frontmatter"
	"github.com/hay-kot/flint/pkgs/set"
)

// RequiredFunc is a builtin rule for flint the allows the user/caller to specify a
// set of Required keys in the frontmatter. These respect dot seperated keys.
// so you can require nested keys by providing "author.name" as a key.
func (b BuiltIns) RequiredFunc(required []string) CheckerFunc {
	requiredSet := set.New(required...)

	return func(fm *frontmatter.FrontMatter) error {
		fmKeys := set.New(fm.Keys()...)
		missing := fmKeys.Missing(requiredSet)

		if missing.Len() > 0 {
			errs := make([]ValueError, 0, missing.Len())
			for _, key := range missing.Slice() {
				errs = append(errs, ValueError{
					Line:  fmtKeyCords(fm.KeyCords(key)),
					Field: key,
				})
			}

			return FieldErrors{
				ID:          b.ID,
				Level:       b.Level,
				Description: b.Description,
				Fields:      errs,
			}
		}

		return nil
	}
}
