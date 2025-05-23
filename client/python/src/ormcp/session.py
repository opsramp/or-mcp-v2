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
from aiohttp_sse_client import client as sse_client

from .exceptions import SessionError, JSONRPCError
from .utils import parse_session_id_from_sse, create_jsonrpc_request, parse_jsonrpc_response

logger = logging.getLogger(__name__)


class AsyncSSEClient:
    """
    An asyncio-based SSE client that uses aiohttp-sse-client.
    
    This maintains a persistent connection to the server and
    processes events in real-time, ensuring that the session
    is properly registered with the server.
    """
    
    def __init__(self, url, headers=None, timeout=10):
        """Initialize the SSE client."""
        self.url = url
        self.headers = headers or {}
        self.timeout = timeout
        self.client = None
        self.session = None
        self.is_connected = False
        self.event_queue = queue.Queue()
        self._keep_running = True
        self._loop = None
        self._task = None
        self._reconnect_count = 0
        self._max_reconnect_attempts = 5
        self._reconnect_delay = 1.0  # Initial reconnect delay in seconds
        self._last_activity = time.time()
        
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
            
            # Create a new asyncio loop in the current thread if one doesn't exist
            try:
                self._loop = asyncio.get_event_loop()
            except RuntimeError:
                self._loop = asyncio.new_event_loop()
                asyncio.set_event_loop(self._loop)
            
            # Reset reconnect count
            self._reconnect_count = 0
            self._last_activity = time.time()
            
            # Run the connection task
            self._task = self._loop.create_task(self._connect_and_process())
            
            # Start a thread to run the asyncio loop if it's not already running
            if not self._loop.is_running():
                def run_loop():
                    try:
                        self._loop.run_until_complete(self._task)
                    except Exception as e:
                        logger.error(f"Error in SSE event loop: {str(e)}")
                
                thread = threading.Thread(target=run_loop, daemon=True)
                thread.start()
            
            # Wait a bit to ensure connection is established
            time.sleep(0.5)
            
            logger.debug("SSE connection requested")
            return True
            
        except Exception as e:
            logger.error(f"Failed to connect to SSE endpoint: {str(e)}", exc_info=True)
            return False
    
    async def _connect_and_process(self):
        """Connect to SSE endpoint and process events."""
        try:
            logger.debug("Starting SSE connection and event processing")
            
            # Setup client session with reconnection
            self.session = aiohttp.ClientSession()
            self.client = sse_client.EventSource(
                self.url,
                session=self.session,
                headers=self.headers,
                timeout=self.timeout,
                reconnection_time=self._reconnect_delay,  # Reconnect after delay if disconnected
                max_connect_retry=self._max_reconnect_attempts
            )
            
            # Connect to the SSE endpoint
            await self.client.connect()
            self.is_connected = True
            self._last_activity = time.time()
            logger.debug("SSE connection established")
            
            # Queue a connection success event
            self.event_queue.put({
                'id': 'connection-established',
                'event': 'connection',
                'data': 'Connection established',
                'timestamp': time.time()
            })
            
            # Start the health checker
            health_check_task = asyncio.create_task(self._check_connection_health())
            
            # Process events
            while self._keep_running:
                try:
                    # Use async iterator instead of get_event()
                    async for event in self.client:
                        if event:
                            # Update last activity timestamp
                            self._last_activity = time.time()
                            
                            # Process the event
                            event_data = {
                                'id': event.last_event_id,
                                'event': event.type or 'message',
                                'data': event.data,
                                'timestamp': time.time()
                            }
                            
                            # Put the event in the queue for consumers
                            self.event_queue.put(event_data)
                            
                            # Log detailed event information at debug level
                            logger.debug(f"SSE event received: {event_data['event']} - {event_data['data'][:200]}")
                            
                        if not self._keep_running:
                            break
                            
                except asyncio.CancelledError:
                    break
                except Exception as e:
                    logger.error(f"Error processing SSE event: {str(e)}")
                    # Mark as disconnected to trigger reconnect
                    self.is_connected = False
                    await asyncio.sleep(0.1)  # Prevent tight loop in case of errors
                    
            # Cancel health check task
            if health_check_task and not health_check_task.done():
                health_check_task.cancel()
                    
        except Exception as e:
            if self._keep_running:  # Only log if not intentionally stopped
                logger.error(f"SSE connection error: {str(e)}", exc_info=True)
                self.is_connected = False
                
                # Queue a connection error event
                self.event_queue.put({
                    'id': 'connection-error',
                    'event': 'error',
                    'data': f"Connection error: {str(e)}",
                    'timestamp': time.time()
                })
                
                # Attempt to reconnect if we haven't exceeded max attempts
                if self._reconnect_count < self._max_reconnect_attempts:
                    self._reconnect_count += 1
                    reconnect_delay = min(30, self._reconnect_delay * (2 ** (self._reconnect_count - 1)))  # Exponential backoff
                    logger.info(f"Attempting to reconnect in {reconnect_delay:.1f} seconds (attempt {self._reconnect_count}/{self._max_reconnect_attempts})")
                    
                    # Queue a reconnection event
                    self.event_queue.put({
                        'id': 'reconnection-attempt',
                        'event': 'reconnecting',
                        'data': f"Reconnection attempt {self._reconnect_count}/{self._max_reconnect_attempts}",
                        'timestamp': time.time()
                    })
                    
                    await asyncio.sleep(reconnect_delay)
                    if self._keep_running:
                        # Create a new task for reconnection
                        asyncio.create_task(self._connect_and_process())
                else:
                    logger.error(f"Max reconnection attempts ({self._max_reconnect_attempts}) reached, giving up")
                    # Queue a connection failed event
                    self.event_queue.put({
                        'id': 'connection-failed',
                        'event': 'error',
                        'data': f"Connection failed after {self._max_reconnect_attempts} attempts",
                        'timestamp': time.time()
                    })
        finally:
            self.is_connected = False
            if self.client and not self._task.cancelled():
                await self._cleanup()
            logger.debug("SSE event processing ended")
    
    async def _check_connection_health(self):
        """Periodically check if the connection is healthy."""
        while self._keep_running and not self._task.cancelled():
            try:
                # Check if we've had activity recently (30 seconds)
                inactive_time = time.time() - self._last_activity
                if inactive_time > 30 and self.is_connected:
                    logger.warning(f"No SSE activity for {inactive_time:.1f} seconds, connection may be stale")
                    # Queue a ping event to test the connection
                    self.event_queue.put({
                        'id': 'health-check',
                        'event': 'ping',
                        'data': 'Connection health check',
                        'timestamp': time.time()
                    })
                    
                    # If it's been too long (60+ seconds), consider the connection dead
                    if inactive_time > 60:
                        logger.error(f"Connection appears dead after {inactive_time:.1f} seconds of inactivity")
                        self.is_connected = False
                        # This will trigger reconnection logic in _connect_and_process when it fails
                        if self.client:
                            await self.client.close()
                
                # Sleep for 5 seconds before checking again
                await asyncio.sleep(5)
            except asyncio.CancelledError:
                break
            except Exception as e:
                logger.error(f"Error in connection health check: {str(e)}")
                await asyncio.sleep(5)
    
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
        logger.debug(f"Waiting for event type: {event_type} with timeout: {timeout}s")
        
        while time.time() - start_time < timeout:
            event = self.get_event(timeout=1)
            
            if event:
                logger.debug(f"Received event while waiting: {event['event']} - data: {event['data'][:100]}...")
                if event_type is None or event['event'] == event_type:
                    logger.debug(f"Event matches requested type: {event_type}")
                    return event
                else:
                    logger.debug(f"Event type mismatch: expected {event_type}, got {event['event']}")
            
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
        
        if self._task and not self._task.done():
            # Cancel the task if it's still running
            self._task.cancel()
            
        self.is_connected = False
        
        # If we have our own loop, close it
        if self._loop and hasattr(self._loop, 'is_running') and self._loop.is_running():
            try:
                asyncio.run_coroutine_threadsafe(self._cleanup(), self._loop)
            except Exception as e:
                logger.warning(f"Error closing SSE connection: {str(e)}")
    
    async def _cleanup(self):
        """Clean up resources."""
        if self.client:
            try:
                await self.client.close()
            except Exception as e:
                logger.warning(f"Error closing SSE client: {str(e)}")
        
        if self.session:
            try:
                await self.session.close()
            except Exception as e:
                logger.warning(f"Error closing aiohttp session: {str(e)}")


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
        self.sse_client: Optional[AsyncSSEClient] = None
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
            logger.debug(f"Establishing async SSE connection to {self._get_sse_url()}")
            self.sse_client = AsyncSSEClient(
                self._get_sse_url(),
                timeout=self.connection_timeout
            )
            
            if not self.sse_client.connect():
                raise SessionError("Failed to connect to SSE endpoint")
            
            # Wait for the endpoint event which contains the session ID
            logger.debug("Waiting for endpoint event with session ID...")
            
            # Wait for any event first and log it to debug
            any_event = self.sse_client.wait_for_event(timeout=self.connection_timeout)
            if any_event:
                logger.debug(f"First received event: {any_event['event']} - {any_event['data']}")
                
                # If this happens to be an endpoint event, use it
                if any_event['event'] == 'endpoint':
                    endpoint_event = any_event
                else:
                    # Otherwise, keep waiting specifically for an endpoint event
                    endpoint_event = self.sse_client.wait_for_event('endpoint', timeout=self.connection_timeout)
            else:
                endpoint_event = None
            
            if not endpoint_event:
                # Try to get a session ID from any event we've received
                # Some servers might not use the 'endpoint' event type
                for attempt in range(3):
                    any_event = self.sse_client.get_event(timeout=2)
                    if any_event:
                        logger.debug(f"Found alternate event: {any_event['event']} - {any_event['data']}")
                        session_id = parse_session_id_from_sse(any_event['data'])
                        if session_id:
                            logger.info(f"Extracted session ID from alternate event: {session_id}")
                            self.session_id = session_id
                            self.is_connected = True
                            self._start_event_processing()
                            return self.session_id
            
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
                    
                    # Call specific handler for this event type
                    if event_type in self._event_handlers:
                        handler = self._event_handlers[event_type]
                        try:
                            handler(event['data'])
                        except Exception as e:
                            logger.error(f"Error in event handler for {event_type}: {str(e)}")
                    
                    # Call wildcard handler if registered
                    if '*' in self._event_handlers:
                        handler = self._event_handlers['*']
                        try:
                            handler(event['data'])
                        except Exception as e:
                            logger.error(f"Error in wildcard event handler: {str(e)}")
                    
                    # Handle special event types
                    if event_type == 'ping':
                        # Server sending ping to keep connection alive
                        logger.debug("Received ping event")
                    elif event_type == 'error':
                        # Handle error events
                        logger.warning(f"Received error event: {event['data']}")
                    
                    # Log the event at debug level
                    logger.debug(f"Processed event: {event_type}")
                
                # Check connection status
                if not self.sse_client.is_connected and self.is_connected:
                    logger.warning("SSE connection lost, session may be invalid")
                    # Don't immediately set is_connected to False
                    # The AsyncSSEClient will attempt to reconnect
                
                # Small sleep to prevent tight loop
                time.sleep(0.1)
            
            # If we exit the loop, connections may have been lost
            if self.is_connected and (not self.sse_client or not self.sse_client.is_connected):
                logger.warning("Event processing stopped due to lost connection")
                
                # Try to reconnect once
                try:
                    logger.info("Attempting to reconnect SSE client...")
                    if self.sse_client:
                        self.sse_client.close()  # Clean up old client
                    
                    # Create new client and connect
                    self.sse_client = AsyncSSEClient(
                        self._get_sse_url(), 
                        timeout=self.connection_timeout
                    )
                    if self.sse_client.connect():
                        logger.info("SSE client reconnected successfully")
                        # Restart event processing
                        thread = threading.Thread(target=process_events, daemon=True)
                        thread.start()
                        return  # Exit this thread, new one will take over
                    else:
                        logger.error("Failed to reconnect SSE client")
                        self.is_connected = False
                except Exception as e:
                    logger.error(f"Error reconnecting SSE client: {str(e)}")
                    self.is_connected = False
            
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
        
        # Check if SSE connection is still active
        if not self.sse_client or not self.sse_client.is_connected:
            logger.warning("SSE connection appears to be lost before sending request")
            
            # Try to reconnect first before failing
            try:
                logger.info("Attempting to reconnect before sending request...")
                if self.sse_client:
                    self.sse_client.close()  # Clean up old client
                
                # Create new client and connect
                self.sse_client = AsyncSSEClient(
                    self._get_sse_url(), 
                    timeout=self.connection_timeout
                )
                if not self.sse_client.connect():
                    raise SessionError("SSE connection lost and reconnection failed")
                logger.info("SSE client reconnected successfully")
                # Restart event processing
                self._start_event_processing()
            except Exception as e:
                logger.error(f"Reconnection failed: {str(e)}")
                self.is_connected = False
                raise SessionError(f"SSE connection lost and reconnection failed: {str(e)}")
        
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
                        if (isinstance(error, dict) and 
                            error.get("code") == -32602 and 
                            "Invalid session ID" in error.get("message", "")):
                            logger.error("Server rejected session ID - this may indicate the SSE connection was not properly recognized")
                            # Try to reconnect
                            self.is_connected = False
                            new_session_id = self.connect()
                            logger.info(f"Successfully reconnected with new session ID: {new_session_id}")
                            # Retry the request once
                            return await self.send_request(method, params, timeout)
                        
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
            event_type: The type of event to handle, or '*' for all events
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