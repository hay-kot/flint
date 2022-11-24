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

type RuleEnum struct {
	Values []string `yaml:"values"`
	Fields []string `yaml:"fields"`
}

type Rule struct {
	Description string    `yaml:"description"`
	Level       string    `yaml:"level"`
	Required    []string  `yaml:"builtin.required"`
	Match       RuleMatch `yaml:"builtin.match"`
	Enum        RuleEnum  `yaml:"builtin.enum"`
}

func (r Rule) Check(id string, fm frontmatter.FrontMatter) error {
	var errors RuleErrors

	check := builtins.New(id, r.Level, r.Description)

	if len(r.Required) > 0 {
		if err := check.Required(fm, set.New(r.Required...)); err != nil {
			errors = append(errors, err)
		}
	}

	if len(r.Match.Re) > 0 {
		if err := check.Match(fm, r.Match.Re, r.Match.Fields); err != nil {
			errors = append(errors, err)
		}
	}

	if len(r.Enum.Values) > 0 {
		if err := check.Enum(fm, r.Enum.Values, r.Enum.Fields); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}
