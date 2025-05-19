#!/bin/bash

# Move to project root (directory of this script's parent)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_ROOT" || exit 1

# Check for config.yaml in project root
if [ ! -f config.yaml ]; then
  echo "ERROR: config.yaml not found in project root ($PROJECT_ROOT). Exiting."
  exit 1
fi

# Function to check if server is ready
check_server_ready() {
  response=$(curl -s http://localhost:8080/health)
  if [ $? -eq 0 ] && [ -n "$response" ]; then
    # Check if response contains "status": "ok"
    if echo "$response" | grep -q '"status":"ok"'; then
      return 0
    fi
  fi
  return 1
}

# Kill any existing server process
echo "Checking for existing server process..."
pkill -f "or-mcp-server" || true

# Start the server in the background
echo "Starting server..."
make run &
SERVER_PID=$!

# Wait for server to be ready (up to 30 seconds)
echo "Waiting for server to be ready..."
for i in {1..30}; do
  if check_server_ready; then
    echo "Server is ready!"
    break
  fi
  if [ $i -eq 30 ]; then
    echo "ERROR: Server failed to start within 30 seconds"
    kill -15 $SERVER_PID 2>/dev/null
    exit 1
  fi
  echo -n "."
  sleep 1
done
echo

# Run the real API integration tests
echo "Running real API integration tests..."
go test -v ./tests/integrations_real_api_test.go

# Store the test result
TEST_RESULT=$?

# Stop the server
echo "Stopping server..."
kill -15 $SERVER_PID 2>/dev/null

# Exit with the test result
if [ $TEST_RESULT -eq 0 ]; then
    echo "✅ Integration tests passed!"
    exit 0
else
    echo "❌ Integration tests failed!"
    exit 1
fi 