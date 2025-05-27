#!/bin/bash

# Test script for Resource Management Real API Testing
# This script runs comprehensive tests against the real OpsRamp API

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}========================================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================================${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Get the script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Change to project root
cd "$PROJECT_ROOT"

print_status "🧪 Resource Management Real API Testing"

# Check if config.yaml exists
if [ ! -f "config.yaml" ]; then
    print_error "config.yaml not found in project root"
    print_warning "Please create config.yaml with your OpsRamp credentials"
    print_warning "You can copy from config.yaml.template and fill in your values"
    exit 1
fi

print_success "Found config.yaml"

# Create output directories
print_status "📁 Creating required directories..."
mkdir -p output/logs
mkdir -p output/test-results
print_success "Directories created"

# Run unit tests first
print_status "🧪 Running unit tests..."
echo "Running resource types tests..."
go test -v ./pkg/types/ -run ".*Resource.*" | tee output/test-results/types-test.log

echo "Running resource tools tests..."
go test -v ./pkg/tools/ -run ".*Resource.*" | tee output/test-results/tools-test.log

print_success "Unit tests completed"

# Run real API tests
print_status "🌐 Running real API tests..."
echo "Testing against real OpsRamp API..."

# Run with timeout to prevent hanging
timeout 300s go test -v ./tests/ -run "TestResourcesRealAPI.*" -timeout=5m | tee output/test-results/real-api-test.log

if [ $? -eq 0 ]; then
    print_success "Real API tests completed successfully"
else
    print_warning "Some real API tests may have failed or been skipped"
    print_warning "Check output/test-results/real-api-test.log for details"
fi

# Run benchmark tests
print_status "⚡ Running performance benchmarks..."
echo "Running resource management benchmarks..."

timeout 120s go test -v ./tests/ -run="^$" -bench="BenchmarkResourcesRealAPI.*" -benchtime=5s | tee output/test-results/benchmark-test.log

if [ $? -eq 0 ]; then
    print_success "Benchmark tests completed"
else
    print_warning "Benchmark tests may have failed or been skipped"
fi

# Test coverage
print_status "📊 Generating test coverage..."
echo "Calculating test coverage for resource management..."

go test -coverprofile=output/test-results/coverage.out ./pkg/types/ ./pkg/tools/ -run ".*Resource.*"
go tool cover -html=output/test-results/coverage.out -o output/test-results/coverage.html

print_success "Coverage report generated: output/test-results/coverage.html"

# Summary
print_status "📋 Test Summary"
echo "Test results saved to:"
echo "  - Types tests: output/test-results/types-test.log"
echo "  - Tools tests: output/test-results/tools-test.log"
echo "  - Real API tests: output/test-results/real-api-test.log"
echo "  - Benchmark tests: output/test-results/benchmark-test.log"
echo "  - Coverage report: output/test-results/coverage.html"
echo ""
echo "Logs saved to: output/logs/or-mcp.log"

# Check if any tests failed
if grep -q "FAIL" output/test-results/*.log; then
    print_warning "Some tests failed. Check the log files for details."
    exit 1
else
    print_success "All tests completed successfully!"
fi 