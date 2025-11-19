package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"to-struct/internal/options"

	"github.com/pelletier/go-toml/v2"
	"github.com/twpayne/go-jsonstruct/v3"
)

type TOMLGenerator struct{}

func (g *TOMLGenerator) Generate(input io.Reader, opts *options.Options) ([]byte, error) {
	data, err := io.ReadAll(input)
	if err != nil {
		return nil, fmt.Errorf("failed to read TOML input: %w", err)
	}

	var tomlData interface{}
	if err := toml.Unmarshal(data, &tomlData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal TOML: %w", err)
	}

	jsonData, err := json.Marshal(tomlData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert TOML to JSON: %w", err)
	}

	generatorOptions := g.buildOptions(opts)
	generator := jsonstruct.NewGenerator(generatorOptions...)

	if err := generator.ObserveJSONReader(bytes.NewReader(jsonData)); err != nil {
		return nil, err
	}

	return generator.Generate()
}

func (g *TOMLGenerator) buildOptions(opts *options.Options) []jsonstruct.GeneratorOption {
	var generatorOptions []jsonstruct.GeneratorOption

	if opts.SkipUnparsable {
		generatorOptions = append(generatorOptions, jsonstruct.WithSkipUnparsableProperties(true))
	}

	if len(opts.Tags) > 0 {
		for _, tag := range opts.Tags {
			generatorOptions = append(generatorOptions, jsonstruct.WithStructTagName(tag))
		}
	} else {
		generatorOptions = append(generatorOptions, jsonstruct.WithStructTagName("toml"))
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

	generatorOptions = append(generatorOptions, jsonstruct.WithTypeName(opts.TypeName))
	generatorOptions = append(generatorOptions, jsonstruct.WithPackageName(opts.PackageName))

	return generatorOptions
}
