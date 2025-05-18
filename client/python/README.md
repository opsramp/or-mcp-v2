# HPE OpsRamp MCP Python Client

A Python client for interacting with the HPE OpsRamp MCP (Model Context Protocol) server.

## Overview

This client provides a reliable way to interact with the HPE OpsRamp MCP server from Python applications. It includes a browser-like SSE client implementation that solves session validation issues with the mark3labs/mcp-go server library.

## Quick Start

### Prerequisites

- Python 3.7 or higher
- A running MCP server (see [server setup](../README.md))

### Installation

Using Makefile (recommended):

```bash
# One-step setup - creates virtual environment and installs dependencies
make

# Activate the virtual environment
source .venv/bin/activate  # On Windows: .venv\Scripts\activate
```

### Running Examples

Using the Makefile:

```bash
# Run the browser-like client example
make run-browser

# Run the example with debug logging enabled
make run-browser ARGS="--debug" 

# Run the integrations example
make run-integrations

# Run any other example
make run-example EXAMPLE=examples/check_server.py
```

### Testing

Using the Makefile:

```bash
# Run all unit tests
make unit-test

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

The `Makefile` provides several useful commands:

```
make                 # Set up development environment
make setup           # Create virtual environment and install dependencies
make run-browser     # Run the browser-like client example
make run-integrations # Run the integrations example
make run-example     # Run a specific example file
make unit-test       # Run unit tests
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
```

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
│   ├── integration/           # Integration tests
│   │   ├── test_browser_like_client.py
│   │   └── test_server_connection.py
│   └── utils/                 # Test utilities
│       ├── server_runner.py   # Server control for tests
│       └── test_config.py     # Test configuration
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

3. Verify your session ID is valid:
   ```bash
   curl "http://localhost:8080/debug?sessionId=<your-session-id>"
   ```

For more help, see the [Troubleshooting Guide](../docs/TROUBLESHOOTING.md).

## License

This project is licensed under the Apache License 2.0 - see the LICENSE file for details. 