package frontmatter

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

var (
	YAMLSeparator      = []byte("---")
	TOMLSeparator      = []byte("+++")
	JSONSeparatorStart = []byte("{")
	JSONSeparatorEnd   = []byte("}")
	NewLine            = []byte("\n")
)

func isSeparator(line []byte, first bool) bool {
	if first {
		return bytes.Equal(line, YAMLSeparator) || bytes.Equal(line, TOMLSeparator) || bytes.Equal(line, JSONSeparatorStart)
	}

	return bytes.Equal(line, YAMLSeparator) || bytes.Equal(line, TOMLSeparator) || bytes.Equal(line, JSONSeparatorEnd)
}

func extractFrontMatter(r io.Reader) ([]byte, format, error) {
	bits := make([]byte, 0)

	first := false
	success := false

	fmFormat := formatUnknown

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()

		if !first && isSeparator(line, true) {
			first = true

			switch {
			case bytes.Equal(line, YAMLSeparator):
				fmFormat = formatYAML
			case bytes.Equal(line, TOMLSeparator):
				fmFormat = formatTOML
			case bytes.Equal(line, JSONSeparatorStart):
				fmFormat = formatJSON
				bits = append(bits, line...)
			}

			continue
		}

		if first && isSeparator(line, false) {
			success = true
			if fmFormat == formatJSON {
				bits = append(bits, line...)
			}
			break
		}

		bits = append(bits, append(line, NewLine...)...)
	}

	if !success {
		return nil, fmFormat, errors.New("frontmatter: no front matter found")
	}

	return bits, fmFormat, nil
}
