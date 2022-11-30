package frontmatter_test

import (
	"io"
	"strings"
	"testing"
	"time"

	"github.com/hay-kot/flint/pkgs/frontmatter"
	"github.com/stretchr/testify/assert"
)

func tomlReader() io.Reader { return strings.NewReader(tomlString) }

var tomlString = `+++
tags = [ "foo", "bar" ]
categories = [ "foo", "bar" ]
title = "Hello World"
date = 2012-12-12T00:00:00.000Z

[before]
title = "Hello World"

[nested]
key = "value"
+++`

func yamlReader() io.Reader { return strings.NewReader(yamlString) }

var yamlString = `---
before:
  title: Hello World
tags:
  - foo
  - bar
categories:
  - foo
  - bar
title: Hello World
date: 2012-12-12
nested:
  key: value
---`

var data = map[string]interface{}{
	"before": map[string]interface{}{
		"title": "Hello World",
	},
	"tags": []interface{}{
		"foo",
		"bar",
	},
	"categories": []interface{}{
		"foo",
		"bar",
	},
	"title": "Hello World",
	"date":  time.Date(2012, 12, 12, 0, 0, 0, 0, time.UTC),
	"nested": map[string]interface{}{
		"key": "value",
	},
}

func Test_FrontMatter_Read(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name     string
		args     args
		wantData map[string]any
		wantErr  bool
	}{
		{
			name: "YAML",
			args: args{
				r: strings.NewReader(yamlString),
			},
			wantData: data,
			wantErr:  false,
		},
		{
			name: "TOML",
			args: args{
				r: strings.NewReader(tomlString),
			},
			wantData: data,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := frontmatter.Read(tt.args.r)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantData, got.Data())
		})
	}
}

func TestFrontMatter_KeyCords(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		reader   io.Reader
		name     string
		args     args
		wantLine int
		wantCol  int
	}{
		{
			reader:   yamlReader(),
			name:     "YAML: title (string)",
			args:     args{key: "title"},
			wantLine: 10,
			wantCol:  1,
		},
		{
			reader:   yamlReader(),
			name:     "YAML: tags (array)",
			args:     args{key: "tags"},
			wantLine: 4,
			wantCol:  1,
		},
		{
			reader:   yamlReader(),
			name:     "YAML: nested.key (string)",
			args:     args{key: "nested.key"},
			wantLine: 13,
			wantCol:  3,
		},
		// {
		// 	reader:   tomlReader(),
		// 	name:     "TOML: title (string)",
		// 	args:     args{key: "title"},
		// 	wantLine: 4,
		// 	wantCol:  1,
		// },
		// {
		// 	reader:   tomlReader(),
		// 	name:     "TOML: tags (array)",
		// 	args:     args{key: "tags"},
		// 	wantLine: 2,
		// 	wantCol:  1,
		// },
		// {
		// 	reader:   tomlReader(),
		// 	name:     "TOML: nested.key (string)",
		// 	args:     args{key: "nested.key"},
		// 	wantLine: 11,
		// 	wantCol:  3,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, _ := frontmatter.Read(tt.reader)
			gotLine, gotCol := fm.KeyCords(tt.args.key)

			assert.Equal(t, tt.wantLine, gotLine, "Row Mismatch: want=%d got=%d", tt.wantLine, gotLine)
			assert.Equal(t, tt.wantCol, gotCol, "Col Mismatch: want=%d got=%d", tt.wantCol, gotCol)
		})
	}
}

func TestFrontMatter_Content(t *testing.T) {
	tests := []struct {
		args io.Reader
		name string
		want []byte
	}{
		{
			name: "YAML",
			args: strings.NewReader(yamlString),
			want: []byte(`before:
  title: Hello World
tags:
  - foo
  - bar
categories:
  - foo
  - bar
title: Hello World
date: 2012-12-12
nested:
  key: value
`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, _ := frontmatter.Read(tt.args)

			got := fm.Content()

			assert.Equal(t, tt.want, got, "Content Mismatch: \n--- want ---\n%s\n--- got ---\n%s\n--- End --", string(tt.want), string(got))
		})
	}
}

func TestFrontMatter_Data(t *testing.T) {

	tests := []struct {
		name string
		args io.Reader
		want map[string]any
	}{
		{
			name: "YAML",
			args: yamlReader(),
			want: data,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, _ := frontmatter.Read(tt.args)
			assert.Equal(t, tt.want, fm.Data())
		})
	}
}

func TestFrontMatter_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		args   args
		want   any
		wantOk bool
	}{
		{
			name:   "title (string)",
			args:   args{key: "title"},
			want:   "Hello World",
			wantOk: true,
		},
		{
			name:   "tags (array)",
			args:   args{key: "tags"},
			want:   []interface{}{"foo", "bar"},
			wantOk: true,
		},
		{
			name:   "nested.key (string)",
			args:   args{key: "nested.key"},
			want:   "value",
			wantOk: true,
		},
		{
			name:   "no key",
			args:   args{key: "no-key"},
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, _ := frontmatter.Read(strings.NewReader(yamlString))
			got, ok := fm.Get(tt.args.key)

			assert.Equal(t, tt.wantOk, ok, "Ok Mismatch: want=%t got=%t", tt.wantOk, ok)
			assert.Equal(t, tt.want, got, "Value Mismatch: want=%v got=%v", tt.want, got)
		})
	}
}

func TestFrontMatter_Has(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "title (string)",
			args: args{key: "title"},
			want: true,
		},
		{
			name: "tags (array)",
			args: args{key: "tags"},
			want: true,
		},
		{
			name: "nested.key (string)",
			args: args{key: "nested.key"},
			want: true,
		},
		{
			name: "no key",
			args: args{key: "no-key"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, _ := frontmatter.Read(strings.NewReader(yamlString))
			assert.Equal(t, tt.want, fm.Has(tt.args.key))
		})
	}
}

func TestFrontMatter_Keys(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "YAML",
			want: []string{
				"categories",
				"title",
				"date",
				"nested.key",
				"tags",
				"before.title",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, _ := frontmatter.Read(strings.NewReader(yamlString))
			assert.ElementsMatch(t, tt.want, fm.Keys())
		})
	}
}
