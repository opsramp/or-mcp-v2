"""
Utility to start and stop the MCP server during tests.
"""

import os
import sys
import subprocess
import time
import signal
import atexit
import requests
from pathlib import Path
import logging

logger = logging.getLogger(__name__)

class ServerRunner:
    """
    Helper class to manage the MCP server process during tests.
    """
    
    def __init__(self, port=8080, debug=True):
        """Initialize the server runner."""
        self.port = port
        self.debug = debug
        self.process = None
        self.server_url = f"http://localhost:{port}"
        
        # Find the repository root
        current_dir = Path(__file__).resolve().parent
        while not (current_dir / "go.mod").exists() and current_dir != current_dir.parent:
            current_dir = current_dir.parent
        
        if not (current_dir / "go.mod").exists():
            raise RuntimeError("Could not find repository root (go.mod file)")
            
        self.repo_root = current_dir
        self.server_cmd = str(current_dir / "cmd" / "server" / "main.go")
        
        # Register the cleanup function
        atexit.register(self.stop)
    
    def start(self, timeout=10):
        """
        Start the MCP server.
        
        Args:
            timeout: How long to wait for the server to start in seconds
            
        Returns:
            True if server started successfully, False otherwise
        """
        # Check if the port is already in use
        if self._check_port_in_use():
            logger.info(f"Port {self.port} already in use, assuming server is running")
            return True
            
        # Start the server
        logger.info(f"Starting MCP server on port {self.port}")
        env = os.environ.copy()
        env["PORT"] = str(self.port)
        
        if self.debug:
            env["DEBUG"] = "true"
            
        try:
            self.process = subprocess.Popen(
                ["go", "run", self.server_cmd],
                cwd=str(self.repo_root),
                env=env,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE
            )
            
            # Wait for the server to start
            start_time = time.time()
            while time.time() - start_time < timeout:
                try:
                    response = requests.get(f"{self.server_url}/health", timeout=1)
                    if response.status_code == 200:
                        logger.info(f"Server started successfully (PID: {self.process.pid})")
                        return True
                except requests.RequestException:
                    pass
                    
                time.sleep(0.5)
                
                # Check if process is still running
                if self.process.poll() is not None:
                    stdout, stderr = self.process.communicate()
                    logger.error(f"Server failed to start, exited with code {self.process.returncode}")
                    logger.error(f"STDOUT: {stdout.decode('utf-8')}")
                    logger.error(f"STDERR: {stderr.decode('utf-8')}")
                    return False
            
            # If we got here, server didn't start in time
            logger.error(f"Server failed to start within {timeout} seconds")
            self.stop()
            return False
            
        except Exception as e:
            logger.error(f"Error starting server: {e}")
            return False
    
    def stop(self):
        """Stop the server if it's running."""
        if self.process and self.process.poll() is None:
            logger.info(f"Stopping MCP server (PID: {self.process.pid})")
            try:
                self.process.terminate()
                self.process.wait(timeout=5)
            except subprocess.TimeoutExpired:
                logger.warning("Server didn't terminate gracefully, killing it")
                self.process.kill()
            
            self.process = None
            logger.info("Server stopped")
    
    def _check_port_in_use(self):
        """Check if the port is already in use."""
        import socket
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
            return s.connect_ex(('localhost', self.port)) == 0
    
    def check_health(self):
        """Check if the server is healthy."""
        try:
            response = requests.get(f"{self.server_url}/health", timeout=2)
            return response.status_code == 200
        except requests.RequestException:
            return False
    
    def __enter__(self):
        """Start the server when entering a context."""
        self.start()
        return self
        
    def __exit__(self, exc_type, exc_val, exc_tb):
        """Stop the server when exiting a context."""
        self.stop()


def get_server_runner(port=8080, debug=True):
    """Get a server runner instance."""
    return ServerRunner(port=port, debug=debug)


if __name__ == "__main__":
    # Simple test of the ServerRunner
    logging.basicConfig(level=logging.INFO)
    runner = ServerRunner()
    
    if runner.start():
        print("Server started successfully")
        print("Health check:", runner.check_health())
        time.sleep(2)
        runner.stop()
    else:
        print("Failed to start server") 