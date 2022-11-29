package flint

import (
	"testing"

	"github.com/hay-kot/flint/flint/builtins"
)

func TestFmtFileErrors(t *testing.T) {
	type args struct {
		path string
		e    []error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{
				path: "test.md",
				e:    []error{},
			},
			want: "",
		},
		{
			name: "single error",
			args: args{
				path: "test.md",
				e: []error{
					builtins.ValueErrors{
						ID:          "TEST001",
						Level:       "error",
						Description: "test description",
						Errors: []builtins.ValueError{
							{
								Line:        "1:1",
								Description: "test description",
								Field:       "test",
							},
						},
					},
				},
			},
			want: `test.md
   1:1    error    test    test description    TEST001
`,
		},
		{
			name: "multiple errors",
			args: args{
				path: "test.md",
				e: []error{
					builtins.ValueErrors{
						ID:          "TEST001",
						Level:       "error",
						Description: "test description",
						Errors: []builtins.ValueError{
							{
								Line:        "1:1",
								Description: "test description",
								Field:       "test",
							},
							{
								Line:        "2:1",
								Description: "test description",
								Field:       "test2",
							},
						},
					},
				},
			},
			want: `test.md
   1:1    error    test     test description    TEST001
   2:1    error    test2    test description    TEST001
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FmtFileErrors(tt.args.path, tt.args.e, WithColor(false))
			want := tt.want

			if got != want {
				t.Errorf("FmtFileErrors()\n-- got --\n%v\n-- want --\n%v\n", got, tt.want)
			}
		})
	}
}
