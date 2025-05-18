package types

// Resource represents an OpsRamp resource
type Resource struct {
	ID                        string                 `json:"id"`
	HostName                  string                 `json:"hostName"`
	IPAddress                 string                 `json:"ipAddress"`
	Identity                  string                 `json:"identity"`
	CreatedDate               string                 `json:"createdDate"`
	UpdatedDate               string                 `json:"updatedDate"`
	Type                      string                 `json:"type"`
	NativeType                string                 `json:"nativeType"`
	State                     string                 `json:"state"`
	Source                    string                 `json:"source"`
	Status                    string                 `json:"status"`
	AliasName                 string                 `json:"aliasName"`
	Tags                      []Tag                  `json:"tags"`
	Name                      string                 `json:"name"`
	ResourceName              string                 `json:"resourceName"`
	Consoles                  []string               `json:"consoles"`
	Properties                map[string]interface{} `json:"properties"`
	ClientUniqueID            string                 `json:"clientUniqueId"`
	ResourceType              string                 `json:"resourceType"`
	AgentInstalled            bool                   `json:"agentInstalled"`
	AgentStatus               string                 `json:"agentStatus"`
	AgentLastConnectedTime    string                 `json:"agentLastConnectedTime,omitempty"`
	Location                  *Location              `json:"location,omitempty"`
	ManagementProfile         *ManagementProfile     `json:"managementProfile,omitempty"`
	DNSName                   string                 `json:"dnsName,omitempty"`
	SerialNumber              string                 `json:"serialNumber,omitempty"`
	Make                      string                 `json:"make,omitempty"`
	Model                     string                 `json:"model,omitempty"`
	SystemUID                 string                 `json:"systemUID,omitempty"`
	ProviderUID               string                 `json:"providerUID,omitempty"`
	ProviderType              string                 `json:"providerType,omitempty"`
	Description               string                 `json:"description,omitempty"`
	OS                        string                 `json:"os,omitempty"`
	Category                  string                 `json:"category,omitempty"`
	InstalledBy               string                 `json:"installedBy,omitempty"`
	InstalledTime             string                 `json:"installedTime,omitempty"`
	ModifiedTime              string                 `json:"modifiedTime,omitempty"`
	ModifiedBy                string                 `json:"modifiedBy,omitempty"`
	AccountLastDiscoveredTime string                 `json:"accountLastDiscoveredTime,omitempty"`
}

// Tag represents a resource tag
type Tag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Location represents a resource location
type Location struct {
	ID   interface{} `json:"id"`
	Name string      `json:"name"`
}

// ManagementProfile represents a resource management profile
type ManagementProfile struct {
	ID   interface{} `json:"id"`
	Name string      `json:"name"`
}

// ResourceSearchParams represents the parameters for searching resources
type ResourceSearchParams struct {
	PageNo                int    `json:"pageNo,omitempty"`
	PageSize              int    `json:"pageSize,omitempty"`
	IsDescendingOrder     bool   `json:"isDescendingOrder,omitempty"`
	SortName              string `json:"sortName,omitempty"`
	QueryString           string `json:"queryString,omitempty"`
	IncludeStatus         bool   `json:"includeStatus,omitempty"`
	HostName              string `json:"hostName,omitempty"`
	DNSName               string `json:"dnsName,omitempty"`
	ResourceName          string `json:"resourceName,omitempty"`
	AliasName             string `json:"aliasName,omitempty"`
	ID                    string `json:"id,omitempty"`
	SerialNumber          string `json:"serialNumber,omitempty"`
	IPAddress             string `json:"ipAddress,omitempty"`
	SystemUID             string `json:"systemUID,omitempty"`
	State                 string `json:"state,omitempty"`
	Type                  string `json:"type,omitempty"`
	DeviceType            string `json:"deviceType,omitempty"`
	ResourceType          string `json:"resourceType,omitempty"`
	StartCreationDate     string `json:"startCreationDate,omitempty"`
	EndCreationDate       string `json:"endCreationDate,omitempty"`
	StartUpdationDate     string `json:"startUpdationDate,omitempty"`
	EndUpdationDate       string `json:"endUpdationDate,omitempty"`
	Tags                  string `json:"tags,omitempty"`
	Template              string `json:"template,omitempty"`
	AgentProfile          string `json:"agentProfile,omitempty"`
	GatewayProfile        string `json:"gatewayProfile,omitempty"`
	InstanceID            string `json:"instanceId,omitempty"`
	AccountNumber         string `json:"accountNumber,omitempty"`
	AgentInstalled        *bool  `json:"agentInstalled,omitempty"`
	DeviceGroup           string `json:"deviceGroup,omitempty"`
	ServiceGroup          string `json:"serviceGroup,omitempty"`
	DeviceLocation        string `json:"deviceLocation,omitempty"`
	IsEquals              string `json:"isEquals,omitempty"`
	AliasIP               string `json:"aliasIp,omitempty"`
	AppRoles              string `json:"appRoles,omitempty"`
	OSArchitecture        string `json:"osArchitecture,omitempty"`
	AssetManagedTime      string `json:"assetManagedTime,omitempty"`
	FirstAssetManagedTime string `json:"firstAssetManagedTime,omitempty"`
	Category              string `json:"category,omitempty"`
	Make                  string `json:"make,omitempty"`
	Model                 string `json:"model,omitempty"`
	ProviderType          string `json:"providerType,omitempty"`
	ProviderUID           string `json:"providerUID,omitempty"`
}

// ResourceSearchResponse represents the response from a resource search
type ResourceSearchResponse struct {
	Results         []Resource `json:"results"`
	TotalResults    int        `json:"totalResults"`
	OrderBy         string     `json:"orderBy"`
	PageNo          int        `json:"pageNo"`
	PageSize        int        `json:"pageSize"`
	TotalPages      int        `json:"totalPages"`
	NextPage        bool       `json:"nextPage"`
	DescendingOrder bool       `json:"descendingOrder"`
}

// DetailedResource represents a detailed view of an OpsRamp resource
type DetailedResource struct {
	Resource
	Components            []string               `json:"components,omitempty"`
	BIOS                  map[string]interface{} `json:"bios,omitempty"`
	CPUs                  []CPU                  `json:"cpus,omitempty"`
	GeneralInfo           map[string]interface{} `json:"generalInfo,omitempty"`
	InstalledApp          map[string]interface{} `json:"installedApp,omitempty"`
	MetricTypes           []MetricType           `json:"metricTypes,omitempty"`
	NetworkCardDetails    []NetworkCard          `json:"networkCardDetails,omitempty"`
	DiscoveryProfile      map[string]interface{} `json:"discoveryProfile,omitempty"`
	AppRoles              []string               `json:"appRoles,omitempty"`
	LogicalDiskDrives     []LogicalDiskDrive     `json:"logicalDiskDrives,omitempty"`
	AvailabilityStatus    string                 `json:"availabilityStatus,omitempty"`
	UpDownSince           string                 `json:"upDownSince,omitempty"`
	LastMetricValue       int                    `json:"lastMetricValue,omitempty"`
	LastMetricUpdatedTime string                 `json:"lastMetricUpdatedTime,omitempty"`
	MetricUnit            string                 `json:"metricUnit,omitempty"`
	DefaultMetric         string                 `json:"defaultMetric,omitempty"`
	Applications          []Application          `json:"applications,omitempty"`
	DiscoveredServices    []DiscoveredService    `json:"discoveredServices,omitempty"`
	Warranty              *Warranty              `json:"warranty,omitempty"`
}

// CPU represents a CPU in a resource
type CPU struct {
	Name         string `json:"name"`
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	Speed        string `json:"speed"`
	Cores        int    `json:"cores"`
	LogicalCores int    `json:"logicalCores"`
	Architecture string `json:"architecture"`
}

// NetworkCard represents a network card in a resource
type NetworkCard struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	MACAddress  string `json:"macAddress"`
	IPAddress   string `json:"ipAddress"`
	SubnetMask  string `json:"subnetMask"`
	Gateway     string `json:"gateway"`
	DHCPEnabled bool   `json:"dhcpEnabled"`
}

// MetricType represents a metric type for a resource
type MetricType struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Unit        string `json:"unit"`
}

// Application represents an application installed on a resource
type Application struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Vendor      string `json:"vendor"`
	InstallDate string `json:"installDate"`
	Path        string `json:"path"`
	Size        int64  `json:"size"`
}

// DiscoveredService represents a service discovered on a resource
type DiscoveredService struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Status      string `json:"status"`
	StartType   string `json:"startType"`
	Path        string `json:"path"`
	Port        int    `json:"port"`
	Protocol    string `json:"protocol"`
}

// Warranty represents warranty information for a resource
type Warranty struct {
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	Status      string `json:"status"`
	Provider    string `json:"provider"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

// ResourceCreateRequest represents the request to create a resource
type ResourceCreateRequest struct {
	AliasName                string             `json:"aliasName,omitempty"`
	AlternateIP              string             `json:"alternateIP,omitempty"`
	ExtResourceID            string             `json:"extResourceId,omitempty"`
	ManagementProfile        string             `json:"managementProfile,omitempty"`
	ResourceNetworkInterface []NetworkInterface `json:"resourceNetworkInterface,omitempty"`
	LogicalDiskDrives        []LogicalDiskDrive `json:"logicalDiskDrives,omitempty"`
	OOBInterfaceCards        []OOBInterfaceCard `json:"oobInterfaceCards,omitempty"`
	ResourceType             string             `json:"resourceType"`
	HostName                 string             `json:"hostName,omitempty"`
	IPAddress                string             `json:"ipAddress,omitempty"`
	DNSName                  string             `json:"dnsName,omitempty"`
	SerialNumber             string             `json:"serialNumber,omitempty"`
	Make                     string             `json:"make,omitempty"`
	Model                    string             `json:"model,omitempty"`
	Description              string             `json:"description,omitempty"`
	OS                       string             `json:"os,omitempty"`
	Category                 string             `json:"category,omitempty"`
	Location                 string             `json:"location,omitempty"`
	Tags                     []Tag              `json:"tags,omitempty"`
	Properties               map[string]any     `json:"properties,omitempty"`
}

// NetworkInterface represents a network interface for a resource
type NetworkInterface struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	MACAddress  string `json:"macAddress"`
	IPAddress   string `json:"ipAddress"`
	SubnetMask  string `json:"subnetMask"`
	Gateway     string `json:"gateway"`
	DHCPEnabled bool   `json:"dhcpEnabled"`
}

// LogicalDiskDrive represents a logical disk drive for a resource
type LogicalDiskDrive struct {
	Name       string `json:"name"`
	FileSystem string `json:"fileSystem"`
	Size       int64  `json:"size"`
	FreeSpace  int64  `json:"freeSpace"`
}

// OOBInterfaceCard represents an out-of-band interface card for a resource
type OOBInterfaceCard struct {
	Type      string `json:"type"`
	IPAddress string `json:"ipAddress"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

// ResourceUpdateRequest represents the request to update a resource
type ResourceUpdateRequest struct {
	AliasName                string             `json:"aliasName,omitempty"`
	AlternateIP              string             `json:"alternateIP,omitempty"`
	ExtResourceID            string             `json:"extResourceId,omitempty"`
	ManagementProfile        string             `json:"managementProfile,omitempty"`
	ResourceNetworkInterface []NetworkInterface `json:"resourceNetworkInterface,omitempty"`
	LogicalDiskDrives        []LogicalDiskDrive `json:"logicalDiskDrives,omitempty"`
	OOBInterfaceCards        []OOBInterfaceCard `json:"oobInterfaceCards,omitempty"`
	ResourceType             string             `json:"resourceType"`
	HostName                 string             `json:"hostName,omitempty"`
	IPAddress                string             `json:"ipAddress,omitempty"`
	DNSName                  string             `json:"dnsName,omitempty"`
	SerialNumber             string             `json:"serialNumber,omitempty"`
	Make                     string             `json:"make,omitempty"`
	Model                    string             `json:"model,omitempty"`
	Description              string             `json:"description,omitempty"`
	OS                       string             `json:"os,omitempty"`
	Category                 string             `json:"category,omitempty"`
	Location                 string             `json:"location,omitempty"`
	Tags                     []Tag              `json:"tags,omitempty"`
	Properties               map[string]any     `json:"properties,omitempty"`
	State                    string             `json:"state,omitempty"`
}

// ResourceBulkUpdateRequest represents a request to update multiple resources
type ResourceBulkUpdateRequest struct {
	ResourceIDs []string               `json:"resourceIds"`
	Updates     map[string]interface{} `json:"updates"`
}

// ResourceBulkDeleteRequest represents a request to delete multiple resources
type ResourceBulkDeleteRequest struct {
	ResourceIDs []string `json:"resourceIds"`
}

// ResourceTypeInfo represents information about a resource type
type ResourceTypeInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

// ResourceStateChangeRequest represents a request to change a resource's state
type ResourceStateChangeRequest struct {
	State string `json:"state"`
}

// ResourceMetricsRequest represents a request to get metrics for a resource
type ResourceMetricsRequest struct {
	MetricNames []string `json:"metricNames"`
	StartTime   string   `json:"startTime"`
	EndTime     string   `json:"endTime"`
	Interval    string   `json:"interval,omitempty"`
}

// ResourceMetricsResponse represents a response containing resource metrics
type ResourceMetricsResponse struct {
	ResourceID string                    `json:"resourceId"`
	Metrics    []ResourceMetricDataPoint `json:"metrics"`
}

// ResourceMetricDataPoint represents a metric data point for a resource
type ResourceMetricDataPoint struct {
	Name      string  `json:"name"`
	Timestamp string  `json:"timestamp"`
	Value     float64 `json:"value"`
	Unit      string  `json:"unit"`
}
