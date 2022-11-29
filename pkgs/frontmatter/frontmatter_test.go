package frontmatter_test

import (
	"io"
	"strings"
	"testing"
	"time"

	"github.com/hay-kot/flint/pkgs/frontmatter"
	"github.com/stretchr/testify/assert"
)

var YAMLFrontMatter = `---
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

func TestRead_YAML(t *testing.T) {
	assert := assert.New(t)

	reader := strings.NewReader(YAMLFrontMatter)

	fm, err := frontmatter.Read(reader)
	assert.NoError(err)

	data := fm.Data()

	assert.Equal("Hello World", data["title"])
	assert.Equal([]interface{}{"foo", "bar"}, data["tags"])
	assert.Equal([]interface{}{"foo", "bar"}, data["categories"])

	dt, _ := time.Parse("2006-01-02", "2012-12-12")
	assert.Equal(dt, data["date"])
}

func TestFrontMatter_KeyCords(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name  string
		args  args
		wantX int
		wantY int
	}{
		{
			name:  "title (string)",
			args:  args{key: "title"},
			wantX: 8,
			wantY: 1,
		},
		{
			name:  "tags (array)",
			args:  args{key: "tags"},
			wantX: 2,
			wantY: 1,
		},
		{
			name:  "nested.key (string)",
			args:  args{key: "nested.key"},
			wantX: 11,
			wantY: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, _ := frontmatter.Read(strings.NewReader(YAMLFrontMatter))
			gotX, gotY := fm.KeyCords(tt.args.key)

			assert.Equal(t, tt.wantX, gotX, "(X) Row Mismatch: want=%d got=%d", tt.wantX, gotX)
			assert.Equal(t, tt.wantY, gotY, "(Y) Col Mismatch: want=%d got=%d", tt.wantY, gotY)
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
			args: strings.NewReader(YAMLFrontMatter),
			want: []byte(`tags:
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
	dt, _ := time.Parse("2006-01-02", "2012-12-12")

	tests := []struct {
		name string
		args io.Reader
		want map[string]any
	}{
		{
			name: "YAML",
			args: strings.NewReader(YAMLFrontMatter),
			want: map[string]any{
				"tags":       []interface{}{"foo", "bar"},
				"categories": []interface{}{"foo", "bar"},
				"title":      "Hello World",
				"date":       dt,
				"nested": map[string]interface{}{
					"key": "value",
				},
			},
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
			fm, _ := frontmatter.Read(strings.NewReader(YAMLFrontMatter))
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
			fm, _ := frontmatter.Read(strings.NewReader(YAMLFrontMatter))
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
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, _ := frontmatter.Read(strings.NewReader(YAMLFrontMatter))
			assert.ElementsMatch(t, tt.want, fm.Keys())
		})
	}
}
