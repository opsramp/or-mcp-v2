# OpsRamp AI Agent - Non-Interactive Testing & Automation

This document provides comprehensive information about using the OpsRamp AI Agent in non-interactive modes with the organized testing framework.

## 🎯 Overview

The OpsRamp AI Agent supports multiple non-interactive modes for automation, testing, and batch processing:

1. **Organized Testing Framework** - Comprehensive test suites for Integration and Resource functionality
2. **Single Prompt Testing** - Process individual queries for validation
3. **Batch Processing** - Process multiple prompts from structured files
4. **Multi-Provider Testing** - Compare performance across AI providers
5. **Evidence Collection** - Automated API payload capture and reporting

## 🚀 Setup

Before using any mode, ensure the testing framework is properly set up:

```bash
# Navigate to the agent directory
cd client/agent

# Setup the testing framework
make setup

# Verify MCP server connectivity
make check-server
```

## 🧪 Organized Testing Framework

### **Integration Management Testing**

Test the `integrations` tool functionality with organized test suites:

```bash
# Basic integration tests (5-10 prompts)
make test-integrations-basic-organized

# Advanced integration tests (20-30 prompts)
make test-integrations-advanced-organized

# All integration tests
make test-integrations-all-organized
```

**Features**:
- Real OpsRamp API integration (no mocks)
- Token-efficient prompts to avoid OpenAI limits
- Complete API payload capture
- Structured output in `tests/integration/output/`

### **Resource Management Testing**

Test the `resources` tool functionality with comprehensive test suites:

```bash
# Basic resource tests (5-10 prompts)
make test-resources-basic-organized

# Comprehensive resource tests (20-30 prompts)
make test-resources-comprehensive-organized

# Ultra-complex resource tests (50+ prompts)
make test-resources-ultra-organized

# All resource tests
make test-resources-all-organized
```

**Features**:
- Real OpsRamp API integration (no mocks)
- Pagination and filtering support
- Performance metrics collection
- Structured output in `tests/resources/output/`

### **Multi-Provider Testing**

Compare AI model performance across different providers:

```bash
# Test all providers (OpenAI, Anthropic, Google)
make test-all-providers-organized

# Test integration functionality across providers
make test-providers-integration-organized

# Test resource functionality across providers
make test-providers-resources-organized
```

**Providers Supported**:
- **OpenAI**: GPT-3.5-Turbo, GPT-4
- **Anthropic**: Claude-3-Haiku, Claude-3-Sonnet
- **Google**: Gemini-1.5-Flash, Gemini-1.5-Pro

### **Master Test Commands**

Run comprehensive test suites:

```bash
# Run complete test suite (all categories)
make test-complete-organized

# Run all tests and generate HTML report
make test-complete-organized && make generate-test-report-html
```

## 📊 Evidence Collection & Reporting

### **Automated Report Generation**

Generate comprehensive reports in multiple formats:

```bash
# Generate HTML report with visual dashboards
make generate-test-report-html

# Generate JSON report for automation
make generate-test-report-json

# Generate text report for simple analysis
make generate-test-report-organized
```

### **Evidence Management**

Manage test evidence and API payloads:

```bash
# Show test evidence summary
make show-test-evidence-organized

# Clean up old test data (30+ days)
make cleanup-test-data-organized

# Show what would be cleaned (dry run)
make cleanup-test-data-dry-organized
```

## 🎯 Single Prompt Testing

For immediate validation and debugging:

```bash
# Test a specific question
make test-single QUESTION="List all integrations"

# Test integration-specific query
make test-single QUESTION="Show me the first 3 integrations with their status"

# Test resource-specific query
make test-single QUESTION="Show me first 5 resources"
```

**Options**:
- Uses the core testing engine for consistency
- Captures API payloads for evidence
- Provides detailed execution logs
- Supports all AI providers

## 📁 Custom Batch Processing

Process custom prompt files with the organized framework:

```bash
# Custom prompts file with specific test count
make test-custom PROMPTS_FILE=my_prompts.txt MAX_TESTS=5

# Use specific test data directory
make test-custom PROMPTS_FILE=tests/integration/test_data/basic_integration_prompts.txt
```

### **Prompt File Format**

Create structured prompt files for testing:

```
# Integration Management Prompts
List all integrations
Show me the first 3 integrations
What integrations are currently enabled?

# Resource Management Prompts  
Show me first 5 resources
List resource types available
Get basic resource inventory
```

**Best Practices**:
- Use token-efficient prompts to avoid provider limits
- Group related prompts by functionality
- Include complexity levels (basic → advanced → ultra)
- Add comments for organization

## 🔧 Advanced Configuration

### **Environment Variables**

Configure the testing framework with environment variables:

```bash
# AI Provider Configuration
export OPENAI_API_KEY="sk-proj-..."
export ANTHROPIC_API_KEY="sk-ant-api03-..."
export GOOGLE_API_KEY="..."

# Model Selection
export OPENAI_MODEL="gpt-4"
export ANTHROPIC_MODEL="claude-3-sonnet-20240229"
export GOOGLE_MODEL="gemini-1.5-pro"

# MCP Server Configuration
export MCP_SERVER_URL="http://localhost:8080"
```

### **Test Configuration Options**

Customize test execution with parameters:

```bash
# Specify maximum number of tests
make test-custom PROMPTS_FILE=prompts.txt MAX_TESTS=10

# Use specific AI provider
OPENAI_MODEL=gpt-3.5-turbo make test-integrations-basic-organized

# Enable debug mode
DEBUG=true make test-resources-basic-organized
```

## 📈 Performance Monitoring

### **Real-Time Analysis**

Monitor test performance and results:

```bash
# Show latest test results
make analyze-results

# Monitor test execution in real-time
tail -f tests/integration/output/logs/*.log

# Check API payload evidence
ls -la tests/evidence/api_payloads/
```

### **Performance Metrics**

The framework tracks comprehensive metrics:

- **Success Rate**: Percentage of successful tests
- **Response Time**: Average API response time
- **Token Efficiency**: Prompts processed per minute
- **Error Rate**: Failed requests and tool calls
- **Provider Comparison**: Cross-platform performance

## 🛠️ Development & Debugging

### **Development Testing**

Use development-focused commands for debugging:

```bash
# Single test for development
make dev-test

# Interactive testing mode
make run-interactive

# Check server connectivity
make check-server
```

### **Debugging Options**

Debug test execution and API interactions:

```bash
# Enable verbose logging
DEBUG=true make test-integrations-basic-organized

# View detailed API payloads
cat tests/integration/output/payloads/*.jsonl | jq .

# Analyze test failures
grep -i error tests/integration/output/logs/*.log
```

## 🎉 Automation Integration

### **CI/CD Pipeline Integration**

Integrate the testing framework into CI/CD pipelines:

```bash
#!/bin/bash
# CI/CD Test Script

# Setup
cd client/agent
make setup

# Verify connectivity
make check-server || exit 1

# Run basic validation tests
make test-integrations-basic-organized || exit 1
make test-resources-basic-organized || exit 1

# Generate evidence report
make generate-test-report-json

# Archive results
cp tests/evidence/test_reports/*.json $CI_ARTIFACTS_DIR/
```

### **Scheduled Testing**

Set up automated testing schedules:

```bash
# Crontab example for daily testing
0 2 * * * cd /path/to/agent && make test-complete-organized && make generate-test-report-html

# Weekly comprehensive testing
0 1 * * 0 cd /path/to/agent && make test-all-providers-organized && make cleanup-test-data-organized
```

## 📋 Output Structure

### **Organized Output Directories**

The framework maintains organized output structure:

```
tests/
├── integration/output/
│   ├── logs/           # Integration test execution logs
│   ├── payloads/       # Integration API evidence
│   └── reports/        # Integration test reports
├── resources/output/
│   ├── logs/           # Resource test execution logs
│   ├── payloads/       # Resource API evidence
│   └── reports/        # Resource test reports
├── multi_provider/output/
│   ├── logs/           # Multi-provider test logs
│   ├── payloads/       # Provider comparison evidence
│   └── reports/        # Provider comparison reports
└── evidence/
    ├── api_payloads/   # Consolidated API evidence
    ├── test_reports/   # Generated reports archive
    └── screenshots/    # Visual evidence (if applicable)
```

### **File Naming Conventions**

Consistent naming for easy identification:

- **Logs**: `test_execution_YYYYMMDD_HHMMSS.log`
- **Payloads**: `api_payloads_YYYYMMDD_HHMMSS.jsonl`
- **Reports**: `test_report_daily_YYYYMMDD_HHMMSS.html`

## 🔍 Troubleshooting

### **Common Issues**

**Issue: "MCP server is not accessible"**
```bash
# Check server status
make check-server

# Start server (from project root)
cd ../../ && make run &
```

**Issue: "Token limit exceeded"**
```bash
# Use basic test variants
make test-integrations-basic-organized

# Switch to Anthropic models
export ANTHROPIC_MODEL="claude-3-haiku-20240307"
make test-resources-basic-organized
```

**Issue: "No test evidence found"**
```bash
# Run tests to generate evidence
make test-resources-basic-organized

# Check evidence directories
make show-test-evidence-organized
```

### **Debug Commands**

Comprehensive debugging toolkit:

```bash
# Check framework status
make check-server
make analyze-results

# View recent logs
tail -20 tests/integration/output/logs/*.log

# Validate test data files
ls -la tests/*/test_data/

# Check API evidence
find tests/ -name "*.jsonl" -exec wc -l {} \;
```

## 📞 Support

For issues with non-interactive modes:

1. **Check Documentation**: Review this guide and `tests/README.md`
2. **Run Diagnostics**: Use `make check-server` and `make analyze-results`
3. **Review Evidence**: Check `tests/evidence/` for API payloads and reports
4. **Generate Reports**: Use `make generate-test-report-html` for detailed analysis
5. **Contact Support**: Provide test logs and evidence files

---

**Ready for automation?** Start with `make test-integrations-basic-organized` to validate your setup, then explore the comprehensive non-interactive testing capabilities!

*Framework Version: 2.0 (Organized)*  
*Status: Production Ready ✅* 