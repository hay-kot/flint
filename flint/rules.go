package flint

import (
	"strings"

	"github.com/hay-kot/flint/flint/builtins"
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

type RuleDateFormat struct {
	Fields []string `yaml:"fields"`
	Format []string `yaml:"format"`
}

type RuleEnum struct {
	Values []string `yaml:"values"`
	Fields []string `yaml:"fields"`
}

type RuleLength struct {
	Min    int      `yaml:"min"`
	Max    int      `yaml:"max"`
	Fields []string `yaml:"fields"`
}

type Rule struct {
	Description string         `yaml:"description"`
	Level       string         `yaml:"level"`
	Required    []string       `yaml:"builtin.required"`
	Match       RuleMatch      `yaml:"builtin.match"`
	Enum        RuleEnum       `yaml:"builtin.enum"`
	DateFormat  RuleDateFormat `yaml:"builtin.date"`
	Length      RuleLength     `yaml:"builtin.length"`
}

func (r Rule) Funcs(id string) []builtins.Checker {
	check := builtins.New(id, r.Level, r.Description)
	var funcs []builtins.Checker

	if len(r.Required) > 0 {
		funcs = append(funcs, check.RequiredFunc(r.Required))
	}

	if len(r.Enum.Fields) > 0 {
		funcs = append(funcs, check.EnumFunc(r.Enum.Values, r.Enum.Fields))
	}

	if len(r.Match.Fields) > 0 {
		funcs = append(funcs, check.MatchFunc(r.Match.Re, r.Match.Fields))
	}

	if len(r.DateFormat.Fields) > 0 {
		funcs = append(funcs, check.DateFormatFunc(r.DateFormat.Format, r.DateFormat.Fields))
	}

	if len(r.Length.Fields) > 0 {
		funcs = append(funcs, check.LengthFunc(r.Length.Min, r.Length.Max, r.Length.Fields))
	}

	return funcs
}
