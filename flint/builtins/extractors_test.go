package builtins

import (
	"sort"
	"testing"

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
