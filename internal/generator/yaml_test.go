package generator

import (
	"strings"
	"testing"

	"to-struct/internal/options"
)

func TestYAMLGenerator_Generate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		opts    *options.Options
		wantErr bool
		checks  []string
	}{
		{
			name: "simple yaml",
			input: `name: John
age: 30
email: john@example.com`,
			opts: &options.Options{
				TypeName:       "User",
				PackageName:    "main",
				SkipUnparsable: true,
			},
			wantErr: false,
			checks:  []string{"type User struct", "Name", "Age", "Email", "yaml:"},
		},
		{
			name: "nested yaml",
			input: `user:
  name: John
  address:
    city: New York
    street: Main St`,
			opts: &options.Options{
				TypeName:       "Config",
				PackageName:    "main",
				SkipUnparsable: true,
			},
			wantErr: false,
			checks:  []string{"type Config struct", "User", "Address", "City", "Street"},
		},
		{
			name: "yaml array",
			input: `items:
  - name: Item1
    price: 100
  - name: Item2
    price: 200`,
			opts: &options.Options{
				TypeName:       "Store",
				PackageName:    "main",
				SkipUnparsable: true,
			},
			wantErr: false,
			checks:  []string{"type Store struct", "Items", "[]struct", "Name", "Price"},
		},
		{
			name: "with custom package",
			input: `name: Test`,
			opts: &options.Options{
				TypeName:       "Model",
				PackageName:    "models",
				SkipUnparsable: true,
			},
			wantErr: false,
			checks:  []string{"package models", "type Model struct"},
		},
		{
			name:    "invalid yaml",
			input:   "invalid: yaml: content: [[[",
			opts:    options.NewOptions(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &YAMLGenerator{}
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
					t.Errorf("Generate() result missing expected string: %q\nGot:\n%s", check, resultStr)
				}
			}
		})
	}
}
