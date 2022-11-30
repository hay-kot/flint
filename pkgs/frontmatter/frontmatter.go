package frontmatter

import (
	"errors"
	"io"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

var (
	ErrNoFrontMatter = errors.New("no front matter found")
)

type format int

const (
	formatUnknown format = iota
	formatTOML
	formatYAML
)

type result struct {
	value any
	ok    bool
}

type FrontMatter struct {
	format  format
	content []byte
	data    map[string]any
	keys    []string
	values  map[string]result

	yamlNode *yaml.Node
}

// Read construct a new FrontMatter from the given reader.
// Supported formats are
//   - TOML
//   - YAML
func Read(r io.Reader) (fm *FrontMatter, err error) {
	content, fmFormat, err := extractFrontMatter(r)
	if err != nil {
		return nil, err
	}

	if content == nil {
		return nil, ErrNoFrontMatter
	}

	data := make(map[string]any)
	yamlNode := yaml.Node{}

	switch fmFormat {
	case formatTOML:
		err = toml.Unmarshal(content, &data)
	case formatYAML:
		err = yaml.Unmarshal(content, &data)
		if err == nil {
			err = yaml.Unmarshal(content, &yamlNode)
		}
	}

	if err != nil {
		return nil, err
	}

	return &FrontMatter{
		format:   fmFormat,
		content:  content,
		data:     data,
		yamlNode: &yamlNode,
	}, nil
}

// KeyCords returns the line and starting column of the given key.
func (fm *FrontMatter) KeyCords(key string) (line int, col int) {
	parts := strings.Split(key, ".")

	const SeparatorOffset = 1

	// Parse the YAML node to find the key
	if fm.yamlNode == nil {
		return -1, -1
	}

	switch fm.format {
	case formatTOML:
		return -1, -1
	case formatYAML:
		line, col = yamlFindLineAndCol(fm.yamlNode, parts)
	}

	if line == -1 {
		return -1, -1
	}

	return line + SeparatorOffset, col
}

// Content returns a copy of the byte slice containing the content of the
// front matter.
func (fm *FrontMatter) Content() []byte {
	c := make([]byte, len(fm.content))
	copy(c, fm.content)
	return c
}

// Data returns a copy of the front matter data.
func (fm *FrontMatter) Data() map[string]any {
	d := make(map[string]any, len(fm.data))
	for k, v := range fm.data {
		d[k] = v
	}
	return d
}

// Get returns the value of the given key. If the key is not found, the
// default value is returned and ok is false.
func (fm *FrontMatter) Get(key string) (any, bool) {
	if fm.values == nil {
		fm.values = make(map[string]result)
	}

	v, ok := fm.values[key]

	if !ok {
		parts := strings.Split(key, ".")
		var val any
		val, ok = get(fm.data, parts)
		v = result{
			value: val,
			ok:    ok,
		}
		fm.values[key] = v
	}

	return v.value, v.ok
}

func get(data map[string]any, parts []string) (any, bool) {
	if len(parts) == 1 {
		v, ok := data[parts[0]]
		if !ok {
			return nil, false
		}

		return v, true
	}

	v, ok := data[parts[0]]
	if !ok {
		return nil, false
	}

	switch v := v.(type) {
	case map[string]any:
		return get(v, parts[1:])
	default:
		return nil, false
	}
}

// Has returns true if the given key exists in the front matter.
func (fm *FrontMatter) Has(key string) bool {
	_, ok := fm.Get(key)
	return ok
}

// Keys returns a slice of all keys in the front matter.
func (fm *FrontMatter) Keys() []string {
	if fm.keys == nil {
		fm.keys = keys(fm.data)
	}

	cp := make([]string, len(fm.keys))
	copy(cp, fm.keys)
	return cp
}

func keys(data map[string]any) []string {
	keyset := make([]string, 0, len(data))

	for k := range data {
		v := data[k]

		switch v := v.(type) {
		case map[string]any:
			for _, key := range keys(v) {
				keyset = append(keyset, k+"."+key)
			}
		default:
			keyset = append(keyset, k)
		}
	}

	return keyset
}
