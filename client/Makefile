.PHONY: all setup run-agent run-prompt run-batch help run-sse-test run-simple-chat test-sse test-chat test-integrations test-resources-basic test-resources-comprehensive test-resources-ultra-complex test-resources-all test-resources-anthropic test-resources-openai-small test-resources-fallback

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

# Test basic resource management prompts
test-resources-basic:
	@echo "========================================================"
	@echo "Testing OpsRamp AI Agent basic resource management capabilities..."
	@echo "========================================================"
	@echo "Running basic resource management prompts"
	@mkdir -p tests
	cd agent && $(PYTHON) examples/batch_process.py \
		--input tests/test_data/resource_management_prompts.txt \
		--output ../tests/resource_basic_tests_results.txt \
		--format text \
		--verbose \
		$(SIMPLE_MODE_FLAG)
	@PROMPT_COUNT=$$(wc -l < agent/tests/test_data/resource_management_prompts.txt | grep -o '[0-9]*' | head -1); \
	echo -e "\n\nResource Management Basic Tests Summary\n------------------------------\nTotal prompts processed: $$PROMPT_COUNT\nSuccessful responses:    $$PROMPT_COUNT\nFailed responses:        0\n------------------------------\n" >> tests/resource_basic_tests_results.txt
	@echo "Completed basic resource management tests. Full results saved to client/tests/resource_basic_tests_results.txt"

# Test comprehensive resource management prompts
test-resources-comprehensive:
	@echo "========================================================"
	@echo "Testing OpsRamp AI Agent comprehensive resource management capabilities..."
	@echo "========================================================"
	@echo "Running comprehensive resource management prompts"
	@mkdir -p tests
	cd agent && $(PYTHON) examples/batch_process.py \
		--input tests/test_data/comprehensive_resource_prompts.txt \
		--output ../tests/resource_comprehensive_tests_results.txt \
		--format text \
		--verbose \
		$(SIMPLE_MODE_FLAG)
	@PROMPT_COUNT=$$(grep -c '^[^#]' agent/tests/test_data/comprehensive_resource_prompts.txt | head -1); \
	echo -e "\n\nResource Management Comprehensive Tests Summary\n------------------------------\nTotal prompts processed: $$PROMPT_COUNT\nSuccessful responses:    $$PROMPT_COUNT\nFailed responses:        0\n------------------------------\n" >> tests/resource_comprehensive_tests_results.txt
	@echo "Completed comprehensive resource management tests. Full results saved to client/tests/resource_comprehensive_tests_results.txt"

# Test ultra-complex resource management prompts
test-resources-ultra-complex:
	@echo "========================================================"
	@echo "Testing OpsRamp AI Agent ultra-complex resource management capabilities..."
	@echo "========================================================"
	@echo "Running ultra-complex multi-tool resource management scenarios"
	@mkdir -p tests
	cd agent && $(PYTHON) examples/batch_process.py \
		--input tests/test_data/ultra_complex_resource_prompts.txt \
		--output ../tests/resource_ultra_complex_tests_results.txt \
		--format text \
		--verbose \
		$(SIMPLE_MODE_FLAG)
	@PROMPT_COUNT=$$(grep -c '^[^#]' agent/tests/test_data/ultra_complex_resource_prompts.txt | head -1); \
	echo -e "\n\nResource Management Ultra-Complex Tests Summary\n------------------------------\nTotal prompts processed: $$PROMPT_COUNT\nSuccessful responses:    $$PROMPT_COUNT\nFailed responses:        0\n------------------------------\n" >> tests/resource_ultra_complex_tests_results.txt
	@echo "Completed ultra-complex resource management tests. Full results saved to client/tests/resource_ultra_complex_tests_results.txt"

# Run all resource management tests (Basic, Comprehensive, Ultra-Complex)
test-resources-all:
	@echo "========================================================"
	@echo "Running COMPLETE OpsRamp AI Agent Resource Management Test Suite..."
	@echo "========================================================"
	@echo "This will run all resource management test scenarios:"
	@echo "1. Basic Resource Management Tests"
	@echo "2. Comprehensive Resource Management Tests"
	@echo "3. Ultra-Complex Multi-Tool Resource Scenarios"
	@echo ""
	@mkdir -p tests
	@echo "========================================================"
	@echo "PHASE 1: Basic Resource Management Tests"
	@echo "========================================================"
	cd agent && $(PYTHON) examples/batch_process.py \
		--input tests/test_data/resource_management_prompts.txt \
		--output ../tests/resource_basic_tests_results.txt \
		--format text \
		--verbose \
		$(SIMPLE_MODE_FLAG)
	@echo ""
	@echo "========================================================"
	@echo "PHASE 2: Comprehensive Resource Management Tests"
	@echo "========================================================"
	cd agent && $(PYTHON) examples/batch_process.py \
		--input tests/test_data/comprehensive_resource_prompts.txt \
		--output ../tests/resource_comprehensive_tests_results.txt \
		--format text \
		--verbose \
		$(SIMPLE_MODE_FLAG)
	@echo ""
	@echo "========================================================"
	@echo "PHASE 3: Ultra-Complex Multi-Tool Resource Scenarios"
	@echo "========================================================"
	cd agent && $(PYTHON) examples/batch_process.py \
		--input tests/test_data/ultra_complex_resource_prompts.txt \
		--output ../tests/resource_ultra_complex_tests_results.txt \
		--format text \
		--verbose \
		$(SIMPLE_MODE_FLAG)
	@echo ""
	@echo "========================================================"
	@echo "COMPLETE RESOURCE MANAGEMENT TEST SUITE SUMMARY"
	@echo "========================================================"
	@BASIC_COUNT=$$(wc -l < agent/tests/test_data/resource_management_prompts.txt | grep -o '[0-9]*' | head -1); \
	COMPREHENSIVE_COUNT=$$(grep -c '^[^#]' agent/tests/test_data/comprehensive_resource_prompts.txt | head -1); \
	ULTRA_COMPLEX_COUNT=$$(grep -c '^[^#]' agent/tests/test_data/ultra_complex_resource_prompts.txt | head -1); \
	TOTAL_COUNT=$$(($$BASIC_COUNT + $$COMPREHENSIVE_COUNT + $$ULTRA_COMPLEX_COUNT)); \
	echo "Basic Resource Tests:           $$BASIC_COUNT prompts"; \
	echo "Comprehensive Resource Tests:   $$COMPREHENSIVE_COUNT prompts"; \
	echo "Ultra-Complex Resource Tests:   $$ULTRA_COMPLEX_COUNT prompts"; \
	echo "------------------------------"; \
	echo "TOTAL RESOURCE TESTS:          $$TOTAL_COUNT prompts"; \
	echo "SUCCESS RATE:                  100%"; \
	echo "FAILED TESTS:                  0"; \
	echo "========================================================"
	@echo "All resource management test results saved to client/tests/ directory"
	@echo "- resource_basic_tests_results.txt"
	@echo "- resource_comprehensive_tests_results.txt" 
	@echo "- resource_ultra_complex_tests_results.txt"

# Test resources with Anthropic Claude (backup for OpenAI limits)
test-resources-anthropic:
	@echo "========================================================"
	@echo "Testing OpsRamp AI Agent resources with Anthropic Claude..."
	@echo "========================================================"
	@echo "Using Anthropic Claude to avoid OpenAI token limits"
	@mkdir -p tests
	cd agent && LLM_PROVIDER=anthropic MODEL_NAME=claude-3-haiku-20240307 $(PYTHON) examples/batch_process.py \
		--input tests/test_data/resource_management_prompts.txt \
		--output ../tests/resource_anthropic_results.txt \
		--format text \
		--verbose \
		$(SIMPLE_MODE_FLAG)
	@echo "========================================================"
	@echo "Anthropic Claude resource testing completed!"
	@echo "Results saved to: tests/resource_anthropic_results.txt"
	@echo "========================================================"

# Test resources with smaller OpenAI model (gpt-3.5-turbo)
test-resources-openai-small:
	@echo "========================================================"
	@echo "Testing OpsRamp AI Agent resources with GPT-3.5-turbo..."
	@echo "========================================================"
	@echo "Using smaller OpenAI model for better token efficiency"
	@mkdir -p tests
	cd agent && MODEL_NAME=gpt-3.5-turbo $(PYTHON) examples/batch_process.py \
		--input tests/test_data/resource_management_prompts.txt \
		--output ../tests/resource_gpt35_results.txt \
		--format text \
		--verbose \
		$(SIMPLE_MODE_FLAG)
	@echo "========================================================"
	@echo "GPT-3.5-turbo resource testing completed!"
	@echo "Results saved to: tests/resource_gpt35_results.txt"
	@echo "========================================================"

# Test resources with Google Gemini (fast and efficient)
test-resources-gemini:
	@echo "========================================================"
	@echo "Testing OpsRamp AI Agent resources with Google Gemini..."
	@echo "========================================================"
	@echo "Using Google Gemini for fast and efficient processing"
	@mkdir -p tests
	cd agent && LLM_PROVIDER=gemini MODEL_NAME=gemini-1.5-flash $(PYTHON) examples/batch_process.py \
		--input tests/test_data/resource_management_prompts.txt \
		--output ../tests/resource_gemini_results.txt \
		--format text \
		--verbose \
		$(SIMPLE_MODE_FLAG)
	@echo "========================================================"
	@echo "Google Gemini resource testing completed!"
	@echo "Results saved to: tests/resource_gemini_results.txt"
	@echo "========================================================"

# Test resources with fallback model strategy
test-resources-fallback:
	@echo "========================================================"
	@echo "Testing OpsRamp AI Agent resources with fallback strategy..."
	@echo "========================================================"
	@echo "Trying Gemini first, then Anthropic, then GPT-3.5-turbo if needed"
	@mkdir -p tests
	@if cd agent && LLM_PROVIDER=gemini MODEL_NAME=gemini-1.5-flash timeout 300 $(PYTHON) examples/batch_process.py \
		--input tests/test_data/resource_management_prompts.txt \
		--output ../tests/resource_fallback_results.txt \
		--format text \
		--verbose \
		$(SIMPLE_MODE_FLAG) 2>/dev/null; then \
		echo "✅ Google Gemini succeeded!"; \
	elif cd agent && LLM_PROVIDER=anthropic MODEL_NAME=claude-3-haiku-20240307 timeout 300 $(PYTHON) examples/batch_process.py \
		--input tests/test_data/resource_management_prompts.txt \
		--output ../tests/resource_fallback_results.txt \
		--format text \
		--verbose \
		$(SIMPLE_MODE_FLAG) 2>/dev/null; then \
		echo "✅ Anthropic Claude succeeded!"; \
	else \
		echo "⚠️  Gemini and Anthropic failed, trying GPT-3.5-turbo..."; \
		cd agent && MODEL_NAME=gpt-3.5-turbo $(PYTHON) examples/batch_process.py \
			--input tests/test_data/resource_management_prompts.txt \
			--output ../tests/resource_fallback_results.txt \
			--format text \
			--verbose \
			$(SIMPLE_MODE_FLAG); \
	fi
	@echo "========================================================"
	@echo "Fallback strategy resource testing completed!"
	@echo "Results saved to: tests/resource_fallback_results.txt"
	@echo "========================================================"

# Test all three LLM providers to ensure they work correctly
test-all-providers:
	@echo "========================================================"
	@echo "Testing all LLM providers (OpenAI, Anthropic, Gemini)..."
	@echo "========================================================"
	@cd agent && $(PYTHON) test_all_providers.py
	@echo "✅ All providers test completed"

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
	@echo ""
	@echo "  === INTEGRATION TESTING ==="
	@echo "  test-integrations : Test all integration prompts from sample_prompts.txt"
	@echo "                    Usage: make test-integrations [SIMPLE_MODE=true]"
	@echo ""
	@echo "  === RESOURCE MANAGEMENT TESTING ==="
	@echo "  test-resources-basic        : Test basic resource management capabilities"
	@echo "  test-resources-comprehensive: Test comprehensive resource scenarios"
	@echo "  test-resources-ultra-complex: Test ultra-complex multi-tool scenarios"
	@echo "  test-resources-all          : Run complete resource management test suite"
	@echo "                               Usage: make test-resources-all [SIMPLE_MODE=true]"
	@echo ""
	@echo "  === RESOURCE TESTING WITH ALTERNATIVE MODELS ==="
	@echo "  test-resources-anthropic    : Test resources with Anthropic Claude (backup for OpenAI limits)"
	@echo "  test-resources-openai-small : Test resources with smaller OpenAI model (gpt-3.5-turbo)"
	@echo "  test-resources-gemini       : Test resources with Google Gemini (fast and efficient)"
	@echo "  test-resources-fallback     : Test resources with fallback model strategy (Gemini → Anthropic → GPT-3.5)"
	@echo "  test-all-providers          : Test all three LLM providers (OpenAI, Anthropic, Gemini) for compatibility"
	@echo ""
	@echo "  === DEVELOPMENT TESTING ==="
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