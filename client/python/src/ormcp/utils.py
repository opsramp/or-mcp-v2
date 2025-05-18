"""
Utility functions for the HPE OpsRamp MCP client.
"""

import json
import uuid
from typing import Dict, Any, Optional

def generate_request_id() -> str:
    """Generate a unique ID for a JSON-RPC request."""
    return str(uuid.uuid4())


def create_jsonrpc_request(method: str, params: Optional[Dict[str, Any]] = None, request_id: Optional[str] = None) -> Dict[str, Any]:
    """
    Create a JSON-RPC 2.0 request object.
    
    Args:
        method: The method to call
        params: The parameters to pass to the method
        request_id: A unique ID for the request (generated if not provided)
        
    Returns:
        A JSON-RPC 2.0 request object
    """
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
    """
    Parse a JSON-RPC 2.0 response.
    
    Args:
        response_text: The JSON-RPC response as text
        
    Returns:
        The parsed JSON-RPC response
        
    Raises:
        ValueError: If the response is not valid JSON
    """
    try:
        return json.loads(response_text)
    except json.JSONDecodeError as e:
        raise ValueError(f"Invalid JSON response: {response_text}") from e


def parse_session_id_from_sse(event_data: str) -> Optional[str]:
    """
    Parse a session ID from an SSE event.
    
    Args:
        event_data: The event data from an SSE connection
        
    Returns:
        The session ID, or None if not found
    """
    if "sessionId=" in event_data:
        # Extract session ID from the format: /message?sessionId=<uuid>
        parts = event_data.split("sessionId=")
        if len(parts) > 1:
            session_id = parts[1].strip()
            return session_id
    return None 