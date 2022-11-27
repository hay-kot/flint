package frontmatter

import (
	"bytes"
	"errors"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	ErrNoFrontMatter = errors.New("no front matter found")
)

type result struct {
	value any
	ok    bool
}

type FrontMatter struct {
	content []byte
	data    map[string]any
	keys    []string
	values  map[string]result
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

func read(r io.Reader) (data map[string]any, content []byte, err error) {
	content, err = extractFrontMatter(r)
	if err != nil {
		return nil, nil, err
	}

	if content == nil {
		return nil, nil, ErrNoFrontMatter
	}

	data = make(map[string]any)
	err = yaml.Unmarshal(content, data)
	if err != nil {
		return nil, nil, err
	}

	return data, content, nil
}

// KeyCords returns the line and starting column of the given key.
func (fm *FrontMatter) KeyCords(key string) (x int, y int) {
	parts := strings.Split(key, ".")
	const offset = 2 // 2 line offset for 0 index and --- separator

	return keyCordsFinder(bytes.Split(fm.content, NewLine), parts, offset)
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
	parts := strings.Split(key, ".")

	if fm.values == nil {
		fm.values = make(map[string]result)
	}

	v, ok := fm.values[key]

	if !ok {
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
			return "", false
		}

		return v, true
	}

	v, ok := data[parts[0]]
	if !ok {
		return "", false
	}

	switch v := v.(type) {
	case map[string]any:
		return get(v, parts[1:])
	default:
		return "", false
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
