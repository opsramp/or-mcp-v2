.PHONY: setup test clean clean-all run-example run-browser run-integrations help

# Variables
PYTHON=python3
VENV=.venv
VENV_ACTIVATE=$(VENV)/bin/activate
PIP=$(VENV)/bin/pip
PYTEST=$(VENV)/bin/pytest
SERVER_PID_FILE=.server.pid

# Default target
all: setup
	@echo "========================================================"
	@echo "✅ Python client setup complete!"
	@echo "   To activate virtual environment: source $(VENV_ACTIVATE)"
	@echo "   To run tests: make test"
	@echo "   To run example: make run-browser"
	@echo "========================================================"

# Setup the virtual environment and install dependencies
setup: $(VENV_ACTIVATE)

$(VENV_ACTIVATE):
	@echo "========================================================"
	@echo "🐍 Setting up Python virtual environment..."
	@echo "========================================================"
	@$(PYTHON) -m venv $(VENV)
	@echo "📦 Installing dependencies..."
	@. $(VENV_ACTIVATE) && $(PIP) install -r requirements.txt
	@echo "✅ Setup complete! Activate with: source $(VENV_ACTIVATE)"

# Run tests (automatically starts the server if needed)
test: setup
	@echo "========================================================"
	@echo "🧪 Running Python client tests (will start server if needed)..."
	@echo "========================================================"
	@./run_tests.sh

# Run unit tests only (no server)
unit-test: setup
	@echo "========================================================"
	@echo "🧪 Running Python client unit tests (no server required)..."
	@echo "========================================================"
	@. $(VENV_ACTIVATE) && $(PYTEST) tests/ -v

# Run a specific example
run-example: setup
	@echo "========================================================"
	@echo "🚀 Running example $(EXAMPLE)..."
	@echo "========================================================"
	@if [ -z "$(EXAMPLE)" ]; then \
		echo "❌ Please specify an example with EXAMPLE=path/to/example.py"; \
		exit 1; \
	fi
	@. $(VENV_ACTIVATE) && $(PYTHON) $(EXAMPLE) $(ARGS)

# Run the browser-like example
run-browser: setup
	@echo "========================================================"
	@echo "🚀 Running browser-like client example..."
	@echo "========================================================"
	@. $(VENV_ACTIVATE) && $(PYTHON) examples/browser_like_example.py $(ARGS)
	@echo "✅ Browser-like example completed"

# Run the integrations example
run-integrations: setup
	@echo "========================================================"
	@echo "🚀 Running integrations client example..."
	@echo "========================================================"
	@. $(VENV_ACTIVATE) && $(PYTHON) examples/call_integrations.py $(ARGS)
	@echo "✅ Integrations example completed"

# Clean Python cache files
clean:
	@echo "========================================================"
	@echo "🧹 Cleaning Python cache files..."
	@echo "========================================================"
	@find . -type d -name "__pycache__" -exec rm -rf {} +
	@find . -type d -name "*.egg-info" -exec rm -rf {} +
	@find . -type f -name "*.pyc" -delete
	@find . -type f -name "*.pyo" -delete
	@find . -type f -name "*.pyd" -delete
	@rm -rf .pytest_cache
	@rm -f server.log
	@echo "✅ Python client cleaned"

# Clean all, including virtual environment
clean-all: clean
	@echo "========================================================"
	@echo "🧹 Cleaning all, including virtual environment..."
	@echo "========================================================"
	@rm -rf $(VENV)
	@rm -rf test_output
	@rm -f $(SERVER_PID_FILE)
	@echo "✅ Clean complete!"

# Help
help:
	@echo "========================================================"
	@echo "HPE OpsRamp MCP Python Client Makefile Help"
	@echo "========================================================"
	@echo "Available targets:"
	@echo "  all (default)    - Setup the virtual environment"
	@echo "  setup            - Setup the virtual environment"
	@echo "  test             - Run all tests (starts server if needed)"
	@echo "  unit-test        - Run unit tests only (no server)"
	@echo "  run-example      - Run a specific example (use EXAMPLE=path/to/example.py)"
	@echo "  run-browser      - Run the browser-like client example"
	@echo "  run-integrations - Run the integrations client example"
	@echo "  clean            - Clean Python cache files"
	@echo "  clean-all        - Clean everything, including virtual environment"
	@echo "  help             - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make run-browser ARGS='--debug'"
	@echo "  make run-example EXAMPLE=examples/check_server.py" 