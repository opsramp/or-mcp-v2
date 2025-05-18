"""
Test configuration for MCP client tests.
"""

import os
import logging
from pathlib import Path

# Base configuration
DEFAULT_SERVER_URL = "http://localhost:8080"
DEFAULT_CONNECTION_TIMEOUT = 10
DEFAULT_REQUEST_TIMEOUT = 30

# Get configuration from environment variables
SERVER_URL = os.environ.get("MCP_SERVER_URL", DEFAULT_SERVER_URL)
DEBUG = os.environ.get("DEBUG", "").lower() in ("true", "1", "yes")
CONNECTION_TIMEOUT = int(os.environ.get("CONNECTION_TIMEOUT", DEFAULT_CONNECTION_TIMEOUT))
REQUEST_TIMEOUT = int(os.environ.get("REQUEST_TIMEOUT", DEFAULT_REQUEST_TIMEOUT))

# Integration test flags
AUTO_START_SERVER = os.environ.get("AUTO_START_SERVER", "").lower() in ("true", "1", "yes")
INTEGRATION_TEST_PORT = int(os.environ.get("INTEGRATION_TEST_PORT", 8080))

# Project paths
REPO_ROOT = Path(__file__).resolve().parents[4]  # Four levels up from here
CLIENT_DIR = REPO_ROOT / "client" / "python"
SERVER_CMD = REPO_ROOT / "cmd" / "server" / "main.go"
OUTPUT_DIR = CLIENT_DIR / "test_output"
LOG_DIR = OUTPUT_DIR / "logs"

# Create output directories
OUTPUT_DIR.mkdir(exist_ok=True)
LOG_DIR.mkdir(exist_ok=True)

# Logging configuration
def configure_logging(level=None):
    """Configure logging for tests."""
    if level is None:
        level = logging.DEBUG if DEBUG else logging.INFO
        
    logging.basicConfig(
        level=level,
        format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
        handlers=[
            logging.StreamHandler(),
            logging.FileHandler(LOG_DIR / "test.log")
        ]
    )
    
    # Silence some verbose loggers
    logging.getLogger('urllib3').setLevel(logging.WARNING)
    logging.getLogger('sseclient').setLevel(logging.INFO) 