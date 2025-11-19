# to-struct

A powerful CLI tool for generating Go struct definitions from various data formats (JSON, YAML, XML, TOML).

## Features

- ✅ **Multiple Format Support**: JSON, YAML, XML, TOML
- ✅ **Smart Type Detection**: Automatically infers Go types from data
- ✅ **Customizable Output**: Control type names, package names, and struct tags
- ✅ **Multiple Tags**: Generate multiple struct tags (json, yaml, xml, toml)
- ✅ **Stdin/File Support**: Read from files or stdin pipe
- ✅ **Nested Structures**: Handles complex nested objects and arrays

## Installation

```bash
go install github.com/netscrawler/to-struct/cmd/to-struct@latest
```

Or build from source:

```bash
git clone https://github.com/netscrawler/to-struct.git
cd to-struct
make build
```

## Usage

### Basic Usage

```bash
# Generate from JSON file
to-struct -f json -i input.json -t User

# Generate from YAML file
to-struct -f yaml -i config.yaml -t Config

# Generate from XML file
to-struct -f xml -i data.xml -t Response

# Generate from TOML file
to-struct -f toml -i config.toml -t Settings
```

### Advanced Usage

```bash
# Read from stdin
cat data.json | to-struct -f json -t User

# Specify package name
to-struct -f json -i input.json -t User -p models

# Output to file
to-struct -f json -i input.json -t User -o user.go

# Generate multiple struct tags
to-struct -f json -i input.json -t User --tags json,yaml,xml

# Add omitempty to all tags
to-struct -f json -i input.json -t User --omitempty
```

## Command Line Options

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--format` | `-f` | Input format: json, yaml, xml, toml | (required) |
| `--input` | `-i` | Input file path (if not specified, reads from stdin) | |
| `--output` | `-o` | Output file path | stdout |
| `--type` | `-t` | Name of the generated type | Root |
| `--package` | `-p` | Package name for generated code | main |
| `--tags` | | Struct tags to generate (comma-separated) | format-specific |
| `--omitempty` | | Add omitempty to all tags | false |
| `--skip-unparsable` | | Skip unparsable properties | true |

## Examples

### Example 1: JSON to Go Struct

Input (`user.json`):
```json
{
  "name": "John Doe",
  "age": 30,
  "email": "john@example.com",
  "address": {
    "street": "123 Main St",
    "city": "New York"
  }
}
```

Command:
```bash
to-struct -f json -i user.json -t User
```

Output:
```go
package main

type User struct {
	Address struct {
		City   string `json:"city"`
		Street string `json:"street"`
	} `json:"address"`
	Age   int    `json:"age"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
```

### Example 2: YAML with Multiple Tags

Input (`config.yaml`):
```yaml
database:
  host: localhost
  port: 5432
server:
  host: 0.0.0.0
  port: 8080
```

Command:
```bash
to-struct -f yaml -i config.yaml -t Config --tags json,yaml
```

Output:
```go
package main

type Config struct {
	Database struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"database"`
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
}
```

### Example 3: TOML with Custom Package

Input (`settings.toml`):
```toml
name = "MyApp"
version = "1.0.0"

[database]
host = "localhost"
port = 5432
```

Command:
```bash
to-struct -f toml -i settings.toml -t Settings -p config -o settings.go
```

Output (`settings.go`):
```go
package config

type Settings struct {
	Database struct {
		Host string `toml:"host"`
		Port int    `toml:"port"`
	} `toml:"database"`
	Name    string `toml:"name"`
	Version string `toml:"version"`
}
```

### Example 4: Using Stdin Pipeline

```bash
curl -s https://api.example.com/user/1 | to-struct -f json -t User -p models

# Or with echo
echo '{"id": 123, "status": "active"}' | to-struct -f json -t Status
```

## How It Works

`to-struct` uses the excellent [go-jsonstruct](https://github.com/twpayne/go-jsonstruct) library under the hood for type inference and struct generation. For XML and TOML formats, the tool first converts them to an intermediate format that go-jsonstruct can process.

## Testing

Run the test suite:

```bash
make test
```

## Building

```bash
# Build for current platform
make build

# Run tests
make test

# Clean build artifacts
make clean
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

## Acknowledgments

- [go-jsonstruct](https://github.com/twpayne/go-jsonstruct) - Core struct generation engine
- [cobra](https://github.com/spf13/cobra) - CLI framework
- [go-toml](https://github.com/pelletier/go-toml) - TOML parser
