#!/bin/bash

# Run the unit tests for the integrations tool
echo "Running unit tests for the integrations tool..."
go test -v ./tests/integrations_test.go

# Optionally run the client test (requires manual intervention)
# echo "Running client test (requires manual intervention)..."
# go test -v ./tests/client_test.go -run TestMCPClient
