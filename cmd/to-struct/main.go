package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"to-struct/internal/generator"
	"to-struct/internal/options"

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
		StringVarP(&opts.Input, "input", "i", "", "Input file path (use --stdin to read from stdin)")
	rootCmd.Flags().
		StringVarP(&opts.Output, "output", "o", "", "Output file path (default: stdout)")
	rootCmd.Flags().StringVarP(&opts.TypeName, "type", "t", "Root", "Name of the generated type")
	rootCmd.Flags().
		StringVarP(&opts.PackageName, "package", "p", "main", "Package name for generated code")
	rootCmd.Flags().BoolVar(&opts.UseStdin, "stdin", false, "Read from stdin instead of file")
	rootCmd.Flags().
		StringSliceVar(&opts.Tags, "tags", []string{}, "Struct tags to generate (comma-separated: json,xml,yaml,toml)")
	rootCmd.Flags().BoolVar(&opts.OmitEmpty, "omitempty", false, "Add omitempty to all tags")
	rootCmd.Flags().
		BoolVar(&opts.SkipUnparsable, "skip-unparsable", true, "Skip unparsable properties")

	rootCmd.MarkFlagRequired("format")
}

func run(cmd *cobra.Command, args []string) error {
	if opts.Format == "" {
		return fmt.Errorf("format is required")
	}

	var input io.Reader
	if opts.Input == "" {
		input = os.Stdin
	} else {
		file, err := os.Open(opts.Input)
		if err != nil {
			return fmt.Errorf("failed to open input file: %w", err)
		}
		defer file.Close()
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
		fmt.Print(string(result))
	} else {
		if err := os.WriteFile(opts.Output, result, 0o644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
