.PHONY: test test-verbose test-coverage test-services test-http test-websocket clean help

# Default target
all: test

# Run all tests
test:
	@echo "Running all tests..."
	@go test ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -func=coverage.out
	@echo ""
	@echo "To view HTML coverage report, run: make coverage-html"

# Generate HTML coverage report
coverage-html: test-coverage
	@echo "Generating HTML coverage report..."
	@go tool cover -html=coverage.out

# Test specific packages
test-services:
	@echo "Testing services package..."
	@go test -v ./services/...

test-http:
	@echo "Testing HTTP client..."
	@go test -v ./http

test-websocket:
	@echo "Testing WebSocket..."
	@go test -v ./websocket

test-errors:
	@echo "Testing errors package..."
	@go test -v ./errors

# Run tests for a specific test function
test-func:
	@echo "Usage: make test-func FUNC=TestName"
	@go test -v -run $(FUNC)

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy

# Clean test artifacts
clean:
	@echo "Cleaning test artifacts..."
	@rm -f coverage.out
	@echo "Done!"

# Build the library
build:
	@echo "Building library..."
	@go build ./...

# Run all checks (format, lint, test)
check: fmt lint test
	@echo "All checks passed!"

# Help
help:
	@echo "BingX Go Library - Makefile Commands"
	@echo ""
	@echo "Available targets:"
	@echo "  test             - Run all tests"
	@echo "  test-verbose     - Run tests with verbose output"
	@echo "  test-coverage    - Run tests with coverage report"
	@echo "  coverage-html    - Generate HTML coverage report"
	@echo "  test-services    - Test services package only"
	@echo "  test-http        - Test HTTP client only"
	@echo "  test-websocket   - Test WebSocket only"
	@echo "  test-errors      - Test errors package only"
	@echo "  test-func        - Run specific test (use FUNC=TestName)"
	@echo "  fmt              - Format code"
	@echo "  lint             - Run linter"
	@echo "  tidy             - Tidy dependencies"
	@echo "  clean            - Clean test artifacts"
	@echo "  build            - Build the library"
	@echo "  check            - Run format, lint, and test"
	@echo "  help             - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make test"
	@echo "  make test-verbose"
	@echo "  make test-coverage"
	@echo "  make test-func FUNC=TestNewClient"
