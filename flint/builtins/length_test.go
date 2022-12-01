package builtins

import (
	"strings"
	"testing"

	"github.com/hay-kot/flint/pkgs/frontmatter"
	"github.com/stretchr/testify/assert"
)

func TestBuiltIns_Length(t *testing.T) {
	type args struct {
		fields []string
		min    int
		max    int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "string",
			args: args{
				min:    1,
				max:    99,
				fields: []string{"title"},
			},
			wantErr: false,
		},
		{
			name: "string - exceed max",
			args: args{
				min:    0,
				max:    1,
				fields: []string{"title"},
			},
			wantErr: true,
		},
		{
			name: "list",
			args: args{
				min:    1,
				max:    3,
				fields: []string{"tags"},
			},
			wantErr: false,
		},
		{
			name: "list - exceed max",
			args: args{
				min:    0,
				max:    2,
				fields: []string{"tags"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New("test", "test", "test")
			fm, _ := frontmatter.Read(strings.NewReader(yml))

			check := b.LengthFunc(tt.args.min, tt.args.max, tt.args.fields)
			err := check(fm)

			switch {
			case tt.wantErr:
				if err == nil {
					t.Errorf("BuiltIns.LengthFunc() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				assert.True(t, IsValueErrors(err))
			case (err != nil) != tt.wantErr:
				t.Errorf("BuiltIns.LengthFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
