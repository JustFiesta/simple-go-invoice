package models

// APIResponse represents a standard REST API response
type APIResponse struct {
	Data  interface{}       `json:"data,omitempty"`
	Links *ResponseLinks    `json:"links,omitempty"`
	Meta  *ResponseMeta     `json:"meta,omitempty"`
	Error *APIError         `json:"error,omitempty"`
}

// ResponseLinks provides HATEOAS links
type ResponseLinks struct {
	Self    string            `json:"self,omitempty"`
	Related map[string]string `json:"related,omitempty"`
}

// ResponseMeta provides pagination and additional metadata
type ResponseMeta struct {
	Page       int   `json:"page,omitempty"`
	Limit      int   `json:"limit,omitempty"`
	Total      int64 `json:"total,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`
}

// APIError represents a standard error response
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// PaginatedResponse creates a paginated response
func PaginatedResponse(data interface{}, page, limit int, total int64, selfLink string) APIResponse {
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	
	return APIResponse{
		Data: data,
		Meta: &ResponseMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
		Links: &ResponseLinks{
			Self: selfLink,
		},
	}
}

// SingleResourceResponse creates a response for single resource
func SingleResourceResponse(data interface{}, selfLink string, related map[string]string) APIResponse {
	return APIResponse{
		Data: data,
		Links: &ResponseLinks{
			Self:    selfLink,
			Related: related,
		},
	}
}

// ErrorResponse creates an error response
func ErrorResponse(code int, message, details string) APIResponse {
	return APIResponse{
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}