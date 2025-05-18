# HPE OpsRamp MCP Testing Guide

This document outlines the testing strategy for the HPE OpsRamp MCP system, including both the Go server and Python client.

## Testing Overview

The testing strategy involves multiple levels of testing:

1. **Unit Tests** - Testing individual components in isolation
2. **Integration Tests** - Testing the interaction between components
3. **End-to-End Tests** - Testing the complete system workflow

## Prerequisites

- Go 1.18+ (for server tests)
- Python 3.7+ (for client tests)
- Make (for running server test scripts)

## Server Testing

### Unit Tests

Server unit tests verify the functionality of individual server components:

```bash
# Run all server unit tests
cd $GOPATH/src/github.com/opsramp/or-mcp-v2
go test ./...

# Run tests for a specific package
go test ./pkg/tools/...
```

### Server Integration Tests

Server integration tests verify the HTTP endpoints and tool functionality:

```bash
# Run the integration tests
cd $GOPATH/src/github.com/opsramp/or-mcp-v2
make integration-test

# Run with debug mode enabled
make integration-test-debug
```

### Server Health Check

To perform a quick health check of the server:

```bash
# Start the server
cd $GOPATH/src/github.com/opsramp/or-mcp-v2
go run cmd/server/main.go

# In another terminal
curl http://localhost:8080/health | python -m json.tool
```

## Python Client Testing

### Client Unit Tests

The Python client has unit tests for its components:

```bash
# Navigate to the Python client directory
cd $GOPATH/src/github.com/opsramp/or-mcp-v2/client/python

# Run unit tests
python -m pytest tests/
```

### Client Integration Tests

The client integration tests verify the client's communication with the server:

```bash
# Make sure the server is running
cd $GOPATH/src/github.com/opsramp/or-mcp-v2
go run cmd/server/main.go

# In another terminal, run the integration tests
cd $GOPATH/src/github.com/opsramp/or-mcp-v2/client/python
python -m pytest tests/integration/
```

### Automated Testing with run_tests.sh

The Python client includes a script to automate testing:

```bash
cd $GOPATH/src/github.com/opsramp/or-mcp-v2/client/python
./run_tests.sh
```

This script:
1. Starts the MCP server (if not already running)
2. Runs unit tests
3. Runs integration tests
4. Stops the server (if it was started by the script)

## End-to-End Testing

For comprehensive end-to-end testing:

```bash
# Navigate to the project root
cd $GOPATH/src/github.com/opsramp/or-mcp-v2

# Run the test_mcp_server.sh script
./test_mcp_server.sh test
```

This script:
1. Starts the MCP server
2. Checks server health
3. Tests the Python client against the server
4. Verifies tool functionality

## Testing the Browser-Like Client

The browser-like client implementation can be tested directly:

```bash
# Start the server
cd $GOPATH/src/github.com/opsramp/or-mcp-v2
go run cmd/server/main.go

# In another terminal, run the browser-like example
cd $GOPATH/src/github.com/opsramp/or-mcp-v2/client/python
python examples/browser_like_example.py --debug
```

## Testing Individual Tools

To test a specific tool like the integrations tool:

```bash
# Start the server
cd $GOPATH/src/github.com/opsramp/or-mcp-v2
go run cmd/server/main.go

# In another terminal, run the integrations example
cd $GOPATH/src/github.com/opsramp/or-mcp-v2/client/python
python examples/call_integrations.py --debug
```

## Test Configuration

Tests can be configured using environment variables:

### Server Test Configuration

- `PORT` - Server port (default: 8080)
- `DEBUG` - Enable debug mode (set to "true")

### Client Test Configuration

- `MCP_SERVER_URL` - Server URL (default: http://localhost:8080)
- `DEBUG` - Enable debug logging (default: false)
- `AUTO_START_SERVER` - Auto-start server for tests (default: false)
- `CONNECTION_TIMEOUT` - Connection timeout in seconds (default: 10)
- `REQUEST_TIMEOUT` - Request timeout in seconds (default: 30)

## Continuous Integration

For CI environments, the following commands are recommended:

```bash
# For server tests
make build && make test && make integration-test-debug

# For client tests
cd client/python && python -m pytest
```

## Troubleshooting Tests

If you encounter test failures:

1. Check that the server is running
2. Verify the port is not already in use
3. Look at the server logs for errors
4. Enable debug mode for more verbose output
5. Check for any session validation issues (see TROUBLESHOOTING.md)

Common test errors and solutions are documented in the [Troubleshooting Guide](./TROUBLESHOOTING.md). 