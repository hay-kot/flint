package builtins

import (
	"reflect"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_checkStruct(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkStruct(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("checkStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ensureExported(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want string
	}{
		{
			name: "lowercase",
			key:  "abcd",
			want: "Abcd",
		},
		{
			name: "mixed case",
			key:  "aBcD",
			want: "ABcD",
		},
		{
			name: "uppercase",
			key:  "Hello",
			want: "Hello",
		},
		{
			name: "empty",
			key:  "",
			want: "",
		},
		{
			name: "single letter",
			key:  "a",
			want: "A",
		},
		{
			name: "z is not uppercase",
			key:  "z",
			want: "Z",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ensureExported(tt.key); got != tt.want {
				t.Errorf("ensureExported() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createStruct(t *testing.T) {
	tests := []struct {
		name   string
		fields map[string]string
		want   reflect.Value
	}{
		{
			name:   "empty",
			fields: map[string]string{},
			want:   reflect.ValueOf(struct{}{}),
		},
		{
			name: "one field",
			fields: map[string]string{
				"Abcd": `validate:"string"`,
			},
			want: reflect.ValueOf(struct {
				Abcd string `validate:"string"`
			}{}),
		},
		{
			name: "two fields",
			fields: map[string]string{
				"abcd": `validate:"string"`,
				"xyz":  `validate:"string"`,
			},
			want: reflect.ValueOf(struct {
				Abcd string `validate:"string"`
				Xyz  string `validate:"string"`
			}{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ptr := createStruct(tt.fields)

			got := ptr.Elem()

			// Compare the keys of the two structs
			{
				gotKeys := make([]string, 0, got.NumField())

				for i := 0; i < got.NumField(); i++ {
					gotKeys = append(gotKeys, got.Type().Field(i).Name)
				}

				wantKeys := make([]string, 0, tt.want.NumField())

				for i := 0; i < tt.want.NumField(); i++ {
					wantKeys = append(wantKeys, tt.want.Type().Field(i).Name)
				}

				sort.Strings(gotKeys)
				sort.Strings(wantKeys)

				if !reflect.DeepEqual(gotKeys, wantKeys) {
					t.Errorf("createStruct() = %v, want %v", gotKeys, wantKeys)
				}
			}

			// Compare the tags of the two structs
			{
				gotTags := make([]string, 0, got.NumField())

				for i := 0; i < got.NumField(); i++ {
					gotTags = append(gotTags, got.Type().Field(i).Tag.Get("validate"))
				}

				wantTags := make([]string, 0, tt.want.NumField())

				for i := 0; i < tt.want.NumField(); i++ {
					wantTags = append(wantTags, tt.want.Type().Field(i).Tag.Get("validate"))
				}

				assert.ElementsMatch(t, gotTags, wantTags)
			}
		})
	}
}

func Test_fillStruct(t *testing.T) {
	type args struct {
		structValue reflect.Value
		fields      map[string]string
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "empty",
			args: args{
				structValue: createStruct(map[string]string{}),
				fields:      map[string]string{},
			},
			want: struct{}{},
		},
		{
			name: "one field",
			args: args{
				structValue: createStruct(map[string]string{
					"Abcd": `validate:"string"`,
				}),
				fields: map[string]string{
					"Abcd": `Abcd Value`,
				},
			},
			want: struct {
				Abcd string `validate:"string"`
			}{
				Abcd: `Abcd Value`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fillStruct(tt.args.structValue, tt.args.fields); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fillStruct() = %v, want %v", got, tt.want)
			}
		})
	}
}
