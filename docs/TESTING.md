# HPE OpsRamp MCP Testing Guide

This document outlines the testing strategy for the HPE OpsRamp MCP system, including both the Go server and Python client.

## Testing Overview

The testing strategy involves multiple levels of testing:

1. **Unit Tests** - Testing individual components in isolation
2. **Integration Tests** - Testing the interaction between components
3. **End-to-End Tests** - Testing the complete system workflow
4. **AI Agent Tests** - Testing the AI agent's integration expertise

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

## AI Agent Testing

The OpsRamp AI Agent provides specialized expertise for handling integrations. Testing this capability is crucial to ensure the agent responds correctly to various integration-related queries.

### Integration Expertise Testing

The `test-integrations` target in the client Makefile runs a comprehensive test of the agent's ability to handle integration-related prompts:

```bash
# Navigate to the client directory
cd $GOPATH/src/github.com/opsramp/or-mcp-v2/client

# Run all integration tests
make test-integrations

# Run tests without connecting to an actual MCP server (mock mode)
make test-integrations SIMPLE_MODE=true
```

This test:
1. Processes 37 different integration-related prompts from `agent/examples/sample_prompts.txt`
2. Verifies the agent can respond appropriately to each type of query
3. Generates detailed test results in `client/tests/integration_tests_results.txt`

### Integration Test Categories

The test prompts cover various aspects of integration management:

1. **Basic integration listing** - Listing and summarizing available integrations
2. **Integration details** - Getting detailed information about specific integrations
3. **Filtering and querying** - Finding integrations by type, status, or category
4. **Integration types** - Explaining different types of integrations 
5. **Integration operations** - Enabling, disabling, creating, updating, and deleting integrations
6. **Advanced queries** - Comparing integrations, analyzing installation history, etc.

### Testing with Mock Data

When run with `SIMPLE_MODE=true`, the tests use pre-defined mock responses to simulate interactions with HPE Alletra storage, Redfish server, and VMware vCenter integrations. This allows testing without an actual MCP server connection.

### Using Persistent SSE Connection Tests

For testing long-running connections with the MCP server, the `persistent_sse_example.py` script can be used:

```bash
# Start the server in one terminal
cd $GOPATH/src/github.com/opsramp/or-mcp-v2
make run-debug

# In another terminal, run the persistent SSE example
cd $GOPATH/src/github.com/opsramp/or-mcp-v2/client/python/examples
python persistent_sse_example.py --run-time 300 --polling-interval 30
```

This script tests:
- Establishing and maintaining a persistent SSE connection
- Handling connection interruptions and recovery
- Periodically checking integrations data over a long-running session
- Using event handlers to process server messages

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
- `SIMPLE_MODE` - Run in simple mode without actual MCP connection (default: false)

## Continuous Integration

For CI environments, the following commands are recommended:

```bash
# For server tests
make build && make test && make integration-test-debug

# For client tests
cd client/python && python -m pytest

# For AI agent integration tests
cd client && make test-integrations SIMPLE_MODE=true
```

## Troubleshooting Tests

If you encounter test failures:

1. Check that the server is running
2. Verify the port is not already in use
3. Look at the server logs for errors
4. Enable debug mode for more verbose output
5. Check for any session validation issues (see TROUBLESHOOTING.md)

Common test errors and solutions are documented in the [Troubleshooting Guide](./TROUBLESHOOTING.md). 