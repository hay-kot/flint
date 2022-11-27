package frontmatter

import (
	"bufio"
	"bytes"
	"io"
)

var (
	YAMLSeparator = []byte("---")
	NewLine       = []byte("\n")
)

func keyCordsFinder(lines [][]byte, keyPath []string, offset int) (x int, y int) {
	line, col, take := -1, 0, 0

	for i, l := range lines {
		key := keyPath[0]
		if bytes.Contains(l, []byte(key+":")) {
			line = (i + offset)
			offset += i
			take = i
			for j, c := range l {
				if c != ' ' {
					col = j + 1

					break
				}
			}
			break
		}
	}

	if line == -1 {
		return -1, -1
	}

	if len(keyPath) == 1 {
		return line, col
	}

	return keyCordsFinder(lines[take:], keyPath[1:], offset)
}

func extractFrontMatter(r io.Reader) ([]byte, error) {
	bits := make([]byte, 0)

	first := false
	success := false

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()

		if !first && bytes.Equal(line, YAMLSeparator) {
			first = true
			continue
		}

		if first && bytes.Equal(line, YAMLSeparator) {
			success = true
			break
		}

		bits = append(bits, append(line, NewLine...)...)
	}

	if !success {
		return nil, nil
	}

	return bits, nil
}
