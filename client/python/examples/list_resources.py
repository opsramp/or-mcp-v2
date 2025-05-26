#!/usr/bin/env python
"""
Example script demonstrating how to use the OpsRamp MCP Python client
to list resources from an MCP server.
"""

import asyncio
import sys
import os
import logging
import json
import time

# Add the parent directory to the Python path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from src.ormcp.client import MCPClient
from src.ormcp.exceptions import MCPError

# Configure logging
logging.basicConfig(
    level=logging.DEBUG,  # Set to DEBUG for more verbose output
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger('mcp_example')


async def main():
    # Server URL (default to localhost:8080)
    server_url = os.environ.get('MCP_SERVER_URL', 'http://localhost:8080')
    
    logger.info(f"Connecting to MCP server at {server_url}")
    
    # Create the client with longer timeout for debugging
    client = MCPClient(server_url, connection_timeout=30)
    
    try:
        # Check server health first
        logger.info("Checking server health...")
        try:
            import requests
            health_response = requests.get(f"{server_url}/health")
            logger.info(f"Server health status: {health_response.status_code}")
            logger.info(f"Health response: {health_response.text}")
        except Exception as e:
            logger.error(f"Error checking server health: {e}")
        
        # Initialize the connection
        logger.info("Initializing connection...")
        await client.initialize(client_name="example-client", client_version="1.0.0")
        
        # List available tools
        logger.info("Listing available tools...")
        tools = await client.list_tools()
        logger.info(f"Available tools: {json.dumps(tools, indent=2)}")
        
        # Check if resources tool is available
        found_resources = False
        for tool in tools:
            if tool.get("name") == "resources":
                found_resources = True
                logger.info("Found resources tool!")
                
                # List resources
                logger.info("Listing resources...")
                resources = await client.call_tool("resources", {"action": "list"})
                logger.info(f"Resources: {json.dumps(resources, indent=2)}")
                
                # Get resource types
                logger.info("Getting resource types...")
                types = await client.call_tool("resources", {"action": "getResourceTypes"})
                logger.info(f"Resource types: {json.dumps(types, indent=2)}")
                
                break
        
        if not found_resources:
            logger.warning("Resources tool not found")
            logger.info("Available tools: " + ", ".join([t.get("name", "unknown") for t in tools]))
        
    except MCPError as e:
        logger.error(f"MCP error: {e}")
        logger.error(f"Error details: {str(e)}")
    except Exception as e:
        logger.error(f"Unexpected error: {e}")
        import traceback
        logger.error(traceback.format_exc())
    finally:
        # Close the connection
        logger.info("Closing connection...")
        await client.close()


if __name__ == "__main__":
    asyncio.run(main()) 