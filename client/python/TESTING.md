# HPE OpsRamp MCP Client Testing Guide

This document outlines the testing strategy for the HPE OpsRamp MCP Python client.

## Testing Structure

Our testing is organized into the following components:

1. **Unit Tests** - Testing individual client components in isolation
2. **Integration Tests** - Testing the client against a running MCP server
3. **Example Scripts** - Demonstrating client usage in real-world scenarios

## Test Scripts

### Unit Tests
Located in the `tests/` directory, these tests verify the functionality of individual components:
- `test_client.py` - Tests for the MCPClient class

### Integration Tests
Located in the `tests/integration/` directory, these tests verify client-server communication:
- `test_server_connection.py` - Tests for establishing connections to the server
- `test_browser_like_client.py` - Tests for the browser-like SSE client

### Test Utilities
- `tests/utils/server_runner.py` - Helper for starting/stopping the Go server during tests
- `tests/utils/test_config.py` - Configuration for tests

## Running Tests

### 1. Unit Tests
```bash
cd client/python
python -m pytest tests/
```

### 2. Integration Tests (Requires running server)
```bash
cd client/python
# Start the server in another terminal first
python -m pytest tests/integration/
```

### 3. Run all tests with server auto-start
```bash
cd client/python
./run_tests.sh  # This will start the server automatically
```

### 4. Manual Testing with Example Scripts
```bash
# Start server
cd $GOPATH/src/github.com/opsramp/or-mcp-v2
go run cmd/server/main.go

# In another terminal, run examples
cd $GOPATH/src/github.com/opsramp/or-mcp-v2/client/python/examples
python browser_like_example.py --debug  # Test browser-like client
python call_integrations.py --debug     # Test calling integrations tool
```

## Test Environment

Tests can be configured using environment variables:
- `MCP_SERVER_URL` - The URL of the MCP server (default: http://localhost:8080)
- `DEBUG` - Enable debug logging (default: false)
- `AUTO_START_SERVER` - Automatically start the server for tests (default: false)
- `CONNECTION_TIMEOUT` - Connection timeout in seconds (default: 10)
- `REQUEST_TIMEOUT` - Request timeout in seconds (default: 30)

## Test Dependencies

The tests require the following Python packages:
- `pytest`
- `pytest-asyncio`
- `pytest-timeout`

These are included in the `requirements.txt` file.

## Troubleshooting Tests

If tests are failing, check:
1. Is the server running?
2. Are you using the correct port?
3. Does the server have the integrations tool registered?
4. Is debug mode enabled for more verbose output?

For more details on troubleshooting, see the project's [Troubleshooting Guide](../docs/TROUBLESHOOTING.md).

## License

This project is licensed under the Apache License 2.0 - see the LICENSE file for details. 