"""
Exception classes for the HPE OpsRamp MCP client.
"""

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


class JSONRPCError(MCPError):
    """Raised when there's a JSON-RPC protocol error."""
    
    # Standard JSON-RPC error codes
    PARSE_ERROR = -32700
    INVALID_REQUEST = -32600
    METHOD_NOT_FOUND = -32601
    INVALID_PARAMS = -32602
    INTERNAL_ERROR = -32603
    
    @classmethod
    def from_response(cls, response):
        """Create an exception from a JSON-RPC error response."""
        error = response.get('error', {})
        message = error.get('message', 'Unknown JSON-RPC error')
        code = error.get('code', 0)
        data = error.get('data')
        return cls(message, code, data)


class ResourceError(MCPError):
    """Raised when there's an error with resource management operations."""
    
    # Resource-specific error codes
    RESOURCE_NOT_FOUND = "RESOURCE_NOT_FOUND"
    RESOURCE_ACCESS_DENIED = "RESOURCE_ACCESS_DENIED"
    RESOURCE_INVALID_PARAMS = "RESOURCE_INVALID_PARAMS"
    RESOURCE_API_ERROR = "RESOURCE_API_ERROR"
    
    def __init__(self, message, code=None, data=None, resource_id=None, action=None):
        super().__init__(message, code, data)
        self.resource_id = resource_id
        self.action = action 