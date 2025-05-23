# HPE OpsRamp MCP Documentation

This directory contains documentation for the HPE OpsRamp MCP (Model Context Protocol) implementation.

## Contents

| Document               | Description                                             |
|------------------------|---------------------------------------------------------|
| [ARCHITECTURE.md](./ARCHITECTURE.md) | System architecture and component overview |
| [TESTING.md](./TESTING.md) | Testing strategy and procedures |
| [TROUBLESHOOTING.md](./TROUBLESHOOTING.md) | Common issues and solutions |
| [MCP_GO_FORK.md](./MCP_GO_FORK.md) | Information about our fork of mark3labs/mcp-go |

## Related Documentation

Additional documentation can be found in specific component directories:

- [Agent Documentation](../client/agent/README.md) - Documentation for the AI agent with integrations expertise
- [Python Client](../client/python/README.md) - Documentation for the Python client library

## Quick Links

- [Main README](../README.md) - Project overview and quick start guide
- [Client Makefile](../client/Makefile) - Client-specific build and run targets
- [Server Setup](../README.md#server-setup-and-running) - Server setup instructions
- [Agent Testing](./TESTING.md#ai-agent-testing) - Testing the agent's integration expertise

## Integration Expert

The AI agent in this project is specialized in OpsRamp integrations management, with expertise in:

1. Listing and categorizing available integrations
2. Providing detailed information about specific integrations
3. Explaining integration operations (create, update, delete, enable, disable)
4. Filtering integrations by type, category, or status
5. Describing integration types and capabilities

For more information on integration testing, see:
- [Integration Testing](./TESTING.md#integration-expertise-testing)
- [Sample Prompts](../client/agent/examples/sample_prompts.txt)
