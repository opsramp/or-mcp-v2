#!/bin/bash

# Test script for MCP Inspector SSE compatibility
set -e

echo "🔍 Testing MCP Inspector SSE Flow"
echo "================================="

BASE_URL="http://localhost:8080"

# Test 1: Connect to SSE endpoint and get message endpoint
echo "1. Connecting to SSE endpoint..."
SSE_OUTPUT=$(timeout 3 curl -s -N -H "Accept: text/event-stream" "$BASE_URL/sse" | head -2)
echo "SSE Response:"
echo "$SSE_OUTPUT"

# Extract session ID from the endpoint
SESSION_ENDPOINT=$(echo "$SSE_OUTPUT" | grep "data:" | cut -d' ' -f2)
SESSION_ID=$(echo "$SESSION_ENDPOINT" | grep -o 'sessionId=[^&]*' | cut -d'=' -f2)

if [ -z "$SESSION_ID" ]; then
    echo "❌ Failed to extract session ID"
    exit 1
fi

echo "✅ Extracted session ID: $SESSION_ID"
echo "✅ Message endpoint: $SESSION_ENDPOINT"

# Test 2: Initialize MCP protocol
echo ""
echo "2. Testing MCP initialization..."
INIT_RESPONSE=$(curl -s -X POST "$BASE_URL$SESSION_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":0,"method":"initialize","params":{"protocolVersion":"2025-03-26","capabilities":{"roots":{"listChanged":true},"sampling":{}},"clientInfo":{"name":"mcp-inspector","version":"0.14.3"}}}')

echo "Initialize response: $INIT_RESPONSE"

if echo "$INIT_RESPONSE" | grep -q "Accepted"; then
    echo "✅ MCP initialization sent (202 Accepted)"
else
    echo "❌ MCP initialization failed"
    echo "Response: $INIT_RESPONSE"
fi

# Test 3: List tools
echo ""
echo "3. Testing tool listing..."
TOOLS_RESPONSE=$(curl -s -X POST "$BASE_URL$SESSION_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list"}')

echo "Tools response: $TOOLS_RESPONSE"

if echo "$TOOLS_RESPONSE" | grep -q "Accepted"; then
    echo "✅ Tools listing sent (202 Accepted)"
else
    echo "❌ Tools listing failed"
    echo "Response: $TOOLS_RESPONSE"
fi

# Test 4: Execute a tool
echo ""
echo "4. Testing tool execution..."
EXEC_RESPONSE=$(curl -s -X POST "$BASE_URL$SESSION_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"integrations","arguments":{"action":"list"}}}')

echo "Tool execution response: $EXEC_RESPONSE"

if echo "$EXEC_RESPONSE" | grep -q "Accepted"; then
    echo "✅ Tool execution sent (202 Accepted)"
else
    echo "❌ Tool execution failed"
    echo "Response: $EXEC_RESPONSE"
fi

echo ""
echo "🎉 MCP Inspector SSE flow test completed!"
echo "   All requests should return 202 Accepted"
echo "   Actual responses are sent via SSE stream"
echo ""
echo "💡 To use with MCP Inspector:"
echo "   1. Start server: DEBUG=true ./server"
echo "   2. Connect MCP Inspector to: $BASE_URL/sse"
echo "   3. MCP Inspector will automatically handle the SSE flow" 