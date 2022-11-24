package builtins

import (
	"sort"
	"strings"
	"testing"

	"github.com/hay-kot/flint/pkgs/frontmatter"
	"github.com/hay-kot/flint/pkgs/set"
	"github.com/stretchr/testify/assert"
)

func Test_extractKeys(t *testing.T) {
	type args struct {
		in map[string]any
	}
	tests := []struct {
		name string
		args args
		want *set.Set[string]
	}{
		{
			name: "empty map",
			args: args{
				in: map[string]any{},
			},
			want: set.New[string](),
		},
		{
			name: "single key",
			args: args{
				in: map[string]any{
					"foo": "bar",
				},
			},
			want: set.New("foo"),
		},
		{
			name: "nested keys",
			args: args{
				in: map[string]any{
					"foo": map[string]any{
						"bar": "baz",
						"qux": "quux",
						"quuz": map[string]any{
							"corge": "grault",
						},
					},
					"quuz": "corge",
				},
			},
			want: set.New("foo.bar", "foo.qux", "quuz", "foo.quuz.corge"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := tt.want.Slice()
			got := extractKeys(tt.args.in).Slice()

			sort.Strings(want)
			sort.Strings(got)
			assert.Equal(t, want, got)
		})
	}
}

var YAMLFrontMatter = `---
title: "Hello World"
date: 2012-12-12
tags:
  - foo
categories:
  - foo
key:
  nested: value
---`

func Test_Required(t *testing.T) {

	type args struct {
		fm   frontmatter.FrontMatter
		keys *set.Set[string]
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		missing []string
	}{
		{
			name: "no keys",
			args: args{
				keys: set.New[string](),
			},
			wantErr: false,
		},
		{
			name: "single key",
			args: args{
				keys: set.New("title"),
			},
			wantErr: false,
		},
		{
			name: "multiple keys",
			args: args{
				keys: set.New("title", "date"),
			},
			wantErr: false,
		},
		{
			name: "missing key",
			args: args{
				keys: set.New("title", "date", "missing"),
			},
			wantErr: true,
			missing: []string{"missing"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error

			tt.args.fm, err = frontmatter.Read(strings.NewReader(YAMLFrontMatter))
			if err != nil {
				t.Fatal(err)
			}

			err = Required("id", "level", "desc", tt.args.fm, tt.args.keys)
			if (err != nil) != tt.wantErr {
				t.Errorf("required() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				ekr := err.(ErrorKeysRequired)

				sort.Strings(ekr.Fields)
				sort.Strings(tt.missing)
				assert.Equal(t, tt.missing, ekr.Fields)
			}
		})
	}
}
