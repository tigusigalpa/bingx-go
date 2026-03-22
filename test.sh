#!/bin/bash

# BingX Go Library Test Runner
# This script runs the test suite with various options

set -e

echo "==================================="
echo "BingX Go Library Test Suite"
echo "==================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
VERBOSE=false
COVERAGE=false
PACKAGE=""

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -c|--coverage)
            COVERAGE=true
            shift
            ;;
        -p|--package)
            PACKAGE="$2"
            shift 2
            ;;
        -h|--help)
            echo "Usage: ./test.sh [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  -v, --verbose    Run tests with verbose output"
            echo "  -c, --coverage   Generate coverage report"
            echo "  -p, --package    Run tests for specific package (e.g., ./services)"
            echo "  -h, --help       Show this help message"
            echo ""
            echo "Examples:"
            echo "  ./test.sh                    # Run all tests"
            echo "  ./test.sh -v                 # Run with verbose output"
            echo "  ./test.sh -c                 # Generate coverage report"
            echo "  ./test.sh -p ./services      # Test only services package"
            echo "  ./test.sh -v -c              # Verbose with coverage"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use -h or --help for usage information"
            exit 1
            ;;
    esac
done

# Build test command
TEST_CMD="go test"

if [ "$PACKAGE" != "" ]; then
    TEST_CMD="$TEST_CMD $PACKAGE"
else
    TEST_CMD="$TEST_CMD ./..."
fi

if [ "$VERBOSE" = true ]; then
    TEST_CMD="$TEST_CMD -v"
fi

if [ "$COVERAGE" = true ]; then
    TEST_CMD="$TEST_CMD -coverprofile=coverage.out -covermode=atomic"
fi

# Run tests
echo "Running: $TEST_CMD"
echo ""

if eval $TEST_CMD; then
    echo ""
    echo -e "${GREEN}✓ All tests passed!${NC}"
    
    # Generate coverage report if requested
    if [ "$COVERAGE" = true ]; then
        echo ""
        echo "Generating coverage report..."
        go tool cover -func=coverage.out
        echo ""
        echo "To view HTML coverage report, run:"
        echo "  go tool cover -html=coverage.out"
    fi
    
    exit 0
else
    echo ""
    echo -e "${RED}✗ Tests failed!${NC}"
    exit 1
fi
