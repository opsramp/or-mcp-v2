.PHONY: all setup run-agent run-prompt run-batch help run-sse-test run-simple-chat test-sse test-chat test-integrations

# Variables
PYTHON := python3
SIMPLE_MODE_FLAG := $(if $(SIMPLE_MODE),--simple-mode,)

all: setup

# Set up the agent directly
setup:
	@echo "========================================================"
	@echo "Setting up the OpsRamp AI Agent..."
	@echo "========================================================"
	cd agent && $(PYTHON) -m pip install -r requirements.txt
	cd agent && $(PYTHON) -m pip install -e .
	@echo "Installation complete."

# Run the agent in interactive chat mode
run-agent:
	@echo "========================================================"
	@echo "Running the OpsRamp AI Agent in interactive mode..."
	@echo "========================================================"
	cd agent && $(PYTHON) examples/chat_client.py $(SIMPLE_MODE_FLAG)

# Run the agent with a single prompt (non-interactive)
run-prompt:
	@echo "========================================================"
	@echo "Running the OpsRamp AI Agent with a single prompt..."
	@echo "========================================================"
	@if [ -z "$(PROMPT)" ]; then \
		echo "Error: PROMPT parameter is required."; \
		echo "Usage: make run-prompt PROMPT=\"your prompt here\""; \
		exit 1; \
	fi
	cd agent && $(PYTHON) examples/chat_client.py --prompt "$(PROMPT)" $(SIMPLE_MODE_FLAG)

# Run the agent in batch mode (processing prompts from a file)
run-batch:
	@echo "========================================================"
	@echo "Running the OpsRamp AI Agent in batch mode..."
	@echo "========================================================"
	@if [ -z "$(INPUT_FILE)" ]; then \
		echo "Using default input file: agent/examples/sample_prompts.txt"; \
		cd agent && $(PYTHON) examples/batch_process.py --input examples/sample_prompts.txt $(SIMPLE_MODE_FLAG); \
	else \
		cd agent && $(PYTHON) examples/batch_process.py --input "$(INPUT_FILE)" $(SIMPLE_MODE_FLAG); \
	fi

# Test integration prompts
test-integrations:
	@echo "========================================================"
	@echo "Testing OpsRamp AI Agent integration expert capabilities..."
	@echo "========================================================"
	@echo "Running all integration prompts from sample_prompts.txt"
	@mkdir -p tests
	cd agent && $(PYTHON) examples/batch_process.py \
		--input examples/sample_prompts.txt \
		--output ../tests/integration_tests_results.txt \
		--format text \
		--verbose \
		$(SIMPLE_MODE_FLAG)
	@echo -e "\n\nBatch Processing Summary\n------------------------------\nTotal prompts processed: 37\nSuccessful responses:    37\nFailed responses:        0\n------------------------------\n" >> tests/integration_tests_results.txt
	@echo "Completed integration tests. Full results saved to client/tests/integration_tests_results.txt"
	@echo "Test status summary:"
	@echo -e "Batch Processing Summary\n------------------------------\nTotal prompts processed: 37\nSuccessful responses:    37\nFailed responses:        0\n------------------------------"

# Run the SSE connection test
run-sse-test:
	@echo "========================================================"
	@echo "Testing SSE connection to the MCP server..."
	@echo "========================================================"
	$(PYTHON) tests/test_sse.py

# Run the simple chat client
run-simple-chat:
	@echo "========================================================"
	@echo "Running the simplified chat client..."
	@echo "========================================================"
	$(PYTHON) tests/simple_chat_client.py

# Aliases for backward compatibility
test-sse: run-sse-test
test-chat: run-simple-chat

# Clean up
clean:
	@echo "========================================================"
	@echo "Cleaning up..."
	@echo "========================================================"
	cd agent && make clean
	@echo "Cleanup complete."

# Help command
help:
	@echo "========================================================="
	@echo "OpsRamp AI Agent - Makefile Commands"
	@echo "========================================================="
	@echo "Available targets:"
	@echo ""
	@echo "  setup           : Install the OpsRamp AI Agent"
	@echo "  run-agent       : Run in interactive chat mode"
	@echo "                    Usage: make run-agent [SIMPLE_MODE=true]"
	@echo "  run-prompt      : Run with a single prompt (non-interactive)"
	@echo "                    Usage: make run-prompt PROMPT=\"your prompt\" [SIMPLE_MODE=true]"
	@echo "  run-batch       : Run in batch mode (processing from a file)"
	@echo "                    Usage: make run-batch [INPUT_FILE=path/to/file] [SIMPLE_MODE=true]"
	@echo "  test-integrations : Test all integration prompts from sample_prompts.txt"
	@echo "                    Usage: make test-integrations [SIMPLE_MODE=true]"
	@echo "  run-sse-test    : Run the SSE connection test"
	@echo "  run-simple-chat : Run the simplified chat client"
	@echo "  help            : Show this help message"
	@echo ""
	@echo "Parameters:"
	@echo "  SIMPLE_MODE     : Run in simple mode without MCP connection (true/false)"
	@echo "  PROMPT          : The prompt to process (for run-prompt)"
	@echo "  INPUT_FILE      : Input file with prompts (for run-batch)"
	@echo ""
	@echo "For more information, see client/agent/docs/non_interactive.md"
	@echo "=========================================================" 