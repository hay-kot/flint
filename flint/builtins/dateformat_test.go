package builtins

import (
	"strings"
	"testing"

	"github.com/hay-kot/flint/pkgs/frontmatter"
	"github.com/stretchr/testify/assert"
)

func TestBuiltIns_DateFormat(t *testing.T) {
	type args struct {
		formats []string
		fields  []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "simple case",
			args: args{
				formats: []string{"2006-01-02"},
				fields:  []string{"date"},
			},
			wantErr: false,
		},
		{
			name: "simple case with error",
			args: args{
				formats: []string{"2006-01-02"},
				fields:  []string{"date_bad"},
			},
			wantErr: false,
		},
		{
			name: "mixed case",
			args: args{
				formats: []string{"2006-01-02T15:04:05Z07:00", "2006-01-02"},
				fields:  []string{"date"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New("test", "test", "test")
			fm, _ := frontmatter.Read(strings.NewReader(yml))

			check := b.DateFormatFunc(tt.args.formats, tt.args.fields)
			err := check(fm)

			switch {
			case tt.wantErr:
				if err == nil {
					t.Errorf("BuiltIns.DateFormatFunc() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				assert.ErrorAs(t, err, &ValueErrors{})
			case (err != nil) != tt.wantErr:
				t.Errorf("BuiltIns.DateFormatFunc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
