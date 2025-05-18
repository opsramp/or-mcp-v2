#!/usr/bin/env python
"""
Example script demonstrating how to use the OpsRamp MCP Python client
to list integrations from an MCP server.
"""

import asyncio
import sys
import os
import logging
import json

# Add the parent directory to the Python path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from src.ormcp.client import MCPClient
from src.ormcp.exceptions import MCPError

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger('mcp_example')


async def main():
    # Server URL (default to localhost:8080)
    server_url = os.environ.get('MCP_SERVER_URL', 'http://localhost:8080')
    
    logger.info(f"Connecting to MCP server at {server_url}")
    
    # Create the client
    client = MCPClient(server_url)
    
    try:
        # Initialize the connection
        logger.info("Initializing connection...")
        await client.initialize(client_name="example-client", client_version="1.0.0")
        
        # List available tools
        logger.info("Listing available tools...")
        tools = await client.list_tools()
        logger.info(f"Available tools: {json.dumps(tools, indent=2)}")
        
        # Check if integrations tool is available
        for tool in tools:
            if tool.get("name") == "integrations":
                logger.info("Found integrations tool!")
                
                # List integrations
                logger.info("Listing integrations...")
                integrations = await client.call_tool("integrations", {"action": "list"})
                logger.info(f"Integrations: {json.dumps(integrations, indent=2)}")
                
                # List integration types
                logger.info("Listing integration types...")
                types = await client.call_tool("integrations", {"action": "listTypes"})
                logger.info(f"Integration types: {json.dumps(types, indent=2)}")
                
                break
        else:
            logger.warning("Integrations tool not found")
        
    except MCPError as e:
        logger.error(f"MCP error: {e}")
    except Exception as e:
        logger.error(f"Unexpected error: {e}")
    finally:
        # Close the connection
        logger.info("Closing connection...")
        await client.close()


if __name__ == "__main__":
    asyncio.run(main()) 