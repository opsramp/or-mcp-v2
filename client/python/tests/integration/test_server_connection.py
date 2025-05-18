"""
Integration tests for server connection.
"""

import pytest
import pytest_asyncio
import logging
import asyncio
import json
from pathlib import Path
import sys
import os

# Add the src directory to the Python path
sys.path.append(os.path.join(os.path.dirname(__file__), '../..'))

from src.ormcp.client import MCPClient
from src.ormcp.exceptions import SessionError, JSONRPCError

from ..utils.server_runner import ServerRunner
from ..utils.test_config import SERVER_URL, AUTO_START_SERVER, configure_logging

logger = logging.getLogger(__name__)

# Configure logging for tests
configure_logging()

@pytest.fixture(scope="module")
def server():
    """Start the server for tests if AUTO_START_SERVER is True."""
    if AUTO_START_SERVER:
        with ServerRunner() as runner:
            if not runner.start():
                pytest.skip("Failed to start server")
            yield runner
    else:
        # Use a context manager to ensure cleanup even if no server is started
        with ServerRunner() as runner:
            # Just check if the server is already running
            if not runner.check_health():
                pytest.skip("Server not running and AUTO_START_SERVER is False")
            yield runner


@pytest_asyncio.fixture
async def client():
    """Create an MCP client for testing."""
    client = MCPClient(SERVER_URL, auto_connect=False)
    yield client
    await client.close()


class TestServerConnection:
    """Test server connection functionality."""
    
    def test_server_health(self, server):
        """Test that the server health endpoint returns 200."""
        assert server.check_health()
    
    @pytest.mark.asyncio
    async def test_client_connect(self, client):
        """Test that the client can connect to the server."""
        # Connect to the server
        session_id = client.connect()
        assert session_id is not None
        assert client.session.is_connected
        
        # Check that we can initialize
        result = await client.initialize(client_name="test-client", client_version="1.0.0")
        assert result is not None
    
    @pytest.mark.asyncio
    async def test_list_tools(self, client):
        """Test that the client can list available tools."""
        # Connect and initialize
        client.connect()
        await client.initialize(client_name="test-client", client_version="1.0.0")
        
        # List tools
        result = await client.list_tools()
        
        # The tools can be returned in different formats depending on the server
        if isinstance(result, list):
            # Handle the case where result is a direct list of tools
            tools = result
            assert len(tools) > 0
        elif isinstance(result, dict) and 'tools' in result:
            # Handle the case where result is a dict with a 'tools' key
            tools = result['tools']
            assert isinstance(tools, list)
            assert len(tools) > 0
        else:
            # Unexpected format
            assert False, f"Unexpected tools format: {result}"
        
        # Log the tools for debugging
        logger.info(f"Available tools: {json.dumps(tools, indent=2)}")
    
    @pytest.mark.skipif(
        os.environ.get("MCP_INTEGRATION_TEST") != "1",
        reason="Set MCP_INTEGRATION_TEST=1 to run integration tests that require a server"
    )
    @pytest.mark.asyncio
    async def test_call_integrations_tool(self, client):
        """Test that the client can call the integrations tool."""
        # Connect and initialize
        client.connect()
        await client.initialize(client_name="test-client", client_version="1.0.0")
        
        # Call the integrations tool
        result = await client.call_tool("integrations", {"action": "list"})
        assert result is not None
        
        # Print result for debugging
        logger.info(f"Integrations tool result: {json.dumps(result, indent=2)}")


if __name__ == "__main__":
    # For manual testing
    logging.basicConfig(level=logging.DEBUG)
    
    async def test():
        client = MCPClient(SERVER_URL)
        try:
            # Connect and initialize
            client.connect()
            await client.initialize(client_name="test-client", client_version="1.0.0")
            
            # List tools
            tools = await client.list_tools()
            print(f"Available tools: {json.dumps(tools, indent=2)}")
            
            # Call the integrations tool
            result = await client.call_tool("integrations", {"action": "list"})
            print(f"Integrations tool result: {json.dumps(result, indent=2)}")
        finally:
            await client.close()
    
    asyncio.run(test()) 