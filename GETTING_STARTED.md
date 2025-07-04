# Getting Started Guide

Complete step-by-step guide to set up and replicate the HPE OpsRamp MCP Server with AI Agent Testing Platform.

## 🎯 What You'll Achieve

By following this guide, you'll have:
- ✅ **Fully functional MCP server** with Integration & Resource Management
- ✅ **AI agent testing platform** with multi-provider LLM support
- ✅ **Zero-vulnerability security** with professional-grade hardening
- ✅ **Comprehensive testing framework** with organized test suites
- ✅ **Real OpsRamp integration** with production API connectivity

## 📋 Prerequisites

### **System Requirements**
- **Operating System**: macOS, Linux, or Windows with WSL2
- **Go**: Version 1.19 or higher
- **Python**: Version 3.8 or higher
- **Git**: For cloning the repository
- **Internet Connection**: For downloading dependencies and API access

### **Account Requirements**
- **HPE OpsRamp Account** with API access
- **LLM Provider Account** (at least one):
  - OpenAI API key (recommended for testing)
  - Anthropic API key (excellent for token efficiency)
  - Google AI API key (alternative option)

### **Knowledge Prerequisites**
- Basic command line usage
- Understanding of API keys and configuration files
- Familiarity with environment variables

## 🚀 Step-by-Step Setup

### **Step 1: Clone and Initial Setup**

```bash
# Clone the repository
git clone https://github.com/opsramp/or-mcp-v2.git
cd or-mcp-v2
```

# Getting Started Guide

Complete step-by-step guide to set up and replicate the HPE OpsRamp MCP Server with AI Agent Testing Platform.

## 🎯 What You'll Achieve

By following this guide, you'll have:
- ✅ **Fully functional MCP server** with Integration & Resource Management
- ✅ **AI agent testing platform** with multi-provider LLM support
- ✅ **Zero-vulnerability security** with professional-grade hardening
- ✅ **Comprehensive testing framework** with organized test suites
- ✅ **Real OpsRamp integration** with production API connectivity

## 📋 Prerequisites

### **System Requirements**
- **Operating System**: macOS (Intel/ARM64), Linux, or Windows with WSL2
- **Go**: Version 1.21 or higher (automatically installed by `make configure`)
- **Python**: Version 3.8 or higher (automatically installed by `make configure`)
- **Git**: For cloning the repository (automatically installed by `make configure`)
- **Internet Connection**: For downloading dependencies and API access

### **Account Requirements**
- **HPE OpsRamp Account** with API access
- **LLM Provider Account** (at least one):
  - OpenAI API key (recommended for testing)
  - Anthropic API key (excellent for token efficiency)
  - Google AI API key (alternative option)

### **Knowledge Prerequisites**
- Basic command line usage
- Understanding of API keys and configuration files
- Familiarity with environment variables

## 🚀 Step-by-Step Setup

### **Step 1: Clone and Initial Setup**

```bash
# Clone the repository
git clone https://github.com/opsramp/or-mcp-v2.git
cd or-mcp-v2
```

### **Step 2: Auto-Configure System & Install Toolchains** ⭐ **NEW!**

**🎉 One-Command Setup for Any Platform!**

The project now includes an intelligent configuration system that automatically detects your OS and CPU architecture, then installs any missing toolchains.

```bash
# 🔧 Auto-configure your system (run this first!)
make configure

# This will:
# ✅ Detect your OS (macOS, Linux) and architecture (x86_64, ARM64)
# ✅ Install Go 1.21+ if missing or outdated
# ✅ Install Python 3.8+ if missing or outdated  
# ✅ Install Git, Make, curl/wget if missing
# ✅ Validate Go compilation for your platform
# ✅ Provide next-step guidance
```

**🍎 ARM64 MacBook Users**: This automatically handles Apple Silicon compatibility!

**🐧 Linux Users**: Automatically detects and uses your package manager (apt, dnf, yum, pacman).

**💡 System Information**: Want to see what's already installed?
```bash
make show-system-info  # Shows OS, architecture, and toolchain status
```

### **Step 3: Verify Auto-Configuration (Optional)**

After running `make configure`, you can verify everything is working:

```bash
# Verify Go installation
go version  # Should show Go 1.21+

# Verify Python installation
python3 --version  # Should show Python 3.8+
```

### **Step 3: Set Up Python Environment (Automated Method)**

The easiest way to set up the Python environment is to use our automated setup target:

```bash
# Run the automated Python setup
make python-setup

# Activate the virtual environment (after setup completes)
source .venv/bin/activate  # On Windows: .venv\Scripts\activate
```

**This will automatically:**
- Create a Python virtual environment (.venv)
- Install and upgrade pip
- Install the agent and all its dependencies
- Install the Python client libraries
- Configure the development environment

### **Step 3 (Alternative): Set Up Python Environment (Manual Method)**

If you prefer to set up the environment manually:

```bash
# Create Python virtual environment
python3 -m venv .venv

# Activate virtual environment
source .venv/bin/activate  # On Windows: .venv\Scripts\activate

# Install agent and dependencies
pip install --upgrade pip
pip install -e client/agent
pip install -e "client/agent[all]"

# Install Python client libraries
pip install -e client/python
```

### **Step 4: Configure OpsRamp Credentials**

```bash
# Copy the configuration template
cp config.yaml.template config.yaml

# Edit the configuration file
nano config.yaml  # Or use your preferred editor
```

**Configure the following in `config.yaml`:**
```yaml
opsramp:
  base_url: "https://your-tenant.api.opsramp.com"
  api_key: "your-opsramp-api-key"
  api_secret: "your-opsramp-api-secret"
  tenant_id: "your-tenant-id"

server:
  port: 8080
  host: "localhost"
  
logging:
  level: "info"
  format: "json"
```

**📋 For detailed configuration instructions, see [CONFIGURATION_GUIDE.md](CONFIGURATION_GUIDE.md)**

### **Step 5: Configure AI Agent**

```bash
# Navigate to agent directory
cd client/agent

# Copy environment template
cp .env.template .env

# Edit the environment file
nano .env  # Or use your preferred editor
```

**Configure the following in `.env`:**
```bash
# OpenAI Configuration (recommended)
OPENAI_API_KEY=your-openai-api-key
OPENAI_MODEL=gpt-4
OPENAI_MAX_TOKENS=4000

# Anthropic Configuration (optional but recommended)
ANTHROPIC_API_KEY=your-anthropic-api-key
ANTHROPIC_MODEL=claude-3-sonnet-20240229

# Google AI Configuration (optional)
GOOGLE_API_KEY=your-google-api-key
GOOGLE_MODEL=gemini-pro

# Server Configuration
MCP_SERVER_URL=http://localhost:8080
```

### **Step 6: Build and Start the MCP Server**

```bash
# Return to project root (if you're in client/agent)
cd ../../

# Build and start the server
# Note: 'make all' already includes python-setup, so Python environment is set up automatically
make all

# Verify server is running
make health-check
```

**Expected output:**
```
✅ Server is running on http://localhost:8080
✅ Health check passed
✅ Tools available: integrations, resources
```

### **Step 7: Verify Agent Connectivity**

```bash
# Make sure your virtual environment is activated
source .venv/bin/activate  # If not already activated

# Navigate to agent directory
cd client/agent

# Check server connectivity
make check-server

# Test basic functionality
make test-basic
```

**Expected output:**
```
✅ Server connectivity verified
✅ Tools discovered: integrations (10 actions), resources (14 actions)
✅ Basic tests passed: 3/3 successful
```

## 🧪 Testing Your Setup

### **Quick Validation Tests**

```bash
# Make sure your virtual environment is activated
source .venv/bin/activate  # If not already activated

# Test integration management
make test-integrations-basic-organized

# Test resource management
make test-resources-basic-organized

# Interactive modes
make chat-interactive     # True interactive chat (best user experience)
make run-interactive      # Test with preset prompts
```

### **Comprehensive Testing**

```bash
# Run complete test suite
make test-complete-organized

# Generate test report
make generate-test-report-html

# View test evidence
make show-test-evidence-organized
```

### **Multi-Provider Testing**

```bash
# Test all configured LLM providers
make test-all-providers-organized

# Compare provider performance
make test-provider-comparison
```

## 🔧 Troubleshooting Common Issues

### **🍎 ARM64 MacBook (Apple Silicon) Issues**

**Issue**: Compilation failures on ARM64 MacBooks

**Solutions**:
```bash
# Use the auto-configure system (recommended)
make configure  # Automatically detects ARM64 and installs compatible toolchains

# Manual verification of architecture
make show-system-info  # Should show "arm64 → arm64"

# If Go compilation still fails, try:
go clean -cache
go clean -modcache
make clean-all
make configure  # Re-run configuration
make all        # Rebuild everything
```

**Issue**: Go not found or wrong architecture

**Solutions**:
```bash
# Remove old Go installation (if any)
sudo rm -rf /usr/local/go

# Install Go for ARM64 automatically
make install-go

# Or install manually via Homebrew (recommended for Mac)
brew install go

# Verify Go is correctly configured for ARM64
go env GOOS GOARCH  # Should show "darwin arm64"
```

### **🐧 Linux Compilation Issues**

**Issue**: Missing build tools or package manager issues

**Solutions**:
```bash
# Use auto-configure to install everything
make configure

# Manual installation for specific distros:
# Ubuntu/Debian:
sudo apt-get update && sudo apt-get install -y golang-go python3 python3-pip git build-essential

# RHEL/CentOS/Fedora:
sudo dnf install -y golang python3 python3-pip git make gcc

# Arch Linux:
sudo pacman -S go python python-pip git base-devel
```

### **Python Environment Issues**

**Issue**: Errors with Python packages or dependencies

**Solutions**:
```bash
# Recreate Python environment from scratch
rm -rf .venv
make python-setup
source .venv/bin/activate

# Verify packages are installed correctly
pip list | grep -E "(openai|anthropic|google|requests)"

# Manual package installation if needed
pip install openai anthropic google-generativeai requests
```

### **Test Data Issues**

**Issue**: Missing test data files for agent tests

**Solutions**:
```bash
# Ensure test data directories exist
mkdir -p client/agent/test_data
mkdir -p client/agent/tests/shared/engines/test_data

# Create prompt files for different test types
# For basic integration tests (simpler version to avoid token limits):
cat > client/agent/tests/shared/engines/test_data/basic_integration_prompts.txt << EOF
List all integrations
Show integration types
What resources do we have?
EOF

# For comprehensive integration tests:
cat > client/agent/tests/shared/engines/test_data/comprehensive_integration_prompts.txt << EOF
List all available integrations
Show me integration types available
How do I create a new integration?
How can I update an existing integration?
How do I view integration details?
EOF

# For interactive mode:
cat > client/agent/tests/shared/engines/test_data/interactive_test_scenarios.txt << EOF
List all integrations
Show me AWS integrations
How do I create a new integration?
List all available resource types
Show me all resources
EOF

# Copy to alternative location that might be checked
cp client/agent/tests/shared/engines/test_data/*.txt client/agent/test_data/

# Verify test data files exist
ls -la client/agent/test_data/
ls -la client/agent/tests/shared/engines/test_data/
```

### **Server Won't Start**

**Issue**: `make all` fails or server doesn't start

**Solutions**:
```bash
# Check Go installation
go version

# Verify configuration
cat config.yaml

# Check for port conflicts
lsof -i :8080

# Try manual build
go build -o build/server cmd/server/main.go
./build/server
```

### **Agent Can't Connect to Server**

**Issue**: `make check-server` fails

**Solutions**:
```bash
# Verify server is running
curl http://localhost:8080/health

# Check firewall settings
# Ensure port 8080 is accessible

# Verify configuration
grep MCP_SERVER_URL .env
```

### **OpsRamp API Authentication Fails**

**Issue**: API calls return authentication errors

**Solutions**:
```bash
# Verify credentials in config.yaml
# Check OpsRamp tenant URL format
# Ensure API key has proper permissions
# Test credentials manually:
curl -H "Authorization: Bearer your-api-key" \
     "https://your-tenant.api.opsramp.com/api/v2/tenants/your-tenant-id/integrations"
```

### **LLM API Errors**

**Issue**: AI agent tests fail with API errors

**Solutions**:
```bash
# Check API key validity
# Verify model names in .env
# Check rate limits and quotas
# Test with different provider:
make test-single PROVIDER=anthropic QUESTION="Test question"
```

### **Token Limit Errors**

**Issue**: "Request too large" or token limit errors

**Solutions**:
```bash
# Use Anthropic (better token efficiency)
export DEFAULT_PROVIDER=anthropic

# Use basic complexity tests
make test-resources-basic-organized

# Reduce test scope
make test-custom MAX_TESTS=3
```

### **Missing MCP-Go Library Issues**

**Issue**: Missing internal/mcp-go directory or "mcp-go package not found" errors

**Solutions**:
```bash
# Verify the vendored library was properly downloaded
ls -la internal/mcp-go

# If the directory is missing, re-clone the repository
git clone https://github.com/opsramp/or-mcp-v2.git
cd or-mcp-v2

# Verify all files are present
ls -la internal/mcp-go/mcp/
```

### **Interactive Chat Issues**

**Issue**: Problems with chat-interactive mode

**Solutions**:
```bash
# Ensure virtual environment is activated
source .venv/bin/activate

# Check server is running
curl http://localhost:8080/health

# Ensure all dependencies are installed
cd client/agent && pip install -e ".[all]"

# Run with debug output
cd client/agent && DEBUG=true make chat-interactive
```

## 🎯 Next Steps After Setup

### **Start Interacting with the Agent**
```bash
# Start an interactive chat session
make chat-interactive

# Ask questions like:
# "What integrations do we have?"
# "Show me all resources with critical status"
# "How do I create a new integration?"
```

### **Explore Integration Management**
```bash
# Learn about integration capabilities
cat INTEGRATIONS.md

# Test integration scenarios
make test-single QUESTION="What integrations do we have?"
make test-single QUESTION="Show me integration types available"
```

### **Explore Resource Management**
```bash
# Learn about resource capabilities
cat RESOURCES.md

# Test resource scenarios
make test-single QUESTION="What resources do we have?"
make test-single QUESTION="Show me server resources with high CPU"
```

### **Build Custom Workflows**
```bash
# Create custom test scenarios
echo "Show me all critical alerts" > my_prompts.txt
echo "List integrations created this week" >> my_prompts.txt
make test-custom PROMPTS_FILE=my_prompts.txt

# Develop automation scripts
python shared/engines/enhanced_real_mcp_integration_test.py \
  --custom-prompts my_prompts.txt \
  --output-format json
```

### **Security Validation**
```bash
# Run comprehensive security scans
make security-scan

# Verify zero vulnerabilities
make security-go
make security-python
```

## 📊 Replication Success Metrics

### **Functional Verification**
- ✅ **Server Health**: `make health-check` passes
- ✅ **Tool Discovery**: 24 total actions available (10 integration + 14 resource)
- ✅ **API Connectivity**: Real OpsRamp data retrieval
- ✅ **LLM Integration**: Successful tool calling with chosen provider

### **Testing Verification**
- ✅ **Basic Tests**: 100% success rate on `make test-basic`
- ✅ **Integration Tests**: Successful integration management operations
- ✅ **Resource Tests**: Successful resource management operations
- ✅ **Interactive Mode**: Responsive to natural language queries

### **Security Verification**
- ✅ **Zero Vulnerabilities**: `make security-scan` shows no issues
- ✅ **Proper Authentication**: OpsRamp API calls authenticated
- ✅ **Secure Configuration**: No credentials in logs or outputs

### **Performance Verification**
- ✅ **Response Times**: Sub-second responses for basic operations
- ✅ **Token Efficiency**: Tests complete within LLM provider limits
- ✅ **Resource Usage**: Reasonable CPU and memory consumption

## 🎉 Success! What You've Accomplished

Congratulations! You now have:

1. **Production-Ready MCP Server** with comprehensive OpsRamp integration
2. **AI Agent Testing Platform** with multi-provider LLM support
3. **Zero-Vulnerability Security** with enterprise-grade hardening
4. **Comprehensive Testing Framework** with organized test suites
5. **Real API Integration** with actual OpsRamp data (no mocks)
6. **True Interactive Chat** mode for direct, intuitive agent interaction

## 📚 Continue Learning

- **[🔗 INTEGRATIONS.md](INTEGRATIONS.md)** - Deep dive into integration management
- **[🖥️ RESOURCES.md](RESOURCES.md)** - Explore resource management capabilities
- **[⚙️ CONFIGURATION_GUIDE.md](CONFIGURATION_GUIDE.md)** - Advanced configuration options
- **[📖 README.md](README.md)** - Project overview and quick reference

## 🤝 Getting Help

If you encounter issues:

1. **Check the troubleshooting section** above
2. **Review the configuration guide** for detailed setup instructions
3. **Run diagnostic commands** to identify specific issues
4. **Check logs** in `output/` directory for detailed error information

---

**🎯 Ready to build amazing AI agents with OpsRamp integration!** Start exploring the interactive testing mode and build your own custom workflows.