# Test Suite

This directory contains all test scripts and test files for the HPE OpsRamp MCP Server project.

## 🧪 Test Structure

### MCP Protocol Tests
- `test_mcp_flow.sh` - Basic MCP protocol flow validation
- `test_mcp_inspector.sh` - Comprehensive MCP Inspector compatibility testing
- `test_mcp_inspector_sse.sh` - SSE-specific MCP Inspector testing
- `test_mcp_inspector_full.sh` - Full MCP Inspector workflow with persistent connections
- `test_mcp_simple.sh` - Simple MCP functionality validation

### Go Tests
- `integrations_real_api_test.go` - Real API integration tests
- `resources_real_api_test.go` - Real API resource tests
- `resources_test.go` - Unit tests for resource functionality
- `logging_test.go` - Logging functionality tests
- `testutils.go` - Test utilities and helpers

### Shell Test Scripts
- `test_integration_server.sh` - Integration server testing
- `run_tests.sh` - Test runner script

### Other Test Files
- `test_client.py` - Python client testing
- `old/` - Legacy test files
- `security/` - Security-related tests

## 🚀 Running Tests

### MCP Protocol Tests
```bash
# Basic MCP flow
./tests/test_mcp_flow.sh

# MCP Inspector compatibility
./tests/test_mcp_inspector.sh

# Full MCP Inspector workflow
./tests/test_mcp_inspector_full.sh
```

### Go Tests
```bash
# Run all Go tests
go test ./tests/...

# Run specific test file
go test ./tests/integrations_real_api_test.go
```

### Makefile Targets
```bash
# Basic testing
make test-basic

# Complete testing
make test-complete-organized

# Multi-provider testing
make test-all-providers-organized
```

## 📋 Test Requirements

- **Server Running**: Most tests require the MCP server to be running (`make run-debug` or `make run`)
- **Configuration**: Proper `config.yaml` setup for API tests
- **Dependencies**: `jq` for JSON processing in shell scripts
- **Permissions**: All `.sh` files should be executable

## 🔍 Test Coverage

The test suite covers:
- ✅ MCP protocol compliance
- ✅ Tool discovery and execution
- ✅ SSE transport layer
- ✅ JSON-RPC message handling
- ✅ Error handling and validation
- ✅ Real API integration
- ✅ Resource management
- ✅ Security scanning 