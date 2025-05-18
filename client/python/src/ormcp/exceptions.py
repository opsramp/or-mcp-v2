"""
Exception classes for the HPE OpsRamp MCP client.
"""

class MCPError(Exception):
    """Base exception for all MCP-related errors."""
    
    def __init__(self, message, code=None, data=None):
        super().__init__(message)
        self.code = code
        self.data = data
        
    def __str__(self):
        if self.code:
            return f"{self.args[0]} (code: {self.code})"
        return self.args[0]


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
    """Raised when there's an error in the JSON-RPC protocol."""
    
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