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
			b := New("test", "test", "test")
			fm, _ := frontmatter.Read(strings.NewReader(yml))

			check := b.MatchFunc(tt.args.re, tt.args.fields)
			err := check(fm)

			switch {
			case tt.wantErr:
				if err == nil {
					t.Errorf("BuiltIns.MatchFunc() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				assert.True(t, IsValueErrors(err))
			case (err != nil) != tt.wantErr:
				t.Errorf("BuiltIns.MatchFunc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
