package builtins

import (
	"strings"
	"testing"

	"github.com/hay-kot/flint/pkgs/frontmatter"
	"github.com/stretchr/testify/assert"
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

			check := b.EnumFunc(tt.args.values, tt.args.fields)
			err := check(fm)

			switch {
			case tt.wantErr:
				if err == nil {
					t.Errorf("BuiltIns.EnumFunc() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				assert.ErrorAs(t, err, &ValueErrors{})
			case (err != nil) != tt.wantErr:
				t.Errorf("BuiltIns.EnumFunc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
