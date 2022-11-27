package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/hay-kot/flint/flint"
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

			err = conf.Run(cwd)

			if err != nil {
				switch {
				case errors.As(err, &flint.FlintErrors{}):
					errs := err.(flint.FlintErrors)
					sorted := make([]string, len(errs))
					i := 0
					for k := range errs {
						sorted[i] = k
						i++
					}
					sort.Strings(sorted)

					for _, fp := range sorted {
						fmt.Println(flint.FmtFileErrors(fp, errs[fp]))
					}
				default:
					return fmt.Errorf("failed to run flint: %w", err)
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
