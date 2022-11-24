package builtins

import (
	"strings"
	"testing"

	"github.com/hay-kot/flint/pkgs/frontmatter"
)

func TestBuiltIns_Enum(t *testing.T) {

	type args struct {
		values []string
		fields []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "simple case",
			args: args{
				values: []string{"value_enum"},
				fields: []string{"key.nested"},
			},
			wantErr: false,
		},
		{
			name: "simple case with error",
			args: args{
				values: []string{"value_enumx"},
				fields: []string{"key.nested"},
			},
			wantErr: true,
		},
		{
			name: "list type error",
			args: args{
				values: []string{"foo", "bar"},
				fields: []string{"key.nested"},
			},
			wantErr: true,
		},
		{
			name: "list type",
			args: args{
				values: []string{"foo", "bar", "baz"},
				fields: []string{"tags"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New("test", "test", "test")

			fm, _ := frontmatter.Read(strings.NewReader(yml))

			if err := b.Enum(fm, tt.args.values, tt.args.fields); (err != nil) != tt.wantErr {
				t.Errorf("BuiltIns.Enum() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
