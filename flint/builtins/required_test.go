package builtins

import (
	"sort"
	"strings"
	"testing"

	"github.com/hay-kot/flint/pkgs/frontmatter"
	"github.com/hay-kot/flint/pkgs/set"
	"github.com/stretchr/testify/assert"
)

func TestBuiltIns_Required(t *testing.T) {
	type args struct {
		fm   frontmatter.FrontMatter
		keys *set.Set[string]
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		missing []string
	}{
		{
			name: "no keys",
			args: args{
				keys: set.New[string](),
			},
			wantErr: false,
		},
		{
			name: "single key",
			args: args{
				keys: set.New("title"),
			},
			wantErr: false,
		},
		{
			name: "multiple keys",
			args: args{
				keys: set.New("title", "date"),
			},
			wantErr: false,
		},
		{
			name: "missing key",
			args: args{
				keys: set.New("title", "date", "missing"),
			},
			wantErr: true,
			missing: []string{"missing"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error

			tt.args.fm, err = frontmatter.Read(strings.NewReader(yml))
			if err != nil {
				t.Fatal(err)
			}

			checks := New("test", "test", "test")

			err = checks.Required(tt.args.fm, tt.args.keys)
			if (err != nil) != tt.wantErr {
				t.Errorf("required() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				ekr := err.(ErrorKeysRequired)

				sort.Strings(ekr.Fields)
				sort.Strings(tt.missing)
				assert.Equal(t, tt.missing, ekr.Fields)
			}
		})
	}
}
