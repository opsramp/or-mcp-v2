package types

import (
	"encoding/json"
	"testing"
)

func TestResourceSerialization(t *testing.T) {
	// Test Resource struct JSON serialization
	resource := Resource{
		ID:           "res-123",
		HostName:     "test-host",
		IPAddress:    "192.168.1.100",
		Name:         "Test Resource",
		ResourceName: "test-resource",
		Type:         "DEVICE",
		ResourceType: "Windows Server",
		State:        "UP",
		Status:       "ACTIVE",
		Location:     &Location{ID: "loc-1", Name: "Data Center 1"},
		Tags: []Tag{
			{Name: "environment", Value: "production"},
			{Name: "team", Value: "infrastructure"},
		},
		CreatedDate: "2023-01-01T00:00:00Z",
		UpdatedDate: "2023-01-02T00:00:00Z",
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(resource)
	if err != nil {
		t.Fatalf("Failed to marshal Resource to JSON: %v", err)
	}

	// Unmarshal back to Resource
	var unmarshaledResource Resource
	err = json.Unmarshal(jsonData, &unmarshaledResource)
	if err != nil {
		t.Fatalf("Failed to unmarshal Resource from JSON: %v", err)
	}

	// Verify key fields
	if unmarshaledResource.ID != resource.ID {
		t.Errorf("Expected ID %s, got %s", resource.ID, unmarshaledResource.ID)
	}
	if unmarshaledResource.HostName != resource.HostName {
		t.Errorf("Expected HostName %s, got %s", resource.HostName, unmarshaledResource.HostName)
	}
	if unmarshaledResource.IPAddress != resource.IPAddress {
		t.Errorf("Expected IPAddress %s, got %s", resource.IPAddress, unmarshaledResource.IPAddress)
	}
	if len(unmarshaledResource.Tags) != len(resource.Tags) {
		t.Errorf("Expected %d tags, got %d", len(resource.Tags), len(unmarshaledResource.Tags))
	}
}

func TestResourceMinimalSerialization(t *testing.T) {
	// Test ResourceMinimal struct
	minimal := ResourceMinimal{
		ID:           "res-456",
		HostName:     "minimal-host",
		IPAddress:    "10.0.0.1",
		Name:         "Minimal Resource",
		ResourceName: "minimal-resource",
		Type:         "SERVICE",
		ResourceType: "Web Service",
		State:        "DOWN",
		Status:       "INACTIVE",
		Location:     &Location{ID: "loc-2", Name: "Cloud Region"},
		Tags: []Tag{
			{Name: "service", Value: "web"},
		},
		UpdatedDate: "2023-01-03T00:00:00Z",
	}

	// Test JSON roundtrip
	jsonData, err := json.Marshal(minimal)
	if err != nil {
		t.Fatalf("Failed to marshal ResourceMinimal: %v", err)
	}

	var unmarshaled ResourceMinimal
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal ResourceMinimal: %v", err)
	}

	if unmarshaled.ID != minimal.ID {
		t.Errorf("Expected ID %s, got %s", minimal.ID, unmarshaled.ID)
	}
}

func TestResourceSearchParams(t *testing.T) {
	// Test ResourceSearchParams with various configurations
	params := ResourceSearchParams{
		PageNo:            1,
		PageSize:          50,
		SortName:          "name",
		IsDescendingOrder: true,
		QueryString:       "production",
		HostName:          "web-server",
		DNSName:           "web.example.com",
		ResourceName:      "web-resource",
		IPAddress:         "172.16.0.1",
		State:             "UP",
		Type:              "DEVICE",
		ResourceType:      "Linux Server",
		Tags:              "environment:production,team:web",
		AgentInstalled:    BoolPtr(true),
		DeviceGroup:       "Web Servers",
		ServiceGroup:      "Frontend Services",
	}

	// Test JSON serialization
	jsonData, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Failed to marshal ResourceSearchParams: %v", err)
	}

	var unmarshaled ResourceSearchParams
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal ResourceSearchParams: %v", err)
	}

	// Verify key fields
	if unmarshaled.PageSize != params.PageSize {
		t.Errorf("Expected PageSize %d, got %d", params.PageSize, unmarshaled.PageSize)
	}
	if unmarshaled.IsDescendingOrder != params.IsDescendingOrder {
		t.Errorf("Expected IsDescendingOrder %v, got %v", params.IsDescendingOrder, unmarshaled.IsDescendingOrder)
	}
	if *unmarshaled.AgentInstalled != *params.AgentInstalled {
		t.Errorf("Expected AgentInstalled %v, got %v", *params.AgentInstalled, *unmarshaled.AgentInstalled)
	}
}

func TestResourceSearchResponse(t *testing.T) {
	// Test ResourceSearchResponse
	response := ResourceSearchResponse{
		Results: []Resource{
			{
				ID:           "res-1",
				HostName:     "host1",
				Name:         "Resource 1",
				ResourceType: "Server",
			},
			{
				ID:           "res-2",
				HostName:     "host2",
				Name:         "Resource 2",
				ResourceType: "Service",
			},
		},
		TotalResults:    2,
		PageNo:          1,
		PageSize:        10,
		NextPage:        false,
		DescendingOrder: false,
	}

	// Test JSON roundtrip
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal ResourceSearchResponse: %v", err)
	}

	var unmarshaled ResourceSearchResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal ResourceSearchResponse: %v", err)
	}

	if len(unmarshaled.Results) != len(response.Results) {
		t.Errorf("Expected %d results, got %d", len(response.Results), len(unmarshaled.Results))
	}
	if unmarshaled.TotalResults != response.TotalResults {
		t.Errorf("Expected TotalResults %d, got %d", response.TotalResults, unmarshaled.TotalResults)
	}
}

func TestDeviceGroup(t *testing.T) {
	// Test DeviceGroup struct
	group := DeviceGroup{
		ID:            "group-123",
		Name:          "Production Servers",
		Description:   "All production servers",
		Type:          "STATIC",
		ParentID:      "parent-group-456",
		CreatedBy:     "admin",
		CreatedDate:   "2023-01-01T00:00:00Z",
		UpdatedBy:     "admin",
		UpdatedDate:   "2023-01-02T00:00:00Z",
		ResourceCount: 25,
	}

	// Test JSON serialization
	jsonData, err := json.Marshal(group)
	if err != nil {
		t.Fatalf("Failed to marshal DeviceGroup: %v", err)
	}

	var unmarshaled DeviceGroup
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal DeviceGroup: %v", err)
	}

	if unmarshaled.ID != group.ID {
		t.Errorf("Expected ID %s, got %s", group.ID, unmarshaled.ID)
	}
	if unmarshaled.Name != group.Name {
		t.Errorf("Expected Name %s, got %s", group.Name, unmarshaled.Name)
	}
	if unmarshaled.ResourceCount != group.ResourceCount {
		t.Errorf("Expected ResourceCount %d, got %d", group.ResourceCount, unmarshaled.ResourceCount)
	}
}

func TestResourceError(t *testing.T) {
	// Test ResourceError creation and methods
	err := NewResourceError(ResourceErrorTypeNotFound, "RESOURCE_NOT_FOUND", "Resource with ID 'test-123' not found")

	if err.Type != ResourceErrorTypeNotFound {
		t.Errorf("Expected error type %s, got %s", ResourceErrorTypeNotFound, err.Type)
	}
	if err.Code != "RESOURCE_NOT_FOUND" {
		t.Errorf("Expected error code RESOURCE_NOT_FOUND, got %s", err.Code)
	}

	// Test Error() method
	errorString := err.Error()
	if errorString == "" {
		t.Errorf("Expected non-empty error string")
	}

	// Test JSON serialization
	jsonData, err2 := json.Marshal(err)
	if err2 != nil {
		t.Fatalf("Failed to marshal ResourceError: %v", err2)
	}

	var unmarshaled ResourceError
	err3 := json.Unmarshal(jsonData, &unmarshaled)
	if err3 != nil {
		t.Fatalf("Failed to unmarshal ResourceError: %v", err3)
	}

	if unmarshaled.Type != err.Type {
		t.Errorf("Expected error type %s, got %s", err.Type, unmarshaled.Type)
	}
}

func TestResourceCreateRequest(t *testing.T) {
	// Test ResourceCreateRequest
	request := ResourceCreateRequest{
		HostName:     "new-server",
		IPAddress:    "192.168.1.200",
		ResourceType: "Linux Server",
		Location:     "Data Center 2",
		Tags: []Tag{
			{Name: "environment", Value: "staging"},
		},
		Properties: map[string]interface{}{
			"os":      "Ubuntu 20.04",
			"memory":  "16GB",
			"storage": "500GB",
		},
	}

	// Test JSON serialization
	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal ResourceCreateRequest: %v", err)
	}

	var unmarshaled ResourceCreateRequest
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal ResourceCreateRequest: %v", err)
	}

	if unmarshaled.HostName != request.HostName {
		t.Errorf("Expected HostName %s, got %s", request.HostName, unmarshaled.HostName)
	}
	if unmarshaled.ResourceType != request.ResourceType {
		t.Errorf("Expected ResourceType %s, got %s", request.ResourceType, unmarshaled.ResourceType)
	}
}

func TestResourceUpdateRequest(t *testing.T) {
	// Test ResourceUpdateRequest
	request := ResourceUpdateRequest{
		HostName: "updated-host",
		Location: "Updated Location",
		Tags: []Tag{
			{Name: "environment", Value: "production"},
			{Name: "updated", Value: "true"},
		},
		Properties: map[string]interface{}{
			"memory": "32GB",
			"notes":  "Resource updated during maintenance",
		},
	}

	// Test JSON serialization
	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal ResourceUpdateRequest: %v", err)
	}

	var unmarshaled ResourceUpdateRequest
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal ResourceUpdateRequest: %v", err)
	}

	if unmarshaled.HostName != request.HostName {
		t.Errorf("Expected HostName %s, got %s", request.HostName, unmarshaled.HostName)
	}
	if len(unmarshaled.Tags) != len(request.Tags) {
		t.Errorf("Expected %d tags, got %d", len(request.Tags), len(unmarshaled.Tags))
	}
}

func TestResourceBulkOperations(t *testing.T) {
	// Test ResourceBulkUpdateRequest
	bulkUpdate := ResourceBulkUpdateRequest{
		ResourceIDs: []string{"res-1", "res-2", "res-3"},
		Updates: map[string]interface{}{
			"location": "Updated Data Center",
			"status":   "ACTIVE",
		},
	}

	jsonData, err := json.Marshal(bulkUpdate)
	if err != nil {
		t.Fatalf("Failed to marshal ResourceBulkUpdateRequest: %v", err)
	}

	var unmarshaled ResourceBulkUpdateRequest
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal ResourceBulkUpdateRequest: %v", err)
	}

	if len(unmarshaled.ResourceIDs) != len(bulkUpdate.ResourceIDs) {
		t.Errorf("Expected %d resource IDs, got %d", len(bulkUpdate.ResourceIDs), len(unmarshaled.ResourceIDs))
	}

	// Test ResourceBulkDeleteRequest
	bulkDelete := ResourceBulkDeleteRequest{
		ResourceIDs: []string{"res-4", "res-5"},
	}

	jsonData, err = json.Marshal(bulkDelete)
	if err != nil {
		t.Fatalf("Failed to marshal ResourceBulkDeleteRequest: %v", err)
	}

	var unmarshaledDelete ResourceBulkDeleteRequest
	err = json.Unmarshal(jsonData, &unmarshaledDelete)
	if err != nil {
		t.Fatalf("Failed to unmarshal ResourceBulkDeleteRequest: %v", err)
	}

	if len(unmarshaledDelete.ResourceIDs) != len(bulkDelete.ResourceIDs) {
		t.Errorf("Expected %d resource IDs, got %d", len(bulkDelete.ResourceIDs), len(unmarshaledDelete.ResourceIDs))
	}
}

func TestTag(t *testing.T) {
	// Test Tag struct
	tag := Tag{
		Name:  "environment",
		Value: "production",
	}

	// Test JSON serialization
	jsonData, err := json.Marshal(tag)
	if err != nil {
		t.Fatalf("Failed to marshal Tag: %v", err)
	}

	var unmarshaled Tag
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal Tag: %v", err)
	}

	if unmarshaled.Name != tag.Name {
		t.Errorf("Expected tag name %s, got %s", tag.Name, unmarshaled.Name)
	}
	if unmarshaled.Value != tag.Value {
		t.Errorf("Expected tag value %s, got %s", tag.Value, unmarshaled.Value)
	}
}

// Helper functions for test data
func BoolPtr(b bool) *bool {
	return &b
}

func IntPtr(i int) *int {
	return &i
}

func StringPtr(s string) *string {
	return &s
}
