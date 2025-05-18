"""
Unit tests for the OpsRamp MCP Python client.
"""

import unittest
import json
import asyncio
from unittest.mock import patch, MagicMock, AsyncMock

import sys
import os

# Add the parent directory to the Python path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from src.ormcp.client import MCPClient, SyncMCPClient
from src.ormcp.exceptions import MCPError, ConnectionError, ToolError, JSONRPCError, SessionError
from src.ormcp.utils import create_jsonrpc_request, parse_jsonrpc_response, parse_session_id_from_sse


class TestUtils(unittest.TestCase):
    """Test the utility functions."""
    
    def test_create_jsonrpc_request(self):
        """Test creating a JSON-RPC request."""
        # Without params
        request = create_jsonrpc_request("test_method", request_id="123")
        self.assertEqual(request["jsonrpc"], "2.0")
        self.assertEqual(request["id"], "123")
        self.assertEqual(request["method"], "test_method")
        self.assertNotIn("params", request)
        
        # With params
        request = create_jsonrpc_request("test_method", {"foo": "bar"}, "123")
        self.assertEqual(request["jsonrpc"], "2.0")
        self.assertEqual(request["id"], "123")
        self.assertEqual(request["method"], "test_method")
        self.assertEqual(request["params"], {"foo": "bar"})
    
    def test_parse_jsonrpc_response(self):
        """Test parsing a JSON-RPC response."""
        response_text = json.dumps({
            "jsonrpc": "2.0",
            "id": "123",
            "result": {"foo": "bar"}
        })
        response = parse_jsonrpc_response(response_text)
        self.assertEqual(response["jsonrpc"], "2.0")
        self.assertEqual(response["id"], "123")
        self.assertEqual(response["result"], {"foo": "bar"})
        
        # Test error response
        response_text = json.dumps({
            "jsonrpc": "2.0",
            "id": "123",
            "error": {"code": -32601, "message": "Method not found"}
        })
        response = parse_jsonrpc_response(response_text)
        self.assertEqual(response["jsonrpc"], "2.0")
        self.assertEqual(response["id"], "123")
        self.assertEqual(response["error"]["code"], -32601)
        self.assertEqual(response["error"]["message"], "Method not found")
        
        # Test invalid JSON
        with self.assertRaises(ValueError):
            parse_jsonrpc_response("not json")
    
    def test_parse_session_id_from_sse(self):
        """Test parsing a session ID from an SSE event."""
        event_data = "/message?sessionId=abc123"
        session_id = parse_session_id_from_sse(event_data)
        self.assertEqual(session_id, "abc123")
        
        # Test no session ID
        event_data = "no session id here"
        session_id = parse_session_id_from_sse(event_data)
        self.assertIsNone(session_id)


class TestMCPClient(unittest.TestCase):
    """Test the MCP client."""
    
    def setUp(self):
        """Set up the test case."""
        # Patch the MCPSession class
        self.session_patcher = patch('src.ormcp.client.MCPSession')
        self.mock_session_class = self.session_patcher.start()
        self.mock_session = self.mock_session_class.return_value
        
        # Set default return values
        self.mock_session.is_connected = True
        self.mock_session.connect.return_value = "fake-session-id"
        
        # Create the client
        self.client = MCPClient("http://example.com", auto_connect=False)
    
    def tearDown(self):
        """Clean up after the test case."""
        self.session_patcher.stop()
    
    def test_connect(self):
        """Test connecting to the MCP server."""
        session_id = self.client.connect()
        self.assertEqual(session_id, "fake-session-id")
        self.mock_session.connect.assert_called_once()
        
        # Test connection error
        self.mock_session.connect.side_effect = Exception("Connection error")
        with self.assertRaises(ConnectionError):
            self.client.connect()

    def test_initialize(self):
        """Test initializing the client."""
        # Mock the async send_request method
        async def mock_send_request(*args, **kwargs):
            return {"result": {"version": "1.0.0"}}
        self.mock_session.send_request = AsyncMock(side_effect=mock_send_request)
        
        # Run the test with an async wrapper
        async def test_init():
            result = await self.client.initialize()
            self.assertEqual(result, {"version": "1.0.0"})
            self.assertTrue(self.client._initialized)
            self.mock_session.send_request.assert_called_with(
                "initialize", 
                {"client": {"name": "python-client", "version": "1.0.0"}},
                timeout=30
            )
            
            # Test initialization with custom name and version
            self.client._initialized = False
            result = await self.client.initialize("test-name", "2.0.0")
            self.mock_session.send_request.assert_called_with(
                "initialize", 
                {"client": {"name": "test-name", "version": "2.0.0"}},
                timeout=30
            )
            
            # Test initialization error
            self.client._initialized = False
            self.mock_session.send_request = AsyncMock(side_effect=Exception("Init error"))
            with self.assertRaises(ConnectionError):
                await self.client.initialize()
                
        asyncio.run(test_init())

    def test_list_tools(self):
        """Test listing tools."""
        # Mock the async send_request method
        mock_tools = [{"name": "tool1"}, {"name": "tool2"}]
        async def mock_send_request(*args, **kwargs):
            return {"result": mock_tools}
        self.mock_session.send_request = AsyncMock(side_effect=mock_send_request)
        
        # Run the test with an async wrapper
        async def test_list():
            # Set up as initialized
            self.client._initialized = True
            
            tools = await self.client.list_tools()
            self.assertEqual(tools, mock_tools)
            self.assertEqual(self.client._available_tools, mock_tools)
            
            # Test not initialized error
            self.client._initialized = False
            with self.assertRaises(MCPError):
                await self.client.list_tools()
                
            # Test error from server
            self.client._initialized = True
            self.mock_session.send_request = AsyncMock(side_effect=Exception("List error"))
            with self.assertRaises(MCPError):
                await self.client.list_tools()
                
        asyncio.run(test_list())

    def test_call_tool(self):
        """Test calling a tool."""
        # Mock the async send_request method
        mock_result = {"status": "success", "data": {"value": 42}}
        async def mock_send_request(*args, **kwargs):
            method, params = args
            if method == "tools/call" and params.get("name") == "test_tool":
                return {"result": mock_result}
            raise Exception(f"Unexpected call: {method} - {params}")
            
        self.mock_session.send_request = AsyncMock(side_effect=mock_send_request)
        
        # Run the test with an async wrapper
        async def test_tool_call():
            # Set up as initialized
            self.client._initialized = True
            
            result = await self.client.call_tool("test_tool", {"param": "value"})
            self.assertEqual(result, mock_result)
            self.mock_session.send_request.assert_called_with(
                "tools/call", 
                {"name": "test_tool", "arguments": {"param": "value"}},
                timeout=60
            )
            
            # Test tool error
            self.mock_session.send_request = AsyncMock(side_effect=Exception("Tool error"))
            with self.assertRaises(ToolError):
                await self.client.call_tool("test_tool", {"param": "value"})
                
        asyncio.run(test_tool_call())
    
    def test_sync_client(self):
        """Test the synchronous client wrapper."""
        # We'll mock the async client's methods that the sync client wraps
        with patch('src.ormcp.client.MCPClient') as mock_async_client_class:
            mock_async_client = mock_async_client_class.return_value
            mock_async_client.connect.return_value = "fake-session-id"
            
            # Mock async methods with AsyncMock to help with later awaits
            mock_async_client.initialize = AsyncMock(return_value={"version": "1.0.0"})
            mock_async_client.list_tools = AsyncMock(return_value=[{"name": "tool1"}])
            mock_async_client.call_tool = AsyncMock(return_value={"result": "value"})
            mock_async_client.close = AsyncMock()
            
            # Create a sync client
            sync_client = SyncMCPClient("http://example.com", auto_connect=False)
            
            # Test connect
            session_id = sync_client.connect()
            self.assertEqual(session_id, "fake-session-id")
            mock_async_client.connect.assert_called_once()
            
            # Test other methods
            result = sync_client.initialize()
            self.assertEqual(result, {"version": "1.0.0"})
            
            result = sync_client.list_tools()
            self.assertEqual(result, [{"name": "tool1"}])
            
            result = sync_client.call_tool("test_tool", {"arg": "value"})
            self.assertEqual(result, {"result": "value"})
            
            # Test close
            sync_client.close()
            mock_async_client.close.assert_called_once()

    def test_close(self):
        """Test closing the client connection."""
        # Run the test with an async wrapper
        async def test_close():
            await self.client.close()
            self.mock_session.close.assert_called_once()
            self.assertFalse(self.client._initialized)
                
        asyncio.run(test_close())


if __name__ == '__main__':
    unittest.main() 