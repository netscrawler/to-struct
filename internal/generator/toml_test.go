package generator

import (
	"strings"
	"testing"

	"to-struct/internal/options"
)

func TestTOMLGenerator_Generate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		opts    *options.Options
		wantErr bool
		checks  []string
	}{
		{
			name: "simple toml",
			input: `name = "John"
age = 30
email = "john@example.com"`,
			opts: &options.Options{
				TypeName:       "User",
				PackageName:    "main",
				SkipUnparsable: true,
			},
			wantErr: false,
			checks:  []string{"type User struct", "Name", "Age", "Email", "toml:"},
		},
		{
			name: "nested toml",
			input: `name = "App"
version = "1.0"

[database]
host = "localhost"
port = 5432`,
			opts: &options.Options{
				TypeName:       "Config",
				PackageName:    "main",
				SkipUnparsable: true,
			},
			wantErr: false,
			checks:  []string{"type Config struct", "Database", "Host", "Port"},
		},
		{
			name: "toml array",
			input: `[[products]]
name = "Product1"
price = 100

[[products]]
name = "Product2"
price = 200`,
			opts: &options.Options{
				TypeName:       "Store",
				PackageName:    "main",
				SkipUnparsable: true,
			},
			wantErr: false,
			checks:  []string{"type Store struct", "Products", "[]struct", "Name", "Price"},
		},
		{
			name:    "invalid toml",
			input:   `invalid toml [[[`,
			opts:    options.NewOptions(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &TOMLGenerator{}
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
