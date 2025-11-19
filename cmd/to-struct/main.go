package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/netscrawler/to-struct/internal/generator"
	"github.com/netscrawler/to-struct/internal/options"
	"github.com/spf13/cobra"
)

var opts *options.Options

var rootCmd = &cobra.Command{
	Use:   "to-struct",
	Short: "Generate Go structs from JSON, YAML, XML, or TOML",
	Long: `to-struct is a CLI tool that generates Go struct definitions from various data formats.
Supports JSON, YAML, XML, and TOML formats with customizable options.`,
	RunE: run,
}

func init() {
	opts = options.NewOptions()

	rootCmd.Flags().
		StringVarP(&opts.Format, "format", "f", "", "Input format: json, yaml, xml, toml (required)")
	rootCmd.Flags().
		StringVarP(&opts.Input, "input", "i", "", "Input file path (default: reads from stdin)")
	rootCmd.Flags().
		StringVarP(&opts.Output, "output", "o", "", "Output file path (default: stdout)")
	rootCmd.Flags().StringVarP(&opts.TypeName, "type", "t", "Root", "Name of the generated type")
	rootCmd.Flags().
		StringVarP(&opts.PackageName, "package", "p", "main", "Package name for generated code")
	rootCmd.Flags().
		StringSliceVar(&opts.Tags, "tags", []string{}, "Struct tags to generate (comma-separated: json,xml,yaml,toml)")
	rootCmd.Flags().BoolVar(&opts.OmitEmpty, "omitempty", false, "Add omitempty to all tags")
	rootCmd.Flags().
		BoolVar(&opts.SkipUnparsable, "skip-unparsable", true, "Skip unparsable properties")

	err := rootCmd.MarkFlagRequired("format")
	if err != nil {
		panic(err)
	}
}

func run(_ *cobra.Command, _ []string) error {
	if opts.Format == "" {
		return errors.New("format is required")
	}

	var input io.Reader

	if opts.Input == "" {
		input = os.Stdin
	} else {
		file, err := os.Open(opts.Input)
		if err != nil {
			return fmt.Errorf("failed to open input file: %w", err)
		}

		defer func() {
			closeErr := file.Close()
			if closeErr != nil {
				fmt.Fprintf(os.Stderr, "warning: failed to close file: %v\n", closeErr)
			}
		}()

		input = file
	}

	factory := generator.NewGeneratorFactory()

	gen := factory.GetGenerator(strings.ToLower(opts.Format))
	if gen == nil {
		return fmt.Errorf("unsupported format: %s", opts.Format)
	}

	result, err := gen.Generate(input, opts)
	if err != nil {
		return fmt.Errorf("failed to generate struct: %w", err)
	}

	if opts.Output == "" {
		_, err = os.Stdout.Write(result)
		if err != nil {
			return fmt.Errorf("failed to write to stdout: %w", err)
		}
	} else {
		err := os.WriteFile(opts.Output, result, 0o600)
		if err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	return nil
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
