"""
Utility functions for the HPE OpsRamp MCP client.
"""

import json
import uuid
import re
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
    # Try to parse JSON first if the event data looks like JSON
    if event_data.strip().startswith('{'):
        try:
            data = json.loads(event_data)
            # Look for common session ID fields
            for field in ['sessionId', 'session_id', 'id']:
                if field in data:
                    return str(data[field])
                    
            # Look in nested objects
            if 'session' in data and isinstance(data['session'], dict):
                session = data['session']
                for field in ['id', 'sessionId', 'session_id']:
                    if field in session:
                        return str(session[field])
        except json.JSONDecodeError:
            pass  # Not JSON, try other formats
    
    # Extract session ID from URL-like format: /message?sessionId=<uuid>
    if "sessionId=" in event_data:
        parts = event_data.split("sessionId=")
        if len(parts) > 1:
            # Get everything up to the next & or space or end of string
            session_id = re.split(r'[&\s]', parts[1])[0].strip()
            return session_id
    
    # Look for UUIDs in the event data
    uuid_pattern = r'[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}'
    uuids = re.findall(uuid_pattern, event_data, re.IGNORECASE)
    if uuids:
        return uuids[0]  # Return the first UUID found
        
    return None 