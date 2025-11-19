package generator

import (
	"testing"
)

func TestGeneratorFactory_GetGenerator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		format   string
		wantType string
		wantNil  bool
	}{
		{
			name:     "json generator",
			format:   "json",
			wantType: "*generator.JSONGenerator",
			wantNil:  false,
		},
		{
			name:     "yaml generator",
			format:   "yaml",
			wantType: "*generator.YAMLGenerator",
			wantNil:  false,
		},
		{
			name:     "yml generator",
			format:   "yml",
			wantType: "*generator.YAMLGenerator",
			wantNil:  false,
		},
		{
			name:     "xml generator",
			format:   "xml",
			wantType: "*generator.XMLGenerator",
			wantNil:  false,
		},
		{
			name:     "toml generator",
			format:   "toml",
			wantType: "*generator.TOMLGenerator",
			wantNil:  false,
		},
		{
			name:    "unsupported format",
			format:  "csv",
			wantNil: true,
		},
		{
			name:    "empty format",
			format:  "",
			wantNil: true,
		},
	}

	factory := NewGeneratorFactory()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := factory.GetGenerator(tt.format)
			if tt.wantNil {
				if got != nil {
					t.Errorf("GetGenerator() = %v, want nil", got)
				}
			} else {
				if got == nil {
					t.Errorf("GetGenerator() = nil, want %s", tt.wantType)
				}
			}
		})
	}
}
