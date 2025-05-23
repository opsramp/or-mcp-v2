import logging
from typing import Any, Dict, List, Optional

logger = logging.getLogger(__name__)

class MCPError(Exception):
    """Base exception for MCP client errors."""
    pass

class MCPClient:
    def __init__(self, session):
        self.session = session
        self._available_tools = []
        self._initialized = False

    def _ensure_initialized(self):
        if not self._initialized:
            raise MCPError("Client not initialized")

    async def list_tools(self, timeout: int = 30) -> list:
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
            # Accept both {"result": {"tools": [...]}} and just a list
            if isinstance(response, list):
                tools = response
            elif isinstance(response, dict):
                tools = response.get("result", {}).get("tools", [])
            else:
                tools = []
            self._available_tools = tools
            logger.debug(f"Received {len(tools)} available tools")
            return tools
        except Exception as e:
            logger.error(f"Failed to list tools: {str(e)}", exc_info=True)
            raise MCPError(f"Failed to list tools: {str(e)}") 