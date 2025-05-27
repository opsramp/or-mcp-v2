# Integration Management Documentation

Complete guide to HPE OpsRamp integration management capabilities through the MCP server and AI agent testing platform.

## 🔗 Overview

The Integration Management system provides comprehensive tools for managing HPE OpsRamp integrations through AI agents. With **10 comprehensive actions**, you can perform complete lifecycle management of integrations.

## 🛠️ Available Integration Actions

### **Core Integration Operations**

#### 1. **`integrations:list`** - List All Integrations
**Purpose**: Retrieve all integrations with optional filtering and search capabilities

**Parameters**:
- `limit` (optional): Maximum number of results to return
- `offset` (optional): Number of results to skip for pagination
- `search` (optional): Search term to filter integrations

**Example Usage**:
```bash
make test-single QUESTION="List all integrations"
make test-single QUESTION="Show me the first 10 integrations"
make test-single QUESTION="Find integrations containing 'monitoring'"
```

**Response**: Array of integration objects with basic information

---

#### 2. **`integrations:get`** - Get Basic Integration Information
**Purpose**: Retrieve basic information about a specific integration

**Parameters**:
- `integrationId` (required): Unique identifier of the integration

**Example Usage**:
```bash
make test-single QUESTION="Get details for integration ID 12345"
make test-single QUESTION="Show me integration information for ID abc-def-123"
```

**Response**: Basic integration object with core properties

---

#### 3. **`integrations:getDetailed`** - Get Comprehensive Integration Details
**Purpose**: Retrieve comprehensive details about a specific integration including configuration, status, and metadata

**Parameters**:
- `integrationId` (required): Unique identifier of the integration

**Example Usage**:
```bash
make test-single QUESTION="Get detailed information for integration 12345"
make test-single QUESTION="Show me comprehensive details for integration abc-def-123"
```

**Response**: Detailed integration object with full configuration and metadata

---

### **Integration Lifecycle Management**

#### 4. **`integrations:create`** - Create New Integration
**Purpose**: Create a new integration with specified configuration

**Parameters**:
- `name` (required): Name of the integration
- `type` (required): Integration type identifier
- `configuration` (required): Integration-specific configuration object
- `description` (optional): Description of the integration

**Example Usage**:
```bash
make test-single QUESTION="Create a new monitoring integration named 'Production Servers'"
make test-single QUESTION="Set up a new integration for AWS CloudWatch"
```

**Response**: Created integration object with assigned ID

---

#### 5. **`integrations:update`** - Update Integration Configuration
**Purpose**: Update an existing integration's configuration or properties

**Parameters**:
- `integrationId` (required): Unique identifier of the integration
- `configuration` (optional): Updated configuration object
- `name` (optional): Updated name
- `description` (optional): Updated description

**Example Usage**:
```bash
make test-single QUESTION="Update integration 12345 with new configuration"
make test-single QUESTION="Change the name of integration abc-def-123 to 'Updated Monitor'"
```

**Response**: Updated integration object

---

#### 6. **`integrations:delete`** - Remove Integration
**Purpose**: Permanently delete an integration

**Parameters**:
- `integrationId` (required): Unique identifier of the integration

**Example Usage**:
```bash
make test-single QUESTION="Delete integration with ID 12345"
make test-single QUESTION="Remove integration abc-def-123"
```

**Response**: Confirmation of deletion

---

### **Integration State Management**

#### 7. **`integrations:enable`** - Activate Integration
**Purpose**: Enable/activate an integration to start data collection

**Parameters**:
- `integrationId` (required): Unique identifier of the integration

**Example Usage**:
```bash
make test-single QUESTION="Enable integration 12345"
make test-single QUESTION="Activate integration abc-def-123"
```

**Response**: Updated integration object with enabled status

---

#### 8. **`integrations:disable`** - Deactivate Integration
**Purpose**: Disable/deactivate an integration to stop data collection

**Parameters**:
- `integrationId` (required): Unique identifier of the integration

**Example Usage**:
```bash
make test-single QUESTION="Disable integration 12345"
make test-single QUESTION="Deactivate integration abc-def-123"
```

**Response**: Updated integration object with disabled status

---

### **Integration Type Management**

#### 9. **`integrations:listTypes`** - List Available Integration Types
**Purpose**: Retrieve all available integration types that can be created

**Parameters**: None

**Example Usage**:
```bash
make test-single QUESTION="What integration types are available?"
make test-single QUESTION="Show me all supported integration types"
```

**Response**: Array of integration type objects with capabilities

---

#### 10. **`integrations:getType`** - Get Integration Type Details
**Purpose**: Retrieve detailed information about a specific integration type

**Parameters**:
- `typeId` (required): Unique identifier of the integration type

**Example Usage**:
```bash
make test-single QUESTION="Get details for integration type 'aws-cloudwatch'"
make test-single QUESTION="Show me configuration options for Azure Monitor integration"
```

**Response**: Detailed integration type object with configuration schema

---

## 🧪 Testing Integration Management

### **Basic Integration Testing**
```bash
# Quick validation of integration capabilities
make test-integrations-basic-organized

# Test specific integration scenarios
cd client/agent
python shared/engines/enhanced_real_mcp_integration_test.py \
  --test-type integration \
  --complexity basic \
  --max-tests 5
```

### **Comprehensive Integration Testing**
```bash
# Full integration testing suite
make test-integrations-comprehensive-organized

# Advanced integration scenarios
python shared/engines/enhanced_real_mcp_integration_test.py \
  --test-type integration \
  --complexity comprehensive \
  --max-tests 20
```

### **Interactive Integration Testing**
```bash
# True interactive chat mode (recommended)
make chat-interactive

# Test with preset prompts
make run-interactive

# In interactive chat, ask questions like:
# "What integrations do we have?"
# "Create a new monitoring integration"
# "Show me details for integration 12345"
# "Disable all inactive integrations"
# "List all AWS integrations in our environment"
```

## 📊 Integration Data Structure

### **Basic Integration Object**
```json
{
  "id": "integration-12345",
  "name": "Production Monitoring",
  "type": "aws-cloudwatch",
  "status": "enabled",
  "created": "2024-01-15T10:30:00Z",
  "lastModified": "2024-01-20T14:45:00Z"
}
```

### **Detailed Integration Object**
```json
{
  "id": "integration-12345",
  "name": "Production Monitoring",
  "type": "aws-cloudwatch",
  "status": "enabled",
  "description": "Monitors production AWS infrastructure",
  "configuration": {
    "region": "us-east-1",
    "accessKey": "AKIA...",
    "secretKey": "***",
    "namespaces": ["AWS/EC2", "AWS/RDS"]
  },
  "metadata": {
    "created": "2024-01-15T10:30:00Z",
    "lastModified": "2024-01-20T14:45:00Z",
    "createdBy": "user@example.com",
    "version": "1.2.0"
  },
  "statistics": {
    "dataPointsCollected": 15420,
    "lastDataCollection": "2024-01-20T14:40:00Z",
    "errorCount": 0
  }
}
```

## 🔧 Technical Implementation

### **Server-Side Implementation**
- **File**: `pkg/tools/integrations.go` - Core integration management logic
- **File**: `pkg/tools/integrations_api.go` - API implementation with 10 actions
- **File**: `pkg/types/integrations.go` - Integration type definitions

### **API Endpoints**
All integration actions are available through the MCP server at:
- **Base URL**: `http://localhost:8080`
- **Protocol**: JSON-RPC 2.0 over HTTP with Server-Sent Events
- **Authentication**: OpsRamp API credentials (configured in `config.yaml`)

### **Error Handling**
- **Validation Errors**: Invalid parameters or missing required fields
- **Authentication Errors**: Invalid OpsRamp credentials
- **Not Found Errors**: Integration or type not found
- **Permission Errors**: Insufficient permissions for operation

## 🎯 Common Use Cases

### **DevOps Scenarios**
```bash
# Monitor new infrastructure
"Create a monitoring integration for our new Kubernetes cluster"

# Audit existing integrations
"List all enabled integrations and their last data collection time"

# Troubleshoot integration issues
"Show me detailed information for all failed integrations"
```

### **Operations Scenarios**
```bash
# Bulk operations
"Disable all integrations that haven't collected data in 30 days"

# Configuration management
"Update the AWS integration to monitor additional services"

# Capacity planning
"Show me all integration types and their resource requirements"
```

### **Compliance Scenarios**
```bash
# Security auditing
"List all integrations with their creation dates and creators"

# Configuration review
"Show me the configuration for all database monitoring integrations"

# Access control
"Which integrations are currently enabled and collecting data?"
```

## 🚀 Getting Started with Integration Management

1. **Setup**: Follow [GETTING_STARTED.md](GETTING_STARTED.md) for initial configuration
2. **Configure**: Set up OpsRamp credentials in `config.yaml`
3. **Test**: Run basic integration tests to verify connectivity
4. **Explore**: Use chat-interactive mode for direct interaction with integrations
5. **Automate**: Build custom workflows using the integration actions

## 📚 Related Documentation

- **[🚀 GETTING_STARTED.md](GETTING_STARTED.md)** - Complete setup guide
- **[⚙️ CONFIGURATION_GUIDE.md](CONFIGURATION_GUIDE.md)** - Configuration details
- **[🖥️ RESOURCES.md](RESOURCES.md)** - Resource management capabilities
- **[📖 README.md](README.md)** - Project overview

---

**Ready to manage integrations?** Start with the [Quick Start guide](README.md#-quick-start) and explore the chat-interactive mode!