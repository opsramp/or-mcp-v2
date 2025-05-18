#!/usr/bin/env python
"""
Server health check utility.

This script checks if the MCP server is running and healthy.
It can be used to verify server availability before running tests.
"""

import os
import sys
import json
import argparse
import logging
import requests
from pathlib import Path

# Add the parent directory to the Python path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from src.ormcp.exceptions import MCPError

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger('check_server')

def check_health(server_url):
    """Check if the server is healthy."""
    try:
        response = requests.get(f"{server_url}/health", timeout=5)
        if response.status_code == 200:
            return response.json()
        else:
            logger.error(f"Server health check failed with status code: {response.status_code}")
            return None
    except requests.RequestException as e:
        logger.error(f"Failed to connect to server: {e}")
        return None

def check_readiness(server_url):
    """Check if the server is ready to accept requests."""
    try:
        response = requests.get(f"{server_url}/readiness", timeout=5)
        if response.status_code == 200:
            return response.json()
        else:
            logger.error(f"Server readiness check failed with status code: {response.status_code}")
            return None
    except requests.RequestException as e:
        logger.error(f"Failed to connect to server: {e}")
        return None

def main():
    """Main function."""
    parser = argparse.ArgumentParser(description="MCP Server Health Check")
    parser.add_argument("--url", default="http://localhost:8080", 
                        help="Server URL (default: http://localhost:8080)")
    parser.add_argument("--debug", action="store_true", help="Enable debug output")
    args = parser.parse_args()
    
    # Set log level based on args
    if args.debug:
        logger.setLevel(logging.DEBUG)
    
    logger.info(f"Checking server health at {args.url}")
    
    # Check basic health
    health_data = check_health(args.url)
    if health_data:
        logger.info("Server is healthy")
        logger.info(f"Server status: {health_data.get('status', 'unknown')}")
        logger.info(f"Server uptime: {health_data.get('uptime', 'unknown')}")
        
        # Show tools if available
        if 'tools' in health_data:
            tools = health_data['tools']
            logger.info(f"Available tools: {', '.join(tools)}")
    else:
        logger.error("Server is not healthy")
        return 1
    
    # Check readiness
    readiness_data = check_readiness(args.url)
    if readiness_data:
        is_ready = readiness_data.get('ready', False)
        if is_ready:
            logger.info("Server is ready to accept requests")
        else:
            logger.warning("Server is not ready to accept requests")
            if 'checks' in readiness_data:
                for check, status in readiness_data['checks'].items():
                    logger.warning(f"Check '{check}': {status}")
    else:
        logger.warning("Could not check server readiness")
    
    return 0 if health_data else 1


if __name__ == "__main__":
    sys.exit(main()) 