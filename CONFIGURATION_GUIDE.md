# OpsRamp AI Agent - Complete Configuration Guide

This guide provides **all configuration settings** needed to get the OpsRamp AI Agent working successfully. Follow these steps to configure both the **MCP Server** and **AI Agent Client** components.

## ⚡ Quick Start (5 Minutes)

**For the impatient - minimal setup to get running:**

1. **Get your credentials:**
   - OpsRamp: auth key, secret, tenant ID from your OpsRamp console
   - LLM: OpenAI API key from [platform.openai.com](https://platform.openai.com/api-keys)

2. **Configure server:**
   ```bash
   cp config.yaml.template config.yaml
   # Edit config.yaml with your OpsRamp credentials
   ```

3. **Configure client:**
   ```bash
   cd client/agent
   cp .env.template .env
   # Edit .env with: OPENAI_API_KEY=sk-proj-your-key
   ```

4. **Start and test:**
   ```bash
   # Terminal 1: Start server (from project root)
   make all
   
   # Terminal 2: Test agent (from client/agent directory)
   cd client/agent && make test-basic
   ```

**For detailed setup instructions, continue reading below.**

## 🏗️ Architecture Overview

```
┌─────────────────────┐    ┌─────────────────────┐    ┌─────────────────────┐
│   AI Agent Client   │───▶│    MCP Server       │───▶│   OpsRamp API       │
│  (Python/LLM)      │    │  (Go/HTTP Server)   │    │  (REST API)         │
└─────────────────────┘    └─────────────────────┘    └─────────────────────┘
```

The agent consists of:
1. **MCP Server**: Go-based HTTP server that connects to OpsRamp APIs
2. **AI Agent Client**: Python-based client that uses LLMs to interact with the MCP server

## 📋 Prerequisites

### Required Software
- **Go 1.19+** (for MCP server)
- **Python 3.8+** (for AI agent client)
- **OpsRamp Account** with API access
- **At least one LLM API key** (OpenAI, Anthropic, or Google)

### Required API Access
- **OpsRamp API credentials** (auth key, secret, tenant ID)
- **LLM Provider API key** (choose one or more):
  - OpenAI API key
  - Anthropic API key  
  - Google Gemini API key

## 🔧 Configuration Files

### 1. MCP Server Configuration (`config.yaml`)

Create `config.yaml` in the project root directory:

```yaml
# MCP Server OpsRamp API Configuration
opsramp:
  # REQUIRED: Your OpsRamp instance URLs
  tenant_url: "https://your-tenant-instance.opsramp.com"
  auth_url: "https://your-tenant-instance.opsramp.com/tenancy/auth/oauth/token"
  
  # REQUIRED: OpsRamp API Credentials
  auth_key: "YOUR_OPSRAMP_AUTH_KEY_HERE"
  auth_secret: "YOUR_OPSRAMP_AUTH_SECRET_HERE"
  tenant_id: "YOUR_OPSRAMP_TENANT_ID_HERE"
  
  # Resource Management Settings
  resources:
    default_page_size: 50        # Default number of resources per page
    max_page_size: 1000         # Maximum resources per page
    cache_ttl: 300              # Cache time-to-live (5 minutes)
    enable_bulk_operations: true # Enable bulk resource operations
    max_bulk_size: 100          # Maximum bulk operation size
    
    # Performance Settings
    request_timeout: 30         # API request timeout (seconds)
    retry_attempts: 3           # Number of retry attempts
    retry_delay: 1000          # Retry delay (milliseconds)
    
    # Monitoring Settings
    enable_metrics: true        # Enable performance metrics
    metrics_interval: 60        # Metrics collection interval (seconds)
```

### 2. AI Agent Client Configuration (`.env`)

Create `.env` file in the `client/agent/` directory:

```bash
# =============================================================================
# LLM PROVIDER CONFIGURATION (Choose at least one)
# =============================================================================

# OpenAI Configuration
OPENAI_API_KEY=sk-proj-your-openai-api-key-here
OPENAI_MODEL=gpt-4

# Anthropic Configuration  
ANTHROPIC_API_KEY=sk-ant-api03-your-anthropic-api-key-here
ANTHROPIC_MODEL=claude-3-sonnet-20240229

# Google Gemini Configuration
GEMINI_API_KEY=your-google-gemini-api-key
GOOGLE_MODEL=gemini-1.5-pro

# =============================================================================
# AGENT CONFIGURATION
# =============================================================================

# Default LLM Provider (openai, anthropic, or gemini)
LLM_PROVIDER=openai

# MCP Server Connection
MCP_SERVER_URL=http://localhost:8080
CONNECTION_TIMEOUT=60
REQUEST_TIMEOUT=30

# =============================================================================
# OPSRAMP CONFIGURATION (Alternative to config.yaml)
# =============================================================================

# OpsRamp API Configuration (overrides config.yaml if set)
OPSRAMP_TENANT_URL=https://your-tenant-instance.opsramp.com
OPSRAMP_AUTH_URL=https://your-tenant-instance.opsramp.com/tenancy/auth/oauth/token
OPSRAMP_AUTH_KEY=your-opsramp-auth-key
OPSRAMP_AUTH_SECRET=your-opsramp-auth-secret
OPSRAMP_TENANT_ID=your-opsramp-tenant-id

# =============================================================================
# SERVER CONFIGURATION (Optional)
# =============================================================================

# MCP Server Settings
PORT=8080                    # Server port (default: 8080)
DEBUG=true                   # Enable debug logging
LOG_LEVEL=debug             # Logging level (debug, info, warn, error)

# =============================================================================
# TESTING CONFIGURATION (Optional)
# =============================================================================

# Testing Settings
SIMPLE_MODE=false           # Set to true for testing without MCP server
MOCK_MODE=false            # Set to true for mock testing
```

## 🚀 Setup Instructions

### Step 1: Get OpsRamp API Credentials

1. **Log into your OpsRamp instance**
2. **Navigate to Setup → Integrations → API Management**
3. **Create or locate your API credentials:**
   - **Auth Key**: Your API authentication key
   - **Auth Secret**: Your API authentication secret
   - **Tenant ID**: Your OpsRamp tenant identifier
   - **Tenant URL**: Your OpsRamp instance URL (e.g., `https://mycompany.opsramp.com`)

### Step 2: Get LLM Provider API Keys

Choose **at least one** LLM provider:

#### OpenAI (Recommended)
1. Visit [OpenAI API Platform](https://platform.openai.com/api-keys)
2. Create an API key
3. Copy the key (starts with `sk-proj-` or `sk-`)

#### Anthropic
1. Visit [Anthropic Console](https://console.anthropic.com/)
2. Create an API key
3. Copy the key (starts with `sk-ant-api03-`)

#### Google Gemini
1. Visit [Google AI Studio](https://aistudio.google.com/app/apikey)
2. Create an API key
3. Copy the key

### Step 3: Configure the MCP Server

1. **Copy the template:**
   ```bash
   cp config.yaml.template config.yaml
   ```

2. **Edit `config.yaml`** with your OpsRamp credentials:
   ```yaml
   opsramp:
     tenant_url: "https://your-actual-tenant.opsramp.com"
     auth_url: "https://your-actual-tenant.opsramp.com/tenancy/auth/oauth/token"
     auth_key: "your-actual-auth-key"
     auth_secret: "your-actual-auth-secret"
     tenant_id: "your-actual-tenant-id"
   ```

### Step 4: Configure the AI Agent Client

1. **Copy the template and create `.env` file**:
   ```bash
   cd client/agent
   cp .env.template .env
   ```

2. **Edit `.env` file** with your actual API keys:
   ```bash
   # Minimal setup - just add your OpenAI key
   OPENAI_API_KEY=sk-proj-your-actual-openai-key
   LLM_PROVIDER=openai
   
   # Or use Anthropic for better token efficiency
   # ANTHROPIC_API_KEY=sk-ant-api03-your-actual-anthropic-key
   # LLM_PROVIDER=anthropic
   
   # Or use Google Gemini
   # GEMINI_API_KEY=your-actual-gemini-key
   # LLM_PROVIDER=gemini
   ```

### Step 5: Install Dependencies

1. **Install Go dependencies:**
   ```bash
   go mod download
   ```

2. **Install Python dependencies:**
   ```bash
   cd client/agent
   pip install -r requirements.txt
   # Or use the setup script
   make setup
   ```

### Step 6: Start the MCP Server

```bash
# Complete build and start (recommended)
make all

# Or just build and run
make run

# Or run in debug mode for troubleshooting
make run-debug
```

**Expected output:**
```
INFO: Starting HPE OpsRamp MCP server
INFO: Registered tool: integrations
INFO: Registered tool: resources
INFO: Startup health check passed: successfully listed X integrations
INFO: Server listening on :8080
```

### Step 7: Test the Configuration

```bash
# Test server health and functionality (from project root)
make health-check
make test
make integration-test

# Test basic agent functionality (from client/agent directory)
cd client/agent
make test-basic

# Test specific tools with organized framework
make test-integrations-basic-organized
make test-resources-basic-organized
```

## 🔍 Configuration Validation

### Verify MCP Server Configuration

1. **Check server health:**
   ```bash
   curl http://localhost:8080/health
   ```
   Expected response: `{"status": "healthy", "uptime": "..."}`

2. **Check server readiness:**
   ```bash
   curl http://localhost:8080/readiness
   ```
   Expected response: `{"status": "ready", "tools": ["integrations", "resources"]}`

3. **Check debug information:**
   ```bash
   curl http://localhost:8080/debug
   ```

### Verify AI Agent Configuration

1. **Test agent connection:**
   ```bash
   cd client/agent
   python examples/chat_client.py --prompt "List available tools"
   ```

2. **Test with specific provider:**
   ```bash
   # Test with OpenAI
   LLM_PROVIDER=openai python examples/chat_client.py --prompt "What integrations do we have?"
   
   # Test with Anthropic
   LLM_PROVIDER=anthropic python examples/chat_client.py --prompt "What resources do we have?"
   ```

## 🛠️ Configuration Options Reference

### MCP Server Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `DEBUG` | `false` | Enable debug logging |
| `LOG_LEVEL` | `info` | Logging level (debug, info, warn, error) |
| `OPSRAMP_TENANT_URL` | - | OpsRamp tenant URL (overrides config.yaml) |
| `OPSRAMP_AUTH_URL` | - | OpsRamp auth URL (overrides config.yaml) |
| `OPSRAMP_AUTH_KEY` | - | OpsRamp auth key (overrides config.yaml) |
| `OPSRAMP_AUTH_SECRET` | - | OpsRamp auth secret (overrides config.yaml) |
| `OPSRAMP_TENANT_ID` | - | OpsRamp tenant ID (overrides config.yaml) |

### AI Agent Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `OPENAI_API_KEY` | - | OpenAI API key |
| `ANTHROPIC_API_KEY` | - | Anthropic API key |
| `GEMINI_API_KEY` | - | Google Gemini API key |
| `LLM_PROVIDER` | `openai` | Default LLM provider (openai, anthropic, gemini) |
| `OPENAI_MODEL` | `gpt-4` | OpenAI model to use |
| `ANTHROPIC_MODEL` | `claude-3-haiku-20240307` | Anthropic model to use |
| `GOOGLE_MODEL` | `gemini-1.5-flash` | Google model to use |
| `MCP_SERVER_URL` | `http://localhost:8080` | MCP server URL |
| `CONNECTION_TIMEOUT` | `60` | Connection timeout (seconds) |
| `REQUEST_TIMEOUT` | `30` | Request timeout (seconds) |
| `SIMPLE_MODE` | `false` | Run without MCP server connection |

### Available LLM Models

#### OpenAI Models
- `gpt-4` (recommended for complex tasks)
- `gpt-4-turbo`
- `gpt-3.5-turbo` (faster, more cost-effective)

#### Anthropic Models
- `claude-3-opus-20240229` (most capable)
- `claude-3-sonnet-20240229` (balanced)
- `claude-3-haiku-20240307` (fastest)

#### Google Gemini Models
- `gemini-1.5-pro` (most capable)
- `gemini-1.5-flash` (faster)

## 🚨 Troubleshooting

### Common Configuration Issues

#### 1. "MCP server is not accessible"
**Cause**: Server not running or wrong URL
**Solution**:
```bash
# Check if server is running (from project root)
make health-check

# Check server logs (from project root)
tail -f output/logs/server.log

# Stop any running server and restart (from project root)
make kill-server
make run

# Or run in debug mode for more information
make run-debug

# Quick build and run
make all
```

#### 2. "Failed to authenticate with OpsRamp API"
**Cause**: Invalid OpsRamp credentials
**Solution**:
- Verify credentials in OpsRamp console
- Check `config.yaml` for typos
- Ensure URLs don't have trailing slashes
- Verify tenant ID is correct

#### 3. "OpenAI API key is required"
**Cause**: Missing or invalid LLM API key
**Solution**:
```bash
# Check .env file exists
ls -la client/agent/.env

# Verify API key format
echo $OPENAI_API_KEY

# Test API key directly
curl -H "Authorization: Bearer $OPENAI_API_KEY" https://api.openai.com/v1/models
```

#### 4. "Request too large for gpt-4"
**Cause**: Token limit exceeded
**Solution**:
- Use Anthropic instead: `LLM_PROVIDER=anthropic`
- Use smaller model: `OPENAI_MODEL=gpt-3.5-turbo`
- Use basic test commands: `make test-basic`

### Configuration Validation Commands

```bash
# Validate server configuration and build (from project root)
make config
make build

# Check server health (from project root)
make health-check

# Run server tests (from project root)
make test
make integration-test

# Validate agent configuration (from client/agent directory)
cd client/agent
python -c "from src.opsramp_agent.utils.config import get_api_keys; print(get_api_keys())"

# Test end-to-end (from client/agent directory)
make test-single QUESTION="What tools are available?"

# Quick functionality test
make test-basic
```

## 📁 File Locations

```
or-mcp-v2/
├── config.yaml.template          # MCP server configuration template
├── config.yaml                   # MCP server configuration (create from template)
├── client/agent/
│   ├── .env.template             # AI agent configuration template
│   ├── .env                      # AI agent configuration (create from template)
│   ├── CONFIGURATION_GUIDE.md    # This configuration guide
│   ├── requirements.txt          # Python dependencies
│   └── src/opsramp_agent/
│       └── utils/config.py       # Configuration utilities
├── cmd/server/main.go            # MCP server entry point
├── common/config.go              # Go configuration utilities
└── output/logs/                  # Server logs
```

### Makefile Structure

The project has **two Makefiles** with different purposes:

1. **Root Makefile** (`/Makefile`): Server build, run, and comprehensive testing
   - **Build & Run**: `make build`, `make run`, `make run-debug`
   - **Testing**: `make test`, `make test-resources-all`, `make integration-test`
   - **Health Checks**: `make health-check`, `make kill-server`
   - **Security**: `make security-full`, `make security-scan`
   - **MCP-GO Library**: `make mcp-go-build`, `make mcp-go-test`
   - **Configuration**: `make config`, `make dirs`, `make clean-all`

2. **Agent Makefile** (`/client/agent/Makefile`): **Primary testing framework** ⭐
   - Organized testing framework
   - Multi-provider LLM testing
   - Comprehensive test suites
   - Interactive testing modes
   - Test reporting and evidence collection

**Use the Agent Makefile** (`client/agent/Makefile`) for all AI agent testing and development.

### Working Directories

**Important**: Commands must be run from the correct directory:

- **Server commands**: Run from **project root** (`/`)
  ```bash
  # Build and start server
  make all          # Complete build process
  make run          # Build and run server
  make run-debug    # Run in debug mode
  
  # Health and management
  make health-check # Quick health check
  make kill-server  # Stop running server
  
  # Testing
  make test                    # Unit tests
  make test-resources-all      # All resource tests
  make integration-test        # Integration tests
  make security-scan          # Security scan
  ```

- **Agent commands**: Run from **client/agent/** directory
  ```bash
  cd client/agent
  
  # All testing commands
  make test-basic
  make test-complete-organized
  make run-interactive
  ```

### Configuration Priority

The system loads configuration in this order (later sources override earlier ones):

1. **Default values** (hardcoded in the application)
2. **config.yaml** (MCP server configuration)
3. **Environment variables** (from shell or .env file)
4. **Command line arguments** (highest priority)

This means you can:
- Use `config.yaml` for OpsRamp settings
- Use `.env` for LLM API keys and agent settings
- Override anything with environment variables
- Override specific settings with command line flags

## ✅ Quick Start Checklist

- [ ] **OpsRamp credentials obtained** (auth key, secret, tenant ID)
- [ ] **LLM API key obtained** (OpenAI, Anthropic, or Google)
- [ ] **`config.yaml` created** with OpsRamp credentials
- [ ] **`.env` file created** in `client/agent/` with LLM API key
- [ ] **Dependencies installed** (`go mod download` and `pip install -r requirements.txt`)
- [ ] **MCP server built and started** (`make all`)
- [ ] **Server health verified** (`make health-check`)
- [ ] **Agent tested** (`cd client/agent && make test-basic`)

## 🎯 Next Steps

Once configured, you can run all commands from the `client/agent/` directory:

### **Basic Testing**
```bash
cd client/agent

# Quick functionality test
make test-basic

# Test specific tools
make test-integrations-basic-organized
make test-resources-basic-organized
```

### **Comprehensive Testing**
```bash
# Complete organized test suite
make test-complete-organized

# Multi-provider testing
make test-all-providers-organized

# Advanced testing by complexity
make test-integrations-advanced-organized
make test-resources-comprehensive-organized
make test-resources-ultra-organized
```

### **Interactive Usage**
```bash
# Enhanced interactive mode (recommended)
make run-interactive

# Single question testing
make test-single QUESTION="What integrations do we have?"

# Custom prompts file
make test-custom PROMPTS_FILE=my_prompts.txt MAX_TESTS=5
```

### **Reporting & Analysis**
```bash
# Generate comprehensive HTML report
make generate-test-report-html

# Generate JSON report for automation
make generate-test-report-json

# Show test evidence summary
make show-test-evidence-organized

# Clean up old test data
make cleanup-test-data-organized
```

### **Development & Debugging**
```bash
# Server health and management (from project root)
make health-check    # Quick server health check
make kill-server     # Stop any running server
make run-debug       # Run server in debug mode

# Agent debugging (from client/agent directory)
make analyze-results # Analyze latest test results
make test-interactive # Interactive testing with scenarios
```

### **Available Makefile Commands**

#### **Main Project Makefile** (`/Makefile`) - Server Operations

| Command | Description |
|---------|-------------|
| `make all` | Complete build process (clean, dirs, config, build) |
| `make build` | Build the server binary |
| `make run` | Build and run the MCP server |
| `make run-debug` | Run server in debug mode |
| `make health-check` | Quick server health check |
| `make kill-server` | Stop any running server |
| `make test` | Run server unit tests |
| `make test-resources-all` | All resource management tests |
| `make integration-test` | Integration tests |
| `make security-full` | Comprehensive security scan |
| `make config` | Setup configuration file |
| `make clean-all` | Clean all artifacts |
| `make help` | Show all server commands |

#### **Agent Makefile** (`/client/agent/Makefile`) - AI Agent Testing

| Command | Description |
|---------|-------------|
| `make setup` | Install dependencies and setup agent |
| `make test-basic` | Quick 3-prompt functionality test |
| `make test-complete-organized` | Complete organized test suite |
| `make test-integrations-basic-organized` | Basic integration tests |
| `make test-resources-basic-organized` | Basic resource tests |
| `make test-all-providers-organized` | Multi-provider comparison tests |
| `make generate-test-report-html` | Generate HTML test report |
| `make run-interactive` | Enhanced interactive mode |
| `make help` | Show all agent commands |

For detailed usage instructions, see the [README.md](README.md) and [testing documentation](tests/README.md). 