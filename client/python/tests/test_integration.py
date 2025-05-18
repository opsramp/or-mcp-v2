"""
Integration tests for the OpsRamp MCP Python client.

These tests require a running MCP server on localhost:8080.
To run the tests:
1. Start the MCP server: `go run cmd/server/main.go`
2. Run the tests: `pytest -xvs client/python/tests/test_integration.py`
"""

import os
import sys
import time
import pytest
import asyncio
from unittest import mock

# Add the parent directory to the Python path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from src.ormcp.client import MCPClient, SyncMCPClient
from src.ormcp.exceptions import MCPError


# Skip the tests if MCP_INTEGRATION_TEST is not set
pytestmark = pytest.mark.skipif(
    not os.environ.get('MCP_INTEGRATION_TEST'),
    reason="Set MCP_INTEGRATION_TEST=1 to run integration tests"
)

# Server URL to use for testing
SERVER_URL = os.environ.get('MCP_SERVER_URL', 'http://localhost:8080')


@pytest.fixture
async def async_client():
    """Fixture to provide an async client for testing."""
    client = MCPClient(SERVER_URL, auto_connect=False)
    yield client
    await client.close()


@pytest.fixture
def sync_client():
    """Fixture to provide a sync client for testing."""
    client = SyncMCPClient(SERVER_URL, auto_connect=False)
    yield client
    client.close()


@pytest.mark.asyncio
async def test_client_connect(async_client):
    """Test connecting to the server."""
    session_id = async_client.connect()
    assert session_id is not None
    assert len(session_id) > 0


@pytest.mark.asyncio
async def test_client_initialize(async_client):
    """Test initializing the connection."""
    # Connect first
    async_client.connect()
    
    # Initialize
    result = await async_client.initialize(
        client_name="integration-test", 
        client_version="1.0.0"
    )
    
    # Verify result
    assert result is not None


@pytest.mark.asyncio
async def test_list_tools(async_client):
    """Test listing available tools."""
    # Connect and initialize
    async_client.connect()
    await async_client.initialize()
    
    # List tools
    tools = await async_client.list_tools()
    
    # Verify tools
    assert isinstance(tools, list)
    assert len(tools) > 0
    
    # Check that each tool is either a string or an object with a name
    for tool in tools:
        if isinstance(tool, str):
            assert len(tool) > 0
        else:
            assert "name" in tool
            assert isinstance(tool["name"], str)
            assert len(tool["name"]) > 0


@pytest.mark.asyncio
async def test_call_integrations_tool(async_client):
    """Test calling the integrations tool."""
    # Connect and initialize
    async_client.connect()
    await async_client.initialize()
    
    # Call the integrations tool directly 
    # We know it's registered on the server from /health endpoint
    result = await async_client.call_tool("integrations", {"action": "list"})
    
    # Verify result
    assert result is not None
    
    # The result is wrapped in a content object with text that contains JSON
    if isinstance(result, dict) and "content" in result:
        content = result["content"]
        assert isinstance(content, list)
        assert len(content) > 0
        
        # Extract and parse the JSON text if available
        text_content = next((item for item in content if item.get("type") == "text"), None)
        if text_content and "text" in text_content:
            import json
            try:
                integrations = json.loads(text_content["text"])
                assert isinstance(integrations, list)
            except json.JSONDecodeError:
                # If it's not valid JSON, that's okay for this test
                pass


def test_sync_client_operations(sync_client):
    """Test synchronous client operations."""
    # Connect
    session_id = sync_client.connect()
    assert session_id is not None
    
    # Initialize
    result = sync_client.initialize(
        client_name="sync-integration-test",
        client_version="1.0.0"
    )
    assert result is not None
    
    # List tools
    tools = sync_client.list_tools()
    assert isinstance(tools, list)
    assert len(tools) > 0
    
    # Call integrations tool directly
    result = sync_client.call_tool("integrations", {"action": "list"})
    assert result is not None
    
    # Verify result format
    if isinstance(result, dict) and "content" in result:
        content = result["content"]
        assert isinstance(content, list)
        assert len(content) > 0


@pytest.mark.asyncio
async def test_client_reconnect(async_client):
    """Test client reconnection capability."""
    # Connect and initialize
    async_client.connect()
    await async_client.initialize()
    
    # Close connection
    await async_client.close()
    
    # Reconnect
    async_client.connect()
    result = await async_client.initialize()
    assert result is not None


@pytest.mark.asyncio
async def test_error_handling(async_client):
    """Test error handling for invalid tool calls."""
    # Connect and initialize
    async_client.connect()
    await async_client.initialize()
    
    # Try to call a non-existent tool
    with pytest.raises(MCPError):
        await async_client.call_tool("non_existent_tool", {})


if __name__ == "__main__":
    # For manual testing
    async def run_tests():
        client = MCPClient(SERVER_URL)
        try:
            client.connect()
            await client.initialize()
            tools = await client.list_tools()
            print(f"Available tools: {tools}")
        finally:
            await client.close()
    
    asyncio.run(run_tests()) 