#!/bin/bash

# This script runs all tests against the real OpsRamp API endpoint.
# Mock tests are no longer used.

# Run the unit tests for the integrations tool
echo "Running unit tests for the integrations tool..."
go test -v ./tests/integrations_test.go

# Run the real API tests for the integrations tool
echo "Running real API tests for the integrations tool..."
go test -v ./tests/integrations_real_api_test.go

# Optionally run the client test (requires manual intervention)
# echo "Running client test (requires manual intervention)..."
# go test -v ./tests/client_test.go -run TestMCPClient
