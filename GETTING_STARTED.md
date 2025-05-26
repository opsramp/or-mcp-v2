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

# Verify Go installation
go version  # Should show Go 1.19+

# Verify Python installation
python3 --version  # Should show Python 3.8+

# Create Python virtual environment
python3 -m venv .venv

# Activate virtual environment
source .venv/bin/activate  # On Windows: .venv\Scripts\activate
```

### **Step 2: Configure OpsRamp Credentials**

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

### **Step 3: Configure AI Agent**

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

### **Step 4: Install Dependencies**

```bash
# Install Python dependencies (from client/agent directory)
make setup

# Verify installation
pip list | grep -E "(openai|anthropic|google|requests)"
```

### **Step 5: Build and Start the MCP Server**

```bash
# Return to project root
cd ../../

# Build and start the server
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

### **Step 6: Verify Agent Connectivity**

```bash
# Navigate back to agent directory
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
# Test integration management
make test-integrations-basic-organized

# Test resource management
make test-resources-basic-organized

# Test interactive mode
make run-interactive
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

## 🎯 Next Steps After Setup

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