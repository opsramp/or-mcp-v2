# HPE OpsRamp MCP Troubleshooting Guide

This document provides solutions to common issues you might encounter when working with the HPE OpsRamp MCP server and client.

## Server Issues

### Server Won't Start

**Symptoms:**
- Error message: `Failed to start server: listen tcp :8080: bind: address already in use`
- Server immediately exits after starting

**Solutions:**
1. Check if another process is using port 8080:
   ```bash
   lsof -i :8080
   # or
   netstat -tuln | grep 8080
   ```

2. Stop the existing process or use a different port:
   ```bash
   # Kill the process using port 8080
   kill <PID>
   
   # Or start the server on a different port
   PORT=8081 go run cmd/server/main.go
   ```

### Server Health Check Fails

**Symptoms:**
- `/health` endpoint returns non-200 status code
- Error in logs about initialization

**Solutions:**
1. Check the server logs for specific errors:
   ```bash
   tail -f output/logs/or-mcp.log
   ```

2. Verify the server has proper permissions to write to log directories:
   ```bash
   # Create required directories if they don't exist
   mkdir -p output/logs
   ```

3. Restart the server with debug mode enabled:
   ```bash
   DEBUG=true go run cmd/server/main.go
   ```

## Client Issues

### Session Validation Errors

**Symptoms:**
- `Invalid session ID` errors in the client
- JSON-RPC requests fail after establishing SSE connection

**Solutions:**
1. Use the BrowserLikeSSEClient implementation:
   ```python
   # Use the browser_like_example.py script as a reference
   # This script properly maintains the SSE connection
   python examples/browser_like_example.py
   ```

2. Check that the session is active on the server (in debug mode):
   ```bash
   # Start server with debug mode
   DEBUG=true go run cmd/server/main.go
   
   # Check debug info with your session ID
   curl "http://localhost:8080/debug?sessionId=<your-session-id>"
   ```

3. Verify that you're using the correct session ID from the SSE connection in your JSON-RPC requests.

### Connection Issues

**Symptoms:**
- Timeout when connecting to the server
- SSE connection closes unexpectedly

**Solutions:**
1. Check that the server is running:
   ```bash
   curl http://localhost:8080/health
   ```

2. Increase connection timeouts:
   ```python
   client = MCPClient("http://localhost:8080", connection_timeout=30)
   ```

3. Check network connectivity and firewalls.

### JSON-RPC Request Failures

**Symptoms:**
- Error: `Method not found`
- Error: `Invalid params`

**Solutions:**
1. Verify you're using the correct method names:
   - Use `initialize` for initialization
   - Use `tools/list` for listing tools
   - Use `tools/call` for calling tools

2. Check the parameters format:
   ```python
   # For initialization
   await client.initialize(client_name="my-client", client_version="1.0.0")
   
   # For tool calls
   await client.call_tool("integrations", {"action": "list"})
   ```

3. Enable debug logging to see the exact request/response:
   ```python
   logging.getLogger('ormcp').setLevel(logging.DEBUG)
   ```

## Testing Issues

### Integration Tests Fail

**Symptoms:**
- Tests timeout
- Session validation errors

**Solutions:**
1. Make sure the server is running before tests:
   ```bash
   # Start the server in debug mode
   DEBUG=true go run cmd/server/main.go
   ```

2. Use the automated test script:
   ```bash
   cd client/python
   ./run_tests.sh
   ```

3. Run tests with debug logging:
   ```bash
   DEBUG=true python -m pytest tests/integration/
   ```

### run_tests.sh Script Fails

**Symptoms:**
- "Port already in use" errors
- Tests fail to start

**Solutions:**
1. Check if the server is already running:
   ```bash
   lsof -i :8080
   ```

2. Kill any existing server processes:
   ```bash
   pkill -f "go run cmd/server/main.go"
   ```

3. Clean up any leftover PID files:
   ```bash
   rm -f client/python/.server.pid
   ```

## Tool-Specific Issues

### Integrations Tool Errors

**Symptoms:**
- `Error calling integrations tool` messages
- Empty responses from the integrations tool

**Solutions:**
1. Verify the tool is registered on the server:
   ```bash
   curl http://localhost:8080/health | grep -o '"tools":\[[^]]*\]'
   ```

2. Check that you're using the correct action:
   ```python
   # List integrations
   result = await client.call_tool("integrations", {"action": "list"})
   
   # Get a specific integration
   result = await client.call_tool("integrations", {"action": "get", "id": "int-001"})
   ```

3. Check the server logs for specific tool errors.

## Common Error Messages

### "Failed to connect to MCP server"

**Likely causes:**
- Server not running
- Network connectivity issues
- Wrong server URL

**Solutions:**
1. Start the server if it's not running
2. Check the server URL (default is http://localhost:8080)
3. Verify network connectivity

### "Server disconnected"

**Likely causes:**
- Server crashed or was stopped
- Network interruption
- SSE connection timeout

**Solutions:**
1. Check if the server is still running
2. Implement reconnection logic in your client code
3. Increase timeouts for long-running connections

### "Session close timed out"

**Likely causes:**
- Server is busy or unresponsive
- Connection already closed

**Solutions:**
1. This is often just a warning and can be ignored if the client operation completed
2. Increase the close timeout:
   ```python
   await client.close(timeout=10)  # 10 seconds timeout
   ```

## Getting More Help

If you continue to experience issues:

1. Enable debug logging on both server and client
2. Check the complete logs for both components
3. Try the example scripts to isolate the problem
4. Create a minimal reproduction of the issue 