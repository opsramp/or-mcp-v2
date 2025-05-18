#!/bin/bash
# Script to run all HPE OpsRamp MCP client tests
# Official repository: https://github.com/opsramp/or-mcp-v2

set -e

# Get the absolute path of the repository root
REPO_ROOT=$(cd "$(dirname "$0")/../../" && pwd)
CLIENT_DIR="$REPO_ROOT/client/python"
SERVER_CMD="$REPO_ROOT/cmd/server/main.go"
SERVER_PORT=8080

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Function to check if port is in use
is_port_in_use() {
  if command -v lsof > /dev/null; then
    lsof -i:$1 > /dev/null 2>&1
    return $?
  elif command -v netstat > /dev/null; then
    netstat -tuln | grep ":$1 " > /dev/null 2>&1
    return $?
  fi
  return 1 # Assume in use if we can't check
}

# Function to start the server
start_server() {
  echo -e "${YELLOW}Starting MCP server...${NC}"
  
  # Check if port is already in use
  if is_port_in_use $SERVER_PORT; then
    echo -e "${RED}Port $SERVER_PORT is already in use.${NC}"
    echo -e "${YELLOW}Assuming server is already running.${NC}"
    return 0
  fi
  
  # Start server in background
  cd "$REPO_ROOT"
  DEBUG=true go run "$SERVER_CMD" > "$CLIENT_DIR/server.log" 2>&1 &
  SERVER_PID=$!
  
  # Wait for server to start
  echo -n "Waiting for server to start"
  for i in {1..10}; do
    if curl -s http://localhost:$SERVER_PORT/health > /dev/null; then
      echo -e "\n${GREEN}Server started successfully (PID: $SERVER_PID)${NC}"
      echo $SERVER_PID > "$CLIENT_DIR/.server.pid"
      return 0
    fi
    echo -n "."
    sleep 1
  done
  
  echo -e "\n${RED}Failed to start server${NC}"
  cat "$CLIENT_DIR/server.log"
  return 1
}

# Function to stop the server
stop_server() {
  if [ -f "$CLIENT_DIR/.server.pid" ]; then
    SERVER_PID=$(cat "$CLIENT_DIR/.server.pid")
    echo -e "${YELLOW}Stopping MCP server (PID: $SERVER_PID)...${NC}"
    kill $SERVER_PID 2>/dev/null || true
    rm "$CLIENT_DIR/.server.pid"
    echo -e "${GREEN}Server stopped${NC}"
  fi
}

# Function to run unit tests
run_unit_tests() {
  echo -e "${YELLOW}Running unit tests...${NC}"
  cd "$CLIENT_DIR"
  python3 -m pytest tests/ -xvs
}

# Function to run integration tests
run_integration_tests() {
  echo -e "${YELLOW}Running integration tests...${NC}"
  cd "$CLIENT_DIR"
  python3 -m pytest tests/integration/ -xvs
}

# Clean up on exit
trap stop_server EXIT

# Make sure we're in the right directory
cd "$CLIENT_DIR"

# Ensure test directories exist
mkdir -p tests/integration
mkdir -p tests/utils

# Start the server
start_server

# Run the tests
run_unit_tests
run_integration_tests

echo -e "${GREEN}All tests completed successfully!${NC}" 