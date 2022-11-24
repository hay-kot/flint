package frontmatter

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	ErrNoFrontMatter = errors.New("no front matter found")
)

type FrontMatter struct {
	content []byte
	data    map[string]interface{}
}

// KeyCords returns the line and starting column of the given key.
func (fm *FrontMatter) KeyCords(key string) (x int, y int) {
	parts := strings.Split(key, ".")
	const offset = 2 // 2 line offset for 0 index and --- separator

	return KeyCordsFinder(bytes.Split(fm.content, NewLine), parts, offset)
}

func KeyCordsFinder(lines [][]byte, keyPath []string, offset int) (x int, y int) {
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

	return KeyCordsFinder(lines[take:], keyPath[1:], offset)
}

func (fm *FrontMatter) LineOfKey(key string) int {
	for i, line := range bytes.Split(fm.content, NewLine) {
		if bytes.HasPrefix(line, []byte(key+":")) {
			return i + 1
		}
	}

	return -1
}

func (fm *FrontMatter) Content() []byte {
	c := make([]byte, len(fm.content))
	copy(c, fm.content)
	return c
}

func (fm *FrontMatter) Data() map[string]interface{} {
	d := make(map[string]interface{}, len(fm.data))
	for k, v := range fm.data {
		d[k] = v
	}
	return d
}

func Read(r io.Reader) (FrontMatter, error) {
	data, content, err := read(r)
	if err != nil {
		return FrontMatter{}, err
	}

	return FrontMatter{
		data:    data,
		content: content,
	}, nil
}

func read(r io.Reader) (data map[string]interface{}, content []byte, err error) {
	content, err = extractFrontMatter(r)
	if err != nil {
		return nil, nil, err
	}

	if content == nil {
		return nil, nil, ErrNoFrontMatter // TODO
	}

	data = make(map[string]interface{})
	err = yaml.Unmarshal(content, data)
	if err != nil {
		return nil, nil, err
	}

	return data, content, nil
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
