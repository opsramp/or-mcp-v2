# Resource Management Tool Design Document

## 📋 **Executive Summary**

This document outlines the comprehensive design and implementation strategy for adding **Resource Management** as the next tool to our HPE OpsRamp MCP server. Resource Management is selected as the optimal next tool because it provides foundational device and infrastructure management capabilities that naturally complement our existing integrations tool.

**Selected API:** [OpsRamp Resource Management API](https://develop.opsramp.com/v2) - Resource Management section

## 🎯 **Tool Selection Rationale**

### **Why Resource Management?**

After thorough analysis of all available OpsRamp API categories, **Resource Management** emerges as the optimal choice for the following strategic reasons:

1. **🏗️ Foundational Synergy**: Resources are the foundation upon which integrations operate - you manage resources first, then install integrations to monitor them
2. **📊 Comprehensive Scope**: 50+ endpoints covering complete resource lifecycle management
3. **🔄 Natural Workflow**: Creates logical progression: Manage Resources → Install Integrations → Monitor Performance
4. **🛠️ Rich Functionality**: Device groups, sites, service groups, availability monitoring, warranty management
5. **👥 Enterprise Value**: Critical for enterprise infrastructure management and organization

### **API Categories Considered**

| Category | Endpoints | Decision | Reasoning |
|----------|-----------|----------|-----------|
| **Resource Management** | 50+ | ✅ **SELECTED** | Foundational, comprehensive, high business value |
| Alerts | 40+ | ⏳ Future | Important but requires resources foundation |
| Automation | 25+ | ⏳ Future | Powerful but complex, better after core tools |
| Monitoring | 20+ | ⏳ Future | Builds on resources + integrations |
| Tickets | 30+ | ⏳ Future | Service management layer |

## 📊 **OpsRamp Resource Management API Analysis**

### **Core API Endpoints (50+ endpoints analyzed)**

Based on the [OpsRamp Resource Management documentation](https://develop.opsramp.com/v2), the Resource Management API provides comprehensive coverage:

#### **🔧 Core Resource Operations**
- `POST /api/v2/tenants/{clientId}/resources` - Create Resource
- `GET /api/v2/tenants/{clientId}/resources/{resourceId}` - Get Resource Details  
- `POST /api/v2/tenants/{clientId}/resources/search` - Search Resources
- `PATCH /api/v2/tenants/{clientId}/resources/{resourceId}` - Update Resource
- `DELETE /api/v2/tenants/{clientId}/resources/{resourceId}` - Delete Resource
- `GET /api/v2/tenants/{clientId}/resources/minimal` - Get Minimal Resource Details

#### **📱 Device Groups Management**
- `POST /api/v2/tenants/{clientId}/deviceGroups` - Create Device Group
- `GET /api/v2/tenants/{clientId}/deviceGroups/{deviceGroupId}` - Get Device Group
- `PATCH /api/v2/tenants/{clientId}/deviceGroups/{deviceGroupId}` - Update Device Group
- `GET /api/v2/tenants/{clientId}/deviceGroups/minimal` - Get Minimal Device Groups
- `GET /api/v2/tenants/{clientId}/deviceGroups/root` - Get Root Level Device Groups

#### **🏢 Sites Management**
- `POST /api/v2/tenants/{clientId}/sites` - Create Site
- `GET /api/v2/tenants/{clientId}/sites/{siteId}` - Get Site Details
- `PATCH /api/v2/tenants/{clientId}/sites/{siteId}` - Update Site
- `POST /api/v2/tenants/{clientId}/sites/search` - Search Sites

#### **⚙️ Service Groups Management**
- `POST /api/v2/tenants/{clientId}/serviceGroups` - Create Service Group
- `GET /api/v2/tenants/{clientId}/serviceGroups/{serviceGroupId}` - Get Service Group
- `PATCH /api/v2/tenants/{clientId}/serviceGroups/{serviceGroupId}` - Update Service Group

#### **📊 Monitoring & Availability**
- `GET /api/v2/tenants/{clientId}/resources/{resourceId}/availability` - Get Resource Availability
- `GET /api/v2/tenants/{clientId}/resources/{resourceId}/applications` - Get Resource Applications
- `GET /api/v2/tenants/{clientId}/resources/{resourceId}/services` - Get Last Discovered Services

#### **🔧 Advanced Operations**
- `POST /api/v2/tenants/{clientId}/resources/{resourceId}/actions` - Resource Management Actions
- `POST /api/v2/tenants/{clientId}/resources/{resourceId}/decommission` - Decommission Resource
- `GET /api/v2/tenants/{clientId}/resources/{resourceId}/warranty` - Get Device Warranty

## 🏗️ **Implementation Architecture**

### **Tool Structure Design**

Following our proven modular architecture from the integrations tool:

```
pkg/tools/
├── resources.go              # Main MCP tool definition
├── resources_api.go          # API client implementation
└── resources_handlers.go     # Action handlers

pkg/types/
├── resources.go              # Resource type definitions
└── resource_groups.go        # Group type definitions

internal/adapters/
└── resources.go              # OpsRamp API adapter
```

### **Tool Actions Design**

Based on comprehensive API analysis, we'll implement **12 core actions**:

#### **🔧 Core Resource Actions**
1. **`list`** - Search and list resources with filtering
2. **`get`** - Get detailed resource information
3. **`create`** - Create new resources
4. **`update`** - Update resource properties
5. **`delete`** - Remove resources
6. **`getMinimal`** - Get minimal resource details for performance

#### **📱 Group Management Actions**  
7. **`listDeviceGroups`** - Manage device groups
8. **`listSites`** - Manage sites and locations
9. **`listServiceGroups`** - Manage service groups

#### **📊 Monitoring Actions**
10. **`getAvailability`** - Get resource availability data
11. **`getApplications`** - Get installed applications
12. **`performAction`** - Execute resource management actions

## 📋 **Phased Implementation Plan**

### **Phase 1: Foundation & Core Operations (Week 1)**

#### **🎯 Objectives**
- Establish basic resource management capabilities
- Implement core CRUD operations
- Create solid foundation for advanced features

#### **📝 Deliverables**
1. **Type Definitions** (`pkg/types/resources.go`)
   - Resource, DeviceGroup, Site, ServiceGroup structs
   - Search parameters and response types
   - Validation and serialization logic

2. **API Client** (`pkg/tools/resources_api.go`)
   - OpsRamp API integration
   - Authentication and error handling
   - HTTP client with retry logic

3. **Core Actions Implementation**
   - `list` - Resource search with comprehensive filtering
   - `get` - Detailed resource retrieval
   - `getMinimal` - Performance-optimized minimal details

4. **Basic Testing**
   - Unit tests for type definitions
   - API client integration tests
   - Core action validation

#### **🧪 Testing Strategy**
```bash
# Unit tests
go test ./pkg/types/...
go test ./pkg/tools/...

# Integration tests with mock server
go test -tags=integration ./internal/adapters/...

# Manual testing with real OpsRamp instance
make test-resources-basic
```

#### **✅ Success Criteria**
- [ ] All type definitions compile and validate
- [ ] API client successfully authenticates with OpsRamp
- [ ] Core actions (list, get, getMinimal) return valid data
- [ ] 100% test coverage for implemented components
- [ ] Zero security vulnerabilities (gosec scan)

### **Phase 2: Advanced Operations & Group Management (Week 2)**

#### **🎯 Objectives**
- Implement full CRUD operations
- Add device groups, sites, and service groups management
- Establish organizational hierarchy support

#### **📝 Deliverables**
1. **CRUD Operations**
   - `create` - Resource creation with validation
   - `update` - Resource modification with conflict handling  
   - `delete` - Safe resource removal with dependency checking

2. **Group Management**
   - `listDeviceGroups` - Device group hierarchy management
   - `listSites` - Site and location management
   - `listServiceGroups` - Service group organization

3. **Advanced Type Support**
   - Hierarchical group structures
   - Parent-child relationships
   - Group membership management

4. **Error Handling Enhancement**
   - Comprehensive error types
   - Retry logic for transient failures
   - Validation and conflict resolution

#### **🧪 Testing Strategy**
```bash
# Group management tests
make test-resources-groups

# CRUD operation tests  
make test-resources-crud

# Error handling validation
make test-resources-errors
```

#### **✅ Success Criteria**
- [ ] All CRUD operations work reliably
- [ ] Group hierarchies are properly managed
- [ ] Parent-child relationships are maintained
- [ ] Error handling covers all edge cases
- [ ] Performance benchmarks meet requirements

### **Phase 3: Monitoring & Advanced Features (Week 3)**

#### **🎯 Objectives**
- Add monitoring and availability features
- Implement resource actions and management operations
- Create comprehensive tool capabilities

#### **📝 Deliverables**
1. **Monitoring Features**
   - `getAvailability` - Resource availability monitoring
   - `getApplications` - Installed application tracking
   - Performance metrics and health indicators

2. **Management Actions**
   - `performAction` - Execute resource management operations
   - Bulk operations support
   - Action scheduling and tracking

3. **Advanced Filtering**
   - Complex search criteria
   - Multi-field filtering
   - Sorting and pagination

4. **Performance Optimization**
   - Caching strategies
   - Bulk operation support
   - Parallel processing where appropriate

#### **🧪 Testing Strategy**
```bash
# Monitoring tests
make test-resources-monitoring

# Advanced feature tests
make test-resources-advanced

# Performance benchmarks
make benchmark-resources
```

#### **✅ Success Criteria**
- [ ] All monitoring features return accurate data
- [ ] Resource actions execute successfully
- [ ] Advanced filtering works correctly
- [ ] Performance meets enterprise requirements
- [ ] All features are well-documented

### **Phase 4: Integration & Comprehensive Testing (Week 4)**

#### **🎯 Objectives**
- Integrate with existing MCP server
- Create comprehensive testing suite
- Establish AI agent testing scenarios

#### **📝 Deliverables**
1. **MCP Server Integration**
   - Tool registration and handler setup
   - JSON-RPC endpoint configuration
   - Session management integration

2. **AI Agent Testing Platform**
   - 50+ resource management test scenarios
   - Integration with existing testing framework
   - Performance and reliability validation

3. **Documentation & Examples**
   - Complete API documentation
   - Usage examples and best practices
   - Troubleshooting guide

4. **Security & Quality Assurance**
   - Security vulnerability scanning
   - Code quality validation
   - Performance benchmarking

#### **🧪 Testing Strategy**
```bash
# Full integration testing
make test-resources-integration

# AI agent testing scenarios
cd client/agent
make test-custom PROMPTS_FILE=test_data/resource_management_prompts.txt

# Security and quality validation
make security-scan
make quality-check
```

#### **✅ Success Criteria**
- [ ] Tool integrates seamlessly with MCP server
- [ ] AI agent testing achieves 100% success rate
- [ ] All documentation is complete and accurate
- [ ] Security scan shows zero vulnerabilities
- [ ] Performance benchmarks meet requirements

## 🧪 **Testing Strategy**

### **Comprehensive Testing Framework**

Following our proven testing methodology from the integrations tool:

#### **🔧 Unit Testing (Automated)**
```bash
# Type validation tests
go test ./pkg/types/resources_test.go

# API client tests  
go test ./pkg/tools/resources_api_test.go

# Handler logic tests
go test ./pkg/tools/resources_handlers_test.go
```

#### **🌐 Integration Testing (Live API)**
```bash
# Real OpsRamp API testing
make test-resources-live

# Error condition testing
make test-resources-errors

# Performance benchmarking
make benchmark-resources
```

#### **🤖 AI Agent Testing (50+ Scenarios)**

Create comprehensive test scenarios similar to our integrations testing:

**Resource Discovery & Listing (12 scenarios)**
- "List all servers in our infrastructure"
- "Show me Windows servers that need updates"
- "Find all resources in the production site"

**Resource Management (10 scenarios)**  
- "Create a new Linux server resource"
- "Update the memory configuration for server-001"
- "Decommission old development servers"

**Group Organization (8 scenarios)**
- "Show me the device group hierarchy"
- "Create a new site for the London office"
- "List all service groups and their members"

**Monitoring & Health (10 scenarios)**
- "Check availability for critical servers"
- "Show me applications installed on web servers"
- "Find resources with availability issues"

**Advanced Operations (10 scenarios)**
- "Perform maintenance restart on server group"
- "Show warranty status for all hardware"
- "Find resources that need patching"

### **Performance Requirements**

#### **Response Time Targets**
- `list` operations: < 2 seconds for 1000+ resources
- `get` operations: < 500ms for detailed resource info
- `create/update` operations: < 1 second
- Bulk operations: < 5 seconds for 100 resources

#### **Reliability Targets**
- 99.9% success rate for all operations
- Graceful handling of OpsRamp API rate limits
- Automatic retry for transient failures
- Circuit breaker for persistent failures

## 🔧 **Technical Implementation Details**

### **Configuration Integration**

Extend existing `config.yaml` structure:

```yaml
opsramp:
  tenant_url: "https://your-instance.opsramp.com"
  auth_url: "https://your-instance.opsramp.com/tenancy/auth/oauth/token"
  auth_key: "YOUR_API_KEY"
  auth_secret: "YOUR_API_SECRET"
  tenant_id: "YOUR_TENANT_ID"
  
  # Resource management specific settings
  resources:
    default_page_size: 50
    max_page_size: 1000
    cache_ttl: 300  # 5 minutes
    enable_bulk_operations: true
```

### **Error Handling Strategy**

Implement comprehensive error handling following enterprise patterns:

```go
type ResourceError struct {
    Code    string    `json:"code"`
    Message string    `json:"message"`
    Type    ErrorType `json:"type"`
    Details map[string]interface{} `json:"details,omitempty"`
}

type ErrorType string

const (
    ErrorTypeValidation   ErrorType = "validation"
    ErrorTypeNotFound     ErrorType = "not_found"
    ErrorTypePermission   ErrorType = "permission"
    ErrorTypeRateLimit    ErrorType = "rate_limit"
    ErrorTypeServerError  ErrorType = "server_error"
)
```

### **Caching Strategy**

Implement intelligent caching for performance:

```go
type ResourceCache struct {
    resources    map[string]*CachedResource
    deviceGroups map[string]*CachedDeviceGroup
    sites        map[string]*CachedSite
    ttl          time.Duration
    mutex        sync.RWMutex
}
```

## 🚀 **Integration with Existing System**

### **MCP Server Integration**

Add to `cmd/server/main.go`:

```go
// Register resources tool
resourcesTool, resourcesHandler := tools.NewResourcesMcpTool()
mcpServer.AddTool(resourcesTool, resourcesHandler)
registeredTools = append(registeredTools, resourcesTool.Name)
```

### **AI Agent Enhancement**

Extend AI agent capabilities to handle resource management queries naturally:

- Resource discovery and inventory management
- Infrastructure organization and hierarchy  
- Capacity planning and resource allocation
- Health monitoring and maintenance coordination

### **Cross-Tool Synergy**

Enable powerful workflows combining resources and integrations:

1. **Discovery Workflow**: "Find all Windows servers and show what integrations are installed"
2. **Planning Workflow**: "List resources in site X and recommend monitoring integrations"  
3. **Maintenance Workflow**: "Check health of resources with VMware integrations"

## 📈 **Success Metrics & KPIs**

### **Functional Metrics**
- ✅ **API Coverage**: 100% of planned endpoints implemented
- ✅ **Test Coverage**: 95%+ code coverage with comprehensive scenarios
- ✅ **AI Agent Success**: 100% success rate on resource management scenarios
- ✅ **Performance**: All response time targets met
- ✅ **Reliability**: 99.9%+ operation success rate

### **Quality Metrics**
- ✅ **Security**: Zero vulnerabilities in security scans
- ✅ **Documentation**: Complete API documentation and examples
- ✅ **Code Quality**: All quality gates passed
- ✅ **Integration**: Seamless integration with existing tools

### **Business Value Metrics**
- ✅ **Workflow Enhancement**: Resource + Integration workflows enabled
- ✅ **Enterprise Readiness**: Production-ready infrastructure management
- ✅ **User Experience**: Natural language resource management queries
- ✅ **System Completeness**: Foundation for additional tools laid

## 🔐 **Security Considerations**

### **Data Protection**
- Sensitive resource information handling
- Credential management for OpsRamp APIs
- Access control and permission validation
- Audit logging for resource management operations

### **API Security**
- Rate limiting and circuit breaker patterns
- Request/response validation and sanitization
- Secure storage of authentication tokens
- Protection against injection attacks

### **Enterprise Security**
- Integration with existing security framework
- Compliance with enterprise security policies
- Secure configuration management
- Encrypted communication channels

## 📚 **Documentation Plan**

### **Technical Documentation**
1. **API Reference**: Complete endpoint documentation
2. **Integration Guide**: Step-by-step integration instructions
3. **Best Practices**: Performance and security guidelines
4. **Troubleshooting**: Common issues and solutions

### **User Documentation**
1. **Getting Started**: Quick start guide for resource management
2. **Use Cases**: Real-world scenarios and examples
3. **AI Agent Guide**: Natural language query examples
4. **Advanced Features**: Complex workflow documentation

## 🎯 **Conclusion**

The Resource Management tool represents a strategic addition that will:

1. **🏗️ Establish Foundation**: Provide essential infrastructure management capabilities
2. **🔄 Enable Workflows**: Create natural progressions from resource management to monitoring
3. **📊 Deliver Value**: Offer immediate business value for enterprise infrastructure management
4. **🚀 Scale Platform**: Lay groundwork for additional monitoring and automation tools

This phased approach ensures systematic, tested, and secure implementation while maintaining our zero-tolerance security standards and 100% success rate achievements.

**Ready to proceed with Phase 1 implementation upon approval.** 🚀 