# Resource Management Documentation

Complete guide to HPE OpsRamp resource management capabilities through the MCP server and AI agent testing platform.

## 🖥️ Overview

The Resource Management system provides comprehensive tools for managing HPE OpsRamp resources through AI agents. With **14 comprehensive actions**, you can perform complete lifecycle management, bulk operations, and detailed analytics on infrastructure resources.

## 🛠️ Available Resource Actions

### **Core Resource Operations**

#### 1. **`resources:search`** - Search Resources with Advanced Filtering
**Purpose**: Search and filter resources using various criteria including type, status, tags, and custom filters

**Parameters**:
- `query` (optional): Search query string
- `resourceType` (optional): Filter by resource type
- `status` (optional): Filter by resource status
- `tags` (optional): Filter by resource tags
- `limit` (optional): Maximum number of results
- `offset` (optional): Pagination offset

**Example Usage**:
```bash
make test-single QUESTION="Search for all server resources"
make test-single QUESTION="Find resources with status 'critical'"
make test-single QUESTION="Show me all Linux servers with monitoring enabled"
```

**Response**: Array of resource objects matching search criteria

---

#### 2. **`resources:get`** - Get Basic Resource Information
**Purpose**: Retrieve basic information about a specific resource

**Parameters**:
- `resourceId` (required): Unique identifier of the resource

**Example Usage**:
```bash
make test-single QUESTION="Get details for resource ID 67890"
make test-single QUESTION="Show me basic information for resource server-001"
```

**Response**: Basic resource object with core properties

---

#### 3. **`resources:getDetailed`** - Get Comprehensive Resource Details
**Purpose**: Retrieve comprehensive details about a specific resource including metrics, configuration, and relationships

**Parameters**:
- `resourceId` (required): Unique identifier of the resource

**Example Usage**:
```bash
make test-single QUESTION="Get detailed information for resource 67890"
make test-single QUESTION="Show me comprehensive details for server-001"
```

**Response**: Detailed resource object with full configuration, metrics, and metadata

---

#### 4. **`resources:getMinimal`** - Get Minimal Resource Information
**Purpose**: Retrieve minimal resource information for efficient bulk operations

**Parameters**:
- `resourceId` (required): Unique identifier of the resource

**Example Usage**:
```bash
make test-single QUESTION="Get minimal info for resource 67890"
make test-single QUESTION="Show me basic status for server-001"
```

**Response**: Minimal resource object with essential properties only

---

### **Resource Lifecycle Management**

#### 5. **`resources:create`** - Create New Resource
**Purpose**: Create a new resource with specified configuration

**Parameters**:
- `name` (required): Name of the resource
- `type` (required): Resource type identifier
- `configuration` (required): Resource-specific configuration object
- `description` (optional): Description of the resource
- `tags` (optional): Resource tags

**Example Usage**:
```bash
make test-single QUESTION="Create a new server resource named 'web-server-01'"
make test-single QUESTION="Add a new database resource for MySQL production"
```

**Response**: Created resource object with assigned ID

---

#### 6. **`resources:update`** - Update Resource Configuration
**Purpose**: Update an existing resource's configuration or properties

**Parameters**:
- `resourceId` (required): Unique identifier of the resource
- `configuration` (optional): Updated configuration object
- `name` (optional): Updated name
- `description` (optional): Updated description

**Example Usage**:
```bash
make test-single QUESTION="Update resource 67890 with new configuration"
make test-single QUESTION="Change the description of server-001"
```

**Response**: Updated resource object

---

#### 7. **`resources:delete`** - Remove Resource
**Purpose**: Permanently delete a resource

**Parameters**:
- `resourceId` (required): Unique identifier of the resource

**Example Usage**:
```bash
make test-single QUESTION="Delete resource with ID 67890"
make test-single QUESTION="Remove resource server-001"
```

**Response**: Confirmation of deletion

---

#### 8. **`resources:changeState`** - Change Resource State
**Purpose**: Change the operational state of a resource (e.g., start, stop, restart)

**Parameters**:
- `resourceId` (required): Unique identifier of the resource
- `state` (required): Target state for the resource

**Example Usage**:
```bash
make test-single QUESTION="Start resource 67890"
make test-single QUESTION="Stop server-001"
make test-single QUESTION="Restart database resource db-prod-01"
```

**Response**: Updated resource object with new state

---

### **Bulk Resource Operations**

#### 9. **`resources:bulkUpdate`** - Bulk Update Multiple Resources
**Purpose**: Update multiple resources simultaneously with the same configuration changes

**Parameters**:
- `resourceIds` (required): Array of resource identifiers
- `configuration` (optional): Configuration updates to apply
- `tags` (optional): Tag updates to apply

**Example Usage**:
```bash
make test-single QUESTION="Update all web servers with new monitoring configuration"
make test-single QUESTION="Apply security patches to all Linux servers"
```

**Response**: Array of updated resource objects

---

#### 10. **`resources:bulkDelete`** - Bulk Delete Multiple Resources
**Purpose**: Delete multiple resources simultaneously

**Parameters**:
- `resourceIds` (required): Array of resource identifiers

**Example Usage**:
```bash
make test-single QUESTION="Delete all test environment resources"
make test-single QUESTION="Remove all decommissioned servers"
```

**Response**: Confirmation of bulk deletion with results

---

### **Resource Analytics & Metadata**

#### 11. **`resources:getMetrics`** - Retrieve Resource Metrics
**Purpose**: Get performance metrics and monitoring data for a resource

**Parameters**:
- `resourceId` (required): Unique identifier of the resource
- `metricType` (optional): Specific metric type to retrieve
- `timeRange` (optional): Time range for metrics

**Example Usage**:
```bash
make test-single QUESTION="Get CPU metrics for server-001"
make test-single QUESTION="Show me memory usage for resource 67890"
make test-single QUESTION="Get all performance metrics for database servers"
```

**Response**: Resource metrics object with performance data

---

#### 12. **`resources:getTags`** - Get Resource Tags
**Purpose**: Retrieve all tags associated with a resource

**Parameters**:
- `resourceId` (required): Unique identifier of the resource

**Example Usage**:
```bash
make test-single QUESTION="Show me all tags for resource 67890"
make test-single QUESTION="Get tags for server-001"
```

**Response**: Array of tag objects

---

#### 13. **`resources:updateTags`** - Update Resource Tags
**Purpose**: Add, update, or remove tags from a resource

**Parameters**:
- `resourceId` (required): Unique identifier of the resource
- `tags` (required): Tag updates to apply
- `operation` (optional): Tag operation (add, update, remove)

**Example Usage**:
```bash
make test-single QUESTION="Add environment tag 'production' to server-001"
make test-single QUESTION="Update tags for resource 67890"
make test-single QUESTION="Remove deprecated tags from all web servers"
```

**Response**: Updated resource object with new tags

---

### **Resource Type Management**

#### 14. **`resources:getResourceTypes`** - List Available Resource Types
**Purpose**: Retrieve all available resource types that can be managed

**Parameters**: None

**Example Usage**:
```bash
make test-single QUESTION="What resource types are available?"
make test-single QUESTION="Show me all supported resource types"
```

**Response**: Array of resource type objects with capabilities

---

## 🧪 Testing Resource Management

### **Basic Resource Testing**
```bash
# Quick validation of resource capabilities
make test-resources-basic-organized

# Test specific resource scenarios
cd client/agent
python shared/engines/enhanced_real_mcp_integration_test.py \
  --test-type resource \
  --complexity basic \
  --max-tests 5
```

### **Comprehensive Resource Testing**
```bash
# Full resource testing suite
make test-resources-comprehensive-organized

# Advanced resource scenarios
python shared/engines/enhanced_real_mcp_integration_test.py \
  --test-type resource \
  --complexity comprehensive \
  --max-tests 20
```

### **Ultra-Complex Resource Testing**
```bash
# Ultra-complex resource scenarios
make test-resources-ultra-organized

# Maximum complexity testing
python shared/engines/enhanced_real_mcp_integration_test.py \
  --test-type resource \
  --complexity ultra \
  --max-tests 30
```

### **Interactive Resource Testing**
```bash
# True interactive chat mode (recommended)
make chat-interactive

# Test with preset prompts
make run-interactive

# In interactive chat, ask questions like:
# "What resources do we have?"
# "Show me all critical servers"
# "Update tags for all production databases"
# "Get metrics for high-CPU resources"
# "Generate a report of all resources"
```

## 📊 Resource Data Structure

### **Basic Resource Object**
```json
{
  "id": "resource-67890",
  "name": "web-server-01",
  "type": "server",
  "status": "running",
  "state": "active",
  "created": "2024-01-15T10:30:00Z",
  "lastModified": "2024-01-20T14:45:00Z"
}
```

### **Detailed Resource Object**
```json
{
  "id": "resource-67890",
  "name": "web-server-01",
  "type": "server",
  "status": "running",
  "state": "active",
  "description": "Production web server",
  "configuration": {
    "os": "Ubuntu 20.04",
    "cpu": "4 cores",
    "memory": "16GB",
    "disk": "500GB SSD",
    "network": "1Gbps"
  },
  "location": {
    "datacenter": "us-east-1",
    "rack": "A-15",
    "position": "U12"
  },
  "tags": [
    {"key": "environment", "value": "production"},
    {"key": "application", "value": "web"},
    {"key": "owner", "value": "devops-team"}
  ],
  "metadata": {
    "created": "2024-01-15T10:30:00Z",
    "lastModified": "2024-01-20T14:45:00Z",
    "createdBy": "admin@example.com",
    "version": "2.1.0"
  },
  "metrics": {
    "cpu": {"current": 45.2, "average": 38.7},
    "memory": {"current": 78.5, "average": 72.1},
    "disk": {"current": 65.3, "average": 60.8}
  },
  "relationships": {
    "parent": "cluster-prod-01",
    "children": ["service-web", "service-api"],
    "dependencies": ["database-01", "cache-01"]
  }
}
```

### **Resource Metrics Object**
```json
{
  "resourceId": "resource-67890",
  "timestamp": "2024-01-20T14:45:00Z",
  "metrics": {
    "cpu": {
      "utilization": 45.2,
      "load1": 1.2,
      "load5": 1.5,
      "load15": 1.8
    },
    "memory": {
      "used": 12.5,
      "free": 3.5,
      "cached": 2.1,
      "utilization": 78.5
    },
    "disk": {
      "used": 326.5,
      "free": 173.5,
      "utilization": 65.3,
      "iops": 1250
    },
    "network": {
      "bytesIn": 1048576,
      "bytesOut": 2097152,
      "packetsIn": 1024,
      "packetsOut": 2048
    }
  }
}
```

## 🔧 Technical Implementation

### **Server-Side Implementation**
- **File**: `pkg/tools/resources.go` - Core resource management logic
- **File**: `pkg/tools/resources_api.go` - API implementation with 14 actions
- **File**: `pkg/types/resources.go` - Resource type definitions

### **API Endpoints**
All resource actions are available through the MCP server at:
- **Base URL**: `http://localhost:8080`
- **Protocol**: JSON-RPC 2.0 over HTTP with Server-Sent Events
- **Authentication**: OpsRamp API credentials (configured in `config.yaml`)

### **Error Handling**
- **Validation Errors**: Invalid parameters or missing required fields
- **Authentication Errors**: Invalid OpsRamp credentials
- **Not Found Errors**: Resource or type not found
- **Permission Errors**: Insufficient permissions for operation
- **State Errors**: Invalid state transitions

## 🎯 Common Use Cases

### **Infrastructure Management**
```bash
# Inventory management
"Show me all servers in the production environment"

# Capacity planning
"Get CPU and memory metrics for all database servers"

# Resource optimization
"Find all resources with high CPU utilization"
```

### **Operations & Monitoring**
```bash
# Health monitoring
"Show me all resources with critical status"

# Performance analysis
"Get detailed metrics for all web servers"

# Troubleshooting
"Find all resources that were modified in the last 24 hours"
```

### **Automation & Bulk Operations**
```bash
# Bulk updates
"Update all test environment resources with new tags"

# State management
"Restart all web servers in the staging environment"

# Cleanup operations
"Delete all resources tagged as 'decommissioned'"
```

### **Compliance & Auditing**
```bash
# Security auditing
"List all resources with their creation dates and owners"

# Configuration compliance
"Show me all servers without proper monitoring tags"

# Change tracking
"Get modification history for all production resources"
```

## 🚀 Getting Started with Resource Management

1. **Setup**: Follow [GETTING_STARTED.md](GETTING_STARTED.md) for initial configuration
2. **Configure**: Set up OpsRamp credentials in `config.yaml`
3. **Test**: Run basic resource tests to verify connectivity
4. **Explore**: Use chat-interactive mode for direct interaction with resources
5. **Automate**: Build custom workflows using the resource actions

## 📚 Related Documentation

- **[🚀 GETTING_STARTED.md](GETTING_STARTED.md)** - Complete setup guide
- **[⚙️ CONFIGURATION_GUIDE.md](CONFIGURATION_GUIDE.md)** - Configuration details
- **[🔗 INTEGRATIONS.md](INTEGRATIONS.md)** - Integration management capabilities
- **[📖 README.md](README.md)** - Project overview

---

**Ready to manage resources?** Start with the [Quick Start guide](README.md#-quick-start) and explore the comprehensive testing framework!