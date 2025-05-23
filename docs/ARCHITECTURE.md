# HPE OpsRamp MCP Architecture

This document describes the architecture of the HPE OpsRamp Model Context Protocol (MCP) implementation.

## High-Level Architecture

```
                 ┌───────────────────┐
                 │                   │
                 │    MCP Server     │
                 │                   │
                 └─────────┬─────────┘
                           │
                           │ HTTP/SSE/JSON-RPC
                           │
                 ┌─────────▼─────────┐
                 │                   │
                 │   Python Client   │
                 │                   │
                 └─────────┬─────────┘
                           │
                           │ Provides tools
                           │
              ┌────────────▼───────────┐
              │                        │
              │   AI Agent with        │
              │   Integration Expertise│
              │                        │
              └────────────────────────┘
```

The system consists of three main components:

1. **Server**: A Go server that implements the MCP specification using our forked version of mark3labs/mcp-go
2. **Client**: A Python client that handles the communication with the server
3. **AI Agent**: An intelligent agent that uses LLMs (either OpenAI or Anthropic) to process natural language queries about OpsRamp integrations

## Server Architecture

### Core Components

The server follows a layered architecture:

1. **HTTP Layer**: Handles HTTP requests, SSE connections, and routing
2. **MCP Protocol Layer**: Implements the MCP protocol specification (session management, message formatting)
3. **Tool Layer**: Tool implementation and registration
4. **OpsRamp API Layer**: Integration with OpsRamp APIs

### MCP-GO Fork

The server relies on a customized fork of the mark3labs/mcp-go library, with several enhancements:

1. **Enhanced SSE Transport**: Improved Server-Sent Events implementation for better connection stability
2. **Extended Client Capabilities**: Additional methods to access transport and capabilities
3. **Custom Testing Framework**: Robust in-process testing utilities
4. **Server Instrumentation**: Added hooks for monitoring and debugging

See [MCP_GO_FORK.md](./MCP_GO_FORK.md) for detailed information about the modifications.

### Tools

The server exposes the following tools:

- **integrations**: Manage OpsRamp integrations
  - List all integrations
  - Get integration details
  - Create/update/delete integrations
  - Enable/disable integrations
  - List integration types
  - Get details about specific integration types

## Client Architecture

The Python client is designed to be modular and extensible:

1. **Transport Layer**: Handles the HTTP and SSE communication
2. **Protocol Layer**: Manages JSON-RPC requests and responses
3. **Session Management**: Handles session creation, maintenance, and reconnection
4. **Tool Interface**: Provides a simple interface to call server tools

## AI Agent Architecture

The AI agent uses large language models to understand and respond to natural language queries about OpsRamp integrations:

```
┌──────────────────────┐     ┌──────────────────────┐     ┌──────────────────────┐
│                      │     │                      │     │                      │
│     User Query       │────►│     LLM Processing   │────►│    Tool Execution    │
│                      │     │                      │     │                      │
└──────────────────────┘     └──────────────────────┘     └──────────────────────┘
                                      │                             │
                                      │                             │
                                      │                             │
                                      ▼                             ▼
                             ┌──────────────────────┐     ┌──────────────────────┐
                             │                      │     │                      │
                             │    Response          │◄────│  Result Processing   │
                             │    Generation        │     │                      │
                             │                      │     │                      │
                             └──────────────────────┘     └──────────────────────┘
```

### OpsRamp Integrations Expertise

The agent's system prompt is specifically designed to provide deep expertise in OpsRamp integrations management:

1. **Comprehensive Tool Knowledge**: The agent understands all integration tool actions (list, get, getDetailed, create, update, delete, enable, disable, listTypes, getType)
2. **Parameter Understanding**: It knows which parameters are required for each action
3. **Integration Categorization**: It can categorize integrations by type, status, and purpose
4. **Lifecycle Expertise**: It understands the full integration lifecycle from discovery to retirement

### AI Agent Components

- **Agent Class**: Main coordinator handling conversation flow
- **System Prompt**: Specialized prompt focusing on integration expertise
- **Tool Handling**: Logic to select appropriate tool based on user request
- **Mock Integration Logic**: Sophisticated mocks for testing without MCP server

### Simple Mode

The agent supports a "simple mode" that doesn't require an MCP server connection:

1. Processes the same queries and commands
2. Returns pre-defined mock responses that simulate actual integration data
3. Follows the same conversational patterns

This mode is useful for:
- Development without a running server
- Testing integration expertise
- Demonstrations of capability

## Communication Flow

1. User sends a question to the AI agent
2. Agent uses LLM to determine which integration action is needed
3. Agent calls the appropriate tool via the Python client
4. Client sends the request to the server via HTTP/JSON-RPC
5. Server processes the request and returns results
6. Client receives the results and passes them back to the agent
7. Agent uses LLM to generate a natural language response
8. Response is presented to the user

## Testing Architecture

The system has a comprehensive testing framework:

1. **Server Tests**: Go unit and integration tests 
2. **Client Tests**: Python unit and integration tests
3. **Agent Tests**: Tests for the AI agent's integration expertise

The integration expertise testing uses:
- 37 diverse integration-related prompts
- Specialized mocks for testing without server connection
- Detailed validation of all response patterns

## Security Considerations

1. **Authentication**: OpsRamp API credentials stored securely
2. **API Keys**: LLM API keys managed via environment variables
3. **No Persistent Storage**: No user data or conversations stored
4. **Configuration Protection**: Sensitive configuration excluded from version control

## Deployment Model

The system is designed for flexible deployment:

1. **Development**: Local development with simple mode
2. **Testing**: CI/CD with mock integration tests
3. **Production**: Deployed with real OpsRamp tenant credentials 