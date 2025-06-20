#!/bin/bash

# Test the complete MCP flow
echo "Testing MCP Server Flow..."

# Step 1: Connect to SSE endpoint and get session
echo "Step 1: Getting SSE session..."
SSE_RESPONSE=$(timeout 5 curl -s -N -H "Accept: text/event-stream" http://localhost:8080/sse | head -2)
echo "$SSE_RESPONSE"

# Extract session ID
SESSION_ID=$(echo "$SSE_RESPONSE" | grep "data:" | sed 's/.*sessionId=//' | tr -d '\r\n')
echo "Session ID: $SESSION_ID"

if [ -z "$SESSION_ID" ]; then
    echo "Error: Could not get session ID"
    exit 1
fi

# Step 2: Send initialize request
echo -e "\nStep 2: Sending initialize request..."
INIT_RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":0,"method":"initialize","params":{"protocolVersion":"2025-03-26","capabilities":{"sampling":{},"roots":{"listChanged":true}},"clientInfo":{"name":"test-client","version":"1.0"}}}' \
  "http://localhost:8080/message?sessionId=$SESSION_ID")

echo "Initialize response: $INIT_RESPONSE"

# Step 3: Send tools/list request
echo -e "\nStep 3: Sending tools/list request..."
TOOLS_RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' \
  "http://localhost:8080/message?sessionId=$SESSION_ID")

echo "Tools response: $TOOLS_RESPONSE"

echo -e "\nMCP Flow test completed!"
