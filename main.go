package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	_ "embed"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/hay-kot/flint/flint"
	"github.com/urfave/cli/v2"
)

//go:embed example.yml
var example []byte

var (
	version = "0.1.0"
	build   = "dev"
	date    = "unknown"

	pattern = "flint.{yml,yaml,toml,json}"
)

func confResolver(absPath string) (path string, err error) {
	fsys := os.DirFS(filepath.Dir(absPath))
	r, err := doublestar.Glob(fsys, pattern)
	if err != nil {
		return "", err
	}

	if len(r) == 0 {
		return "", fmt.Errorf("no config file found")
	}

	return r[0], nil
}

func pathResolver(cwd string, path string) (string, error) {
	if !filepath.IsAbs(path) {
		path = filepath.Join(cwd, path)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		match, err := confResolver(path)
		if err != nil {
			return "", err
		}

		return match, nil
	}

	return path, nil
}

func main() {
	app := &cli.App{
		Name:    "flint",
		Version: fmt.Sprintf("%s (%s), built at %s", version, build, date),
		Usage:   "extensible frontmatter linter",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "path to config file",
				Value:   "flint.yml",
			},
		},
		Action: func(c *cli.Context) error {
			var err error
			start := time.Now()

			cwd := c.Args().Get(0)

			if cwd == "" {
				return errors.New("no path provided")
			}

			if err != nil {
				return fmt.Errorf("failed to get current working directory: %w", err)
			}

			confpath, err := pathResolver(cwd, c.String("config"))
			if err != nil {
				return fmt.Errorf("failed to resolve config file: %w", err)
			}

			conffile, err := os.OpenFile(confpath, os.O_RDONLY, 0)
			if err != nil {
				return fmt.Errorf("failed to find config file: %w", err)
			}

			ext := filepath.Ext(confpath)

			var format flint.ConfigFormat
			switch ext {
			case ".json":
				format = flint.JSON
			case ".toml":
				format = flint.TOML
			case ".yml", ".yaml":
				format = flint.YAML
			default:
				return fmt.Errorf("unsupported config file format: %s", ext)
			}

			conf, err := flint.ReadConfig(conffile, format)
			conffile.Close()
			if err != nil {
				return fmt.Errorf("failed to find config file: %w", err)
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
		Commands: []*cli.Command{
			{
				Name:  "init",
				Usage: "create a flint.yml file in the current working directory",
				Action: func(c *cli.Context) error {
					target := c.String("config")

					if _, err := os.Stat(target); !os.IsNotExist(err) {
						return fmt.Errorf("flint.yml already exists in current working directory")
					}

					err := os.WriteFile(target, example, 0x644)
					if err != nil {
						return fmt.Errorf("failed to create flint.yml: %w", err)
					}

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
