#!/usr/bin/env python3
"""
Test script for MCP server connection and tool availability.
"""

import asyncio
import logging
from ormcp.client import MCPClient

# Configure logging
logging.basicConfig(level=logging.DEBUG)
logger = logging.getLogger(__name__)

async def test_mcp_server():
    """Test MCP server connection and tool availability."""
    # Create client
    client = MCPClient("http://localhost:8080", auto_connect=True)
    
    try:
        # Initialize connection
        logger.info("Initializing MCP connection...")
        await client.initialize()
        
        # List available tools
        logger.info("Listing available tools...")
        tools = await client.list_tools()
        logger.info(f"Available tools: {tools}")
        
        # If integrations tool is available, try to list integrations
        if any(tool.get("name") == "integrations" for tool in tools):
            logger.info("Testing integrations tool...")
            result = await client.call_tool("integrations", {"action": "list"})
            logger.info(f"Integrations result: {result}")
        else:
            logger.warning("Integrations tool not found in available tools")
            
    except Exception as e:
        logger.error(f"Error during MCP server test: {str(e)}", exc_info=True)
    finally:
        # Close the connection
        await client.close()

if __name__ == "__main__":
    asyncio.run(test_mcp_server()) 