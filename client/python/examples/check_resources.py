#!/usr/bin/env python
"""
Simple script to test server connection and session handling
"""

import asyncio
import sys
import os
import logging
import json
import requests

# Configure logging
logging.basicConfig(
    level=logging.DEBUG,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger('mcp_connection_test')

async def main():
    server_url = "http://localhost:8080"
    
    # 1. Test basic HTTP connection
    logger.info(f"Testing connection to {server_url}")
    try:
        response = requests.get(f"{server_url}/health")
        logger.info(f"Health check status: {response.status_code}")
        logger.info(f"Health response: {response.text}")
        
        health_data = response.json()
        if "tools" in health_data:
            logger.info(f"Available tools: {health_data['tools']}")
            
            if "resources" in health_data["tools"]:
                logger.info("✅ Resources tool is registered on the server")
            else:
                logger.error("❌ Resources tool is not registered on the server")
    except Exception as e:
        logger.error(f"Error connecting to server: {e}")
        return
    
    # 2. Test SSE connection
    logger.info("Testing SSE connection")
    try:
        # Import here to ensure the path is set correctly
        sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))
        from src.ormcp.session import MCPSession
        
        session = MCPSession(server_url, connection_timeout=15)
        logger.info("Created session object")
        
        # Just test the connection establishment
        logger.info("Connecting to SSE endpoint...")
        session_id = session.connect()
        logger.info(f"✅ Successfully connected with session ID: {session_id}")
        
        # Test sending a basic message
        logger.info("Testing message sending...")
        response = session.send_message("echo", {"message": "Hello, server!"})
        logger.info(f"Message response: {response}")
        
        # Clean up
        logger.info("Cleaning up session...")
        session.close()
        logger.info("Session closed")
        
    except Exception as e:
        logger.error(f"Error with SSE connection: {e}")
        import traceback
        logger.error(traceback.format_exc())

if __name__ == "__main__":
    asyncio.run(main()) 