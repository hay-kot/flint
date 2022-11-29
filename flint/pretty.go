package flint

import (
	"errors"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/hay-kot/flint/flint/builtins"
)

var (
	StyleFilePath   = lipgloss.NewStyle().Bold(true).Underline(true)
	StyleLightGray  = lipgloss.NewStyle().Foreground(lipgloss.Color("#52545A")).Bold(true)
	StyleLineNumber = lipgloss.NewStyle().Foreground(lipgloss.Color("#52545A")).Bold(true).MarginLeft(3)
	StyleSuccess    = lipgloss.NewStyle().Foreground(lipgloss.Color("#4CAF50"))
	StyleInfo       = lipgloss.NewStyle().Foreground(lipgloss.Color("#2196F3"))
	StyleError      = lipgloss.NewStyle().Foreground(lipgloss.Color("#AB3D30"))
	StyleWarning    = lipgloss.NewStyle().Foreground(lipgloss.Color("#F25D18"))
)

func or[T comparable](a, b T) T {
	var zero T

	if a != zero {
		return a
	}
	return b
}

type prettyOptions struct {
	color bool
}

type prettyOptionFunc func(*prettyOptions)

func WithColor(v bool) prettyOptionFunc {
	return func(o *prettyOptions) {
		o.color = v
	}
}

// FmtFileErrors takes in a map of filepaths to errors and returns a string
// of the the formatted errors output.
//
// Supported Error Types
//   - builtins.FieldErrors
//   - builtins.ValueErrors
//   - RuleErrors (unwrapped into individual errors)
func FmtFileErrors(path string, e []error, optfn ...prettyOptionFunc) string {
	if len(e) == 0 {
		return ""
	}

	opts := prettyOptions{
		color: true,
	}

	for _, opt := range optfn {
		opt(&opts)
	}

	bldr := strings.Builder{}

	cols := [][]string{}

	for _, err := range e {
		all := make([]error, 0)

		if errors.As(err, &RuleErrors{}) {
			for _, ruleError := range err.(RuleErrors) {
				all = append(all, ruleError)
			}
		} else {
			all = append(all, err)
		}

		for _, e := range all {
			switch err := e.(type) {
			case builtins.FieldErrors:
				for _, key := range err.Fields {
					cols = append(cols, []string{
						or(key.Line, "0:0"),
						err.Level,
						key.Field,
						or(key.Description, err.Description),
						err.ID,
					})
				}
			case builtins.ValueErrors:
				for _, m := range err.Errors {
					cols = append(cols, []string{
						m.Line,
						err.Level,
						m.Field,
						or(m.Description, err.Description),
						err.ID,
					})
				}
			}
		}
	}

	if opts.color {
		bldr.WriteString(StyleFilePath.Render(path))
	} else {
		bldr.WriteString(path)
	}
	bldr.WriteString("\n")
	bldr.WriteString(fileErrTable(cols, opts.color))
	return bldr.String()
}

// fileErrTable takes in a 2D array of strings and returns a string of a table.
// including the header. It returns an evenly spaced table for neatly printing
// tables with a consistent and simple look.
func fileErrTable(rows [][]string, color bool) string {
	table := strings.Builder{}
	cols := len(rows[0])

	// Find longest string in each column
	max := make([]int, cols)
	for _, row := range rows {
		for i, s := range row {
			if len(s) > max[i] {
				max[i] = len(s)
			}
		}
	}

	for _, row := range rows {
		last := len(row) - 1
		for j, s := range row {
			spaces := strings.Repeat(" ", max[j]-len(s)+4)

			if color {
				switch j {
				case 1: // Level
					switch RuleLevel(strings.TrimSpace(s)) {
					case LevelError:
						s = StyleError.Render(s)
					case LevelWarning:
						s = StyleWarning.Render(s)
					case LevelInfo:
						s = StyleInfo.Render(s)
					}
				case 2: // Field
					break
				default:
					s = StyleLineNumber.Render(s)
				}
			}

			table.WriteString(s)

			if j != last {
				table.WriteString(spaces)
			}
		}
		table.WriteString("\n")
	}

	return table.String()
}
