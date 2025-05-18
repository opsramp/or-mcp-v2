package types

// FilterCriterion represents a single filter condition for integration config
 type FilterCriterion struct {
	Attribute string `json:"attribute"`
	Operator  string `json:"operator"`
	Value     string `json:"value"`
}

// IntegrationConfig is a generic config for all integration types
 type IntegrationConfig struct {
	Host              string                 `json:"host,omitempty"`
	Username          string                 `json:"username,omitempty"`
	Password          string                 `json:"password,omitempty"`
	CollectorProfile  string                 `json:"collectorProfile,omitempty"`
	FilterCriteria    []FilterCriterion      `json:"filterCriteria,omitempty"`
	DiscoverySchedule string                 `json:"discoverySchedule,omitempty"`
	AttributeMap      map[string]string      `json:"attributeMap,omitempty"`
	EventMap          map[string]string      `json:"eventMap,omitempty"`
	EntityMap         map[string]string      `json:"entityMap,omitempty"`
	CustomAttributes  []string               `json:"customAttributes,omitempty"`
	Direction         string                 `json:"direction,omitempty"`
	Enabled           bool                   `json:"enabled,omitempty"`
	Extra             map[string]interface{} `json:"extra,omitempty"` // for extensibility
}

// Integration represents an installed integration instance
 type Integration struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Status      string            `json:"status"`
	Config      IntegrationConfig `json:"config"`
	CreatedBy   string            `json:"createdBy"`
	CreatedTime string            `json:"createdTime"`
	UpdatedBy   string            `json:"updatedBy"`
	UpdatedTime string            `json:"updatedTime"`
}

// IntegrationType represents the schema for a supported integration type
 type IntegrationType struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Category    string                 `json:"category"`
	ConfigSchema map[string]interface{} `json:"configSchema"`
}
