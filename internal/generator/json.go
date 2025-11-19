package generator

import (
	"io"

	"github.com/netscrawler/to-struct/internal/options"
	"github.com/twpayne/go-jsonstruct/v3"
)

type JSONGenerator struct{}

func (g *JSONGenerator) Generate(input io.Reader, opts *options.Options) ([]byte, error) {
	generatorOptions := g.buildOptions(opts)
	generator := jsonstruct.NewGenerator(generatorOptions...)

	err := generator.ObserveJSONReader(input)
	if err != nil {
		return nil, err
	}

	return generator.Generate()
}

func (g *JSONGenerator) buildOptions(opts *options.Options) []jsonstruct.GeneratorOption {
	var generatorOptions []jsonstruct.GeneratorOption

	if opts.SkipUnparsable {
		generatorOptions = append(generatorOptions, jsonstruct.WithSkipUnparsableProperties(true))
	}

	if len(opts.Tags) > 0 {
		for _, tag := range opts.Tags {
			generatorOptions = append(generatorOptions, jsonstruct.WithStructTagName(tag))
		}
	} else {
		generatorOptions = append(generatorOptions, jsonstruct.WithStructTagName("json"))
	}

	generatorOptions = append(generatorOptions, jsonstruct.WithGoFormat(true))

	if opts.OmitEmpty {
		generatorOptions = append(
			generatorOptions,
			jsonstruct.WithOmitEmptyTags(jsonstruct.OmitEmptyTagsAlways),
		)
	} else {
		generatorOptions = append(generatorOptions, jsonstruct.WithOmitEmptyTags(jsonstruct.OmitEmptyTagsAuto))
	}

	generatorOptions = append(generatorOptions,
		jsonstruct.WithTypeName(opts.TypeName),
		jsonstruct.WithPackageName(opts.PackageName))

	return generatorOptions
}
