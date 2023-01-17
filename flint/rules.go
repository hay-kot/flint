package flint

import (
	"strings"

	"github.com/hay-kot/flint/flint/builtins"
)

type RuleLevel string

const (
	LevelInfo    RuleLevel = "info"
	LevelError   RuleLevel = "error"
	LevelWarning RuleLevel = "warning"
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
	Re     []string `yaml:"re" json:"re" toml:"re"`
	Fields []string `yaml:"fields" json:"fields" toml:"fields"`
}

type RuleDateFormat struct {
	Fields []string `yaml:"fields" json:"fields" toml:"fields"`
	Format []string `yaml:"format" json:"format" toml:"format"`
}

type RuleEnum struct {
	Values []string `yaml:"values" json:"values" toml:"values"`
	Fields []string `yaml:"fields" json:"fields" toml:"fields"`
}

type RuleLength struct {
	Min    int      `yaml:"min" json:"min" toml:"min"`
	Max    int      `yaml:"max" json:"max" toml:"max"`
	Fields []string `yaml:"fields" json:"fields" toml:"fields"`
}

type AssetRules struct {
	Sources []string `yaml:"sources" json:"sources" toml:"sources"`
	Fields  []string `yaml:"fields" json:"fields" toml:"fields"`
}

type TypeRule struct {
	Type   string   `yaml:"name" json:"name" toml:"name"`
	Fields []string `yaml:"fields" json:"fields" toml:"fields"`
}

type Rule struct {
	Description string         `yaml:"description" json:"description" toml:"description"`
	Level       RuleLevel      `yaml:"level" json:"level" toml:"level"`
	Required    []string       `yaml:"required" json:"required" toml:"required"`
	Match       RuleMatch      `yaml:"match" json:"match" toml:"match"`
	Enum        RuleEnum       `yaml:"enum" json:"enum" toml:"enum"`
	DateFormat  RuleDateFormat `yaml:"date" json:"date" toml:"date"`
	Length      RuleLength     `yaml:"length" json:"length" toml:"length"`
	Disallowed  []string       `yaml:"disallowed" json:"disallowed" toml:"disallowed"`
	Assets      AssetRules     `yaml:"assets" json:"assets" toml:"assets"`
	Type        TypeRule       `yaml:"type" json:"type" toml:"type"`
}

func (r Rule) Funcs(id string, types map[string]TypeDef) []builtins.Checker {
	check := builtins.New(id, string(r.Level), r.Description)
	var funcs []builtins.Checker

	if len(r.Required) > 0 {
		funcs = append(funcs, check.RequiredFunc(r.Required))
	}

	if len(r.Disallowed) > 0 {
		funcs = append(funcs, check.DisallowedFunc(r.Disallowed))
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

	if len(r.Assets.Fields) > 0 {
		funcs = append(funcs, check.AssetsFunc(r.Assets.Fields, r.Assets.Sources))
	}

	if len(r.Type.Fields) > 0 {
		funcs = append(funcs, check.TypeCheck(r.Type.Fields, types[r.Type.Type]))
	}

	return funcs
}
