"""
HPE OpsRamp MCP Client - A Python client for the HPE OpsRamp MCP (Model Context Protocol) server.
"""

__version__ = '1.0.0'
__author__ = 'HPE OpsRamp'

from .client import MCPClient
from .exceptions import MCPError, ConnectionError, SessionError, ToolError

__all__ = ['MCPClient', 'MCPError', 'ConnectionError', 'SessionError', 'ToolError'] 