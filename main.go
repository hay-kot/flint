package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/hay-kot/flint/flint"
	"github.com/hay-kot/flint/pkgs/frontmatter"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "flint",
		Usage: "extensible frontmatter linter",
		Action: func(c *cli.Context) error {
			var err error
			start := time.Now()

			cwd := c.Args().Get(0)

			if cwd == "" {
				cwd, err = os.Getwd()
				if err != nil {
					return err
				}
			}

			if err != nil {
				return fmt.Errorf("failed to get current working directory: %w", err)
			}

			confpath := filepath.Join(cwd, "flint.yml")
			if _, err := os.Stat(confpath); os.IsNotExist(err) {
				return fmt.Errorf("failed to find flint.yml in current working directory: %w", err)
			}

			conffile, err := os.OpenFile(confpath, os.O_RDONLY, 0)
			if err != nil {
				return fmt.Errorf("failed to open flint.yml: %w", err)
			}

			conf, err := flint.ReadConfig(conffile)
			if err != nil {
				return fmt.Errorf("failed to read flint.yml: %w", err)
			}

			errors := make(map[string][]error)

			for _, c := range conf.Content {
				matches := []string{}

				for _, p := range c.Paths {
					root, p := doublestar.SplitPattern(p)
					root = filepath.Join(cwd, root)
					fsys := os.DirFS(root)

					relmatches, err := doublestar.Glob(fsys, p)
					if err != nil {
						return fmt.Errorf("failed to glob %s: %w", p, err)
					}

					for _, m := range relmatches {
						matches = append(matches, filepath.Join(root, m))
					}
				}

				for _, m := range matches {
					f, err := os.OpenFile(m, os.O_RDONLY, 0)
					if err != nil {
						return fmt.Errorf("failed to open %s: %w", m, err)
					}

					fm, err := frontmatter.Read(f)
					if err != nil {
						return fmt.Errorf("failed to read frontmatter from %s: %w", m, err)
					}

					for _, r := range c.Rules {
						rule := conf.Rules[r]
						err := rule.Check(r, fm)
						if err != nil {
							errors[m] = append(errors[m], err)
						}
					}
				}
			}

			if len(errors) > 0 {
				sorted := make([]string, len(errors))
				i := 0
				for k := range errors {
					sorted[i] = k
					i++
				}
				sort.Strings(sorted)

				for _, fp := range sorted {
					fmt.Println(flint.FmtFileErrors(fp, errors[fp]))
				}
			} else {
				fmt.Println(flint.StyleSuccess.Render("\n✓ No errors found"))
			}

			fmt.Println(flint.StyleLightGray.Render((fmt.Sprintf("\n✨ flint took %s\n", time.Since(start)))))

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
