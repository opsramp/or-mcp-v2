#!/usr/bin/env python
"""
Example script demonstrating how to call the integrations tool with the OpsRamp MCP client.

This script demonstrates:
1. Connection establishment with browser-like SSE
2. Using debug info to inspect the server state
3. Directly calling the integrations tool
"""

import asyncio
import os
import sys
import logging
import json
import argparse

# Add the parent directory to the Python path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from src.ormcp.client import MCPClient
from src.ormcp.exceptions import MCPError, ToolError

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger('call_integrations')


def setup_argparse():
    """Set up command line arguments."""
    parser = argparse.ArgumentParser(
        description="Call OpsRamp MCP Integrations Tool Example"
    )
    parser.add_argument(
        "--server-url", 
        default=os.environ.get('MCP_SERVER_URL', 'http://localhost:8080'),
        help="MCP server URL (default: http://localhost:8080 or MCP_SERVER_URL env var)"
    )
    parser.add_argument(
        "--debug", 
        action="store_true", 
        help="Enable debug logging"
    )
    parser.add_argument(
        "--action",
        default="list",
        choices=["list", "get"],
        help="Action to perform on the integrations tool"
    )
    parser.add_argument(
        "--id",
        help="Integration ID for the get action"
    )
    return parser.parse_args()


async def main():
    """Run the tool call example."""
    # Parse command line arguments
    args = setup_argparse()
    
    # Set debug logging if requested
    if args.debug:
        logger.setLevel(logging.DEBUG)
        logging.getLogger('ormcp').setLevel(logging.DEBUG)
    
    logger.info(f"Connecting to MCP server at {args.server_url}")
    
    # Create the client
    client = MCPClient(args.server_url)
    
    try:
        # Connect to the server (this happens automatically with auto_connect=True)
        session_id = client.session.session_id
        logger.info(f"Connected with session ID: {session_id}")
        
        # Check server health
        try:
            import requests
            health_response = requests.get(f"{args.server_url}/health", timeout=5)
            if health_response.status_code == 200:
                health_data = health_response.json()
                logger.info(f"Server health: {health_data.get('status')}")
                if 'tools' in health_data:
                    logger.info(f"Server tools: {', '.join(health_data['tools'])}")
        except Exception as e:
            logger.warning(f"Could not check server health: {e}")
        
        # Initialize the connection
        logger.info("Initializing connection...")
        result = await client.initialize(
            client_name="integrations-client", 
            client_version="1.0.0"
        )
        logger.info(f"Initialization result: {json.dumps(result, indent=2)}")
        
        # Build arguments based on the action
        tool_args = {"action": args.action}
        if args.action == "get" and args.id:
            tool_args["id"] = args.id
        
        # Call the integrations tool directly
        logger.info(f"Calling integrations tool with arguments: {tool_args}")
        try:
            result = await client.call_tool("integrations", tool_args)
            logger.info(f"Integrations result: {json.dumps(result, indent=2)}")
        except ToolError as e:
            logger.error(f"Tool error: {e}")
        except MCPError as e:
            logger.error(f"MCP error: {e}")
        
    except Exception as e:
        logger.error(f"Unexpected error: {e}", exc_info=True)
    finally:
        # Close the connection
        logger.info("Closing connection...")
        await client.close()
        logger.info("Connection closed")


if __name__ == "__main__":
    asyncio.run(main()) 