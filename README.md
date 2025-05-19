# HPE OpsRamp MCP (Model Context Protocol) Server

A Go-based implementation of the MCP (Model Context Protocol) server for HPE OpsRamp with a Python client.

## Overview

This project provides the official Model Context Protocol implementation for HPE OpsRamp, consisting of:

1. A Go-based MCP server that exposes HPE OpsRamp integrations as tools
2. A Python client library for interacting with the server

## Prerequisites

Before you begin, you'll need:

1. **Go 1.18 or higher** - For building and running the server
2. **Python 3.7 or higher** - For running the client
3. **OpsRamp Credentials** - You'll need the following from your OpsRamp instance:
   - Tenant URL (e.g., `https://your-instance.opsramp.com`)
   - Auth URL (e.g., `https://your-instance.opsramp.com/tenancy/auth/oauth/token`)
   - Auth Key (API key for authentication)
   - Auth Secret (API secret for authentication)
   - Tenant ID (Your OpsRamp tenant identifier)

## Configuration

⚠️ **IMPORTANT:** Before running the server, you **MUST** create a valid configuration file with your OpsRamp credentials:

1. Copy the template configuration file to create your local config:
   ```bash
   cp config.yaml.template config.yaml
   ```

2. Edit `config.yaml` and replace ALL placeholder values with your actual OpsRamp credentials:
   ```yaml
   opsramp:
     tenant_url: "https://your-tenant-instance.opsramp.com"  # Replace with your actual tenant URL
     auth_url: "https://your-tenant-instance.opsramp.com/tenancy/auth/oauth/token"  # Replace with your actual auth URL
     auth_key: "YOUR_AUTH_KEY_HERE"  # Replace with your actual auth key
     auth_secret: "YOUR_AUTH_SECRET_HERE"  # Replace with your actual auth secret
     tenant_id: "YOUR_TENANT_ID_HERE"  # Replace with your actual tenant ID
   ```

3. Make sure your credentials are correct. The server will not function properly without valid OpsRamp credentials.

> **Security Note:** The `config.yaml` file is excluded from Git to protect sensitive information. Never commit this file to the repository. Each developer needs to create their own local copy with their credentials.

## Quick Start Guide

### Complete Setup with Single Makefile

The project is designed to be built, run, and tested using a single Makefile at the root level:

```bash
# Clone the repository
git clone https://github.com/opsramp/or-mcp-v2.git
cd or-mcp-v2

# REQUIRED: Create your configuration file with actual credentials
cp config.yaml.template config.yaml
# ⚠️ YOU MUST edit config.yaml with your valid OpsRamp credentials before proceeding ⚠️

# Clean any previous builds and build everything fresh
make clean-all
make all

# Kill any running server instances 
make kill-server

# Run the server and client tests together
make test-with-client

# OR run the server and connect with client examples
make run-debug           # Run server in debug mode in the background
make client-run-browser  # Run client browser example in another terminal
make client-run-integrations  # Run integrations example
make kill-server         # Shut down the server when done
```

### Unified Makefile Targets

The root Makefile now provides targets to:

1. Build and run the server
2. Run Python client tests
3. Run client examples
4. Clean up all artifacts

Use `make help` to see all available targets:

```bash
make help
```

### Individual Component Setup

If you prefer to work with the server and client separately:

#### 1. Server Setup and Running

```bash
# Build and start the server
make
make run

# Run with debug mode enabled
make run-debug

# Run quick health check
make health-check
```

#### 2. Python Client Setup

```bash
# Navigate to the Python client directory
cd client/python

# Setup the Python environment
make setup

# Run unit tests (no server required)
make unit-test

# Run the browser example (requires server running)
make run-browser
```

## Server Features

- HTTP-based communication
- SSE (Server-Sent Events) for real-time updates
- JSON-RPC 2.0 for message exchange
- RESTful endpoints for health and status monitoring
- Tools for managing HPE OpsRamp integrations

## Client Features

- Asynchronous API using Python's asyncio
- Browser-like SSE client for reliable connections
- Error handling and retries
- Event processing
- Tool discovery and invocation

## Requirements

- Go 1.18 or higher (for server)
- Python 3.7 or higher (for client)
- Access to HPE OpsRamp APIs (for production use)

## Server Endpoints

- `/sse` - SSE connection endpoint
- `/message` - JSON-RPC message endpoint
- `/health` - Health check endpoint
- `/readiness` - Readiness check endpoint
- `/debug` - Debug information endpoint (in debug mode)

## Configuration

### Server Configuration

The server can be configured via environment variables or through Makefile targets:

```bash
# Change the port when running directly with make
PORT=8090 make run

# Run in debug mode
DEBUG=true make run

# Run with debug logging
LOG_LEVEL=debug make run
```

### Client Configuration

The client can be configured when running examples with Makefile:

```bash
# Run with additional arguments
make client-run-browser ARGS="--debug --server-url=http://localhost:8090"

# Run client integration tests with MCP server
MCP_INTEGRATION_TEST=1 make client-test
```

## MCP Protocol Implementation

This project implements the Model Context Protocol, which defines:

1. A session-based communication model
2. JSON-RPC 2.0 for request/response
3. SSE for server-to-client events
4. Tool discovery and invocation mechanisms

## Tools

Currently implemented tools:

- `integrations` - Manage HPE OpsRamp integrations
  - List integrations
  - Get integration details
  - Create/update integrations
  - List integration types

## Documentation

For more detailed information, see:

- [Architecture Documentation](./docs/ARCHITECTURE.md)
- [Testing Guide](./docs/TESTING.md)
- [Troubleshooting Guide](./docs/TROUBLESHOOTING.md)
- [Python Client Documentation](./client/python/README.md)

## License

This project is licensed under the Apache License 2.0 - see the LICENSE file for details. 