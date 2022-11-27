package builtins

import (
	"strings"
	"testing"

	"github.com/hay-kot/flint/pkgs/frontmatter"
	"github.com/stretchr/testify/assert"
)

func TestBuiltIns_Required(t *testing.T) {
	type args struct {
		required []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "basic key match",
			args: args{
				required: []string{"date"},
			},
			wantErr: false,
		},
		{
			name: "nested match",
			args: args{
				required: []string{"key.nested", "title", "nested.slug"},
			},
			wantErr: false,
		},
		{
			name: "missing key",
			args: args{
				required: []string{"banana"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New("test", "test", "test")
			fm, _ := frontmatter.Read(strings.NewReader(yml))

			check := b.RequiredFunc(tt.args.required)
			err := check(fm)

			switch {
			case tt.wantErr:
				if err == nil {
					t.Errorf("BuiltIns.RequiredFunc() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				assert.ErrorAs(t, err, &FieldErrors{})
			case (err != nil) != tt.wantErr:
				t.Errorf("BuiltIns.RequiredFunc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
