# MCP-GO Fork Modifications

This document describes the fork of the [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) library that is used in this project, including the modifications made, reasons for forking, and how these changes benefit our implementation.

## Overview

The Model Context Protocol (MCP) specification is implemented in Go using the mark3labs/mcp-go library. However, to meet the specific requirements of the HPE OpsRamp MCP server, we have created a fork with several enhancements and modifications.

The fork is located in `internal/mcp-go` and is referenced in the project's `go.mod` file using the replace directive:

```go
replace github.com/mark3labs/mcp-go => ./internal/mcp-go
```

## Reasons for Forking

1. **Improved SSE Transport Layer**: The original library had limitations in its Server-Sent Events (SSE) implementation that affected connection stability and recovery.
2. **Extended Client/Server Capabilities**: We needed additional methods to expose transport interfaces and configuration options.
3. **Enhanced Testing Framework**: Our testing requirements necessitated additional utilities for in-process testing.
4. **Server Instrumentation**: We required additional hooks for monitoring and debugging server activity.

## Key Modifications

### 1. Enhanced SSE Transport Layer

The original SSE implementation had several limitations that we addressed:

```go
// Added improved connection recovery
func (s *SSE) reconnect(ctx context.Context) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    // Enhanced reconnection logic
    if s.closed.Load() {
        return fmt.Errorf("client is closed")
    }
    
    // Exponential backoff retry mechanism
    // ...
}
```

Additional improvements include:
- Robust error handling for network interruptions
- Automatic session recovery
- Improved message tracking and correlation
- More efficient event handling

### 2. Extended Client Interface

We added methods to provide better access to client internals:

```go
// GetTransport returns the transport interface.
func (c *Client) GetTransport() transport.Interface {
    return c.transport
}

// GetServerCapabilities returns the server capabilities.
func (c *Client) GetServerCapabilities() mcp.ServerCapabilities {
    return c.serverCapabilities
}

// GetClientCapabilities returns the client capabilities.
func (c *Client) GetClientCapabilities() mcp.ClientCapabilities {
    return c.clientCapabilities
}
```

These methods allow the OpsRamp implementation to better integrate with and configure the MCP system.

### 3. Custom Testing Framework

We added a more robust testing framework in `mcptest` package:

```go
// Server encapsulates an MCP server and manages resources like pipes and context.
type Server struct {
    name  string
    tools []server.ServerTool

    ctx    context.Context
    cancel func()

    serverReader *io.PipeReader
    serverWriter *io.PipeWriter
    clientReader *io.PipeReader
    clientWriter *io.PipeWriter

    logBuffer bytes.Buffer

    transport transport.Interface
    client    *client.Client

    wg sync.WaitGroup
}
```

This framework provides:
- In-process testing without network overhead
- Simplified tool testing
- Automated session management
- Better test isolation

### 4. Server Instrumentation

Added server hooks and instrumentation for monitoring and debugging:

```go
// Hooks contains hooks to be called at various points in the server execution.
type Hooks struct {
    BeforeRequest  BeforeRequestHook
    AfterRequest   AfterRequestHook
    OnError        ErrorHook
    
    // Method-specific hooks
    Initialize     InitializeHook
    Ping           PingHook
    // ...and more
}
```

These hooks enable:
- Request/response monitoring
- Error tracking
- Performance measurement
- Debugging support

## Benefits of the Fork

1. **Improved Reliability**: Better connection handling and recovery for SSE connections
2. **Extended Functionality**: Additional methods to configure and monitor the MCP system
3. **Better Testing**: Comprehensive testing framework for all components
4. **Enhanced Monitoring**: Server instrumentation for operational visibility

## Usage in the Project

The fork is used throughout the project for:

1. The main MCP server implementation
2. Tool registration and execution
3. Client-server communication
4. Testing infrastructure

The modified library maintains compatibility with the original MCP specification while providing the enhancements needed for the HPE OpsRamp implementation.

## Maintenance Strategy

We plan to:
1. Keep our fork synchronized with upstream changes when appropriate
2. Contribute non-OpsRamp-specific improvements back to the upstream project
3. Maintain a clear separation between core MCP functionality and OpsRamp-specific extensions

This approach allows us to benefit from community improvements while maintaining the specialized functionality required for HPE OpsRamp. 