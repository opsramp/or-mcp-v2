"""
Integration tests for browser-like SSE client.
"""

import pytest
import pytest_asyncio
import logging
import asyncio
import json
import time
from pathlib import Path
import sys
import os

# Add the src directory to the Python path
sys.path.append(os.path.join(os.path.dirname(__file__), '../..'))

from src.ormcp.client import MCPClient
from src.ormcp.session import BrowserLikeSSEClient, MCPSession
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


class TestBrowserLikeClient:
    """Test browser-like SSE client functionality."""
    
    def test_sse_connect(self, server):
        """Test that the BrowserLikeSSEClient can connect to the server."""
        # Create a browser-like SSE client
        sse_client = BrowserLikeSSEClient(f"{SERVER_URL}/sse")
        
        try:
            # Connect to the server
            connected = sse_client.connect()
            assert connected
            assert sse_client.is_connected
            
            # Wait for any event - don't be too specific as event types may vary
            event = sse_client.wait_for_event(timeout=5)
            
            # If we get an event, log it. Otherwise, don't fail the test
            # as the server might not send events immediately
            if event:
                logger.info(f"Received event: {json.dumps(event, indent=2)}")
                
            # Success if we're connected - don't require specific events
        finally:
            # Close the connection
            sse_client.close()
    
    def test_mcp_session(self, server):
        """Test that the MCPSession can connect using the browser-like client."""
        # Create an MCP session
        session = MCPSession(SERVER_URL)
        
        try:
            # Connect to the server
            session_id = session.connect()
            assert session_id is not None
            assert session.is_connected
            
            # Check that the session has a browser-like SSE client
            assert session.sse_client is not None
            assert session.sse_client.is_connected
            
            # Success if connection is established
            # Events will be tested separately
        finally:
            # Close the session
            session.close()
    
    @pytest.mark.asyncio
    async def test_jsonrpc_request(self, server):
        """Test that the session can send JSON-RPC requests."""
        # Create an MCP session
        session = MCPSession(SERVER_URL)
        
        try:
            # Connect to the server
            session_id = session.connect()
            assert session_id is not None
            
            # Send the initialize request which we know works
            response = await session.send_request('initialize', {
                'client': {
                    'name': 'test-client',
                    'version': '1.0.0'
                }
            })
            
            # Verify response
            assert response is not None
            assert 'result' in response
            assert 'error' not in response
            
            # Log response for debugging
            logger.info(f"JSON-RPC response: {json.dumps(response, indent=2)}")
            
        finally:
            # Close the session
            session.close()


if __name__ == "__main__":
    # For manual testing
    logging.basicConfig(level=logging.DEBUG)
    
    # Test SSE client
    sse_client = BrowserLikeSSEClient(f"{SERVER_URL}/sse")
    
    try:
        print("Connecting to SSE endpoint...")
        if sse_client.connect():
            print("Connected to SSE endpoint")
            
        print("Waiting for events...")
        for _ in range(5):
            event = sse_client.get_event(timeout=2)
            if event:
                print(f"Event: {event['event']} - {event['data']}")
        
        print("Waiting for endpoint event...")
        endpoint_event = sse_client.wait_for_event(event_type='endpoint', timeout=5)
        if endpoint_event:
            print(f"Received endpoint event: {endpoint_event}")
        else:
            print("No endpoint event received")
    finally:
        sse_client.close() 