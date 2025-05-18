# HPE OpsRamp MCP (Model Context Protocol) Server and Client

This project implements the official MCP server and client for HPE OpsRamp using the mark3labs/mcp-go library and a custom Python client. The project is hosted at [github.com/opsramp/or-mcp-v2](https://github.com/opsramp/or-mcp-v2).

## Documentation Structure

- [Project README](../README.md) - Main project documentation and server setup
- [Python Client README](../client/python/README.md) - Python client documentation
- [Architecture](./ARCHITECTURE.md) - System architecture and component design
- [Testing Guide](./TESTING.md) - Comprehensive testing strategy
- [Troubleshooting](./TROUBLESHOOTING.md) - Common issues and solutions

## Project Structure

- `cmd/` — Server entrypoints (main.go)
- `pkg/` — Core packages (tools, utilities) 
- `internal/` — Server internal logic
- `common/` — Shared utilities (logging, configurations)
- `client/python/` — Python client implementation
- `tests/` — Server test helpers and integration tests
- `docs/` — Project documentation
