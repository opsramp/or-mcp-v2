#!/usr/bin/env python
"""
Example script demonstrating how to use the OpsRamp MCP Python client 
with a browser-like SSE connection.

This script demonstrates:
1. Connection establishment with browser-like SSE
2. Session ID validation
3. JSON-RPC request/response
4. Tool discovery and invocation
5. Event handling
"""

import asyncio
import os
import sys
import logging
import json
import argparse
from datetime import datetime

# Add the parent directory to the Python path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from src.ormcp.client import MCPClient
from src.ormcp.exceptions import MCPError, SessionError, JSONRPCError

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger('browser_like_example')

# Silence some verbose loggers
logging.getLogger('urllib3').setLevel(logging.WARNING)
logging.getLogger('sseclient').setLevel(logging.INFO)


def setup_argparse():
    """Set up command line arguments."""
    parser = argparse.ArgumentParser(
        description="OpsRamp MCP Browser-like Client Example"
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
        "--wait-events", 
        type=int, 
        default=3,
        help="Number of seconds to wait for events after initialization"
    )
    parser.add_argument(
        "--client-name", 
        default="browser-like-client",
        help="Client name to use for initialization"
    )
    parser.add_argument(
        "--client-version", 
        default="1.0.0",
        help="Client version to use for initialization"
    )
    return parser.parse_args()


async def main():
    """Run the browser-like client example."""
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
        # Connect to the server
        # The connection happens automatically with auto_connect=True (default)
        # The session ID should now be valid on the server side
        session_id = client.session.session_id
        logger.info(f"Connected with session ID: {session_id}")
        
        # Initialize the connection
        logger.info(f"Initializing connection as {args.client_name} v{args.client_version}...")
        result = await client.initialize(
            client_name=args.client_name, 
            client_version=args.client_version
        )
        logger.info(f"Initialization result: {json.dumps(result, indent=2)}")
        
        # Register an event handler to show SSE events
        def event_handler(event_data):
            logger.debug(f"Event received: {event_data}")
        
        client.session.register_event_handler('*', event_handler)
        
        # List available tools
        logger.info("Listing available tools...")
        tools = await client.list_tools()
        
        if not tools:
            logger.warning("No tools available")
        else:
            logger.info(f"Found {len(tools)} available tools:")
            for tool in tools:
                if isinstance(tool, str):
                    logger.info(f"- {tool}")
                else:
                    tool_name = tool.get('name', 'unknown')
                    tool_desc = tool.get('description', 'No description')
                    logger.info(f"- {tool_name}: {tool_desc}")
        
        # Try to call the integrations tool directly
        logger.info("Attempting to call the integrations tool directly...")
        try:
            # List integrations
            integrations = await client.call_tool("integrations", {"action": "list"})
            logger.info(f"Integrations result: {json.dumps(integrations, indent=2)}")
        except MCPError as e:
            logger.error(f"Error calling integrations tool: {e}")
        
        # Wait for events to be received
        if args.wait_events > 0:
            logger.info(f"Waiting {args.wait_events}s for events...")
            await asyncio.sleep(args.wait_events)
        
        # Check events received during the connection
        events = client.session.get_received_events()
        logger.info(f"Received {len(events)} events during the session")
        
        # Print event types and counts
        event_types = {}
        for event in events:
            event_type = event.get('event', 'unknown')
            event_types[event_type] = event_types.get(event_type, 0) + 1
        
        for event_type, count in event_types.items():
            logger.info(f"- {event_type}: {count} events")
        
    except SessionError as e:
        logger.error(f"Session error: {e}")
    except JSONRPCError as e:
        logger.error(f"JSON-RPC error: {e}")
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