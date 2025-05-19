"""
OpsRamp MCP client for interacting with an MCP server.
"""

import asyncio
import logging
from typing import Dict, List, Any, Optional, Union

from .exceptions import MCPError, ConnectionError, ToolError
from .session import MCPSession

logger = logging.getLogger(__name__)


class MCPClient:
    """
    Client for interacting with an OpsRamp MCP server.
    
    This class provides a high-level API for:
    - Connecting to an MCP server
    - Initializing the connection
    - Discovering available tools
    - Calling tools with arguments
    - Managing the session
    """
    
    def __init__(self, server_url: str, auto_connect: bool = True, connection_timeout: int = 10):
        """
        Initialize the MCP client.
        
        Args:
            server_url: The URL of the MCP server
            auto_connect: Whether to automatically connect to the server
            connection_timeout: Timeout in seconds for connections
        """
        self.server_url = server_url
        self.connection_timeout = connection_timeout
        self.session = MCPSession(server_url, connection_timeout=connection_timeout)
        self._initialized = False
        self._available_tools = []
        
        if auto_connect:
            self.connect()
    
    def connect(self) -> str:
        """
        Connect to the MCP server.
        
        Returns:
            The session ID
            
        Raises:
            ConnectionError: If connection fails
        """
        try:
            logger.debug(f"Connecting to MCP server at {self.server_url}")
            return self.session.connect()
        except Exception as e:
            logger.error(f"Failed to connect to MCP server: {str(e)}", exc_info=True)
            raise ConnectionError(f"Failed to connect to MCP server: {str(e)}")
    
    async def initialize(self, client_name: str = "python-client", client_version: str = "1.0.0", timeout: int = 30) -> Dict[str, Any]:
        """
        Initialize the connection with the MCP server.
        
        Args:
            client_name: The name of the client
            client_version: The version of the client
            timeout: Timeout in seconds for the request
            
        Returns:
            The initialization response
            
        Raises:
            ConnectionError: If the initialization fails
        """
        if not self.session.is_connected:
            logger.error("Cannot initialize: Not connected")
            raise ConnectionError("Cannot initialize: Not connected")
        
        try:
            response = await self.session.send_request(
                "initialize",
                {
                    "client": {
                        "name": client_name,
                        "version": client_version
                    }
                },
                timeout=timeout
            )
            
            self._initialized = True
            self.session.is_initialized = True
            logger.info("MCP connection initialized")
            return response.get("result", {})
            
        except Exception as e:
            logger.error(f"Failed to initialize MCP connection: {str(e)}", exc_info=True)
            raise ConnectionError(f"Failed to initialize MCP connection: {str(e)}")
    
    async def list_tools(self, timeout: int = 30) -> List[Dict[str, Any]]:
        """
        List the available tools on the MCP server.
        
        Args:
            timeout: Timeout in seconds for the request
            
        Returns:
            A list of available tools
            
        Raises:
            MCPError: If not initialized or the request fails
        """
        self._ensure_initialized()
        
        try:
            logger.debug("Requesting list of available tools")
            response = await self.session.send_request("tools/list", {}, timeout=timeout)
            if isinstance(response, str):
                import json
                response = json.loads(response)
            
            # Handle both response formats:
            # 1. Direct list: [tool1, tool2, ...]
            # 2. Wrapped in result: {"result": {"tools": [tool1, tool2, ...]}}
            if isinstance(response, list):
                tools = response
            else:
                tools = response.get("result", {}).get("tools", [])
            
            self._available_tools = tools
            logger.debug(f"Received {len(tools)} available tools")
            return tools
            
        except Exception as e:
            logger.error(f"Failed to list tools: {str(e)}", exc_info=True)
            raise MCPError(f"Failed to list tools: {str(e)}")
    
    async def call_tool(self, tool_name: str, arguments: Dict[str, Any], timeout: int = 60) -> Any:
        """
        Call a tool on the MCP server.
        
        Args:
            tool_name: The name of the tool to call
            arguments: The arguments to pass to the tool
            timeout: Timeout in seconds for the request
            
        Returns:
            The result of the tool call
            
        Raises:
            ToolError: If the tool call fails
        """
        self._ensure_initialized()
        
        try:
            logger.debug(f"Calling tool '{tool_name}' with arguments: {arguments}")
            response = await self.session.send_request(
                "tools/call",
                {
                    "name": tool_name,
                    "arguments": arguments
                },
                timeout=timeout
            )
            
            if isinstance(response, str):
                import json
                response = json.loads(response)
            
            result = response.get("result")
            logger.debug(f"Tool '{tool_name}' call successful")
            return result
            
        except Exception as e:
            logger.error(f"Failed to call tool '{tool_name}': {str(e)}", exc_info=True)
            raise ToolError(f"Failed to call tool '{tool_name}': {str(e)}")
    
    def _ensure_initialized(self):
        """Ensure the client is initialized."""
        if not self._initialized:
            logger.error("Client not initialized. Call initialize() first.")
            raise MCPError("Client not initialized. Call initialize() first.")
    
    async def close(self, timeout: int = 5):
        """
        Close the connection to the MCP server.
        
        Args:
            timeout: Time to wait for graceful shutdown in seconds
        """
        logger.debug("Closing MCP client connection")
        
        # Create a task with a timeout to close the session
        try:
            close_task = asyncio.create_task(self._close_session())
            await asyncio.wait_for(close_task, timeout=timeout)
        except asyncio.TimeoutError:
            logger.warning(f"Session close timed out after {timeout}s")
        except Exception as e:
            logger.error(f"Error during session close: {str(e)}", exc_info=True)
        finally:
            self._initialized = False
            logger.info("MCP client closed")
    
    async def _close_session(self):
        """Helper method to close the session asynchronously."""
        # Since session.close is not async, run it in a thread
        loop = asyncio.get_event_loop()
        await loop.run_in_executor(None, self.session.close)


class SyncMCPClient:
    """
    Synchronous wrapper around the asynchronous MCP client.
    
    This class provides a synchronous API for applications that don't use asyncio.
    """
    
    def __init__(self, server_url: str, auto_connect: bool = True, connection_timeout: int = 10):
        """
        Initialize the synchronous MCP client.
        
        Args:
            server_url: The URL of the MCP server
            auto_connect: Whether to automatically connect to the server
            connection_timeout: Timeout in seconds for connections
        """
        self._async_client = MCPClient(server_url, auto_connect, connection_timeout)
        self._loop = asyncio.new_event_loop()
    
    def connect(self) -> str:
        """
        Connect to the MCP server.
        
        Returns:
            The session ID
        """
        return self._async_client.connect()
    
    def initialize(self, client_name: str = "python-client", client_version: str = "1.0.0", timeout: int = 30) -> Dict[str, Any]:
        """
        Initialize the connection with the MCP server.
        
        Args:
            client_name: The name of the client
            client_version: The version of the client
            timeout: Timeout in seconds for the request
            
        Returns:
            The initialization response
        """
        return self._run_async(self._async_client.initialize(client_name, client_version, timeout))
    
    def list_tools(self, timeout: int = 30) -> List[Dict[str, Any]]:
        """
        List the available tools on the MCP server.
        
        Args:
            timeout: Timeout in seconds for the request
            
        Returns:
            A list of available tools
        """
        return self._run_async(self._async_client.list_tools(timeout))
    
    def call_tool(self, tool_name: str, arguments: Dict[str, Any], timeout: int = 60) -> Any:
        """
        Call a tool on the MCP server.
        
        Args:
            tool_name: The name of the tool to call
            arguments: The arguments to pass to the tool
            timeout: Timeout in seconds for the request
            
        Returns:
            The result of the tool call
        """
        return self._run_async(self._async_client.call_tool(tool_name, arguments, timeout))
    
    def close(self, timeout: int = 5):
        """
        Close the connection to the MCP server.
        
        Args:
            timeout: Time to wait for graceful shutdown in seconds
        """
        try:
            self._run_async(self._async_client.close(timeout))
        finally:
            self._loop.close()
    
    def _run_async(self, coroutine):
        """Run an asynchronous coroutine in the event loop."""
        return self._loop.run_until_complete(coroutine) 