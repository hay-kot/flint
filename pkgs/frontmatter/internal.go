package frontmatter

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

var (
	YAMLSeparator = []byte("---")
	TOMLSeparator = []byte("+++")
	NewLine       = []byte("\n")
)

func isSeparator(line []byte) bool {
	return bytes.Equal(line, YAMLSeparator) || bytes.Equal(line, TOMLSeparator)
}

func extractFrontMatter(r io.Reader) ([]byte, format, error) {
	bits := make([]byte, 0)

	first := false
	success := false

	fmFormat := formatUnknown

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()

		if !first && isSeparator(line) {
			first = true

			switch {
			case bytes.Equal(line, YAMLSeparator):
				fmFormat = formatYAML
			case bytes.Equal(line, TOMLSeparator):
				fmFormat = formatTOML
			}

			continue
		}

		if first && isSeparator(line) {
			success = true
			break
		}

		bits = append(bits, append(line, NewLine...)...)
	}

	if !success {
		return nil, fmFormat, errors.New("frontmatter: no front matter found")
	}

	return bits, fmFormat, nil
}
