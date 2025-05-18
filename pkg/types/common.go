package types

// Common types for OpsRamp MCP integration

// Standard ID type
 type ID string

// Standard timestamp type
 type Timestamp string

// Standard error response
 type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Pagination info
 type Pagination struct {
	Page     int `json:"page"`
	PerPage  int `json:"perPage"`
	Total    int `json:"total"`
}
