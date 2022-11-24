package builtins

import (
	"strings"
	"testing"

	"github.com/hay-kot/flint/pkgs/frontmatter"
	"github.com/stretchr/testify/assert"
)

func TestBuiltIns_Match(t *testing.T) {
	type args struct {
		re     []string
		fields []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "literal match",
			args: args{
				re:     []string{"Hello World"},
				fields: []string{"title"},
			},
			wantErr: false,
		},
		{
			name: "nested match",
			args: args{
				re:     []string{"foo"},
				fields: []string{"nested.slug"},
			},
			wantErr: false,
		},
		{
			name: "does not match",
			args: args{
				re:     []string{"foo"},
				fields: []string{"title"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, _ := frontmatter.Read(strings.NewReader(yml))

			checks := New("test", "error", "test")

			err := checks.Match(fm, tt.args.re, tt.args.fields)

			if tt.wantErr {
				assert.ErrorAs(t, err, &ErrGroup{})
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
