#!/usr/bin/env python
"""
Test script for both integrations and resources tools
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
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger('mcp_tools_test')


async def test_tools():
    """Test both integrations and resources tools"""
    # Server URL (default to localhost:8080)
    server_url = os.environ.get('MCP_SERVER_URL', 'http://localhost:8080')
    
    logger.info(f"Connecting to MCP server at {server_url}")
    
    # Create the client
    client = MCPClient(server_url, connection_timeout=15)
    
    try:
        # First check server health
        logger.info("Checking server health...")
        import requests
        health_response = requests.get(f"{server_url}/health")
        logger.info(f"Server health status: {health_response.status_code}")
        health_data = health_response.json()
        logger.info(f"Registered tools: {', '.join([t['name'] for t in health_data.get('tools', [])])}")
        
        # Initialize the connection
        logger.info("Initializing connection...")
        await client.initialize(client_name="tools-test-client", client_version="1.0.0")
        
        # List available tools
        logger.info("Listing available tools...")
        tools = await client.list_tools()
        tool_names = [tool.get("name") for tool in tools]
        logger.info(f"Available tools: {', '.join(tool_names)}")
        
        # Test integrations tool
        if "integrations" in tool_names:
            logger.info("=== Testing integrations tool ===")
            try:
                # List integrations
                logger.info("Listing integrations...")
                integrations = await client.call_tool("integrations", {"action": "list"})
                count = len(integrations.get("results", []))
                logger.info(f"Found {count} integrations")
                
                # List first few integrations
                for i, integration in enumerate(integrations.get("results", [])[:3]):
                    logger.info(f"Integration {i+1}: {integration.get('displayName')}, Type: {integration.get('category')}")
            except Exception as e:
                logger.error(f"Error testing integrations tool: {e}")
        
        # Test resources tool
        if "resources" in tool_names:
            logger.info("=== Testing resources tool ===")
            try:
                # List resources
                logger.info("Listing resources...")
                resources = await client.call_tool("resources", {"action": "list"})
                count = len(resources.get("results", []))
                logger.info(f"Found {count} resources")
                
                # List first few resources
                for i, resource in enumerate(resources.get("results", [])[:3]):
                    logger.info(f"Resource {i+1}: {resource.get('name')}, Type: {resource.get('type')}")
                
                # Get resource types
                logger.info("Getting resource types...")
                types = await client.call_tool("resources", {"action": "getResourceTypes"})
                type_count = len(types)
                logger.info(f"Found {type_count} resource types")
            except Exception as e:
                logger.error(f"Error testing resources tool: {e}")
        
    except MCPError as e:
        logger.error(f"MCP error: {e}")
    except Exception as e:
        logger.error(f"Unexpected error: {e}")
    finally:
        # Close the connection
        logger.info("Closing connection...")
        await client.close()


if __name__ == "__main__":
    asyncio.run(test_tools()) 