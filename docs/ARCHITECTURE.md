# HPE OpsRamp MCP Architecture

This document outlines the architecture of the HPE OpsRamp MCP system, consisting of a Go-based server and a Python client.

## System Overview

The HPE OpsRamp MCP system implements the Model Context Protocol, a communication protocol designed for AI agents to interact with tools. The system has two main components:

1. **MCP Server** - A Go application that exposes tools via HTTP endpoints
2. **MCP Client** - A Python library for connecting to and interacting with the server

```
┌───────────────┐     HTTP/SSE     ┌────────────────┐
│  Python       │  <-------------> │  Go MCP        │
│  MCP Client   │  JSON-RPC 2.0    │  Server        │
└───────────────┘                  └────────────────┘
```

## Server Architecture

The server is built on top of the [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) library and follows a modular design:

```
┌─────────────────────────────────────────────────┐
│                   HTTP Server                    │
├─────────────┬─────────────────────┬─────────────┤
│  SSE Server │  JSON-RPC Handlers  │  Health API │
├─────────────┴─────────────────────┴─────────────┤
│                   MCP Server                     │
├─────────────────────────────────────────────────┤
│                   Tool Registry                  │
└─────────────────────────────────────────────────┘
```

Key components:

1. **HTTP Server** - Handles incoming HTTP requests
2. **SSE Server** - Manages Server-Sent Events connections for real-time communication
3. **JSON-RPC Handlers** - Process JSON-RPC 2.0 requests
4. **Health API** - Provides health and readiness endpoints
5. **MCP Server** - Core implementation of the MCP protocol
6. **Tool Registry** - Manages the available tools and their handlers

## Client Architecture

The Python client is designed with an asynchronous-first approach, with a synchronous wrapper for convenience:

```
┌─────────────────────────────────────────────────┐
│                   MCPClient                      │
├─────────────────────────────────────────────────┤
│                   MCPSession                     │
├─────────────────┬─────────────────┬─────────────┤
│ BrowserLikeSSE  │  JSON-RPC       │  Event      │
│ Client          │  Client         │  Handlers   │
└─────────────────┴─────────────────┴─────────────┘
```

Key components:

1. **MCPClient** - High-level client interface with methods for tool discovery and invocation
2. **MCPSession** - Manages the connection to the server and handles session state
3. **BrowserLikeSSEClient** - Maintains persistent SSE connection with browser-like behavior
4. **JSON-RPC Client** - Formats and sends JSON-RPC 2.0 requests
5. **Event Handlers** - Processes events from the SSE connection

## Communication Flow

1. **Connection Establishment**:
   - Client opens an SSE connection to the server via `/sse` endpoint
   - Server assigns a session ID and returns it in an `endpoint` event

2. **Initialization**:
   - Client sends an `initialize` request with client metadata
   - Server validates the session and responds with protocol capabilities

3. **Tool Discovery**:
   - Client requests available tools with `tools/list`
   - Server responds with list of registered tools

4. **Tool Invocation**:
   - Client calls a tool using `tools/call` with tool name and arguments
   - Server processes the tool request and returns the result

5. **Real-time Events**:
   - Server can send events to the client at any time via the SSE connection
   - Client processes events and makes them available to the application

## Implementation Details

### Server Implementation

The server is implemented in Go using:
- Standard library `net/http` for HTTP serving
- [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) for MCP protocol implementation
- Custom logging and configuration

### Client Implementation

The client is implemented in Python using:
- `aiohttp` for asynchronous HTTP requests
- `sseclient-py` as a base for the BrowserLikeSSEClient
- `asyncio` for asynchronous programming

## Session Management

Session management is a critical aspect of the MCP protocol:

1. The server creates unique session IDs for each client connection
2. These IDs must be used in subsequent JSON-RPC requests
3. The BrowserLikeSSEClient ensures the session remains active
4. Session validation is handled by the mark3labs/mcp-go library

## Deployment Architecture

For production deployments, the recommended architecture is:

```
┌──────────────┐      ┌──────────────┐      ┌──────────────┐
│  AI Agent    │ ---> │  MCP Server  │ ---> │  HPE OpsRamp API │
│  with MCP    │      │  (Go)        │      │  Services    │
│  Client      │ <--- │              │ <--- │              │
└──────────────┘      └──────────────┘      └──────────────┘
```

Where:
- The AI Agent uses the MCP Client to communicate with the server
- The MCP Server acts as a gateway to HPE OpsRamp services
- The server can be deployed behind a load balancer for scaling 