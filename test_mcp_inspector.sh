#!/bin/bash

# Test script for MCP Inspector compatibility
set -e

echo "🔍 Testing MCP Inspector Compatibility"
echo "======================================"

BASE_URL="http://localhost:8080"
SESSION_ID="test-mcp-inspector-$(date +%s)"

# Test 1: Health check
echo "1. Testing health endpoint..."
curl -s "$BASE_URL/health" | jq -r '.status' | grep -q "ok" && echo "✅ Health check passed" || echo "❌ Health check failed"

# Test 2: SSE endpoint connection
echo "2. Testing SSE endpoint..."
curl -s -N -H "Accept: text/event-stream" "$BASE_URL/sse" > /tmp/sse_test.txt &
SSE_PID=$!
sleep 2
kill $SSE_PID 2>/dev/null || true
if grep -q "event: endpoint" /tmp/sse_test.txt; then
    echo "✅ SSE endpoint working"
else
    echo "❌ SSE endpoint failed"
fi
rm -f /tmp/sse_test.txt

# Test 3: MCP Inspector message endpoint - Initialize
echo "3. Testing MCP Inspector initialization..."
INIT_RESPONSE=$(curl -s -X POST "$BASE_URL/message?sessionId=$SESSION_ID" \
  -H "Content-Type: application/json" \
  -H "Accept: text/event-stream" \
  -d '{"jsonrpc":"2.0","id":0,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{"roots":{"listChanged":true},"sampling":{}},"clientInfo":{"name":"MCP-Inspector","version":"1.0.0"}}}')

if echo "$INIT_RESPONSE" | grep -q '"result"'; then
    echo "✅ MCP Inspector initialization successful"
else
    echo "❌ MCP Inspector initialization failed"
    echo "Response: $INIT_RESPONSE"
fi

# Test 4: Send acknowledgment
echo "4. Testing MCP Inspector acknowledgment..."
ACK_RESPONSE=$(curl -s -X POST "$BASE_URL/message?sessionId=$SESSION_ID" \
  -H "Content-Type: application/json" \
  -H "Accept: text/event-stream" \
  -d '{"jsonrpc":"2.0","id":1,"result":{}}')

if echo "$ACK_RESPONSE" | grep -q '"method":"initialized"'; then
    echo "✅ MCP Inspector acknowledgment successful"
else
    echo "❌ MCP Inspector acknowledgment failed"
    echo "Response: $ACK_RESPONSE"
fi

# Test 5: List tools
echo "5. Testing tool listing..."
TOOLS_RESPONSE=$(curl -s -X POST "$BASE_URL/message?sessionId=$SESSION_ID" \
  -H "Content-Type: application/json" \
  -H "Accept: text/event-stream" \
  -d '{"jsonrpc":"2.0","id":2,"method":"tools/list"}')

if echo "$TOOLS_RESPONSE" | grep -q '"tools"' && echo "$TOOLS_RESPONSE" | grep -q '"integrations"' && echo "$TOOLS_RESPONSE" | grep -q '"resources"'; then
    echo "✅ Tool listing successful"
    echo "   Found tools: integrations, resources"
else
    echo "❌ Tool listing failed"
    echo "Response: $TOOLS_RESPONSE"
fi

# Test 6: Test tool execution
echo "6. Testing tool execution..."
TOOL_RESPONSE=$(curl -s -X POST "$BASE_URL/message?sessionId=$SESSION_ID" \
  -H "Content-Type: application/json" \
  -H "Accept: text/event-stream" \
  -d '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"integrations","arguments":{"action":"list"}}}')

if echo "$TOOL_RESPONSE" | grep -q '"result"'; then
    echo "✅ Tool execution successful"
else
    echo "❌ Tool execution failed"
    echo "Response: $TOOL_RESPONSE"
fi

echo ""
echo "🎉 MCP Inspector compatibility test completed!"
echo "   Your MCP Inspector should now work properly with:"
echo "   - SSE endpoint: $BASE_URL/sse"
echo "   - Message endpoint: $BASE_URL/message"
echo ""
echo "💡 To use with MCP Inspector:"
echo "   1. Start server: DEBUG=true ./server"
echo "   2. Connect MCP Inspector to: $BASE_URL/sse"
echo "   3. All tools should be discoverable and executable" 