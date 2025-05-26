.PHONY: all setup test clean run-example help test-basic test-complex test-comprehensive test-all analyze-results run-interactive test-interactive test-single chat-interactive

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
	cd $(TEST_DIR) && OPENAI_MODEL=gpt-4-turbo $(PYTHON) shared/engines/enhanced_real_mcp_integration_test.py \
		--prompts-file test_data/basic_integration_prompts.txt \
		--max-tests 3

# Medium test with 10 prompts  
test-medium:
	@echo "$(GREEN)Running medium integration tests (10 prompts)...$(NC)"
	@mkdir -p $(OUTPUT_DIR)
	cd $(TEST_DIR) && $(PYTHON) shared/engines/enhanced_real_mcp_integration_test.py --max-tests 10

# Ultra-complex scenarios
test-complex:
	@echo "$(GREEN)Running ultra-complex integration tests...$(NC)"
	@mkdir -p $(OUTPUT_DIR)
	cd $(TEST_DIR) && $(PYTHON) shared/engines/enhanced_real_mcp_integration_test.py \
		--prompts-file test_data/ultra_complex_integration_prompts.txt \
		--max-tests 5

# Full comprehensive test suite
test-comprehensive:
	@echo "$(GREEN)Running FULL comprehensive integration test suite...$(NC)"
	@mkdir -p $(OUTPUT_DIR)
	cd $(TEST_DIR) && $(PYTHON) shared/engines/enhanced_real_mcp_integration_test.py \
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
	cd $(TEST_DIR) && $(PYTHON) shared/engines/enhanced_real_mcp_integration_test.py \
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
	
# True interactive mode - Chat directly with the agent
chat-interactive: setup
	@echo "$(GREEN)Starting interactive chat with AI agent...$(NC)"
	@echo "$(YELLOW)Type your questions directly and the agent will respond in real-time$(NC)"
	@echo "$(YELLOW)Press Ctrl+C to exit the chat session$(NC)"
	@echo "$(YELLOW)Example questions to try:$(NC)"
	@echo "  - \"List all integrations in our environment\""
	@echo "  - \"Show me all resources with critical status\""
	@echo "  - \"Generate a report of our infrastructure\""
	@echo ""
	$(PYTHON) examples/chat_client.py

# Development and debugging
dev-test:
	@echo "$(GREEN)Running single test for development...$(NC)"
	@mkdir -p $(OUTPUT_DIR)
	cd $(TEST_DIR) && $(PYTHON) shared/engines/enhanced_real_mcp_integration_test.py --max-tests 1

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
	cd $(TEST_DIR) && OPENAI_MODEL=gpt-4-turbo $(PYTHON) shared/engines/enhanced_real_mcp_integration_test.py \
		--prompts-file test_data/interactive_test_scenarios.txt \
		--max-tests 5

# Interactive testing with predefined scenarios
test-interactive:
	@echo "$(GREEN)Running interactive test scenarios...$(NC)"
	@mkdir -p $(OUTPUT_DIR)
	@echo "Testing common integration queries..."
	cd $(TEST_DIR) && $(PYTHON) shared/engines/enhanced_real_mcp_integration_test.py \
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
	cd $(TEST_DIR) && $(PYTHON) shared/engines/enhanced_real_mcp_integration_test.py \
		--prompts-file ../$(OUTPUT_DIR)/single_question.txt \
		--max-tests 1

# Test individual integrations functionality
test-integrations-basic:
	@echo "Testing basic integrations functionality..."
	@mkdir -p tests/$(OUTPUT_DIR)
	@echo "What integrations do we have?" > tests/$(OUTPUT_DIR)/integrations_basic_test.txt
	@cd tests && python3 shared/engines/enhanced_real_mcp_integration_test.py --max-tests 1 --prompts-file $(OUTPUT_DIR)/integrations_basic_test.txt
	@echo "✅ Basic integrations test completed"

# Test individual resources functionality  
test-resources-basic:
	@echo "Testing basic resources functionality..."
	@mkdir -p tests/$(OUTPUT_DIR)
	@echo "What resources do we have?" > tests/$(OUTPUT_DIR)/resources_basic_test.txt
	@cd tests && python3 shared/engines/enhanced_real_mcp_integration_test.py --max-tests 1 --prompts-file $(OUTPUT_DIR)/resources_basic_test.txt
	@echo "✅ Basic resources test completed"

# Test resources with dedicated prompts file
test-resources-comprehensive:
	@echo "$(GREEN)Running comprehensive resources tests...$(NC)"
	@mkdir -p $(OUTPUT_DIR)
	cd $(TEST_DIR) && $(PYTHON) shared/engines/enhanced_real_mcp_integration_test.py \
		--prompts-file test_data/resource_management_prompts.txt \
		--max-tests 10

# Test both integrations and resources together
test-all-tools:
	@echo "$(GREEN)Testing all tools (integrations + resources)...$(NC)"
	@$(MAKE) test-integrations-basic
	@$(MAKE) test-resources-basic
	@echo "$(GREEN)All tools tested successfully!$(NC)"

# Resources-specific test targets
test-resources-all: test-resources-basic test-resources-comprehensive
	@echo "$(GREEN)All resources tests completed!$(NC)"

# ORGANIZED TEST COMMANDS - Integration and Resources Testing

# Integration Testing (Organized)
test-integrations-basic-organized:
	@echo "$(GREEN)Running organized basic integration tests...$(NC)"
	@cd tests && python integration/scripts/run_integration_tests.py --complexity basic

test-integrations-advanced-organized:
	@echo "$(GREEN)Running organized advanced integration tests...$(NC)"
	@cd tests && python integration/scripts/run_integration_tests.py --complexity advanced

test-integrations-all-organized:
	@echo "$(GREEN)Running all organized integration tests...$(NC)"
	@cd tests && python integration/scripts/run_integration_tests.py --complexity all

# Resource Testing (Organized)
test-resources-basic-organized:
	@echo "$(GREEN)Running organized basic resource tests...$(NC)"
	@cd tests && python resources/scripts/run_resource_tests.py --complexity basic

test-resources-comprehensive-organized:
	@echo "$(GREEN)Running organized comprehensive resource tests...$(NC)"
	@cd tests && python resources/scripts/run_resource_tests.py --complexity comprehensive

test-resources-ultra-organized:
	@echo "$(GREEN)Running organized ultra-complex resource tests...$(NC)"
	@cd tests && python resources/scripts/run_resource_tests.py --complexity ultra

test-resources-all-organized:
	@echo "$(GREEN)Running all organized resource tests...$(NC)"
	@cd tests && python resources/scripts/run_resource_tests.py --complexity all

# Multi-Provider Testing (Organized)
test-all-providers-organized:
	@echo "$(GREEN)Running organized multi-provider tests...$(NC)"
	@cd tests && python multi_provider/scripts/test_all_providers.py

test-providers-integration-organized:
	@echo "$(GREEN)Testing integration functionality across all providers...$(NC)"
	@cd tests && python multi_provider/scripts/test_all_providers.py --functionality integration

test-providers-resources-organized:
	@echo "$(GREEN)Testing resource functionality across all providers...$(NC)"
	@cd tests && python multi_provider/scripts/test_all_providers.py --functionality resources

# Master Test Commands (Organized)
test-complete-organized:
	@echo "$(GREEN)Running complete organized test suite...$(NC)"
	@$(MAKE) test-integrations-basic-organized
	@$(MAKE) test-resources-basic-organized
	@$(MAKE) test-all-providers-organized
	@echo "$(GREEN)Complete organized test suite finished!$(NC)"

# Test Data Management (Organized)
cleanup-test-data-organized:
	@echo "$(YELLOW)Cleaning up old test data...$(NC)"
	@cd tests && python shared/utilities/cleanup_test_data.py

cleanup-test-data-dry-organized:
	@echo "$(YELLOW)Showing what test data would be cleaned...$(NC)"
	@cd tests && python shared/utilities/cleanup_test_data.py --dry-run

# Evidence and Reporting (Organized)
generate-test-report-organized:
	@echo "$(GREEN)Generating comprehensive test report...$(NC)"
	@cd tests && python shared/utilities/generate_test_report.py

generate-test-report-html:
	@echo "$(GREEN)Generating HTML test report...$(NC)"
	@cd tests && python shared/utilities/generate_test_report.py --format html

generate-test-report-json:
	@echo "$(GREEN)Generating JSON test report...$(NC)"
	@cd tests && python shared/utilities/generate_test_report.py --format json

show-test-evidence-organized:
	@echo "$(YELLOW)Organized Test Evidence Summary:$(NC)"
	@echo "Integration Evidence:"
	@ls -la tests/integration/output/payloads/ 2>/dev/null || echo "  No integration payloads found"
	@echo "Resource Evidence:"  
	@ls -la tests/resources/output/payloads/ 2>/dev/null || echo "  No resource payloads found"
	@echo "Multi-Provider Evidence:"
	@ls -la tests/multi_provider/output/payloads/ 2>/dev/null || echo "  No multi-provider payloads found"
	@echo "Consolidated Evidence:"
	@ls -la tests/evidence/ 2>/dev/null || echo "  No consolidated evidence found"

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
	@echo "$(YELLOW)Tool-Specific Testing:$(NC)"
	@echo "  test-integrations-basic     - Test integrations tool basic functionality"
	@echo "  test-resources-basic        - Test resources tool basic functionality"
	@echo "  test-resources-comprehensive - Test resources with dedicated prompts"
	@echo "  test-resources-all          - Run all resources tests"
	@echo "  test-all-tools              - Test both integrations and resources"
	@echo ""
	@echo "$(YELLOW)Organized Testing (Recommended):$(NC)"
	@echo "  test-integrations-basic-organized     - Organized integration basic tests"
	@echo "  test-integrations-advanced-organized  - Organized integration advanced tests"
	@echo "  test-integrations-all-organized       - All organized integration tests"
	@echo "  test-resources-basic-organized        - Organized resource basic tests"
	@echo "  test-resources-comprehensive-organized - Organized resource comprehensive tests"
	@echo "  test-resources-ultra-organized        - Organized resource ultra-complex tests"
	@echo "  test-resources-all-organized          - All organized resource tests"
	@echo "  test-all-providers-organized          - Organized multi-provider tests"
	@echo "  test-complete-organized               - Complete organized test suite"
	@echo ""
	@echo "$(YELLOW)Analysis Commands:$(NC)"
	@echo "  analyze-results - Show latest test results summary"
	@echo "  check-server    - Check if MCP server is running"
	@echo "  generate-test-report-organized - Generate comprehensive test report"
	@echo "  generate-test-report-html      - Generate HTML test report"
	@echo "  generate-test-report-json      - Generate JSON test report"
	@echo "  show-test-evidence-organized   - Show organized test evidence summary"
	@echo "  cleanup-test-data-organized    - Clean up old test data"
	@echo ""
	@echo "$(YELLOW)Interactive Commands:$(NC)"
	@echo "  chat-interactive - Start true interactive chat with the agent (RECOMMENDED)"
	@echo "  run-interactive  - Run enhanced interactive testing with predefined prompts"
	@echo "  run-example      - Run original chat example"
	@echo "  run-prompt       - Run with single prompt"
	@echo "  run-batch        - Run batch processor example"
	@echo ""
	@echo "$(YELLOW)Interactive Testing Commands:$(NC)"
	@echo "  test-interactive  - Interactive test with predefined scenarios"
	@echo "  test-single       - Single question interactive test"
	@echo ""
	@echo "$(YELLOW)Examples:$(NC)"
	@echo "  make chat-interactive                               # True interactive chat mode"
	@echo "  make run-interactive                                # Enhanced interactive mode with preset prompts"
	@echo "  make test-basic                                     # Quick 3-prompt test"
	@echo "  make test-custom PROMPTS_FILE=my_prompts.txt        # Custom prompts"
	@echo "  make test-single QUESTION=\"Show all AWS resources\" # Test a specific question"
