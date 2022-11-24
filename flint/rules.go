package flint

import (
	"strings"

	"github.com/hay-kot/flint/flint/builtins"
	"github.com/hay-kot/flint/pkgs/frontmatter"
	"github.com/hay-kot/flint/pkgs/set"
)

type RuleErrors []error

func (re RuleErrors) Error() string {
	var errors []string
	for _, err := range re {
		errors = append(errors, err.Error())
	}

	return strings.Join(errors, ", ")
}

type RuleMatch struct {
	Re     []string `yaml:"re"`
	Fields []string `yaml:"fields"`
}

type Rule struct {
	Description string    `yaml:"description"`
	Level       string    `yaml:"level"`
	Required    []string  `yaml:"builtin.required"`
	Match       RuleMatch `yaml:"builtin.match"`
}

func (r Rule) Check(id string, fm frontmatter.FrontMatter) error {
	var errors RuleErrors

	if len(r.Required) > 0 {
		if err := builtins.Required(id, r.Level, r.Description, fm, set.New(r.Required...)); err != nil {
			errors = append(errors, err)
		}
	}

	if len(r.Match.Re) > 0 {
		if err := builtins.Match(id, r.Level, r.Description, fm, r.Match.Re, r.Match.Fields); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}
