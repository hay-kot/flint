// Package flint contains the core logic for flint.
package flint

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/hay-kot/flint/flint/builtins"
	"github.com/hay-kot/flint/pkgs/frontmatter"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

type ConfigFormat string

const (
	JSON ConfigFormat = "json"
	TOML ConfigFormat = "toml"
	YAML ConfigFormat = "yaml"
)

type TypeDef = map[string]string

type Content struct {
	Name  string   `yaml:"name" toml:"name" json:"name"`
	Paths []string `yaml:"paths" toml:"paths" json:"paths"`
	Rules []string `yaml:"rules" toml:"rules" json:"rules"`
}

type Config struct {
	Types   map[string]TypeDef `yaml:"types" toml:"types" json:"types"`
	Rules   map[string]Rule    `yaml:"rules" toml:"rules" json:"rules"`
	Content []Content          `yaml:"content" toml:"content" json:"content"`
}

func ReadConfig(r io.Reader, format ConfigFormat) (*Config, error) {
	var c Config
	var err error

	switch format {
	case JSON:
		err = json.NewDecoder(r).Decode(&c)
	case TOML:
		err = toml.NewDecoder(r).Decode(&c)
	case YAML:
		err = yaml.NewDecoder(r).Decode(&c)
	}

	return &c, err
}

// Run is the core data processing pipeline like function that performs the following
// for each conf.Content:
//
//  1. Compiles a list of all files that match the glob patterns
//  2. Transforms that list into a list of Absolute Paths
//  3. Compiles the list of rules to be applied to each file
//
// For Each File:
//  1. Read the frontmatter file into a frontmatter.FrontMatter struct
//  2. Apply the rules to the frontmatter.FrontMatter struct
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

			matches = make([]string, len(relmatches))

			for i, m := range relmatches {
				matches[i] = filepath.Join(root, m)
			}
		}

		allChecks := make([]builtins.Checker, 0, len(c.Rules))

		for _, r := range c.Rules {
			allChecks = append(allChecks, conf.Rules[r].Funcs(r, conf.Types)...)
		}

		for _, m := range matches {
			f, err := os.Open(m)
			if err != nil {
				errs[m] = append(errs[m], FileError{
					Path: m,
					Err:  err,
				})
				continue
			}

			fm, err := frontmatter.Read(f)
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

			err = f.Close()
			if err != nil {
				errs[m] = append(errs[m], FileError{
					Path: m,
					Err:  fmt.Errorf("failed to close file: %w", err),
				})
				continue
			}

			errs[m] = append(errs[m], Apply(fm, allChecks)...)
		}
	}

	if len(errs) > 0 {
		return total, errs
	}

	return total, nil
}

func Apply(fm *frontmatter.FrontMatter, rules []builtins.Checker) []error {
	var errs []error
	for _, rule := range rules {
		err := rule(fm)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}
