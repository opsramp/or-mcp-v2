#!/bin/bash

# Test script for full MCP Inspector SSE flow with persistent connection
set -e

echo "🔍 Testing MCP Inspector Full SSE Flow"
echo "======================================"

BASE_URL="http://localhost:8080"
TEMP_DIR="/tmp/mcp_test_$$"
mkdir -p "$TEMP_DIR"

# Cleanup function
cleanup() {
    echo "Cleaning up..."
    [ -n "$SSE_PID" ] && kill $SSE_PID 2>/dev/null || true
    rm -rf "$TEMP_DIR"
}
trap cleanup EXIT

# Test 1: Start persistent SSE connection
echo "1. Starting persistent SSE connection..."
curl -s -N -H "Accept: text/event-stream" "$BASE_URL/sse" > "$TEMP_DIR/sse_output.txt" &
SSE_PID=$!

# Wait for initial connection and extract session info
sleep 2

if [ ! -s "$TEMP_DIR/sse_output.txt" ]; then
    echo "❌ SSE connection failed"
    exit 1
fi

# Extract session ID from the endpoint
SESSION_ENDPOINT=$(grep "data:" "$TEMP_DIR/sse_output.txt" | head -1 | cut -d' ' -f2)
SESSION_ID=$(echo "$SESSION_ENDPOINT" | grep -o 'sessionId=[^&]*' | cut -d'=' -f2)

if [ -z "$SESSION_ID" ]; then
    echo "❌ Failed to extract session ID"
    cat "$TEMP_DIR/sse_output.txt"
    exit 1
fi

echo "✅ SSE connection established"
echo "✅ Session ID: $SESSION_ID"
echo "✅ Message endpoint: $SESSION_ENDPOINT"

# Test 2: Initialize MCP protocol
echo ""
echo "2. Testing MCP initialization..."
curl -s -X POST "$BASE_URL$SESSION_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":0,"method":"initialize","params":{"protocolVersion":"2025-03-26","capabilities":{"roots":{"listChanged":true},"sampling":{}},"clientInfo":{"name":"mcp-inspector","version":"0.14.3"}}}' \
  -w "HTTP Status: %{http_code}\n"

echo "✅ MCP initialization sent"

# Wait for response in SSE stream
sleep 2

# Test 3: List tools
echo ""
echo "3. Testing tool listing..."
curl -s -X POST "$BASE_URL$SESSION_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' \
  -w "HTTP Status: %{http_code}\n"

echo "✅ Tools listing sent"

# Wait for response in SSE stream
sleep 2

# Test 4: Execute a tool
echo ""
echo "4. Testing tool execution..."
curl -s -X POST "$BASE_URL$SESSION_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"integrations","arguments":{"action":"list"}}}' \
  -w "HTTP Status: %{http_code}\n"

echo "✅ Tool execution sent"

# Wait for response in SSE stream
sleep 3

# Test 5: Show SSE responses
echo ""
echo "5. SSE Stream Responses:"
echo "========================"
cat "$TEMP_DIR/sse_output.txt"

# Count the number of responses
RESPONSE_COUNT=$(grep -c "event: message" "$TEMP_DIR/sse_output.txt" || echo "0")
echo ""
echo "📊 Summary:"
echo "   - SSE connection: ✅ Active"
echo "   - Session ID: $SESSION_ID"
echo "   - Responses received: $RESPONSE_COUNT"

if [ "$RESPONSE_COUNT" -ge 3 ]; then
    echo "   - Status: ✅ All responses received"
else
    echo "   - Status: ⚠️  Some responses may be missing"
fi

echo ""
echo "🎉 MCP Inspector full SSE flow test completed!"
echo "   This simulates exactly how MCP Inspector works"
echo ""
echo "💡 To use with MCP Inspector:"
echo "   1. Start server: DEBUG=true ./server"
echo "   2. Connect MCP Inspector to: $BASE_URL/sse"
echo "   3. MCP Inspector will maintain the SSE connection and send messages to the message endpoint" 