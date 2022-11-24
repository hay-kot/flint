package builtins

import (
	"strings"
	"testing"

	"github.com/hay-kot/flint/pkgs/frontmatter"
)

func TestBuiltIns_DateFormat(t *testing.T) {

	type args struct {
		format []string
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
				format: []string{"2006-01-02"},
				fields: []string{"date"},
			},
			wantErr: false,
		},
		{
			name: "simple case with error",
			args: args{
				format: []string{"2006-01-02"},
				fields: []string{"date_bad"},
			},
			wantErr: false,
		},
		{
			name: "mixed case",
			args: args{
				format: []string{"2006-01-02T15:04:05Z07:00", "2006-01-02"},
				fields: []string{"date"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New("test", "test", "test")

			fm, _ := frontmatter.Read(strings.NewReader(yml))

			if err := b.DateFormat(fm, tt.args.format, tt.args.fields); (err != nil) != tt.wantErr {
				t.Errorf("BuiltIns.DateFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
