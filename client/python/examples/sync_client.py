#!/usr/bin/env python
"""
Example script demonstrating how to use the synchronous OpsRamp MCP Python client
to interact with an MCP server without using asyncio.
"""

import sys
import os
import logging
import json

# Add the parent directory to the Python path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from src.ormcp.client import SyncMCPClient
from src.ormcp.exceptions import MCPError

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger('mcp_sync_example')


def main():
    # Server URL (default to localhost:8080)
    server_url = os.environ.get('MCP_SERVER_URL', 'http://localhost:8080')
    
    logger.info(f"Connecting to MCP server at {server_url}")
    
    # Create the synchronous client
    client = SyncMCPClient(server_url)
    
    try:
        # Initialize the connection
        logger.info("Initializing connection...")
        client.initialize(client_name="example-sync-client", client_version="1.0.0")
        
        # List available tools
        logger.info("Listing available tools...")
        tools = client.list_tools()
        logger.info(f"Available tools: {json.dumps(tools, indent=2)}")
        
        # Try calling the integrations tool if available
        for tool in tools:
            if tool.get("name") == "integrations":
                logger.info("Found integrations tool!")
                
                # Get a specific integration
                integration_id = "int-001"  # Example ID, adjust as needed
                logger.info(f"Getting integration with ID: {integration_id}")
                
                integration = client.call_tool("integrations", {
                    "action": "get",
                    "id": integration_id
                })
                
                logger.info(f"Integration details: {json.dumps(integration, indent=2)}")
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
        client.close()


if __name__ == "__main__":
    main() 