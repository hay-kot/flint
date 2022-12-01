package builtins

import (
	"fmt"

	"github.com/hay-kot/flint/pkgs/frontmatter"
	"github.com/hay-kot/flint/pkgs/set"
)

func (b BuiltIns) DisallowedFunc(disallowed []string) Checker {
	disallowedSet := set.New(disallowed...)

	return func(fm *frontmatter.FrontMatter) error {
		fmKeys := set.New(fm.Keys()...)
		has := fmKeys.Intersection(disallowedSet)

		if has.Len() > 0 {
			errs := make([]FieldError, 0, has.Len())
			for _, key := range has.Slice() {
				errs = append(errs, FieldError{
					Line:        fmtKeyCords(fm.KeyCords(key)),
					Field:       key,
					Description: fmt.Sprintf("disallowed field %q", key),
				})
			}

			return &FieldErrors{
				ID:          b.ID,
				Level:       b.Level,
				Description: b.Description,
				Fields:      errs,
			}
		}

		return nil
	}
}
