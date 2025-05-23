#!/usr/bin/env python
"""
Example script demonstrating how to use the OpsRamp MCP Python client 
with persistent SSE connections.

This script shows:
1. How to establish and maintain a persistent SSE connection
2. How to handle connection interruptions and recover
3. How to periodically check integrations data over a long-running session
4. How to use event handlers to process server messages

Usage:
    python persistent_sse_example.py [--server-url URL] [--run-time SECONDS] [--polling-interval SECONDS]
"""

import asyncio
import os
import sys
import logging
import json
import argparse
import signal
import time
from datetime import datetime, timedelta

# Add the parent directory to the Python path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from src.ormcp.client import MCPClient
from src.ormcp.exceptions import MCPError, SessionError, JSONRPCError, ConnectionError, ToolError

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger('persistent_sse_example')

# Set log levels for other modules
logging.getLogger('urllib3').setLevel(logging.WARNING)


def setup_argparse():
    """Set up command line arguments."""
    parser = argparse.ArgumentParser(
        description="OpsRamp MCP Persistent SSE Connection Example"
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
        "--run-time", 
        type=int,
        default=300,  # 5 minutes
        help="How long to run the example in seconds (default: 300 seconds)"
    )
    parser.add_argument(
        "--polling-interval", 
        type=int,
        default=30,
        help="How often to poll for integrations in seconds (default: 30 seconds)"
    )
    return parser.parse_args()


class PersistentSSEClient:
    """
    Demonstrates a persistent SSE client that maintains a long-running connection
    and periodically checks for integrations.
    """
    
    def __init__(self, server_url, run_time, polling_interval):
        """Initialize the persistent client."""
        self.server_url = server_url
        self.run_time = run_time
        self.polling_interval = polling_interval
        self.client = None
        self.running = False
        self.received_events = []
        self.last_event_time = time.time()
        self.end_time = None
        self.integration_counts = {}
        
        # Set up signal handlers for graceful shutdown
        signal.signal(signal.SIGINT, self._signal_handler)
        signal.signal(signal.SIGTERM, self._signal_handler)
    
    def _signal_handler(self, sig, frame):
        """Handle Ctrl+C and termination signals."""
        logger.info("Shutdown signal received, closing connection gracefully...")
        self.running = False
    
    async def start(self):
        """Start the persistent client."""
        logger.info(f"Starting persistent SSE client for {self.run_time} seconds")
        logger.info(f"Will poll for integrations every {self.polling_interval} seconds")
        
        self.running = True
        self.end_time = datetime.now() + timedelta(seconds=self.run_time)
        
        # Connect to the server
        try:
            self.client = MCPClient(self.server_url)
            
            # Register event handler
            self.client.session.register_event_handler('*', self._handle_event)
            
            # Initialize the connection
            await self.client.initialize(
                client_name="persistent-sse-client",
                client_version="1.0.0"
            )
            
            logger.info(f"Connected and initialized, session will run until {self.end_time}")
            
            # Start the main loop
            await self._run_main_loop()
            
        except ConnectionError as e:
            logger.error(f"Connection error: {e}")
        except Exception as e:
            logger.error(f"Unexpected error: {e}", exc_info=True)
        finally:
            await self._cleanup()
    
    async def _run_main_loop(self):
        """Run the main loop, polling for integrations periodically."""
        last_poll_time = 0
        
        while self.running and datetime.now() < self.end_time:
            try:
                # Check if it's time to poll for integrations
                current_time = time.time()
                if current_time - last_poll_time >= self.polling_interval:
                    await self._poll_integrations()
                    last_poll_time = current_time
                
                # Check for event activity
                inactive_time = current_time - self.last_event_time
                if inactive_time > 60:  # No events for 60+ seconds
                    logger.warning(f"No events received for {inactive_time:.1f} seconds")
                
                # Sleep a short time to prevent tight loop
                await asyncio.sleep(1)
                
            except ConnectionError as e:
                logger.error(f"Connection error during main loop: {e}")
                # Wait before retrying
                await asyncio.sleep(5)
            except Exception as e:
                logger.error(f"Error in main loop: {e}", exc_info=True)
                await asyncio.sleep(5)
    
    async def _poll_integrations(self):
        """Poll for integrations and log results."""
        try:
            # List available tools first
            tools = await self.client.list_tools()
            tool_names = [t["name"] if isinstance(t, dict) else t for t in tools]
            logger.info(f"Available tools: {', '.join(tool_names)}")
            
            # Check if integrations tool is available
            if "integrations" in tool_names:
                logger.info("Polling for integrations...")
                
                # Call the integrations tool to list integrations
                integrations = await self.client.call_tool(
                    "integrations", 
                    {"action": "list"}
                )
                
                if isinstance(integrations, list):
                    # Group by integration type
                    type_counts = {}
                    for integration in integrations:
                        if isinstance(integration, dict):
                            int_type = integration.get("type", "unknown")
                            type_counts[int_type] = type_counts.get(int_type, 0) + 1
                    
                    # Log the results
                    logger.info(f"Found {len(integrations)} total integrations:")
                    for int_type, count in type_counts.items():
                        logger.info(f"  - {int_type}: {count}")
                    
                    # Store for comparison
                    self.integration_counts = type_counts
                else:
                    logger.warning(f"Unexpected integrations response format: {type(integrations)}")
            else:
                logger.warning("Integrations tool not available")
                
        except ToolError as e:
            logger.error(f"Error polling integrations: {e}")
        except Exception as e:
            logger.error(f"Unexpected error polling integrations: {e}", exc_info=True)
    
    def _handle_event(self, event_data):
        """Handle SSE events."""
        # Update the last event time
        self.last_event_time = time.time()
        
        # Store the event
        self.received_events.append({
            'time': datetime.now().isoformat(),
            'data': event_data
        })
        
        # Truncate event history if it gets too large
        if len(self.received_events) > 100:
            self.received_events = self.received_events[-50:]
    
    async def _cleanup(self):
        """Clean up resources."""
        if self.client:
            logger.info("Closing connection...")
            try:
                await self.client.close()
            except Exception as e:
                logger.error(f"Error during cleanup: {e}")
        
        # Print session summary
        event_types = {}
        for event in self.received_events:
            # Try to parse the event data to extract the event type
            try:
                if isinstance(event['data'], str):
                    data = json.loads(event['data'])
                    if isinstance(data, dict) and 'event' in data:
                        event_type = data['event']
                    else:
                        event_type = 'unknown'
                else:
                    event_type = 'unknown'
            except (json.JSONDecodeError, TypeError):
                event_type = 'unknown'
                
            event_types[event_type] = event_types.get(event_type, 0) + 1
        
        logger.info(f"Session summary:")
        logger.info(f"  - Received {len(self.received_events)} events")
        for event_type, count in event_types.items():
            logger.info(f"  - {event_type}: {count} events")
        
        if self.integration_counts:
            logger.info(f"  - Found {sum(self.integration_counts.values())} integrations:")
            for int_type, count in self.integration_counts.items():
                logger.info(f"    - {int_type}: {count}")


async def main():
    """Run the persistent SSE example."""
    # Parse command line arguments
    args = setup_argparse()
    
    # Set debug logging if requested
    if args.debug:
        logger.setLevel(logging.DEBUG)
        logging.getLogger('ormcp').setLevel(logging.DEBUG)
    
    # Create and start the persistent client
    client = PersistentSSEClient(
        server_url=args.server_url,
        run_time=args.run_time,
        polling_interval=args.polling_interval
    )
    
    await client.start()


if __name__ == "__main__":
    asyncio.run(main()) 