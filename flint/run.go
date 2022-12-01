package flint

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/hay-kot/flint/flint/builtins"
	"github.com/hay-kot/flint/pkgs/frontmatter"
)

type FileError struct {
	Path string
	Err  error
}

func (fe FileError) Error() string {
	return fmt.Sprintf("%s: %s", fe.Path, fe.Err)
}

func IsFileError(err error) bool {
	var e FileError
	return errors.As(err, &e)
}

type FlintErrors map[string][]error

func (fe FlintErrors) Error() string {
	return "flint errors"
}

func (conf *Config) Run(cwd string) (int, error) {
	errs := make(FlintErrors)

	total := 0

	for _, c := range conf.Content {
		var matches []string

		for _, p := range c.Paths {
			root, p := doublestar.SplitPattern(p)
			root = filepath.Join(cwd, root)
			fsys := os.DirFS(root)

			relmatches, err := doublestar.Glob(fsys, p)
			total += len(relmatches)
			if err != nil {
				return 0, fmt.Errorf("failed to glob %s: %w", p, err)
			}

			matches = make([]string, 0, len(relmatches))

			for _, m := range relmatches {
				matches = append(matches, filepath.Join(root, m))
			}
		}

		allChecks := make([]builtins.Checker, 0, len(c.Rules))

		for _, r := range c.Rules {
			allChecks = append(allChecks, conf.Rules[r].Funcs(r)...)
		}

		for _, m := range matches {
			f, err := os.OpenFile(m, os.O_RDONLY, 0x0)
			if err != nil {
				return 0, fmt.Errorf("failed to open %s: %w", m, err)
			}

			fm, err := frontmatter.Read(f)
			f.Close()

			if err != nil {
				if errors.Is(err, frontmatter.ErrNoFrontMatter) {
					errs[m] = append(errs[m], FileError{
						Path: m,
						Err:  err,
					})
					continue
				}

				return 0, fmt.Errorf("failed to read frontmatter unknown error: %w", err)
			}

			for _, check := range allChecks {
				err := check(fm)
				if err != nil {
					errs[m] = append(errs[m], err)
				}
			}

		}
	}

	if len(errs) > 0 {
		return total, errs
	}

	return total, nil
}
