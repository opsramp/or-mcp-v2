# OpsRamp AI Agent - Enhanced Integration Testing Platform

A sophisticated AI-powered agent with comprehensive OpsRamp integrations management and advanced testing capabilities.

## 🎯 Features

- **Real MCP Integration** - No mocks, connects directly to OpsRamp MCP server
- **Comprehensive Testing** - 121 test scenarios covering all integration use cases
- **Advanced Analytics** - Performance metrics, complexity scoring, failure analysis
- **Multi-Tool Orchestration** - Complex scenarios with multiple integration tool calls
- **Complete Evidence Trail** - Structured logs for verification and debugging
- **Interactive Chat Interface** - Natural language queries for integrations management
- **Organized Architecture** - Clean separation of tests, data, and output

## 📁 Project Structure

```
client/agent/
├── tests/                              # 🧪 All test-related files
│   ├── enhanced_real_mcp_integration_test.py  # Advanced testing engine
│   └── test_data/                      # 📊 Test input files
│       ├── comprehensive_integration_prompts.txt    # 90 comprehensive scenarios
│       └── ultra_complex_integration_prompts.txt    # 31 ultra-complex scenarios
├── output/                             # 📁 All test results (auto-created)
│   ├── enhanced_integration_test_*.log       # Execution logs
│   ├── request_response_payloads_*.jsonl     # Request/response data  
│   └── comprehensive_test_analysis_*.json    # Test analytics
├── src/                                # 💻 Source code
│   └── opsramp_agent/                  # Main package
├── examples/                           # 📖 Example scripts
├── docs/                               # 📚 Documentation
└── Makefile                            # 🛠️ Build and test automation
```

## 🚀 Quick Start

### Prerequisites

- Python 3.8+
- Access to an OpsRamp MCP server running on `http://localhost:8080`
- OpenAI API key for LLM functionality

### Installation

```bash
# Install dependencies
make setup

# Create .env file with your API key
echo "OPENAI_API_KEY=your_openai_key_here" > .env
```

### Basic Testing

```bash
# Quick test (3 prompts)
make test-basic

# Check results
make analyze-results
```

## 🧪 Comprehensive Testing

### Test Suites Available

| Command | Description | Prompts | Duration | Use Case |
|---------|-------------|---------|-----------|----------|
| `make test-basic` | Quick validation | 3 | ~1 min | Development/CI |
| `make test-medium` | Medium coverage | 10 | ~3 min | Regular testing |
| `make test-complex` | Ultra-complex scenarios | 5 | ~2 min | Advanced testing |
| `make test-comprehensive` | Full test suite | 90 | ~15 min | Complete validation |
| `make test-all` | Basic + Complex | 8 | ~3 min | Standard workflow |

### Advanced Testing Options

```bash
# Custom prompts file
make test-custom PROMPTS_FILE=my_prompts.txt MAX_TESTS=5

# Development testing
make dev-test

# Check MCP server status
make check-server

# Clean test outputs
make clean-output
```

## 📊 Test Categories Covered

### **Basic Integration Management (15 scenarios)**
- Discovery & inventory
- Integration status checks
- Basic configuration queries

### **Advanced Analysis (45 scenarios)**
- Comparative analysis across integrations
- Detailed troubleshooting workflows
- Security and compliance analysis
- Capacity planning and expansion

### **Complex Multi-Tool Scenarios (31 scenarios)**
- Ultra-deep comparative analysis
- Multi-step diagnostic workflows
- Business intelligence reporting
- Risk assessment and mitigation

### **Real-World Operations (30 scenarios)**
- Incident response workflows
- Lifecycle management
- Advanced monitoring setup
- Strategic planning

## 🎯 Integration Tool Capabilities

The agent interacts with the OpsRamp integrations tool with 10 core actions:

| Action | Description | Parameters |
|--------|-------------|------------|
| `list` | List all integrations | None |
| `get` | Get basic integration info | id |
| `getDetailed` | Get comprehensive details | id |
| `create` | Create new integration | config |
| `update` | Update existing integration | id, config |
| `delete` | Delete integration | id |
| `enable` | Enable disabled integration | id |
| `disable` | Disable active integration | id |
| `listTypes` | List available integration types | None |
| `getType` | Get integration type details | id |

## 📈 Analytics & Reporting

### Automated Analysis

Every test run generates comprehensive analytics:

- **Success Rate** - Percentage of successful tests
- **Performance Metrics** - Duration, tool calls, complexity scores
- **Category Analysis** - Performance by test category
- **Tool Usage Statistics** - Most frequently used integrations actions
- **Failure Analysis** - Pattern recognition for failed tests

### Output Files

```bash
output/
├── enhanced_integration_test_[session].log     # Detailed execution log
├── request_response_payloads_[session].jsonl  # Structured API trace
└── comprehensive_test_analysis_[session].json # Complete analytics
```

## 🔧 Verified Integrations

The system has been tested with real OpsRamp integrations:

- **HPE Alletra Storage** (`INTG-2ed93041-eb92-40e9-b6b4-f14ad13e54fc`)
- **Redfish Server** (`INTG-f9e5d2aa-ee17-4e32-9251-493566ebdfca`)  
- **VMware vCenter** (`INTG-00ee85e2-1f84-4fc1-8cf8-5277ae6980dd`)

All authentication configurations, profiles, and metadata are real and verified.

## 💬 Interactive Usage

### Chat Interface

```bash
# Interactive mode
make run-example

# Single prompt
make run-prompt

# Batch processing
make run-batch
```

### Example Queries

- "List all integrations and tell me which ones need updates"
- "Show me comprehensive details on the HPE Alletra integration"
- "Compare our VMware and storage integrations"
- "What integration types are available that we haven't implemented?"
- "Walk me through troubleshooting the down integration"

## 🛠️ Development

### Configuration

Environment variables (`.env` file):
```
OPENAI_API_KEY=your_openai_key_here
ANTHROPIC_API_KEY=your_anthropic_key_here  # Optional
MCP_SERVER_URL=http://localhost:8080       # Optional
```

### Adding Custom Tests

1. Create your prompts file in `tests/test_data/`
2. Run with: `make test-custom PROMPTS_FILE=tests/test_data/my_prompts.txt`

### Debugging

```bash
# Single test for debugging
make dev-test

# Check server connectivity
make check-server

# View detailed logs
tail -f output/enhanced_integration_test_*.log
```

## 📚 Documentation

- [Project Summary](PROJECT_SUMMARY.md) - Organized overview
- [Non-Interactive Modes](docs/non_interactive.md) - Automation guides
- [Configuration Options](docs/configuration.md) - Advanced settings

## 🎉 Success Stories

Recent test results demonstrate:
- **100% Success Rate** on comprehensive integration testing
- **Real MCP Server Integration** with no mocks or simulations
- **Advanced Multi-Tool Orchestration** handling complex scenarios
- **Complete Evidence Trail** for verification and debugging

## 🔄 Continuous Integration

The test suite is designed for CI/CD integration:

```bash
# In your CI pipeline
make setup
make test-basic
make analyze-results
```

## 📋 Makefile Reference

Run `make help` to see all available commands with descriptions and examples.

## 🤝 Contributing

1. Add test scenarios to `tests/test_data/`
2. Run comprehensive tests to verify
3. Update documentation as needed
4. Submit pull request with test evidence

## 📄 License

[License information]
