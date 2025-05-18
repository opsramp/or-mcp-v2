#!/bin/bash
# Comprehensive MCP server test script for the HPE OpsRamp MCP server
# Official repository: https://github.com/opsramp/or-mcp-v2
#
# This script tests the MCP server by:
# 1. Starting the server if not already running
# 2. Testing server health endpoints
# 3. Running Python browser-like client tests
# 4. Testing integration with the client

set -e

# Get the absolute path of the repository root
REPO_ROOT=$(cd "$(dirname "$0")" && pwd)
CLIENT_DIR="$REPO_ROOT/client/python"
SERVER_CMD="$REPO_ROOT/cmd/server/main.go"
SERVER_PORT=8080
SERVER_URL="http://localhost:$SERVER_PORT"

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
    echo -e "${YELLOW}Checking if it's our server...${NC}"
    
    # Try to check if it's our server by checking the health endpoint
    if curl -s "$SERVER_URL/health" > /dev/null; then
      echo -e "${GREEN}MCP server is already running and healthy.${NC}"
      return 0
    else
      echo -e "${RED}Port $SERVER_PORT is in use but doesn't seem to be our MCP server.${NC}"
      echo -e "${RED}Please free up the port and try again.${NC}"
      return 1
    fi
  fi
  
  # Create output directory if it doesn't exist
  mkdir -p "$REPO_ROOT/output/logs"
  
  # Start server in background
  echo -e "${YELLOW}Starting server in debug mode...${NC}"
  cd "$REPO_ROOT"
  DEBUG=true go run "$SERVER_CMD" > "$REPO_ROOT/output/logs/server.log" 2>&1 &
  SERVER_PID=$!
  
  # Wait for server to start
  echo -n "Waiting for server to start"
  for i in {1..10}; do
    if curl -s "$SERVER_URL/health" > /dev/null; then
      echo -e "\n${GREEN}Server started successfully (PID: $SERVER_PID)${NC}"
      echo $SERVER_PID > "$REPO_ROOT/.server.pid"
      return 0
    fi
    echo -n "."
    sleep 1
  done
  
  echo -e "\n${RED}Failed to start server${NC}"
  cat "$REPO_ROOT/output/logs/server.log"
  return 1
}

# Function to stop the server
stop_server() {
  if [ -f "$REPO_ROOT/.server.pid" ]; then
    SERVER_PID=$(cat "$REPO_ROOT/.server.pid")
    echo -e "${YELLOW}Stopping MCP server (PID: $SERVER_PID)...${NC}"
    kill $SERVER_PID 2>/dev/null || true
    rm "$REPO_ROOT/.server.pid"
    echo -e "${GREEN}Server stopped${NC}"
  fi
}

# Function to check server health
check_server_health() {
  echo -e "${YELLOW}Checking server health...${NC}"
  
  # Try to get the health endpoint
  HEALTH_DATA=$(curl -s "$SERVER_URL/health")
  if [ $? -ne 0 ]; then
    echo -e "${RED}Failed to connect to server health endpoint${NC}"
    return 1
  fi
  
  # Parse the health data
  echo "$HEALTH_DATA" | python -m json.tool
  
  # Check if the server reports it's healthy
  if echo "$HEALTH_DATA" | grep -q '"status":"ok"'; then
    echo -e "${GREEN}Server reports healthy status${NC}"
    
    # Extract uptime and tools
    UPTIME=$(echo "$HEALTH_DATA" | grep -o '"uptime":"[^"]*"' | cut -d'"' -f4)
    TOOLS=$(echo "$HEALTH_DATA" | grep -o '"tools":\[[^]]*\]' | cut -d':' -f2)
    
    echo -e "${GREEN}Server uptime: $UPTIME${NC}"
    echo -e "${GREEN}Available tools: $TOOLS${NC}"
    
    return 0
  else
    echo -e "${RED}Server does not report healthy status${NC}"
    return 1
  fi
}

# Function to run the Python client test
run_python_client_test() {
  echo -e "${YELLOW}Running Python client test...${NC}"
  
  # Check if the Python client directory exists
  if [ ! -d "$CLIENT_DIR" ]; then
    echo -e "${RED}Python client directory not found: $CLIENT_DIR${NC}"
    return 1
  fi
  
  # Create a virtual environment if not exists
  if [ ! -d "$CLIENT_DIR/.venv" ]; then
    echo -e "${YELLOW}Creating Python virtual environment...${NC}"
    cd "$CLIENT_DIR"
    python -m venv .venv
    source .venv/bin/activate
    pip install -r requirements.txt
  else
    cd "$CLIENT_DIR"
    source .venv/bin/activate
  fi
  
  # Run the browser-like example
  echo -e "${YELLOW}Running browser-like client example...${NC}"
  cd "$CLIENT_DIR"
  python examples/browser_like_example.py --debug
  RESULT=$?
  
  deactivate
  
  if [ $RESULT -eq 0 ]; then
    echo -e "${GREEN}Browser-like client example completed successfully${NC}"
    return 0
  else
    echo -e "${RED}Browser-like client example failed${NC}"
    return 1
  fi
}

# Function to run tests
run_tests() {
  echo -e "${YELLOW}Starting comprehensive MCP server tests...${NC}"
  
  # Start the server if needed
  start_server
  
  # Check server health
  check_server_health
  
  # Run Python client test
  run_python_client_test
  
  echo -e "${GREEN}All tests completed successfully!${NC}"
}

# Function to show usage
show_usage() {
  echo "Usage: $0 [COMMAND]"
  echo ""
  echo "Commands:"
  echo "  start      Start the MCP server"
  echo "  stop       Stop the MCP server"
  echo "  health     Check server health"
  echo "  test       Run client tests"
  echo "  all        Start server and run all tests (default)"
  echo ""
}

# Clean up on exit if started by this script
trap 'if [ -f "$REPO_ROOT/.server.pid" ]; then stop_server; fi' EXIT

# Parse command line arguments
COMMAND=${1:-all}

case $COMMAND in
  start)
    start_server
    echo -e "${GREEN}Server started. Use '$0 stop' to stop it.${NC}"
    trap - EXIT  # Disable auto-stop on exit
    ;;
  stop)
    stop_server
    ;;
  health)
    check_server_health
    ;;
  test)
    run_python_client_test
    ;;
  all)
    run_tests
    ;;
  *)
    show_usage
    exit 1
    ;;
esac 