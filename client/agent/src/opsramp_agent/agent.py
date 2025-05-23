"""
OpsRamp AI Agent - Main Agent Implementation
"""

import os
import json
import asyncio
import logging
import aiohttp
import uuid
import time
import requests
import threading
import queue
from typing import Dict, List, Any, Optional, Union, Callable
from aiohttp_sse_client import client as sse_client

# LLM providers (choose one based on availability)
try:
    import openai
    HAS_OPENAI = True
except ImportError:
    HAS_OPENAI = False

try:
    from anthropic import Anthropic
    HAS_ANTHROPIC = True
except ImportError:
    HAS_ANTHROPIC = False

from .utils.config import load_env_from_file, get_api_keys

logger = logging.getLogger(__name__)

# MCP Client utilities
class MCPError(Exception):
    """Base exception for MCP client errors."""
    def __init__(self, message, code=None, data=None):
        self.message = message
        self.code = code
        self.data = data
        super().__init__(message)

class ConnectionError(MCPError):
    """Raised when there's an error connecting to the MCP server."""
    pass

class SessionError(MCPError):
    """Raised when there's an error with the session management."""
    pass

class ToolError(MCPError):
    """Raised when there's an error invoking a tool."""
    pass

def generate_request_id() -> str:
    """Generate a unique ID for a JSON-RPC request."""
    return str(uuid.uuid4())

def create_jsonrpc_request(method: str, params: Optional[Dict[str, Any]] = None, request_id: Optional[str] = None) -> Dict[str, Any]:
    """Create a JSON-RPC 2.0 request object."""
    if request_id is None:
        request_id = generate_request_id()
        
    request = {
        "jsonrpc": "2.0",
        "id": request_id,
        "method": method
    }
    
    if params is not None:
        request["params"] = params
        
    return request

def parse_jsonrpc_response(response_text: str) -> Dict[str, Any]:
    """Parse a JSON-RPC 2.0 response."""
    try:
        return json.loads(response_text)
    except json.JSONDecodeError as e:
        raise ValueError(f"Invalid JSON response: {response_text}") from e

def parse_session_id_from_sse(event_data: str) -> Optional[str]:
    """Parse a session ID from an SSE event."""
    if "sessionId=" in event_data:
        # Extract session ID from the format: /message?sessionId=<uuid>
        parts = event_data.split("sessionId=")
        if len(parts) > 1:
            session_id = parts[1].strip()
            return session_id
    return None

# Direct SSE client using the proven working method from chat_client.py
class AsyncSSEClient:
    """
    SSE client that uses the proven working method from chat_client.py.
    Uses direct requests with streaming instead of aiohttp-sse-client.
    """
    
    def __init__(self, url, headers=None, timeout=60):
        """Initialize the SSE client."""
        self.url = url
        self.headers = headers or {}
        self.timeout = timeout
        self.is_connected = False
        self.session_id = None
        
        # Add standard SSE headers
        if 'Accept' not in self.headers:
            self.headers['Accept'] = 'text/event-stream'
        if 'Cache-Control' not in self.headers:
            self.headers['Cache-Control'] = 'no-cache'

    def connect(self):
        """Connect using the proven working method from chat_client.py."""
        if self.is_connected:
            return True
        
        logger.debug(f"Connecting to SSE endpoint using proven method: {self.url}")
        
        try:
            # Use the exact method from chat_client.py that works
            import requests
            response = requests.get(self.url, stream=True, headers=self.headers, timeout=self.timeout)
            
            if response.status_code != 200:
                logger.error(f"Failed to connect to SSE endpoint with status code {response.status_code}")
                return False
            
            logger.debug("Connected to SSE endpoint, waiting for session ID from server...")
            
            # Process the SSE stream to find the endpoint event
            current_event_type = None
            current_event_data = ""
            
            for line in response.iter_lines(decode_unicode=True):
                if line is None:
                    continue
                    
                line = line.strip()
                if not line:
                    # Empty line indicates end of event
                    if current_event_type == "endpoint" and current_event_data:
                        self.session_id = parse_session_id_from_sse(current_event_data)
                        if self.session_id:
                            logger.debug(f"Successfully extracted session ID: {self.session_id}")
                            self.is_connected = True
                            return True
                    
                    # Reset for next event
                    current_event_type = None
                    current_event_data = ""
                    continue
                
                if line.startswith("event:"):
                    current_event_type = line[6:].strip()
                elif line.startswith("data:"):
                    current_event_data = line[5:].strip()
            
            logger.error("Could not extract session ID from SSE stream")
            return False
            
        except Exception as e:
            logger.error(f"Error establishing SSE connection: {str(e)}")
            return False

    def get_event(self, timeout=1):
        """Not used in the proven working method."""
        return None

    def wait_for_event(self, event_type=None, timeout=60):
        """Session ID is already extracted during connect."""
        if event_type == 'endpoint' and self.session_id:
            return {
                'event': 'endpoint',
                'data': f'/message?sessionId={self.session_id}'
            }
        return None

    def close(self):
        """Close the SSE connection."""
        self.is_connected = False

# MCP Session management
class MCPSession:
    """Manages a session with an MCP server."""
    
    def __init__(self, base_url: str, connection_timeout: int = 60):
        """Initialize the session."""
        self.base_url = base_url.rstrip('/')
        self.connection_timeout = connection_timeout
        self.session_id: Optional[str] = None
        self.sse_client: Optional[AsyncSSEClient] = None
        self.is_connected = False
        self.is_initialized = False
        self._event_handlers: Dict[str, Callable] = {}
        self._received_events: List[Dict] = []
        self._start_event_processing()
    
    def _get_message_url(self) -> str:
        """Get the URL for sending messages."""
        if not self.session_id:
            raise SessionError("No active session")
        return f"{self.base_url}/message?sessionId={self.session_id}"
    
    def _get_sse_url(self) -> str:
        """Get the URL for SSE connection."""
        return f"{self.base_url}/sse"
    
    def connect(self) -> str:
        """Connect to the MCP server and get a session ID using the proven working method."""
        if self.is_connected and self.session_id:
            return self.session_id
        
        try:
            # Create and connect SSE client using proven working method
            logger.debug(f"Establishing SSE connection to {self._get_sse_url()}")
            self.sse_client = AsyncSSEClient(
                self._get_sse_url(), 
                timeout=self.connection_timeout
            )
            
            if not self.sse_client.connect():
                raise SessionError("Failed to connect to SSE endpoint")
            
            # Get the session ID directly from the client
            self.session_id = self.sse_client.session_id
            
            if not self.session_id:
                raise SessionError("Could not obtain session ID from SSE connection")
            
            logger.debug(f"Session established with ID: {self.session_id}")
            self.is_connected = True
            return self.session_id
            
        except Exception as e:
            logger.error(f"Failed to establish session: {str(e)}", exc_info=True)
            self.is_connected = False
            raise SessionError(f"Failed to establish session: {str(e)}")
    
    def _start_event_processing(self):
        """Start a background thread to process events from the SSE client."""
        def process_events():
            while True:
                if not self.sse_client or not self.sse_client.is_connected:
                    time.sleep(1.0)  # Wait and retry if no connection
                    continue
                
                event = self.sse_client.get_event(timeout=1.0)
                if not event:
                    continue  # No event received, continue waiting
                
                # Store the event for later retrieval
                self._received_events.append(event)
                
                # Call any registered handlers for this event type
                event_type = event.get('event')
                if event_type in self._event_handlers:
                    try:
                        self._event_handlers[event_type](event)
                    except Exception as e:
                        logger.error(f"Error in event handler for {event_type}: {str(e)}", exc_info=True)
        
        # Start event processing in a daemon thread
        thread = threading.Thread(target=process_events, daemon=True)
        thread.start()
    
    async def send_request(self, method: str, params: Optional[Dict[str, Any]] = None, timeout: Optional[int] = None) -> Dict[str, Any]:
        """Send a JSON-RPC request to the MCP server."""
        if not self.is_connected:
            logger.error("Cannot send request: Not connected")
            raise ConnectionError("Cannot send request: Not connected")
        
        # Use default timeout from client if not specified
        if timeout is None:
            # Try to get timeout from self.request_timeout if available
            timeout = getattr(self, 'request_timeout', 30)
        
        # Create message URL with session ID
        url = self._get_message_url()
        
        # Create JSON-RPC request
        rpcRequest = create_jsonrpc_request(method, params)
        
        try:
            # Make the request
            async with aiohttp.ClientSession() as session:
                async with session.post(url, json=rpcRequest, timeout=timeout) as response:
                    # Check status code
                    if response.status != 200:
                        error_text = await response.text()
                        logger.error(f"Request failed: {response.status}, {error_text}")
                        raise MCPError(f"Request failed: {response.status}", code=response.status)
                    
                    # Parse response
                    response_json = await response.json()
                    
                    # Check for JSON-RPC error
                    if "error" in response_json:
                        error = response_json["error"]
                        logger.error(f"JSON-RPC error: {error}")
                        raise MCPError(
                            error.get("message", "Unknown error"),
                            code=error.get("code", -1),
                            data=error.get("data", None)
                        )
                    
                    return response_json
        except asyncio.TimeoutError:
            logger.error(f"Request timed out after {timeout} seconds")
            raise MCPError(f"Request timed out after {timeout} seconds", code=408)
        except aiohttp.ClientError as e:
            logger.error(f"Network error: {str(e)}")
            raise ConnectionError(f"Network error: {str(e)}")
    
    def register_event_handler(self, event_type: str, handler: Callable):
        """Register a handler function for a specific event type."""
        self._event_handlers[event_type] = handler
        logger.debug(f"Registered handler for event type: {event_type}")
    
    def get_received_events(self, event_type: Optional[str] = None) -> List[Dict]:
        """Get all received events, optionally filtered by event type."""
        if event_type is None:
            return self._received_events.copy()
        else:
            return [e for e in self._received_events if e.get('event') == event_type]
    
    def close(self):
        """Close the session."""
        logger.debug("Closing MCP session")
        if self.sse_client:
            self.sse_client.close()
        self.is_connected = False
        self.is_initialized = False

# MCP Client
class MCPClient:
    """Client for interacting with an OpsRamp MCP server."""
    
    def __init__(self, server_url: str, auto_connect: bool = True, connection_timeout: int = 10, request_timeout: int = 30):
        """Initialize the MCP client."""
        self.server_url = server_url
        self.connection_timeout = connection_timeout
        self.request_timeout = request_timeout
        self.session = MCPSession(server_url, connection_timeout)
        self.is_initialized = False
        self._available_tools = []
        
        if auto_connect:
            self.connect()
    
    def connect(self) -> str:
        """Connect to the MCP server."""
        try:
            logger.debug(f"Connecting to MCP server at {self.server_url}")
            return self.session.connect()
        except Exception as e:
            logger.error(f"Failed to connect to MCP server: {str(e)}", exc_info=True)
            raise ConnectionError(f"Failed to connect to MCP server: {str(e)}")
    
    async def initialize(self, client_name: str = "python-client", client_version: str = "1.0.0", timeout: int = 30) -> Dict[str, Any]:
        """Initialize the connection with the MCP server."""
        if not self.session.is_connected:
            logger.error("Cannot initialize: Not connected")
            raise ConnectionError("Cannot initialize: Not connected")
        
        try:
            # Skip the initialize method as our server doesn't support it
            # Just mark as initialized and proceed
            self.is_initialized = True
            self.session.is_initialized = True
            logger.info("MCP connection initialized (skipped initialize method)")
            return {}
            
        except Exception as e:
            logger.error(f"Failed to initialize MCP connection: {str(e)}", exc_info=True)
            raise ConnectionError(f"Failed to initialize MCP connection: {str(e)}")
    
    async def list_tools(self, timeout: int = 30) -> List[Dict[str, Any]]:
        """List the available tools on the MCP server."""
        if not self.is_initialized:
            logger.error("Cannot list tools: Not initialized")
            raise ConnectionError("Cannot list tools: Not initialized")
        
        try:
            response = await self.session.send_request("tools/list", {}, timeout=timeout)
            # Handle both possible response formats
            if isinstance(response.get("result"), list):
                tools = response.get("result", [])
            else:
                tools = response.get("result", {}).get("tools", [])
            self._available_tools = tools
            logger.debug(f"Retrieved {len(tools)} tools from server")
            return tools
        except Exception as e:
            logger.error(f"Failed to list tools: {str(e)}", exc_info=True)
            raise MCPError(f"Failed to list tools: {str(e)}")
    
    async def call_tool(self, tool_name: str, arguments: Dict[str, Any], timeout: int = 60) -> Any:
        """Call a tool on the MCP server."""
        if not self.is_initialized:
            logger.error("Cannot call tool: Not initialized")
            raise ConnectionError("Cannot call tool: Not initialized")
        
        # Create the request parameters
        params = {
            "name": tool_name,
            "arguments": arguments
        }
        
        # Send the request
        try:
            response = await self.session.send_request("callTool", params, timeout=timeout)
            
            # Extract the result
            if "result" in response:
                return response["result"]
            else:
                # Return an empty array as a fallback for empty responses
                logger.warning(f"Tool call returned no result: {tool_name}")
                return []
                
        except Exception as e:
            logger.error(f"Error calling tool {tool_name}: {str(e)}")
            raise ToolError(f"Error calling tool {tool_name}: {str(e)}")
    
    async def close(self, timeout: int = 5):
        """Close the connection to the MCP server."""
        logger.debug("Closing MCP client connection")
        
        # Create a task with a timeout to close the session
        try:
            close_task = asyncio.create_task(self._close_session())
            await asyncio.wait_for(close_task, timeout=timeout)
        except asyncio.TimeoutError:
            logger.warning(f"Session close timed out after {timeout}s")
        except Exception as e:
            logger.error(f"Error during session close: {str(e)}", exc_info=True)
        finally:
            self.is_initialized = False
            logger.info("MCP client closed")
    
    async def _close_session(self):
        """Helper method to close the session asynchronously."""
        # Since session.close is not async, run it in a thread
        loop = asyncio.get_event_loop()
        await loop.run_in_executor(None, self.session.close)


class Agent:
    """
    OpsRamp AI Agent that uses LLM to understand requests and interact with MCP tools.
    """
    
    def __init__(
        self, 
        server_url: str, 
        llm_provider: str = "openai",
        openai_api_key: Optional[str] = None,
        anthropic_api_key: Optional[str] = None,
        model: Optional[str] = None,
        connection_timeout: int = 60,
        env_file: Optional[str] = None,
        simple_mode: bool = False,
        request_timeout: int = 30
    ):
        """
        Initialize the OpsRamp agent.
        
        Args:
            server_url: The URL of the MCP server
            llm_provider: The LLM provider to use ('openai' or 'anthropic')
            openai_api_key: OpenAI API key (can also be set via OPENAI_API_KEY env var or .env file)
            anthropic_api_key: Anthropic API key (can also be set via ANTHROPIC_API_KEY env var or .env file)
            model: The model to use (defaults to gpt-4 for OpenAI and claude-2 for Anthropic)
            connection_timeout: Connection timeout in seconds
            env_file: Path to .env file containing config variables
            simple_mode: Whether to run in simple mode without MCP connection
            request_timeout: Timeout in seconds for JSON-RPC requests
        """
        self.server_url = server_url
        self.llm_provider = llm_provider.lower()
        self.connection_timeout = connection_timeout
        self.request_timeout = request_timeout
        self.simple_mode = simple_mode
        self._initialized = False
        self.tools = []
        
        # Load environment variables from file if provided
        if env_file:
            load_env_from_file(env_file)
        
        # Get API keys from arguments, environment variables or .env file
        self.openai_api_key, self.anthropic_api_key = get_api_keys(openai_api_key, anthropic_api_key)
        
        # Set default model based on provider
        if model is None:
            if self.llm_provider == "openai":
                self.model = "gpt-4"
            elif self.llm_provider == "anthropic":
                self.model = "claude-2"
            else:
                raise ValueError(f"Unsupported LLM provider: {self.llm_provider}")
        else:
            self.model = model
            
        # Initialize LLM clients
        self.openai_client = None
        self.anthropic_client = None
        
        if self.llm_provider == "openai":
            if not self.openai_api_key:
                raise ValueError("OpenAI API key is required for OpenAI provider")
            try:
                import openai
                self.openai_client = openai.OpenAI(api_key=self.openai_api_key)
            except ImportError:
                raise ValueError("openai package is required for OpenAI provider")
                
        elif self.llm_provider == "anthropic":
            if not self.anthropic_api_key:
                raise ValueError("Anthropic API key is required for Anthropic provider")
            try:
                import anthropic
                self.anthropic_client = anthropic.Anthropic(api_key=self.anthropic_api_key)
            except ImportError:
                raise ValueError("anthropic package is required for Anthropic provider")
        
        # Create MCP client if not in simple mode
        if not self.simple_mode:
            self.mcp_client = MCPClient(
                server_url=self.server_url, 
                auto_connect=False,
                connection_timeout=self.connection_timeout,
                request_timeout=self.request_timeout
            )
        else:
            self.mcp_client = None
        
        self.conversation_history = []
        
        # In simple mode, initialize with mock tools immediately
        if self.simple_mode:
            self.tools = [
                {
                    "name": "integrations",
                    "description": "Manage OpsRamp integrations",
                    "parameters": {
                        "action": "The action to perform (list, get, create, update, delete, enable, disable)",
                        "id": "Integration ID for get, update, delete, enable, disable actions"
                    }
                }
            ]
            self._initialized = True
    
    async def connect(self) -> None:
        """
        Connect to the MCP server and initialize the client.
        """
        if self.simple_mode:
            logger.info("Running in simple mode, skipping MCP connection")
            return
            
        try:
            self.mcp_client.connect()
            await self.mcp_client.initialize(
                client_name="opsramp-ai-agent", 
                client_version="1.0.0",
                timeout=self.request_timeout
            )
            
            # Try to get tools but don't fail if listTools is not available
            try:
                self.tools = await self.mcp_client.list_tools(timeout=self.request_timeout)
                self._initialized = True
                logger.info(f"Connected to MCP server with {len(self.tools)} tools available")
            except MCPError as e:
                # If listTools fails, set default tools based on server capability
                logger.warning(f"Failed to list tools: {str(e)}")
                logger.info("Continuing with default tools")
                self._initialized = True
                
                # Set default tools
                self.tools = [{
                    "name": "integrations",
                    "description": self.integrations_tool_description,
                    "parameters": {
                        "action": {
                            "type": "string",
                            "description": "The action to perform: list, get, getDetailed, create, update, delete, enable, disable, listTypes, getType"
                        },
                        "id": {
                            "type": "string",
                            "description": "The ID of the integration"
                        },
                        "type": {
                            "type": "string",
                            "description": "The type of integration"
                        },
                        "filter": {
                            "type": "object",
                            "description": "Filter criteria for listing integrations"
                        }
                    }
                }]
                
        except MCPError as e:
            logger.error(f"Connection error: {str(e)}")
            raise
    
    async def direct_call_tool(self, tool_name: str, arguments: Dict[str, Any]) -> Any:
        """
        Directly call a tool without relying on list_tools.
        Useful when we know the tool exists but list_tools is not supported.
        
        Args:
            tool_name: Name of the tool to call
            arguments: Tool arguments
            
        Returns:
            Tool execution result
        """
        try:
            logger.info(f"Direct tool call: {tool_name} with args: {arguments}")
            
            # If in simple mode, return mock data
            if self.simple_mode:
                if tool_name == "integrations":
                    return await self._mock_integrations_call(arguments)
                else:
                    return {"result": f"Mock result for {tool_name}", "arguments": arguments}
                
            # Call the tool with the MCP client
            if self.mcp_client and self.mcp_client.is_initialized:
                return await self.mcp_client.call_tool(
                    tool_name=tool_name, 
                    arguments=arguments,
                    timeout=self.request_timeout
                )
            else:
                # For integrations, we can call directly
                if tool_name == "integrations":
                    return await self.direct_call_integrations(arguments)
                else:
                    raise MCPError(f"Tool not available: {tool_name}")
                    
        except Exception as e:
            logger.error(f"Error calling tool {tool_name}: {str(e)}", exc_info=True)
            raise
        
    async def _mock_integrations_call(self, arguments: Dict[str, Any]) -> Any:
        """
        Provide mock responses for integrations tool in simple mode.
        
        Args:
            arguments: Integration tool arguments
            
        Returns:
            Mock integration data
        """
        action = arguments.get("action", "")
        if not action:
            raise ValueError("Action is required for integrations tool")
        
        # Mock integrations for testing
        mock_integrations = [
            {
                "id": "INTG-2ed93041-eb92-40e9-b6b4-f14ad13e54fc",
                "displayName": "hpe-alletra-LabRat",
                "category": "SDK APP",
                "status": "Installed",
                "app": "hpe-alletra",
                "version": "7.0.0",
                "updateAvailable": True,
                "state": "Deployed",
                "installedBy": "user-XXXXX@example.com",
                "installedTime": "2025-02-18T15:34:32+0000",
                "modifiedTime": "2025-02-18T17:02:06+0000"
            },
            {
                "id": "INTG-f9e5d2aa-ee17-4e32-9251-493566ebdfca",
                "displayName": "redfish-server-LabRat",
                "category": "SDK APP",
                "status": "Installed",
                "app": "redfish-server",
                "version": "7.0.0",
                "updateAvailable": True,
                "state": "Deployed",
                "installedBy": "user-XXXXX@example.com",
                "installedTime": "2025-02-18T15:25:00+0000",
                "modifiedTime": "2025-02-18T15:25:45+0000"
            },
            {
                "id": "INTG-00ee85e2-1f84-4fc1-8cf8-5277ae6980dd",
                "displayName": "vcenter-58.51",
                "category": "COMPUTE_INTEGRATION",
                "status": "enabled",
                "ipAddress": "10.54.58.51",
                "installedBy": "user-XXXXX@example.com",
                "installedTime": "2025-02-18T15:40:23+0000",
                "modifiedTime": "2025-02-18T16:50:52+0000"
            }
        ]
        
        if action == "list":
            return mock_integrations
        elif action == "get" and "id" in arguments:
            integration_id = arguments["id"]
            for integration in mock_integrations:
                if integration["id"] == integration_id:
                    return integration
            raise ValueError(f"Integration with ID {integration_id} not found")
        else:
            return {"action": action, "status": "success", "message": f"Mock integration {action} operation"}
    
    async def direct_call_integrations(self, arguments: Dict[str, Any]) -> Any:
        """
        Make a direct HTTP call to the integrations endpoint.
        
        Args:
            arguments: Integration tool arguments
            
        Returns:
            Integration API result
        """
        import aiohttp
        
        action = arguments.get("action", "")
        if not action:
            raise ValueError("Action is required for integrations tool")
            
        # Session ID is required for the request
        if not self.mcp_client.session.session_id:
            raise MCPError("No active session")
            
        # Construct the URL for the integrations tool
        url = f"{self.mcp_client.server_url}/integrations"
        session_id = self.mcp_client.session.session_id
        url = f"{url}?sessionId={session_id}"
        
        logger.info(f"Making direct call to integrations API: {url} with action {action}")
        
        try:
            # Different endpoints based on the action
            async with aiohttp.ClientSession() as session:
                if action == "list":
                    async with session.get(url, timeout=30) as response:
                        if response.status != 200:
                            response_text = await response.text()
                            raise MCPError(f"Request failed with status {response.status}: {response_text}")
                        
                        response_data = await response.json()
                        logger.info(f"Received direct integrations list response: {str(response_data)[:100]}...")
                        return response_data
                
                elif action == "get" and "id" in arguments:
                    integration_id = arguments["id"]
                    get_url = f"{url}/{integration_id}"
                    async with session.get(get_url, timeout=30) as response:
                        if response.status != 200:
                            response_text = await response.text()
                            raise MCPError(f"Request failed with status {response.status}: {response_text}")
                        
                        response_data = await response.json()
                        return response_data
                
                else:
                    raise ValueError(f"Unsupported action for direct integrations call: {action}")
                
        except asyncio.TimeoutError:
            raise MCPError(f"Request timed out after 30s: integrations {action}")
        except Exception as e:
            logger.error(f"Error in direct integrations call: {str(e)}", exc_info=True)
            raise MCPError(f"Error in direct integrations call: {str(e)}")

    async def chat(self, message: str) -> str:
        """
        Process a user message through the LLM and execute any tool calls.
        
        Args:
            message: The user's message
            
        Returns:
            The agent's response
        """
        if not self._initialized:
            await self.connect()
        
        # Add user message to history
        self.conversation_history.append({"role": "user", "content": message})
        
        # Fast path for common queries that we can handle directly
        lower_msg = message.lower()
        if "list all integration" in lower_msg or "show me the integration" in lower_msg or "what integration" in lower_msg:
            try:
                # Direct call to integrations tool with list action
                result = await self.direct_call_tool("integrations", {"action": "list"})
                
                # Format the result
                integrations = []
                if isinstance(result, list):
                    integrations = result
                elif isinstance(result, dict) and "results" in result:
                    integrations = result["results"]
                
                if integrations:
                    # Create a user-friendly response with integration details
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
                else:
                    response = "No integrations found."
                
                self.conversation_history.append({"role": "assistant", "content": response})
                return response
            except Exception as e:
                logger.error(f"Fast path for integrations list failed: {str(e)}")
                # Continue with normal flow if fast path fails
        
        # Create system prompt with available tools
        system_prompt = self._create_system_prompt()
        
        # Get LLM response
        llm_response = await self._get_llm_response(system_prompt)
        
        # Process tool calls if any
        if self._has_tool_calls(llm_response):
            tool_results = await self._process_tool_calls(llm_response)
            
            # Get final response
            final_response = await self._get_final_response(tool_results)
            self.conversation_history.append({"role": "assistant", "content": final_response})
            return final_response
        else:
            # No tool calls, just return the LLM response
            content = self._extract_content(llm_response)
            self.conversation_history.append({"role": "assistant", "content": content})
            return content
    
    def _create_system_prompt(self) -> str:
        """Create the system prompt with available tools."""
        tools_json = json.dumps(self.tools, indent=2)
        
        # Define integration tool description with supported actions
        integrations_tool_description = """COMPREHENSIVE INTEGRATIONS TOOL EXPERTISE:

The "integrations" tool is your primary interface for managing HPE OpsRamp integrations. It supports these actions:

=== DISCOVERY & LISTING OPERATIONS ===
1. "list" - Lists all integrations in the environment
   Use for: "show me all integrations", "what integrations do we have", "list our integrations"
   Example: {"name": "integrations", "arguments": {"action": "list"}}

2. "listTypes" - Lists all available integration types/templates
   Use for: "what integration types are available", "show integration templates", "what can I integrate with"
   Example: {"name": "integrations", "arguments": {"action": "listTypes"}}

=== DETAILED INSPECTION OPERATIONS ===
3. "get" - Get basic info about a specific integration (name, status, type)
   Required: "id" (integration ID)
   Use for: Quick status checks, basic integration info
   Example: {"name": "integrations", "arguments": {"action": "get", "id": "INTG-2ed93041-eb92-40e9-b6b4-f14ad13e54fc"}}

4. "getDetailed" - Get comprehensive integration details (resources, metrics, alerts, discovery runs, full config)
   Required: "id" (integration ID) 
   Use for: Deep troubleshooting, complete integration analysis, resource inventory
   Example: {"name": "integrations", "arguments": {"action": "getDetailed", "id": "INTG-2ed93041-eb92-40e9-b6b4-f14ad13e54fc"}}

5. "getType" - Get detailed info about an integration type/template
   Required: "id" (integration type ID like "vcenter", "hpe-alletra", "redfish-server")
   Use for: Understanding integration capabilities, configuration requirements
   Example: {"name": "integrations", "arguments": {"action": "getType", "id": "vcenter"}}

=== LIFECYCLE MANAGEMENT OPERATIONS ===
6. "create" - Create a new integration
   Required: "config" (object with name, type, credentials, settings)
   Use for: Adding new integrations to monitor additional infrastructure
   Example: {"name": "integrations", "arguments": {"action": "create", "config": {"name": "New VMware", "type": "vcenter", "host": "10.1.1.100"}}}

7. "update" - Update existing integration configuration
   Required: "id" (integration ID) and "config" (updated configuration object)
   Use for: Modifying credentials, changing settings, updating configurations
   Example: {"name": "integrations", "arguments": {"action": "update", "id": "INTG-123", "config": {"credentials": "new_password"}}}

8. "delete" - Permanently remove an integration
   Required: "id" (integration ID)
   Use for: Decommissioning integrations, cleanup
   Example: {"name": "integrations", "arguments": {"action": "delete", "id": "INTG-123"}}

=== OPERATIONAL CONTROL ===
9. "enable" - Activate a disabled integration (starts monitoring)
   Required: "id" (integration ID)
   Use for: Resuming monitoring, bringing integrations back online
   Example: {"name": "integrations", "arguments": {"action": "enable", "id": "INTG-123"}}

10. "disable" - Deactivate an integration (stops monitoring but preserves config)
    Required: "id" (integration ID)
    Use for: Temporarily stopping monitoring, maintenance mode
    Example: {"name": "integrations", "arguments": {"action": "disable", "id": "INTG-123"}}

=== SMART QUERY INTERPRETATION ===
When users ask questions, intelligently choose the right action:

• "What do we have?" → use "list"
• "Tell me about X integration" → use "getDetailed" for comprehensive info
• "What's the status of X?" → use "get" for quick status
• "Show me VMware integrations" → use "list" then filter results
• "What can I integrate with?" → use "listTypes"
• "How do I set up X?" → use "getType" to show capabilities
• "Turn off X integration" → use "disable"
• "Start monitoring Y again" → use "enable"

=== CRITICAL: USER AND INSTALLATION DATA ===
When users ask about WHO installed integrations, user emails, or installation information:
• ALWAYS use "getDetailed" to get complete integration details including installedBy, installedTime, etc.
• NEVER give privacy/policy responses - the integrations tool contains real operational data that's meant to be accessible
• Example: "Who installed integrations?" → use "list" first, then "getDetailed" for each to get installedBy info
• Example: "Show me user emails" → use "getDetailed" to extract installedBy/modifiedBy fields

=== COMMON INTEGRATION TYPES ===
- vcenter: VMware vCenter Server monitoring
- hpe-alletra: HPE Alletra storage arrays
- redfish-server: Redfish-compatible servers
- aws-ec2: Amazon EC2 instances
- azure-vm: Microsoft Azure VMs
- And many more available via "listTypes"

=== BEST PRACTICES ===
- Always use "getDetailed" for troubleshooting and comprehensive analysis
- Use "list" + filtering for category-based queries
- Check integration status before enable/disable operations
- Use "getType" to understand requirements before creating integrations"""
        
        return f"""You are an expert AI assistant for OpsRamp IT Operations Management, specializing in HPE OpsRamp Integrations.

You have deep expertise in HPE OpsRamp integrations and can help users manage their integration ecosystem. You understand the complete lifecycle of integrations including discovery, configuration, monitoring, and troubleshooting.

You have access to the following tools through the OpsRamp MCP server:

{tools_json}

{integrations_tool_description}

When the user asks you to perform an action related to integrations, you should:
1. Identify which integration action is appropriate (list, get, getDetailed, etc.)
2. Determine the correct parameters needed for that action
3. Call the integrations tool with the appropriate action and parameters

For tool calls, use the following format:
```tool
{{"name": "tool_name", "arguments": {{"param1": "value1", "param2": "value2"}}}}
```

For example, if the user asks about what integrations they have, respond with:
```tool
{{"name": "integrations", "arguments": {{"action": "list"}}}}
```

If a user mentions a specific integration by ID or name, consider using the "getDetailed" action instead of "get" to provide more comprehensive information.

If you don't need to call a tool, just respond normally with your knowledge of OpsRamp integrations. Maintain a professional, expert tone focused on helping users manage their OpsRamp integration environment effectively.
"""
    
    async def _get_llm_response(self, system_prompt: str) -> Any:
        """Get a response from the LLM."""
        messages = [
            {"role": "system", "content": system_prompt},
            *self.conversation_history
        ]
        
        if self.llm_provider == "openai":
            # Define OpenAI function tools for integrations
            openai_tools = [
                {
                    "type": "function",
                    "function": {
                        "name": "integrations",
                        "description": "Manage OpsRamp integrations with comprehensive actions for discovery, configuration, and lifecycle management",
                        "parameters": {
                            "type": "object",
                            "properties": {
                                "action": {
                                    "type": "string",
                                    "enum": ["list", "get", "getDetailed", "create", "update", "delete", "enable", "disable", "listTypes", "getType"],
                                    "description": "The action to perform: list (all integrations), get (basic info), getDetailed (comprehensive info), create (new), update (modify), delete (remove), enable (activate), disable (deactivate), listTypes (available types), getType (type details)"
                                },
                                "id": {
                                    "type": "string",
                                    "description": "Integration ID - required for get, getDetailed, update, delete, enable, disable actions"
                                },
                                "config": {
                                    "type": "object",
                                    "description": "Configuration object - required for create and update actions"
                                },
                                "filter": {
                                    "type": "object", 
                                    "description": "Filter criteria for listing integrations"
                                }
                            },
                            "required": ["action"]
                        }
                    }
                }
            ]
            
            response = await asyncio.to_thread(
                self.openai_client.chat.completions.create,
                model=self.model,
                messages=messages,
                tools=openai_tools,
                tool_choice="auto"  # Let OpenAI decide when to use tools
            )
            return response
        
        elif self.llm_provider == "anthropic":
            response = await asyncio.to_thread(
                self.anthropic_client.messages.create,
                model=self.model,
                system=system_prompt,
                messages=self.conversation_history
            )
            return response
    
    def _has_tool_calls(self, llm_response: Any) -> bool:
        """Check if the LLM response contains tool calls."""
        if self.llm_provider == "openai":
            return hasattr(llm_response.choices[0].message, 'tool_calls') and llm_response.choices[0].message.tool_calls
        
        elif self.llm_provider == "anthropic":
            # For Anthropic, we need to check for the tool call format in the text
            content = llm_response.content[0].text
            return "```tool" in content
    
    def _extract_tool_calls(self, llm_response: Any) -> List[Dict[str, Any]]:
        """Extract tool calls from the LLM response."""
        tool_calls = []
        
        if self.llm_provider == "openai":
            tool_calls = [
                {
                    "name": tool_call.function.name,
                    "arguments": json.loads(tool_call.function.arguments)
                }
                for tool_call in llm_response.choices[0].message.tool_calls
            ]
        
        elif self.llm_provider == "anthropic":
            content = llm_response.content[0].text
            
            # Extract tool calls from the markdown blocks
            for part in content.split("```tool"):
                if "```" in part:
                    tool_json = part.split("```")[0].strip()
                    try:
                        tool_call = json.loads(tool_json)
                        tool_calls.append(tool_call)
                    except json.JSONDecodeError:
                        pass
        
        # Validate integrations tool parameters
        validated_calls = []
        for call in tool_calls:
            if call["name"] == "integrations":
                arguments = call["arguments"]
                action = arguments.get("action", "")
                
                # Ensure required parameters are present for each action
                if action in ["get", "getDetailed", "update", "delete", "enable", "disable", "getType"]:
                    if "id" not in arguments or not arguments["id"]:
                        logger.warning(f"Missing required 'id' parameter for integrations {action} action")
                        continue
                
                if action in ["create", "update"]:
                    if "config" not in arguments or not arguments["config"]:
                        logger.warning(f"Missing required 'config' parameter for integrations {action} action")
                        continue
                
                # All validations passed
                validated_calls.append(call)
            else:
                # Pass through non-integrations tool calls without validation
                validated_calls.append(call)
        
        return validated_calls
    
    def _extract_content(self, llm_response: Any) -> str:
        """Extract the text content from the LLM response."""
        if self.llm_provider == "openai":
            return llm_response.choices[0].message.content or ""
        
        elif self.llm_provider == "anthropic":
            content = llm_response.content[0].text
            
            # Remove tool call blocks
            for part in content.split("```tool"):
                if "```" in part:
                    content = content.replace(f"```tool{part.split('```')[0]}```", "")
            
            return content.strip()
    
    async def _process_tool_calls(self, llm_response: Any) -> List[Dict[str, Any]]:
        """Process tool calls from the LLM response."""
        tool_calls = self._extract_tool_calls(llm_response)
        results = []
        
        for tool_call in tool_calls:
            name = tool_call["name"]
            arguments = tool_call["arguments"]
            
            try:
                logger.info(f"Calling tool: {name} with arguments: {arguments}")
                
                # Use direct_call_tool for better reliability
                result = await self.direct_call_tool(name, arguments)
                
                results.append({
                    "name": name,
                    "arguments": arguments,
                    "result": result,
                    "status": "success"
                })
            except Exception as e:
                logger.error(f"Tool call failed: {name} - {str(e)}")
                results.append({
                    "name": name,
                    "arguments": arguments,
                    "error": str(e),
                    "status": "error"
                })
        
        return results
    
    async def _get_final_response(self, tool_results: List[Dict[str, Any]]) -> str:
        """Get a final response from the LLM after tool execution."""
        # Add tool results to history
        for result in tool_results:
            result_str = json.dumps(result, indent=2)
            self.conversation_history.append({"role": "system", "content": f"Tool result: {result_str}"})
        
        # Custom system prompt based on which tool was called
        system_prompt = """Based on the tool execution results, provide a clear and helpful response to the user.
Explain what was done and summarize the results in a user-friendly way."""

        # Check if this was an integrations tool call
        for result in tool_results:
            if result["name"] == "integrations":
                action = result["arguments"].get("action", "")
                
                # Enhanced prompt specifically for integrations tool
                system_prompt = """You are an expert in HPE OpsRamp integrations management. 
                
Based on the integration tool execution results, provide a detailed, expert analysis of the results.

For list operations:
- Categorize integrations by type/status where appropriate
- Highlight important metrics like total count, active vs. inactive
- Point out any notable patterns or issues

For detailed integration information:
- Highlight key operational metrics and status
- Summarize resources, metrics, and alerts if present
- Explain the integration's role in the IT ecosystem

Use appropriate technical terminology and maintain a professional, expert tone.
Respond as a knowledgeable HPE OpsRamp integrations specialist would."""
                break
        
        llm_response = await self._get_llm_response(system_prompt)
        return self._extract_content(llm_response)
    
    async def close(self) -> None:
        """Close the agent and the MCP client connection."""
        try:
            await self.mcp_client.close()
        except Exception as e:
            logger.error(f"Error closing MCP client: {str(e)}") 