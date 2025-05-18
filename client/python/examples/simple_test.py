#!/usr/bin/env python
"""
A simple script to test the OpsRamp MCP client manually.
With special handling for the known session validation issue.
"""

import asyncio
import os
import sys
import time
import logging
import traceback
import signal
import requests
import json
import sseclient

# Add parent directory to Python path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from src.ormcp.client import MCPClient
from src.ormcp.exceptions import MCPError, ConnectionError, SessionError, ToolError, JSONRPCError

# Configure logging
logging.basicConfig(
    level=logging.DEBUG,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)

# Silence some verbose loggers
logging.getLogger('urllib3').setLevel(logging.WARNING)
logging.getLogger('sseclient').setLevel(logging.INFO)
logger = logging.getLogger('simple_test')


# Set up signal handler for clean exits
def handle_signal(signum, frame):
    logger.info(f"Received signal {signum}, exiting...")
    sys.exit(1)

signal.signal(signal.SIGINT, handle_signal)
signal.signal(signal.SIGTERM, handle_signal)


def get_server_info(url, timeout=2):
    """Get server info directly via the health endpoint."""
    try:
        response = requests.get(f"{url}/health", timeout=timeout)
        if response.status_code == 200:
            return response.json()
    except Exception as e:
        logger.error(f"Error getting server info: {e}")
    return None


def get_session_id_directly(url, timeout=5):
    """Get a session ID directly via SSE connection."""
    logger.info(f"Establishing direct SSE connection to {url}/sse...")
    try:
        session = requests.Session()
        session.headers.update({
            'Accept': 'text/event-stream',
            'Cache-Control': 'no-cache'
        })
        
        response = session.get(f"{url}/sse", stream=True, timeout=timeout)
        client = sseclient.SSEClient(response)
        
        # Get the first event
        for event in client.events():
            if event.event == "endpoint":
                logger.debug(f"Received endpoint event: {event.data}")
                # Extract session ID
                if "sessionId=" in event.data:
                    parts = event.data.split("sessionId=")
                    if len(parts) > 1:
                        session_id = parts[1].strip()
                        logger.info(f"Got session ID: {session_id}")
                        return session_id
            break  # Only process the first event
            
        logger.error("No valid session ID received from SSE connection")
        return None
    except Exception as e:
        logger.error(f"Error establishing SSE connection: {e}")
        return None


def send_jsonrpc_request(url, session_id, method, params=None, timeout=5):
    """Send a JSON-RPC request directly."""
    message_url = f"{url}/message?sessionId={session_id}"
    payload = {
        "jsonrpc": "2.0",
        "id": str(time.time()),
        "method": method
    }
    if params:
        payload["params"] = params
        
    logger.debug(f"Sending request to {message_url}: {payload}")
    
    try:
        response = requests.post(
            message_url,
            json=payload,
            headers={"Content-Type": "application/json"},
            timeout=timeout
        )
        
        if response.status_code != 200:
            logger.warning(f"Request returned status code {response.status_code}: {response.text}")
            
        if "Invalid session ID" in response.text:
            logger.warning("Server reported invalid session ID - this is a known issue with the mark3labs/mcp-go library")
            
        return response.json() if response.text else None
    except Exception as e:
        logger.error(f"Error sending request: {e}")
        return None


async def main():
    # Parse command line arguments
    connection_timeout = int(os.environ.get('CONNECTION_TIMEOUT', '10'))
    request_timeout = int(os.environ.get('REQUEST_TIMEOUT', '30'))
    server_url = os.environ.get('MCP_SERVER_URL', 'http://localhost:8080')
    
    logger.info(f"Testing connection to {server_url}")
    logger.info(f"Using connection timeout: {connection_timeout}s, request timeout: {request_timeout}s")
    
    # First check if server is up via health endpoint
    server_info = get_server_info(server_url, timeout=connection_timeout)
    if not server_info:
        logger.error("Server is not responding to health checks")
        return 1
        
    logger.info(f"Server is running: {server_info}")
    logger.info(f"Available tools: {', '.join(server_info.get('tools', []))}")
    
    # Try to get a session ID directly
    session_id = get_session_id_directly(server_url, timeout=connection_timeout)
    if not session_id:
        logger.error("Could not obtain a session ID")
        return 1
    
    # Try direct initialization
    init_result = send_jsonrpc_request(
        server_url, 
        session_id, 
        "initialize", 
        {"client": {"name": "direct-test", "version": "1.0.0"}},
        timeout=request_timeout
    )
    
    if init_result:
        logger.info(f"Initialization result: {init_result}")
        if "error" in init_result and "Invalid session ID" in str(init_result["error"]):
            logger.warning("Known issue: Server rejected session ID due to mark3labs/mcp-go library limitation")
            logger.info("Continuing with tests to demonstrate connection was established")
        elif "result" in init_result:
            logger.info("Successfully initialized connection")
    
    # Try to list tools
    tools_result = send_jsonrpc_request(
        server_url,
        session_id,
        "list_tools",
        {},
        timeout=request_timeout
    )
    
    if tools_result:
        logger.info(f"Tools listing result: {tools_result}")
    
    # Try client via standard API as well
    try:
        client = MCPClient(server_url, auto_connect=False, connection_timeout=connection_timeout)
        
        # Connect
        logger.info("Trying connection via MCPClient...")
        sid = client.connect()
        logger.info(f"Connected with session ID: {sid}")
        
        # Check for integrations tool in server info
        if "integrations" in server_info.get("tools", []):
            logger.info("'integrations' tool is available according to health check")
        else:
            logger.warning("'integrations' tool not found in health check")
        
        # Close client
        await client.close(timeout=5)
        logger.info("Client closed")
        
    except MCPError as e:
        logger.error(f"Error using MCPClient: {e}")
    
    logger.info("Testing completed")
    return 0


if __name__ == "__main__":
    exit_code = asyncio.run(main())
    sys.exit(exit_code) 