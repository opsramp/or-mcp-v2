"""
Session management for the HPE OpsRamp MCP client.
"""

import asyncio
import logging
import threading
import time
import queue
from typing import Optional, Dict, Any, Callable, List

import aiohttp
import requests
import sseclient

from .exceptions import SessionError, JSONRPCError
from .utils import parse_session_id_from_sse, create_jsonrpc_request, parse_jsonrpc_response

logger = logging.getLogger(__name__)


class BrowserLikeSSEClient:
    """
    An SSE client that behaves like a browser's EventSource.
    
    This maintains a persistent connection to the server and
    processes events in real-time, ensuring that the session
    is properly registered with the server.
    """
    
    def __init__(self, url, headers=None, timeout=10):
        """Initialize the SSE client."""
        self.url = url
        self.headers = headers or {}
        self.timeout = timeout
        self.session = requests.Session()
        self.response = None
        self.is_connected = False
        self.event_queue = queue.Queue()
        self._thread = None
        self._keep_running = True
        
        # Add standard SSE headers
        if 'Accept' not in self.headers:
            self.headers['Accept'] = 'text/event-stream'
        if 'Cache-Control' not in self.headers:
            self.headers['Cache-Control'] = 'no-cache'
    
    def connect(self):
        """
        Establish connection to the SSE endpoint.
        
        Returns:
            True if connected, False otherwise
        """
        if self.is_connected:
            return True
            
        try:
            logger.debug(f"Connecting to SSE endpoint: {self.url}")
            # Open a persistent connection that won't close after receiving data
            self.response = self.session.get(
                self.url, 
                stream=True, 
                headers=self.headers,
                timeout=self.timeout
            )
            
            if self.response.status_code != 200:
                logger.error(f"Failed to connect to SSE endpoint: {self.response.status_code}")
                return False
                
            self.is_connected = True
            
            # Start event processing thread
            self._thread = threading.Thread(
                target=self._process_events,
                daemon=True
            )
            self._keep_running = True
            self._thread.start()
            
            logger.debug("SSE connection established")
            return True
            
        except Exception as e:
            logger.error(f"Failed to connect to SSE endpoint: {str(e)}", exc_info=True)
            return False
    
    def _process_events(self):
        """Process SSE events in a background thread."""
        try:
            logger.debug("Starting SSE event processing")
            client = sseclient.SSEClient(self.response)
            
            for event in client.events():
                if not self._keep_running:
                    break
                
                # Process the event
                event_data = {
                    'id': event.id,
                    'event': event.event,
                    'data': event.data,
                    'timestamp': time.time()
                }
                
                # Put the event in the queue for consumers
                self.event_queue.put(event_data)
                
                # Log detailed event information at debug level
                logger.debug(f"SSE event received: {event.event} - {event.data}")
                
        except Exception as e:
            if self._keep_running:  # Only log if not intentionally stopped
                logger.error(f"SSE connection error: {str(e)}", exc_info=True)
                self.is_connected = False
        finally:
            logger.debug("SSE event processing ended")
    
    def get_event(self, timeout=1):
        """
        Get next event from the queue.
        
        Args:
            timeout: How long to wait for an event
            
        Returns:
            Event dict or None if timeout
        """
        try:
            return self.event_queue.get(timeout=timeout)
        except queue.Empty:
            return None
    
    def wait_for_event(self, event_type=None, timeout=10):
        """
        Wait for a specific type of event.
        
        Args:
            event_type: Type of event to wait for, or None for any event
            timeout: Maximum time to wait in seconds
            
        Returns:
            Event dict or None if timeout
        """
        start_time = time.time()
        
        while time.time() - start_time < timeout:
            event = self.get_event(timeout=1)
            
            if event:
                if event_type is None or event['event'] == event_type:
                    return event
            
            # Check if we're still connected
            if not self.is_connected:
                logger.error("Lost SSE connection while waiting for event")
                return None
            
            # Small sleep to prevent tight loop
            time.sleep(0.1)
            
        logger.warning(f"Timeout waiting for event type: {event_type}")
        return None
    
    def close(self):
        """Close the SSE connection."""
        self._keep_running = False
        
        if self.response:
            try:
                self.response.close()
            except Exception as e:
                logger.warning(f"Error closing SSE response: {str(e)}")
            
        self.is_connected = False
        
        if self._thread and self._thread.is_alive():
            self._thread.join(timeout=2.0)
            if self._thread.is_alive():
                logger.warning("SSE thread did not terminate gracefully")


class MCPSession:
    """
    Manages a session with an MCP server.
    
    This class handles:
    - Establishing an SSE connection using browser-like behavior
    - Getting and maintaining a session ID
    - Sending JSON-RPC requests
    - Processing responses
    """
    
    def __init__(self, base_url: str, connection_timeout: int = 10):
        """
        Initialize the session.
        
        Args:
            base_url: The base URL of the MCP server
            connection_timeout: Timeout in seconds for connections
        """
        self.base_url = base_url.rstrip('/')
        self.connection_timeout = connection_timeout
        self.session_id: Optional[str] = None
        self.sse_client: Optional[BrowserLikeSSEClient] = None
        self.is_connected = False
        self.is_initialized = False
        self._event_handlers: Dict[str, Callable] = {}
        self._received_events: List[Dict] = []
    
    def _get_message_url(self) -> str:
        """Get the URL for sending messages."""
        if not self.session_id:
            raise SessionError("No active session")
        return f"{self.base_url}/message?sessionId={self.session_id}"
    
    def _get_sse_url(self) -> str:
        """Get the URL for SSE connection."""
        return f"{self.base_url}/sse"
    
    def connect(self) -> str:
        """
        Connect to the MCP server and get a session ID using browser-like behavior.
        
        Returns:
            The session ID
            
        Raises:
            SessionError: If connection fails
        """
        if self.is_connected and self.session_id:
            return self.session_id
        
        try:
            # Create and connect browser-like SSE client
            logger.debug(f"Establishing browser-like SSE connection to {self._get_sse_url()}")
            self.sse_client = BrowserLikeSSEClient(
                self._get_sse_url(), 
                timeout=self.connection_timeout
            )
            
            if not self.sse_client.connect():
                raise SessionError("Failed to connect to SSE endpoint")
            
            # Wait for the endpoint event which contains the session ID
            logger.debug("Waiting for endpoint event with session ID...")
            endpoint_event = self.sse_client.wait_for_event('endpoint', timeout=self.connection_timeout)
            
            if not endpoint_event:
                raise SessionError(f"No endpoint event received within {self.connection_timeout}s")
            
            # Parse session ID from the event data
            event_data = endpoint_event['data']
            logger.debug(f"Received endpoint event: {event_data}")
            
            session_id = parse_session_id_from_sse(event_data)
            if not session_id:
                raise SessionError("Failed to extract session ID from event data")
            
            self.session_id = session_id
            self.is_connected = True
            
            # Register event handlers for processing other events
            self._start_event_processing()
            
            logger.info(f"Connected with session ID: {self.session_id}")
            return self.session_id
            
        except Exception as e:
            logger.error(f"Connection error: {str(e)}", exc_info=True)
            self.close()  # Clean up any partial connections
            raise SessionError(f"Failed to connect to MCP server: {str(e)}")
    
    def _start_event_processing(self):
        """Start processing events from the SSE client."""
        def process_events():
            """Background thread to process events."""
            while self.is_connected and self.sse_client and self.sse_client.is_connected:
                event = self.sse_client.get_event(timeout=1)
                if event:
                    # Store event for later processing
                    self._received_events.append(event)
                    
                    # Call event handlers if registered
                    event_type = event['event']
                    if event_type in self._event_handlers:
                        handler = self._event_handlers[event_type]
                        try:
                            handler(event['data'])
                        except Exception as e:
                            logger.error(f"Error in event handler for {event_type}: {str(e)}")
                    
                    # Log the event at debug level
                    logger.debug(f"Processed event: {event_type}")
                
                # Small sleep to prevent tight loop
                time.sleep(0.1)
            
            logger.debug("Event processing stopped")
        
        # Start in background thread
        thread = threading.Thread(target=process_events, daemon=True)
        thread.start()
    
    async def send_request(self, method: str, params: Optional[Dict[str, Any]] = None, timeout: int = 30) -> Dict[str, Any]:
        """
        Send a JSON-RPC request to the MCP server.
        
        Args:
            method: The method to call
            params: The parameters to pass to the method
            timeout: Timeout in seconds for the request
            
        Returns:
            The JSON-RPC response
            
        Raises:
            SessionError: If not connected or session error
            JSONRPCError: If the response contains an error
        """
        if not self.is_connected:
            raise SessionError("Not connected to MCP server")
        
        request = create_jsonrpc_request(method, params)
        message_url = self._get_message_url()
        
        logger.debug(f"Sending request to {message_url}: {request}")
        
        async with aiohttp.ClientSession() as session:
            try:
                # Use timeout for the request
                async with session.post(
                    message_url,
                    json=request,
                    headers={"Content-Type": "application/json"},
                    timeout=aiohttp.ClientTimeout(total=timeout)
                ) as response:
                    response_text = await response.text()
                    logger.debug(f"Received response: {response_text}")
                    
                    response_data = parse_jsonrpc_response(response_text)
                    
                    # Check for errors in the response
                    if "error" in response_data:
                        # Check for "Invalid session ID" error specifically
                        error = response_data.get("error", {})
                        if error.get("code") == -32602 and "Invalid session ID" in error.get("message", ""):
                            logger.error("Server rejected session ID - this may indicate the SSE connection was not properly recognized")
                            # You might try to reconnect or other recovery here
                        
                        raise JSONRPCError.from_response(response_data)
                    
                    return response_data
                
            except asyncio.TimeoutError:
                logger.error(f"Request timed out after {timeout}s: {method}")
                raise SessionError(f"Request timed out after {timeout}s: {method}")
            except aiohttp.ClientError as e:
                logger.error(f"Failed to send request: {str(e)}")
                raise SessionError(f"Failed to send request: {str(e)}")
    
    def register_event_handler(self, event_type: str, handler: Callable):
        """
        Register a handler for SSE events.
        
        Args:
            event_type: The type of event to handle
            handler: The function to call when the event is received
        """
        self._event_handlers[event_type] = handler
    
    def get_received_events(self, event_type: Optional[str] = None) -> List[Dict]:
        """
        Get all received events, optionally filtered by type.
        
        Args:
            event_type: The type of events to return, or None for all events
            
        Returns:
            List of event dictionaries
        """
        if event_type is None:
            return self._received_events.copy()
        else:
            return [e for e in self._received_events if e['event'] == event_type]
    
    def close(self):
        """Close the session and clean up resources."""
        logger.debug("Closing session")
        self.is_connected = False
        self.is_initialized = False
        
        if self.sse_client:
            try:
                logger.debug("Closing SSE client")
                self.sse_client.close()
            except Exception as e:
                logger.warning(f"Error closing SSE client: {str(e)}")
            self.sse_client = None
        
        logger.info("Session closed") 