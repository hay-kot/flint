package builtins

import (
	"strings"
	"testing"

	"github.com/hay-kot/flint/pkgs/frontmatter"
	"github.com/stretchr/testify/assert"
)

func TestBuiltIns_DisallowedFunc(t *testing.T) {
	type args struct {
		disallowed []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "basic key match",
			args: args{
				disallowed: []string{"date"},
			},
			wantErr: true,
		},
		{
			name: "nested match",
			args: args{
				disallowed: []string{"key.nested", "title", "nested.slug"},
			},
			wantErr: true,
		},
		{
			name: "missing key",
			args: args{
				disallowed: []string{"banana"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New("test", "test", "test")
			fm, _ := frontmatter.Read(strings.NewReader(yml))

			check := b.DisallowedFunc(tt.args.disallowed)
			err := check(fm)

			switch tt.wantErr {
			case true:
				assert.Error(t, err)
				assert.True(t, IsFieldErrors(err))
			case false:
				assert.NoError(t, err)
			}
		})
	}
}
