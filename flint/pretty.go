package flint

import (
	"errors"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/hay-kot/flint/flint/builtins"
)

var (
	StyleFilePath   = lipgloss.NewStyle().Bold(true).Underline(true)
	StyleLightGray  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#52545A"))
	StyleSuccess    = lipgloss.NewStyle().Foreground(lipgloss.Color("#4CAF50"))
	StyleLineNumber = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#52545A")).MarginLeft(3)
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

func FmtFileErrors(path string, e []error) string {
	bldr := strings.Builder{}
	bldr.WriteString(StyleFilePath.Render(path))
	bldr.WriteString("\n")

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
			case builtins.ErrorKeysRequired:
				for _, key := range err.Fields {
					cols = append(cols, []string{
						"0:0",
						err.Level,
						key,
						err.Description,
						err.ID,
					})
				}
			case builtins.ErrGroup:
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

	bldr.WriteString(fileErrTable(cols))
	return bldr.String()
}

// fileErrTable takes in a 2D array of strings and returns a string of a table.
// including the header. It returns an evenly spaced table for neatly printing
// tables with a consistent and simple look.
func fileErrTable(rows [][]string) string {
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
		for j, s := range row {
			spaces := strings.Repeat(" ", max[j]-len(s)+4)

			switch j {
			case 1:
				if strings.Contains(s, "error") {
					s = StyleError.Render(s)
				} else if strings.Contains(s, "warning") {
					s = StyleWarning.Render(s)
				}
			case 2:
				break
			default:
				s = StyleLineNumber.Render(s)
			}

			table.WriteString(s)
			table.WriteString(spaces)
		}
		table.WriteString("\n")
	}

	return table.String()
}
