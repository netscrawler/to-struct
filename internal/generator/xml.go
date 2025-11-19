package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/clbanning/mxj/v2"
	"github.com/netscrawler/to-struct/internal/options"
	"github.com/twpayne/go-jsonstruct/v3"
)

type XMLGenerator struct{}

func (g *XMLGenerator) Generate(input io.Reader, opts *options.Options) ([]byte, error) {
	data, err := io.ReadAll(input)
	if err != nil {
		return nil, fmt.Errorf("failed to read XML input: %w", err)
	}

	mv, err := mxj.NewMapXml(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	jsonData, err := json.Marshal(mv)
	if err != nil {
		return nil, fmt.Errorf("failed to convert XML to JSON: %w", err)
	}

	generatorOptions := g.buildOptions(opts)
	generator := jsonstruct.NewGenerator(generatorOptions...)

	err = generator.ObserveJSONReader(bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}

	return generator.Generate()
}

func (g *XMLGenerator) buildOptions(opts *options.Options) []jsonstruct.GeneratorOption {
	var generatorOptions []jsonstruct.GeneratorOption

	if opts.SkipUnparsable {
		generatorOptions = append(generatorOptions, jsonstruct.WithSkipUnparsableProperties(true))
	}

	if len(opts.Tags) > 0 {
		for _, tag := range opts.Tags {
			generatorOptions = append(generatorOptions, jsonstruct.WithStructTagName(tag))
		}
	} else {
		generatorOptions = append(generatorOptions, jsonstruct.WithStructTagName("xml"))
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
