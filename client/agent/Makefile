.PHONY: all setup test clean run-example help test-basic test-complex test-comprehensive test-all analyze-results run-interactive test-interactive test-single

# Variables
PYTHON := python3
PIP := $(PYTHON) -m pip
PYTEST := $(PYTHON) -m pytest
TEST_DIR := tests
OUTPUT_DIR := output

# Colors for output
GREEN := \033[0;32m
YELLOW := \033[1;33m
RED := \033[0;31m
NC := \033[0m # No Color

all: setup

setup:
	@echo "$(GREEN)Setting up OpsRamp AI Agent...$(NC)"
	$(PIP) install -e .
	$(PIP) install -e ".[all]"
	@echo "$(GREEN)Installation complete.$(NC)"

# Basic unit tests (if any exist)
test:
	@echo "$(YELLOW)Running basic unit tests...$(NC)"
	@if [ -f "$(TEST_DIR)/test_*.py" ]; then \
		$(PYTEST) -xvs $(TEST_DIR)/; \
	else \
		echo "$(YELLOW)No unit tests found. Use comprehensive integration tests instead.$(NC)"; \
	fi

# Quick test with 3 prompts
test-basic:
	@echo "$(GREEN)Running basic integration tests (3 prompts)...$(NC)"
	@mkdir -p $(OUTPUT_DIR)
	cd $(TEST_DIR) && $(PYTHON) enhanced_real_mcp_integration_test.py --max-tests 3

# Medium test with 10 prompts  
test-medium:
	@echo "$(GREEN)Running medium integration tests (10 prompts)...$(NC)"
	@mkdir -p $(OUTPUT_DIR)
	cd $(TEST_DIR) && $(PYTHON) enhanced_real_mcp_integration_test.py --max-tests 10

# Ultra-complex scenarios
test-complex:
	@echo "$(GREEN)Running ultra-complex integration tests...$(NC)"
	@mkdir -p $(OUTPUT_DIR)
	cd $(TEST_DIR) && $(PYTHON) enhanced_real_mcp_integration_test.py \
		--prompts-file test_data/ultra_complex_integration_prompts.txt \
		--max-tests 5

# Full comprehensive test suite
test-comprehensive:
	@echo "$(GREEN)Running FULL comprehensive integration test suite...$(NC)"
	@mkdir -p $(OUTPUT_DIR)
	cd $(TEST_DIR) && $(PYTHON) enhanced_real_mcp_integration_test.py \
		--prompts-file test_data/comprehensive_integration_prompts.txt

# Run all test suites
test-all: test-basic test-complex
	@echo "$(GREEN)All test suites completed!$(NC)"
	@$(MAKE) analyze-results

# Analyze test results  
analyze-results:
	@echo "$(YELLOW)Test Results Summary:$(NC)"
	@echo "================================"
	@if [ -d "$(OUTPUT_DIR)" ]; then \
		echo "📁 Output files:"; \
		ls -la $(OUTPUT_DIR)/ | grep -E '\.(log|json|jsonl)$$' || echo "No output files found"; \
		echo ""; \
		echo "📊 Latest test analysis:"; \
		latest_analysis=$$(ls -t $(OUTPUT_DIR)/comprehensive_test_analysis_*.json 2>/dev/null | head -1); \
		if [ -n "$$latest_analysis" ]; then \
			echo "$$latest_analysis"; \
			$(PYTHON) -c "import json; data=json.load(open('$$latest_analysis')); summary=data['test_summary']; print(f\"✅ Tests: {summary['successful_tests']}/{summary['total_tests']} ({summary['success_rate']:.1f}%)\"); print(f\"⏱️  Duration: {summary['total_duration']:.1f}s\"); print(f\"🔧 Tool Calls: {summary['total_tool_calls']}\")"; \
		else \
			echo "No analysis files found"; \
		fi \
	else \
		echo "No output directory found - run tests first"; \
	fi

# Custom test with specific prompts file
test-custom:
	@echo "$(GREEN)Running custom integration tests...$(NC)"
	@echo "Usage: make test-custom PROMPTS_FILE=path/to/prompts.txt MAX_TESTS=5"
	@if [ -z "$(PROMPTS_FILE)" ]; then \
		echo "$(RED)Error: PROMPTS_FILE is required$(NC)"; \
		echo "Example: make test-custom PROMPTS_FILE=test_data/ultra_complex_integration_prompts.txt MAX_TESTS=3"; \
		exit 1; \
	fi
	@mkdir -p $(OUTPUT_DIR)
	cd $(TEST_DIR) && $(PYTHON) enhanced_real_mcp_integration_test.py \
		--prompts-file $(PROMPTS_FILE) \
		$(if $(MAX_TESTS),--max-tests $(MAX_TESTS))

clean:
	@echo "$(YELLOW)Cleaning up...$(NC)"
	rm -rf build/
	rm -rf dist/
	rm -rf *.egg-info/
	rm -rf src/*.egg-info/
	find . -name __pycache__ -exec rm -rf {} +
	find . -name "*.pyc" -delete
	@echo "$(GREEN)Cleanup complete.$(NC)"

clean-output:
	@echo "$(YELLOW)Cleaning test output...$(NC)"
	rm -rf $(OUTPUT_DIR)/*
	@echo "$(GREEN)Output cleanup complete.$(NC)"

run-example:
	@echo "$(GREEN)Running chat example...$(NC)"
	$(PYTHON) examples/chat_client.py

run-batch:
	@echo "$(GREEN)Running batch processor example...$(NC)"
	$(PYTHON) examples/batch_process.py --input-file examples/sample_prompts.txt --output-file output.txt

run-prompt:
	@echo "$(GREEN)Running single prompt example...$(NC)"
	$(PYTHON) examples/chat_client.py --prompt "List all integrations"

# Development and debugging
dev-test:
	@echo "$(GREEN)Running single test for development...$(NC)"
	@mkdir -p $(OUTPUT_DIR)
	cd $(TEST_DIR) && $(PYTHON) enhanced_real_mcp_integration_test.py --max-tests 1

# Check MCP server connection
check-server:
	@echo "$(YELLOW)Checking MCP server connection...$(NC)"
	@curl -s http://localhost:8080/health > /dev/null && \
		echo "$(GREEN)✅ MCP server is running$(NC)" || \
		echo "$(RED)❌ MCP server is not accessible$(NC)"

# Enhanced Interactive Testing - Uses our proven comprehensive testing engine
run-interactive:
	@echo "$(GREEN)Running enhanced interactive testing mode...$(NC)"
	@echo "$(YELLOW)This uses our proven comprehensive testing engine for reliable tool usage$(NC)"
	@mkdir -p $(OUTPUT_DIR)
	cd $(TEST_DIR) && $(PYTHON) enhanced_real_mcp_integration_test.py \
		--prompts-file test_data/interactive_test_scenarios.txt \
		--max-tests 5

# Interactive testing with predefined scenarios
test-interactive:
	@echo "$(GREEN)Running interactive test scenarios...$(NC)"
	@mkdir -p $(OUTPUT_DIR)
	@echo "Testing common integration queries..."
	cd $(TEST_DIR) && $(PYTHON) enhanced_real_mcp_integration_test.py \
		--prompts-file test_data/comprehensive_integration_prompts.txt \
		--max-tests 5

# Single question test for immediate verification
test-single:
	@echo "$(GREEN)Running single integration question test...$(NC)"
	@mkdir -p $(OUTPUT_DIR)
	@if [ -z "$(QUESTION)" ]; then \
		echo "$(YELLOW)Usage: make test-single QUESTION='your question here'$(NC)"; \
		echo "$(YELLOW)Example: make test-single QUESTION='what are the emails of users who installed integrations?'$(NC)"; \
		exit 1; \
	fi
	@echo "$(QUESTION)" > $(OUTPUT_DIR)/single_question.txt
	cd $(TEST_DIR) && $(PYTHON) enhanced_real_mcp_integration_test.py \
		--prompts-file ../$(OUTPUT_DIR)/single_question.txt \
		--max-tests 1

help:
	@echo "$(GREEN)OpsRamp AI Agent - Enhanced Makefile$(NC)"
	@echo ""
	@echo "$(YELLOW)Setup Commands:$(NC)"
	@echo "  setup           - Install the agent and dependencies"
	@echo "  clean           - Clean up build artifacts"
	@echo "  clean-output    - Clean test output files"
	@echo ""
	@echo "$(YELLOW)Testing Commands:$(NC)"
	@echo "  test-basic      - Quick test with 3 prompts"
	@echo "  test-medium     - Medium test with 10 prompts"
	@echo "  test-complex    - Ultra-complex scenarios (5 prompts)"
	@echo "  test-comprehensive - Full test suite (90 prompts)"
	@echo "  test-all        - Run basic + complex tests"
	@echo "  test-custom     - Custom test (requires PROMPTS_FILE=...)"
	@echo "  dev-test        - Single test for development"
	@echo ""
	@echo "$(YELLOW)Analysis Commands:$(NC)"
	@echo "  analyze-results - Show latest test results summary"
	@echo "  check-server    - Check if MCP server is running"
	@echo ""
	@echo "$(YELLOW)Interactive Commands:$(NC)"
	@echo "  run-interactive - Run enhanced interactive testing (RECOMMENDED)"
	@echo "  run-example     - Run original chat example"
	@echo "  run-prompt      - Run with single prompt"
	@echo "  run-batch       - Run batch processor example"
	@echo ""
	@echo "$(YELLOW)Interactive Testing Commands:$(NC)"
	@echo "  test-interactive  - Interactive test with predefined scenarios"
	@echo "  test-single       - Single question interactive test"
	@echo ""
	@echo "$(YELLOW)Examples:$(NC)"
	@echo "  make run-interactive                               # Enhanced interactive mode"
	@echo "  make test-interactive                              # Interactive testing"
	@echo "  make test-basic                                    # Quick 3-prompt test"
	@echo "  make test-custom PROMPTS_FILE=my_prompts.txt       # Custom prompts"
	@echo "  make test-custom PROMPTS_FILE=test_data/ultra_complex_integration_prompts.txt MAX_TESTS=2"
