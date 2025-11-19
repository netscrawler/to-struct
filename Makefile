.PHONY: build test clean install run-examples help

BINARY_NAME=to-struct
BUILD_DIR=.
CMD_DIR=./cmd/to-struct

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -f $(BUILD_DIR)/$(BINARY_NAME)
	rm -f coverage.out coverage.html
	@echo "Clean complete"

install: ## Install the binary to GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	go install $(CMD_DIR)
	@echo "Install complete"

run-examples: build ## Run examples with all formats
	@echo "\n=== JSON Example ==="
	./$(BINARY_NAME) -f json -i examples/example.json -t User
	@echo "\n=== YAML Example ==="
	./$(BINARY_NAME) -f yaml -i examples/example.yaml -t Config
	@echo "\n=== XML Example ==="
	./$(BINARY_NAME) -f xml -i examples/example.xml -t Person
	@echo "\n=== TOML Example ==="
	./$(BINARY_NAME) -f toml -i examples/example.toml -t Settings

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...

lint: ## Run linters
	@echo "Running linters..."
	golangci-lint run ./...

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

.DEFAULT_GOAL := help
