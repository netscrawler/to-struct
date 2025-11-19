package generator

import (
	"io"

	"to-struct/internal/options"
)

type Generator interface {
	Generate(input io.Reader, opts *options.Options) ([]byte, error)
}

type GeneratorFactory struct{}

func NewGeneratorFactory() *GeneratorFactory {
	return &GeneratorFactory{}
}

func (f *GeneratorFactory) GetGenerator(format string) Generator {
	switch format {
	case "json":
		return &JSONGenerator{}
	case "yaml", "yml":
		return &YAMLGenerator{}
	case "xml":
		return &XMLGenerator{}
	case "toml":
		return &TOMLGenerator{}
	default:
		return nil
	}
}
