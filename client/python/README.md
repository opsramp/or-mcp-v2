# HPE OpsRamp MCP Python Client

A Python client for interacting with the HPE OpsRamp MCP (Model Context Protocol) server.

## Overview

This client provides a reliable way to interact with the HPE OpsRamp MCP server from Python applications. It includes a browser-like SSE client implementation that solves session validation issues with the mark3labs/mcp-go server library.

## Quick Start

### Prerequisites

- Python 3.7 or higher
- A running MCP server (see [server setup](../README.md))

### Installation and Setup

#### Using Root Level Makefile (Recommended)

You can build, run, and test the entire project (both server and client) using the Makefile at the repository root:

```bash
# From the repository root directory
cd /path/to/or-mcp-v2

# Build everything
make all

# Run the server in debug mode (background)
make run-debug

# Run the Python client browser example
make client-run-browser

# Run the Python client integrations example
make client-run-integrations

# Run all client tests (with a running server)
make client-test

# Run client integration tests
MCP_INTEGRATION_TEST=1 make client-test

# Stop the server when done
make kill-server
```

#### Using Python Client Makefile

If you prefer to work directly with the Python client:

```bash
# Navigate to the Python client directory
cd client/python

# One-step setup - creates virtual environment and installs dependencies
make setup

# Activate the virtual environment
source .venv/bin/activate  # On Windows: .venv\Scripts\activate

# Run examples, tests, etc.
make run-browser
```

### Running Examples

Using the Makefile:

```bash
# From repository root:
make client-run-browser
make client-run-integrations

# OR from client/python directory:
make run-browser
make run-integrations

# Run the example with debug logging enabled
make run-browser ARGS="--debug" 

# Run any other example
make run-example EXAMPLE=examples/check_server.py
```

### Testing

Using the Makefile:

```bash
# From repository root:
make client-test
MCP_INTEGRATION_TEST=1 make client-test

# OR from client/python directory:
make unit-test
MCP_INTEGRATION_TEST=1 make test

# Clean up temporary files
make clean

# Clean everything including virtual environment
make clean-all
```

### Simple Usage Example

```python
import asyncio
from src.ormcp.client import MCPClient

async def main():
    # Create client (automatically connects to the server)
    client = MCPClient("http://localhost:8080")
    
    try:
        # Initialize connection
        await client.initialize(client_name="my-client", client_version="1.0.0")
        
        # List available tools
        tools = await client.list_tools()
        print(f"Available tools: {tools}")
        
        # Call the integrations tool
        result = await client.call_tool("integrations", {"action": "list"})
        print(f"Integrations: {result}")
        
    finally:
        # Always close the client when done
        await client.close()

# Run the async example
asyncio.run(main())
```

## Makefile Commands

The client's `Makefile` provides several useful commands:

```
make setup           # Create virtual environment and install dependencies
make run-browser     # Run the browser-like client example
make run-integrations # Run the integrations example
make run-example     # Run a specific example file
make unit-test       # Run unit tests
make test            # Run all tests (including integration tests if MCP_INTEGRATION_TEST=1)
make clean           # Remove temporary files
make clean-all       # Remove all generated files and virtual environment
make help            # Show all available commands
```

Additional options:

```bash
# Run example with custom arguments
make run-browser ARGS="--debug --server-url=http://custom-server:8080"

# Run a specific example file
make run-example EXAMPLE=examples/custom_example.py

# Run integration tests (requires server running)
MCP_INTEGRATION_TEST=1 make test
```

## Recent Improvements

- Updated `list_tools` method to handle both response formats (direct list and nested structure)
- Fixed integration tests with proper async fixtures using pytest-asyncio
- Added proper error handling for server responses
- Enhanced the client to parse different server response formats

## Key Features

- **Browser-Like SSE Client**: Maintains persistent SSE connections like a browser
- **Session Management**: Handles session IDs and validation automatically
- **Tool Discovery**: Easy access to available MCP tools
- **Asynchronous API**: Built with asyncio for efficient I/O operations
- **Error Handling**: Comprehensive error types and recovery mechanisms
- **Event Processing**: Captures and processes server-sent events

## API Reference

### MCPClient

The main client class for interacting with the MCP server.

```python
from src.ormcp.client import MCPClient

# Create with auto-connect (default)
client = MCPClient("http://localhost:8080")

# Create without auto-connect
client = MCPClient("http://localhost:8080", auto_connect=False)

# Create with custom timeout
client = MCPClient("http://localhost:8080", connection_timeout=30)
```

#### Key Methods:

- `connect()` - Connect to the server
- `initialize(client_name, client_version)` - Initialize the connection
- `list_tools()` - Get list of available tools
- `call_tool(tool_name, arguments)` - Call a tool with arguments
- `close()` - Close the connection

### SyncMCPClient

A synchronous wrapper around the asynchronous client for non-async environments.

```python
from src.ormcp.client import SyncMCPClient

# Create client
client = SyncMCPClient("http://localhost:8080")

# Initialize
client.initialize(client_name="sync-client", client_version="1.0.0")

# List tools
tools = client.list_tools()

# Call a tool
result = client.call_tool("integrations", {"action": "list"})

# Close connection
client.close()
```

## Project Structure

```
client/python/
├── examples/                  # Example scripts
│   ├── browser_like_example.py  # Browser-like client example
│   ├── call_integrations.py     # Tool calling example
│   └── check_server.py          # Server health check
├── src/                       # Source code
│   └── ormcp/                 # Main package
│       ├── client.py          # MCPClient implementation
│       ├── session.py         # Session management
│       ├── exceptions.py      # Custom exceptions
│       └── utils.py           # Utility functions
├── tests/                     # Test directory
│   ├── test_client.py         # Client unit tests
│   ├── test_integration.py    # Integration tests with real server
├── Makefile                   # Makefile for common tasks
├── requirements.txt           # Dependencies
└── README.md                  # This file
```

## Test Configuration

Tests can be configured with environment variables:

- `MCP_SERVER_URL` - Server URL (default: http://localhost:8080)
- `DEBUG` - Enable debug logging (default: false)
- `MCP_INTEGRATION_TEST` - Enable integration tests (default: disabled)

## Known Issues

### Session Validation

This client includes a browser-like SSE client implementation that solves session validation issues with the mark3labs/mcp-go server library. Without this implementation, the server would reject JSON-RPC requests with "Invalid session ID" errors.

The solution involves:
1. Maintaining a persistent SSE connection
2. Properly processing server events
3. Using the valid session ID for all requests

See the [browser_like_example.py](./examples/browser_like_example.py) script for a complete implementation example.

### Timeout Handling

For long-running operations, you may need to adjust timeouts:

```python
# Extend connection timeout
client = MCPClient(server_url, connection_timeout=30)

# Extend request timeout for a specific call
await client.call_tool("long_running_tool", args, timeout=120)

# Extend close timeout
await client.close(timeout=10)
```

## Troubleshooting

If you encounter issues:

1. Enable debug logging:
   ```python
   import logging
   logging.getLogger('ormcp').setLevel(logging.DEBUG)
   ```

2. Check server health:
   ```bash
   python examples/check_server.py --debug
   ```

3. Make sure the server is running and accessible:
   ```bash
   # From repository root
   make kill-server  # Stop any running servers
   make run-debug    # Start a fresh server instance
   ```

4. If tests are failing, try:
   ```bash
   # Start fresh with a clean environment
   cd or-mcp-v2
   make clean-all
   make all
   make test-with-client
   ```

For more help, see the [Troubleshooting Guide](../docs/TROUBLESHOOTING.md).

## License

This project is licensed under the Apache License 2.0 - see the LICENSE file for details. 