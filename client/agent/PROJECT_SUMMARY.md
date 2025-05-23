# ENHANCED MCP INTEGRATION TESTING - ORGANIZED PROJECT

## 📁 Directory Structure

```
client/agent/
├── tests/                              # Test scripts and test infrastructure
│   ├── enhanced_real_mcp_integration_test.py  # Main testing engine
│   ├── test_agent.py                   # Agent unit tests
│   └── test_data/                      # Test input files
│       ├── comprehensive_integration_prompts.txt    # 90 comprehensive scenarios
│       └── ultra_complex_integration_prompts.txt    # 31 ultra-complex scenarios
├── output/                             # All test results and logs
│   ├── enhanced_integration_test_XXXXXXX.log       # Execution logs
│   ├── request_response_payloads_XXXXXXX.jsonl     # Request/response data
│   └── comprehensive_test_analysis_XXXXXXX.json    # Test analysis
├── src/                                # Source code
├── examples/                           # Example scripts and prompts
├── docs/                               # Documentation
└── PROJECT_SUMMARY.md                  # This file
```

## 🎯 Core Components

### **Main Testing Engine**
- **`tests/enhanced_real_mcp_integration_test.py`** - Advanced testing engine
  - Reads prompts from test_data/ directory
  - Makes real calls to MCP server at http://localhost:8080
  - Writes all logs to output/ directory
  - Supports complex multi-tool scenarios

### **Test Data**
- **`tests/test_data/comprehensive_integration_prompts.txt`** - 90 comprehensive scenarios (15 categories)
- **`tests/test_data/ultra_complex_integration_prompts.txt`** - 31 ultra-complex scenarios (10 categories)

### **Test Output**
All test results automatically go to the `output/` directory:
- **Execution logs** - Detailed test execution with timestamps
- **Payload logs** - Structured request/response data (JSONL format)
- **Analysis reports** - Complete test metrics and analytics (JSON format)

## 🚀 How to Use

### **Run Basic Test (5 prompts)**
```bash
cd tests
python enhanced_real_mcp_integration_test.py --max-tests 5
```

### **Run Ultra-Complex Scenarios**
```bash
cd tests
python enhanced_real_mcp_integration_test.py \
  --prompts-file test_data/ultra_complex_integration_prompts.txt \
  --max-tests 3
```

### **Run All Comprehensive Tests**
```bash
cd tests
python enhanced_real_mcp_integration_test.py \
  --prompts-file test_data/comprehensive_integration_prompts.txt
```

### **Custom Prompts File**
```bash
cd tests
python enhanced_real_mcp_integration_test.py \
  --prompts-file /path/to/your/prompts.txt \
  --server-url http://your-server:8080
```

## 📊 What You Get

1. **Real MCP Server Integration** - No mocks, connects to http://localhost:8080
2. **Organized Output** - All logs automatically saved to output/ directory
3. **Full Request/Response Logging** - Every call traced and logged in structured format
4. **Advanced Analytics** - Performance metrics, complexity scoring, failure analysis
5. **Multi-Tool Orchestration** - Complex scenarios with multiple integration tool calls
6. **Complete Evidence Trail** - Structured logs for verification and debugging

## ✅ Verified Data

The system has been tested and verified with real OpsRamp integrations:
- **HPE Alletra** (INTG-2ed93041-eb92-40e9-b6b4-f14ad13e54fc)
- **Redfish Server** (INTG-f9e5d2aa-ee17-4e32-9251-493566ebdfca)  
- **VMware vCenter** (INTG-00ee85e2-1f84-4fc1-8cf8-5277ae6980dd)

All authentication configurations, profiles, and metadata are real and verified.

## 🧹 Clean Organization

- **No clutter** - All files in their proper directories
- **Clear separation** - Tests, data, and output are separated
- **Easy maintenance** - Simple to find and manage files
- **Scalable structure** - Easy to add new tests or data files 