package generator

import (
	"strings"
	"testing"

	"github.com/netscrawler/to-struct/internal/options"
)

func TestJSONGenerator_Generate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		opts    *options.Options
		wantErr bool
		checks  []string
	}{
		{
			name: "simple object",
			input: `{
				"name": "John",
				"age": 30
			}`,
			opts: &options.Options{
				TypeName:       "User",
				PackageName:    "main",
				SkipUnparsable: true,
			},
			wantErr: false,
			checks:  []string{"type User struct", "Name", "Age", "json:"},
		},
		{
			name: "nested object",
			input: `{
				"user": {
					"name": "John",
					"email": "john@example.com"
				}
			}`,
			opts: &options.Options{
				TypeName:       "Response",
				PackageName:    "main",
				SkipUnparsable: true,
			},
			wantErr: false,
			checks:  []string{"type Response struct", "User", "Name", "Email"},
		},
		{
			name: "array of objects",
			input: `{
				"users": [
					{"name": "John", "age": 30},
					{"name": "Jane", "age": 25}
				]
			}`,
			opts: &options.Options{
				TypeName:       "Data",
				PackageName:    "main",
				SkipUnparsable: true,
			},
			wantErr: false,
			checks:  []string{"type Data struct", "Users", "[]struct"},
		},
		{
			name: "with omitempty",
			input: `{
				"name": "John",
				"age": 30
			}`,
			opts: &options.Options{
				TypeName:       "User",
				PackageName:    "main",
				OmitEmpty:      true,
				SkipUnparsable: true,
			},
			wantErr: false,
			checks:  []string{"omitempty"},
		},
		{
			name: "custom tags",
			input: `{
				"name": "John"
			}`,
			opts: &options.Options{
				TypeName:       "User",
				PackageName:    "main",
				Tags:           []string{"json", "yaml"},
				SkipUnparsable: true,
			},
			wantErr: false,
			checks:  []string{"yaml:"},
		},
		{
			name:    "invalid json",
			input:   `{invalid json}`,
			opts:    options.NewOptions(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			g := &JSONGenerator{}
			result, err := g.Generate(strings.NewReader(tt.input), tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if err != nil {
				return
			}

			resultStr := string(result)

			for _, check := range tt.checks {
				if !strings.Contains(resultStr, check) {
					t.Errorf(
						"Generate() result missing expected string: %q\nGot:\n%s",
						check,
						resultStr,
					)
				}
			}
		})
	}
}

func TestJSONGenerator_BuildOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		opts *options.Options
		want int
	}{
		{
			name: "default options",
			opts: options.NewOptions(),
			want: 6,
		},
		{
			name: "with multiple tags",
			opts: &options.Options{
				TypeName:       "Test",
				PackageName:    "main",
				Tags:           []string{"json", "yaml", "xml"},
				SkipUnparsable: true,
			},
			want: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			g := &JSONGenerator{}

			result := g.buildOptions(tt.opts)
			if len(result) != tt.want {
				t.Errorf("buildOptions() returned %d options, want %d", len(result), tt.want)
			}
		})
	}
}
