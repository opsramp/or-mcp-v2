# OpsRamp AI Agent - Comprehensive Testing Framework

This directory contains a comprehensive testing framework for the OpsRamp AI Agent, organized to test both **Integration** and **Resources** functionality with real API evidence collection.

## 🏗️ Directory Structure

```
tests/
├── integration/                    # Integration Management Testing
│   ├── test_data/
│   │   ├── basic_integration_prompts.txt
│   │   ├── advanced_integration_prompts.txt
│   │   └── integration_scenarios.json
│   ├── scripts/
│   │   ├── run_integration_tests.py
│   │   └── validate_integration_results.py
│   └── output/
│       ├── logs/                   # Test execution logs
│       ├── payloads/              # Real API payloads
│       └── reports/               # Test reports
├── resources/                      # Resource Management Testing
│   ├── test_data/
│   │   ├── basic_resource_prompts.txt
│   │   ├── comprehensive_resource_prompts.txt
│   │   ├── ultra_complex_resource_prompts.txt
│   │   └── resource_scenarios.json
│   ├── scripts/
│   │   ├── run_resource_tests.py
│   │   └── validate_resource_results.py
│   └── output/
│       ├── logs/                   # Test execution logs
│       ├── payloads/              # Real API payloads
│       └── reports/               # Test reports
├── multi_provider/                 # Multi-Provider Testing
│   ├── test_data/
│   │   └── provider_scenarios.json
│   ├── scripts/
│   │   └── test_all_providers.py
│   └── output/
│       ├── logs/
│       ├── payloads/
│       └── reports/
├── evidence/                       # Consolidated Evidence
│   ├── api_payloads/              # All API evidence
│   ├── test_reports/              # Generated reports
│   └── screenshots/               # Visual evidence
└── shared/                        # Shared Testing Components
    ├── engines/                   # Core testing engines
    │   └── enhanced_real_mcp_integration_test.py  # Core test engine
    └── utilities/                 # Shared utilities
        ├── generate_test_report.py    # Report generator
        └── cleanup_test_data.py       # Data cleanup
```

## 🚀 Quick Start

### Basic Testing Commands

```bash
# Test integrations functionality
make test-integrations-basic-organized

# Test resources functionality  
make test-resources-basic-organized

# Run complete test suite
make test-complete-organized

# Generate test report
make generate-test-report-html
```

### Advanced Testing Commands

```bash
# Comprehensive integration tests
make test-integrations-advanced-organized

# Ultra-complex resource tests
make test-resources-ultra-organized

# Multi-provider comparison tests
make test-all-providers-organized

# Show test evidence
make show-test-evidence-organized
```

## 📊 Test Categories

### 1. Integration Management Tests

**Purpose**: Test the `integrations` tool functionality for managing OpsRamp integrations.

**Test Levels**:
- **Basic**: Simple integration queries (5-10 prompts)
- **Advanced**: Complex integration management scenarios (20-30 prompts)
- **All**: Complete integration test suite

**Key Features**:
- Real OpsRamp API integration (no mocks)
- Token-efficient prompts to avoid OpenAI limits
- Comprehensive API payload capture
- Integration status monitoring
- Configuration management testing

### 2. Resource Management Tests

**Purpose**: Test the `resources` tool functionality for managing OpsRamp resources.

**Test Levels**:
- **Basic**: Simple resource queries (5-10 prompts)
- **Comprehensive**: Detailed resource analysis (20-30 prompts)
- **Ultra**: Complex resource scenarios (50+ prompts)
- **All**: Complete resource test suite

**Key Features**:
- Real OpsRamp API integration (no mocks)
- Pagination and filtering support
- Resource discovery and inventory
- Performance metrics collection
- Hardware and software resource management

### 3. Multi-Provider Tests

**Purpose**: Compare AI model performance across different providers.

**Providers Tested**:
- OpenAI (GPT-3.5, GPT-4)
- Anthropic (Claude-3-Haiku, Claude-3-Sonnet)
- Google (Gemini-1.5-Flash, Gemini-1.5-Pro)

**Comparison Metrics**:
- Response accuracy
- Token efficiency
- Response time
- Error handling
- API call success rate

## 🔧 Test Configuration

### Environment Setup

Ensure your `.env` file contains:
```bash
# OpenAI Configuration
OPENAI_API_KEY=sk-proj-...
OPENAI_MODEL=gpt-4

# Anthropic Configuration  
ANTHROPIC_API_KEY=sk-ant-api03-...
ANTHROPIC_MODEL=claude-3-sonnet-20240229

# Google Configuration
GOOGLE_API_KEY=...
GOOGLE_MODEL=gemini-1.5-pro

# OpsRamp Configuration
OPSRAMP_API_KEY=...
OPSRAMP_TENANT_ID=...
OPSRAMP_BASE_URL=...
```

### Test Data Files

#### Integration Test Data
- `basic_integration_prompts.txt`: Token-efficient integration queries
- `advanced_integration_prompts.txt`: Complex integration scenarios
- `integration_scenarios.json`: Structured test scenarios

#### Resource Test Data
- `basic_resource_prompts.txt`: Simple resource queries
- `comprehensive_resource_prompts.txt`: Detailed resource analysis
- `ultra_complex_resource_prompts.txt`: Advanced resource scenarios
- `resource_scenarios.json`: Structured resource test scenarios

## 📈 Evidence Collection

### API Payload Evidence

All tests capture real API calls to OpsRamp:
- **Request payloads**: Complete API requests with headers
- **Response payloads**: Full API responses with data
- **Timing information**: Request/response timing
- **Error handling**: Failed requests and error responses

### Test Reports

Comprehensive reports are generated in multiple formats:
- **HTML**: Visual reports with charts and metrics
- **JSON**: Machine-readable test data
- **Text**: Simple text summaries

### Performance Metrics

- **Success Rate**: Percentage of successful tests
- **Response Time**: Average API response time
- **Token Efficiency**: Prompts processed per minute
- **Error Rate**: Failed requests and tool calls

## 🎯 Testing Best Practices

### Token Management

To avoid OpenAI token limits:
1. Use pagination in resource queries
2. Request specific data subsets
3. Implement token-efficient prompts
4. Use alternative models for large datasets

### Real API Testing

All tests use real OpsRamp APIs:
- No mock data or simulated responses
- Actual production environment testing
- Real integration and resource data
- Authentic API error handling

### Evidence Preservation

- All API payloads are preserved as evidence
- Test logs include complete execution traces
- Reports show real performance metrics
- Screenshots capture visual evidence

## 🔍 Troubleshooting

### Common Issues

1. **Token Limit Exceeded**
   - Use basic test variants
   - Implement pagination
   - Switch to Anthropic models

2. **API Connection Failures**
   - Check OpsRamp credentials
   - Verify network connectivity
   - Review API endpoint configuration

3. **Test Script Errors**
   - Check Python dependencies
   - Verify file permissions
   - Review log files for details

### Debug Commands

```bash
# Check server connectivity
make check-server

# Run single test for debugging
make test-single QUESTION="List first 3 integrations"

# Clean up test data
make cleanup-test-data-dry-organized

# Show detailed evidence
make show-test-evidence-organized
```

## 📋 Test Execution Examples

### Integration Testing

```bash
# Basic integration test
cd tests && python integration/scripts/run_integration_tests.py --complexity basic

# Advanced integration test with specific model
cd tests && python integration/scripts/run_integration_tests.py --complexity advanced --model claude-3-sonnet-20240229
```

### Resource Testing

```bash
# Basic resource test
cd tests && python resources/scripts/run_resource_tests.py --complexity basic

# Comprehensive resource test with pagination
cd tests && python resources/scripts/run_resource_tests.py --complexity comprehensive --max-resources 50
```

### Report Generation

```bash
# Generate HTML report
cd tests && python scripts/generate_test_report.py --format html --period daily

# Generate JSON report for automation
cd tests && python scripts/generate_test_report.py --format json --period weekly
```

## 🎉 Success Metrics

### Phase 1 Completion Criteria

- ✅ **Integration Tool**: Fully functional with real API
- ✅ **Resources Tool**: Fully functional with real API  
- ✅ **Test Framework**: Comprehensive testing infrastructure
- ✅ **Evidence Collection**: Real API payload capture
- ✅ **Multi-Provider Support**: OpenAI, Anthropic, Google
- ✅ **Token Management**: Efficient prompt handling
- ✅ **Report Generation**: Automated test reporting

### Quality Metrics

- **Success Rate**: > 95% for basic tests
- **Performance**: > 5 prompts/minute efficiency
- **Coverage**: Both integration and resource functionality
- **Evidence**: Complete API payload collection
- **Reliability**: Consistent test execution

## 🔄 Continuous Integration

The testing framework supports:
- Automated test execution
- Scheduled test runs
- Performance monitoring
- Evidence archival
- Report generation
- Multi-provider comparison

## 📞 Support

For issues or questions:
1. Check the troubleshooting section
2. Review test logs in `output/logs/`
3. Examine API payloads in `evidence/api_payloads/`
4. Generate diagnostic reports
5. Contact the development team

---

**Note**: This testing framework uses real OpsRamp APIs and collects actual production data. All evidence is preserved for verification and compliance purposes. 