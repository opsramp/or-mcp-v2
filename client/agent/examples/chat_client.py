#!/usr/bin/env python3
"""
Interactive chat client for the OpsRamp AI Agent.
"""

import os
import sys
import asyncio
import logging
import argparse
import json
import requests
import time
from typing import Optional

# Add parent directory to path for local development
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from src.opsramp_agent.agent import Agent, MCPError, AsyncSSEClient, parse_session_id_from_sse

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
)
logger = logging.getLogger(__name__)

DEFAULT_SERVER_URL = 'http://localhost:8080'
DEFAULT_LLM_PROVIDER = 'openai'
DEFAULT_ENV_FILE = None

# Set to True to use mock mode (no actual MCP or LLM connections)
MOCK_MODE = False


def establish_sse_connection(server_url):
    """
    Establish an SSE connection directly and get a session ID from the MCP server.
    
    Args:
        server_url: The URL of the MCP server
        
    Returns:
        Session ID if successful, None otherwise
    """
    try:
        print(f"Establishing SSE connection to {server_url}/sse...")
        
        # Set up headers and timeout
        headers = {
            'Accept': 'text/event-stream',
            'Cache-Control': 'no-cache'
        }
        timeout = 60  # seconds
        
        print(f"Using extended timeout of {timeout} seconds...")
        
        # Make direct request without using AsyncSSEClient
        print("Attempting to connect to SSE endpoint...")
        response = requests.get(f"{server_url}/sse", stream=True, headers=headers, timeout=timeout)
        
        if response.status_code != 200:
            print(f"Failed to connect to SSE endpoint with status code {response.status_code}")
            return None
        
        print("Connected to SSE endpoint, waiting for session ID from server...")
        print(f"This may take up to {timeout} seconds...")
        
        # Variables for parsing the SSE stream
        current_event_type = None
        current_event_data = None
        buffer = ""
        session_id = None
        start_time = time.time()
        
        # Process the SSE stream character by character for more reliable parsing
        for chunk in response.iter_content(chunk_size=1, decode_unicode=True):
            # Check for timeout
            if time.time() - start_time > timeout:
                print(f"Timeout after {timeout} seconds waiting for session ID")
                return None
            
            # Decode the chunk if needed
            if isinstance(chunk, bytes):
                try:
                    chunk = chunk.decode('utf-8')
                except UnicodeDecodeError:
                    continue
            
            # Process line endings
            if chunk == '\n':
                line = buffer.strip()
                buffer = ""
                
                if not line:  # Empty line marks the end of an event
                    # Process complete event
                    if current_event_type == 'endpoint' and current_event_data:
                        print(f"Received complete endpoint event: {current_event_data}")
                        session_id = parse_session_id_from_sse(current_event_data)
                        if session_id:
                            print(f"Extracted session ID: {session_id}")
                            break  # Found the session ID, exit the loop
                    
                    current_event_type = None
                    current_event_data = None
                    continue
                
                # Parse event type and data
                if line.startswith('event:'):
                    current_event_type = line[6:].strip()
                    print(f"Event type: {current_event_type}")
                elif line.startswith('data:'):
                    current_event_data = line[5:].strip()
                    print(f"Event data: {current_event_data}")
            else:
                buffer += chunk
        
        # Close the response to free resources
        response.close()
        
        if not session_id:
            print("Could not extract a session ID from the SSE stream")
            return None
        
        print(f"Successfully established session with ID: {session_id}")
        return session_id
        
    except requests.exceptions.Timeout:
        print(f"Connection timed out after {timeout} seconds")
        return None
    except Exception as e:
        print(f"Error establishing SSE connection: {str(e)}")
        import traceback
        traceback.print_exc()
        return None


def get_integrations_direct(server_url, session_id):
    """
    Get integrations data directly from the MCP server using the session ID.
    
    Args:
        server_url: The URL of the MCP server
        session_id: The session ID from the SSE connection
        
    Returns:
        List of integrations if successful, None otherwise
    """
    try:
        url = f"{server_url}/message?sessionId={session_id}"
        # Use the correct format for calling a tool
        payload = {
            "jsonrpc": "2.0",
            "id": "1",
            "method": "callTool",
            "params": {
                "name": "integrations",
                "arguments": {
                    "action": "list"
                }
            }
        }
        
        print(f"Making direct request to {url}")
        response = requests.post(url, json=payload, timeout=30)
        
        if response.status_code != 200:
            print(f"Error: Request failed with status {response.status_code}")
            print(f"Response: {response.text}")
            return None
        
        data = response.json()
        if "error" in data:
            print(f"Error: {data['error']['message']}")
            return None
        
        if "result" in data:
            return data["result"]
        
        return None
    
    except Exception as e:
        print(f"Error making direct request: {str(e)}")
        return None


async def interactive_chat(
    server_url: str,
    llm_provider: str,
    env_file: Optional[str] = None,
    simple_mode: bool = False
):
    """
    Run an interactive chat session with the OpsRamp AI Agent.
    
    Args:
        server_url: URL of the MCP server
        llm_provider: LLM provider to use ('openai' or 'anthropic')
        env_file: Path to .env file containing configuration
        simple_mode: Whether to run in simple mode without MCP connection
    """
    if simple_mode:
        print("\n======================================")
        print("Starting OpsRamp AI Agent in simple mode")
        print("======================================\n")
        
        while True:
            try:
                # Get user query
                user_query = input("\nYou: ")
                if user_query.lower() in ['exit', 'quit', 'q']:
                    print("\nExiting chat session.")
                    break
                    
                print("\nAI Assistant: This is a simple response. In a real session, this would connect to the OpsRamp MCP.")
                
            except KeyboardInterrupt:
                print("\n\nExiting chat session.")
                break
            except Exception as e:
                print(f"\nError: {str(e)}")
        
        return
    
    # First establish direct SSE connection to get session ID
    session_id = establish_sse_connection(server_url)
    if not session_id:
        print("Failed to establish SSE connection with the MCP server.")
        print("Trying to proceed with agent initialization anyway...")
    else:
        print("SSE connection successfully established.")
        # Test the connection with a direct call
        integrations = get_integrations_direct(server_url, session_id)
        if integrations:
            print(f"✅ Connection test successful! Found {len(integrations)} integrations.")
        else:
            print("⚠️ Connection test failed. Proceeding with agent anyway...")
    
    # Initialize the Agent
    try:
        print(f"Connecting to MCP server at {server_url}...")
        agent = Agent(server_url, llm_provider=llm_provider, env_file=env_file)
        
        # Skip the connection step since we already have a session ID
        if session_id:
            print(f"Using existing session ID: {session_id}")
            # Manually set the session ID in the agent
            agent.mcp_client.session.session_id = session_id
            agent.mcp_client.session.is_connected = True
            # Skip the connect() call that would attempt to establish a new connection
        else:
            # Only try to connect if we don't have a valid session ID
            await agent.connect()
        
        print("\n======================================")
        print("Welcome to the OpsRamp AI Agent")
        print("======================================")
        print(f"Connected to server: {server_url}")
        print(f"Available tools: integrations")
        print("Type 'exit', 'quit', or 'q' to end the session.")
        print("======================================\n")
        
        # Main chat loop
        while True:
            try:
                # Get user query
                user_query = input("\nYou: ")
                if user_query.lower() in ['exit', 'quit', 'q']:
                    print("\nExiting chat session.")
                    break
                
                # Special handling for integration listing
                if "list" in user_query.lower() and "integration" in user_query.lower():
                    if session_id:
                        integrations = get_integrations_direct(server_url, session_id)
                        if integrations:
                            # Format the response
                            response = "Here are the OpsRamp integrations:\n\n"
                            for i, integration in enumerate(integrations, 1):
                                display_name = integration.get("displayName", "Unknown")
                                id_value = integration.get("id", "Unknown ID")
                                status = integration.get("status", "Unknown status")
                                category = integration.get("category", "")
                                app = integration.get("app", "")
                                version = integration.get("version", "")
                                
                                response += f"{i}. {display_name} (ID: {id_value})\n"
                                response += f"   Status: {status}, Category: {category}\n"
                                if app and version:
                                    response += f"   App: {app}, Version: {version}\n"
                                response += "\n"
                            
                            print(f"\nAI Assistant: {response}")
                            continue
                
                # Get a response
                response = await agent.chat(user_query)
                
                # Display response
                print(f"\nAI Assistant: {response}")
                
            except KeyboardInterrupt:
                print("\n\nExiting chat session.")
                break
            except Exception as e:
                print(f"\nError: {str(e)}")
        
    except MCPError as e:
        logger.error(f"MCP Error: {str(e)}")
        print(f"Error: {str(e)}")
    except Exception as e:
        logger.error(f"Unexpected error: {str(e)}", exc_info=True)
        print(f"Error: {str(e)}")
    finally:
        # Clean up
        if 'agent' in locals():
            try:
                await agent.close()
            except Exception as e:
                logger.error(f"Error closing agent: {str(e)}")


async def process_single_prompt(
    prompt: str,
    server_url: str,
    llm_provider: str,
    env_file: Optional[str] = None,
    output_json: bool = False,
    simple_mode: bool = False,
    connection_timeout: int = 60,
    request_timeout: int = 15
):
    """
    Process a single prompt and return the response.
    
    Args:
        prompt: The user prompt to process
        server_url: URL of the MCP server
        llm_provider: LLM provider to use ('openai' or 'anthropic')
        env_file: Path to .env file containing configuration
        output_json: Whether to output in JSON format
        simple_mode: Whether to run in simple mode without MCP connection
        connection_timeout: Connection timeout in seconds
        request_timeout: Request timeout in seconds for JSON-RPC requests
    """
    if simple_mode:
        # In simple mode, provide mocked responses for specific queries
        lower_prompt = prompt.lower()
        response = "This is a simple mode response. The agent is not connected to the OpsRamp MCP."
        
        # Special handling for common queries
        if "list all integration" in lower_prompt or "show integration" in lower_prompt:
            response = """Here are the OpsRamp integrations:

1. hpe-alletra-LabRat (ID: INTG-2ed93041-eb92-40e9-b6b4-f14ad13e54fc)
   Status: Installed, Category: SDK APP
   App: hpe-alletra, Version: 7.0.0

2. redfish-server-LabRat (ID: INTG-f9e5d2aa-ee17-4e32-9251-493566ebdfca)
   Status: Installed, Category: SDK APP
   App: redfish-server, Version: 7.0.0

3. vcenter-58.51 (ID: INTG-00ee85e2-1f84-4fc1-8cf8-5277ae6980dd)
   Status: enabled, Category: COMPUTE_INTEGRATION
"""
            
        if output_json:
            import json
            result = {
                "prompt": prompt,
                "response": response,
                "status": "success"
            }
            print(json.dumps(result, indent=2))
        else:
            print(f"Prompt: {prompt}")
            print(f"\nResponse: {response}")
        return
    
    # First establish direct SSE connection to get session ID
    session_id = establish_sse_connection(server_url)
    if session_id:
        print("SSE connection successfully established.")
        
        # Handle special cases directly
        lower_prompt = prompt.lower()
        if "list all integration" in lower_prompt or "show integration" in lower_prompt:
            integrations = get_integrations_direct(server_url, session_id)
            if integrations:
                # Format the response
                response = "Here are the OpsRamp integrations:\n\n"
                for i, integration in enumerate(integrations, 1):
                    display_name = integration.get("displayName", "Unknown")
                    id_value = integration.get("id", "Unknown ID")
                    status = integration.get("status", "Unknown status")
                    category = integration.get("category", "")
                    app = integration.get("app", "")
                    version = integration.get("version", "")
                    
                    response += f"{i}. {display_name} (ID: {id_value})\n"
                    response += f"   Status: {status}, Category: {category}\n"
                    if app and version:
                        response += f"   App: {app}, Version: {version}\n"
                    response += "\n"
                
                if output_json:
                    result = {
                        "prompt": prompt,
                        "response": response,
                        "status": "success"
                    }
                    print(json.dumps(result, indent=2))
                else:
                    print(f"Prompt: {prompt}")
                    print(f"\nResponse: {response}")
                return
    
    try:
        # Initialize the Agent
        agent = Agent(
            server_url=server_url,
            llm_provider=llm_provider,
            env_file=env_file,
            connection_timeout=connection_timeout,
            simple_mode=simple_mode,
            request_timeout=request_timeout
        )
        
        # Connect to the server and initialize the session
        await agent.connect()
        
        # Process the prompt
        response = await agent.chat(prompt)
        
        # Output the result
        if output_json:
            import json
            result = {
                "prompt": prompt,
                "response": response,
                "status": "success"
            }
            print(json.dumps(result, indent=2))
        else:
            print(f"Prompt: {prompt}")
            print(f"\nResponse: {response}")
            
    except Exception as e:
        if output_json:
            import json
            result = {
                "prompt": prompt,
                "error": str(e),
                "status": "error"
            }
            print(json.dumps(result, indent=2))
        else:
            print(f"Error processing prompt: {str(e)}")
    finally:
        # Clean up
        if 'agent' in locals():
            await agent.close()


def setup_argparse():
    """Set up command line arguments."""
    parser = argparse.ArgumentParser(
        description="OpsRamp AI Agent Chat Client"
    )
    parser.add_argument(
        "--server-url", 
        default=DEFAULT_SERVER_URL,
        help=f"OpsRamp MCP server URL (default: {DEFAULT_SERVER_URL})"
    )
    parser.add_argument(
        "--prompt", 
        help="Process a single prompt and exit (non-interactive mode)"
    )
    parser.add_argument(
        "--debug", 
        action="store_true", 
        help="Enable debug logging"
    )
    parser.add_argument(
        "--simple-mode",
        action="store_true",
        help="Run in simple mode without connecting to MCP server"
    )
    parser.add_argument(
        "--env-file",
        default=DEFAULT_ENV_FILE,
        help="Path to .env file containing config variables"
    )
    parser.add_argument(
        "--llm-provider",
        choices=["openai", "anthropic"],
        default=DEFAULT_LLM_PROVIDER,
        help=f"LLM provider to use (default: {DEFAULT_LLM_PROVIDER})"
    )
    parser.add_argument(
        "--connection-timeout",
        type=int,
        default=60,
        help="Connection timeout in seconds (default: 60)"
    )
    parser.add_argument(
        "--request-timeout",
        type=int,
        default=15,
        help="Request timeout in seconds for JSON-RPC requests (default: 15)"
    )
    return parser.parse_args()


def main():
    """Main entry point for the chat client."""
    args = setup_argparse()
    
    # Use simple mode if MOCK_MODE is set to True
    simple_mode = args.simple_mode or MOCK_MODE
    
    if simple_mode:
        print("WARNING: Running in simple mode without connecting to MCP server")
    
    # Run in single prompt mode or interactive mode
    if args.prompt:
        asyncio.run(process_single_prompt(
            args.prompt,
            args.server_url,
            args.llm_provider,
            args.env_file,
            False,
            simple_mode,
            args.connection_timeout,
            args.request_timeout
        ))
    else:
        asyncio.run(interactive_chat(
            args.server_url,
            args.llm_provider,
            args.env_file,
            simple_mode
        ))


if __name__ == '__main__':
    main() 