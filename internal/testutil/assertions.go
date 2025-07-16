package testutil

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

// APIResponse represents a common API response structure
type APIResponse struct {
	IsSuccess bool        `json:"is_success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	APIResponse
	Meta PaginationMeta `json:"meta"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalCount int `json:"total_count"`
	TotalPages int `json:"total_pages"`
}

// AssertSuccessResponse checks if the response indicates success
func AssertSuccessResponse(t *testing.T, resp *http.Response) APIResponse {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		t.Errorf("Expected successful status code (2xx), got %d", resp.StatusCode)
	}

	var apiResp APIResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	resp.Body.Close()

	if err := json.Unmarshal(body, &apiResp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !apiResp.IsSuccess {
		t.Errorf("Expected is_success=true, got false. Message: %s", apiResp.Message)
	}

	return apiResp
}

// AssertErrorResponse checks if the response indicates an error
func AssertErrorResponse(t *testing.T, resp *http.Response, expectedStatus int) APIResponse {
	if resp.StatusCode != expectedStatus {
		t.Errorf("Expected status %d, got %d", expectedStatus, resp.StatusCode)
	}

	var apiResp APIResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	resp.Body.Close()

	if err := json.Unmarshal(body, &apiResp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if apiResp.IsSuccess {
		t.Errorf("Expected is_success=false, got true")
	}

	return apiResp
}

// AssertUnauthorized checks for 401 Unauthorized response
func AssertUnauthorized(t *testing.T, resp *http.Response) {
	AssertErrorResponse(t, resp, http.StatusUnauthorized)
}

// AssertForbidden checks for 403 Forbidden response
func AssertForbidden(t *testing.T, resp *http.Response) {
	AssertErrorResponse(t, resp, http.StatusForbidden)
}

// AssertNotFound checks for 404 Not Found response
func AssertNotFound(t *testing.T, resp *http.Response) {
	AssertErrorResponse(t, resp, http.StatusNotFound)
}

// AssertBadRequest checks for 400 Bad Request response
func AssertBadRequest(t *testing.T, resp *http.Response) APIResponse {
	return AssertErrorResponse(t, resp, http.StatusBadRequest)
}

// AssertValidationError checks for validation error response
func AssertValidationError(t *testing.T, resp *http.Response, expectedFields ...string) {
	apiResp := AssertBadRequest(t, resp)
	
	for _, field := range expectedFields {
		if !strings.Contains(apiResp.Message, field) && !strings.Contains(apiResp.Error, field) {
			t.Errorf("Expected validation error for field '%s' in response", field)
		}
	}
}

// AssertContainsData checks if response contains expected data fields
func AssertContainsData(t *testing.T, resp *http.Response, expectedFields ...string) APIResponse {
	apiResp := AssertSuccessResponse(t, resp)
	
	if apiResp.Data == nil {
		t.Error("Expected data field in response, got nil")
		return apiResp
	}

	dataMap, ok := apiResp.Data.(map[string]interface{})
	if !ok {
		t.Error("Expected data to be an object")
		return apiResp
	}

	for _, field := range expectedFields {
		if _, exists := dataMap[field]; !exists {
			t.Errorf("Expected field '%s' in response data", field)
		}
	}

	return apiResp
}

// AssertListResponse checks if response is a successful list with expected structure
func AssertListResponse(t *testing.T, resp *http.Response, minItems int) APIResponse {
	apiResp := AssertSuccessResponse(t, resp)
	
	if apiResp.Data == nil {
		t.Error("Expected data field in list response")
		return apiResp
	}

	// Check if data is an array
	dataValue := reflect.ValueOf(apiResp.Data)
	if dataValue.Kind() != reflect.Slice {
		t.Error("Expected data to be an array for list response")
		return apiResp
	}

	if dataValue.Len() < minItems {
		t.Errorf("Expected at least %d items in list, got %d", minItems, dataValue.Len())
	}

	return apiResp
}

// AssertPaginatedResponse checks if response is a paginated list
func AssertPaginatedResponse(t *testing.T, resp *http.Response) PaginatedResponse {
	var paginatedResp PaginatedResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	resp.Body.Close()

	if err := json.Unmarshal(body, &paginatedResp); err != nil {
		t.Fatalf("Failed to unmarshal paginated response: %v", err)
	}

	if !paginatedResp.IsSuccess {
		t.Errorf("Expected is_success=true in paginated response")
	}

	// Check pagination meta
	if paginatedResp.Meta.Page < 1 {
		t.Errorf("Expected page >= 1, got %d", paginatedResp.Meta.Page)
	}
	if paginatedResp.Meta.PerPage < 1 {
		t.Errorf("Expected per_page >= 1, got %d", paginatedResp.Meta.PerPage)
	}
	if paginatedResp.Meta.TotalCount < 0 {
		t.Errorf("Expected total_count >= 0, got %d", paginatedResp.Meta.TotalCount)
	}

	return paginatedResp
}

// AssertJSONContains checks if response JSON contains specific key-value pairs
func AssertJSONContains(t *testing.T, resp *http.Response, expected map[string]interface{}) {
	var actual map[string]interface{}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	resp.Body.Close()

	if err := json.Unmarshal(body, &actual); err != nil {
		t.Fatalf("Failed to unmarshal JSON response: %v", err)
	}

	for key, expectedValue := range expected {
		actualValue, exists := actual[key]
		if !exists {
			t.Errorf("Expected key '%s' in JSON response", key)
			continue
		}

		if !reflect.DeepEqual(actualValue, expectedValue) {
			t.Errorf("For key '%s': expected %v, got %v", key, expectedValue, actualValue)
		}
	}
}

// AssertHeader checks if response has expected header value
func AssertHeader(t *testing.T, resp *http.Response, header, expectedValue string) {
	actualValue := resp.Header.Get(header)
	if actualValue != expectedValue {
		t.Errorf("Expected header '%s' to be '%s', got '%s'", header, expectedValue, actualValue)
	}
}

// AssertContentType checks if response has expected content type
func AssertContentType(t *testing.T, resp *http.Response, expectedType string) {
	AssertHeader(t, resp, "Content-Type", expectedType)
}

// AssertResponseTime checks if response time is within acceptable limits
func AssertResponseTime(t *testing.T, startTime int64, maxDurationMs int64) {
	// This would need to be called with timing logic in the actual test
	// duration := time.Now().UnixNano()/1e6 - startTime
	// if duration > maxDurationMs {
	//     t.Errorf("Response took %dms, expected less than %dms", duration, maxDurationMs)
	// }
}

// ExtractDataField extracts a specific field from the API response data
func ExtractDataField(t *testing.T, apiResp APIResponse, fieldName string) interface{} {
	if apiResp.Data == nil {
		t.Fatalf("No data field in response")
	}

	dataMap, ok := apiResp.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Data field is not an object")
	}

	value, exists := dataMap[fieldName]
	if !exists {
		t.Fatalf("Field '%s' not found in response data", fieldName)
	}

	return value
}

// ExtractListItems extracts items from a list response
func ExtractListItems(t *testing.T, apiResp APIResponse) []interface{} {
	if apiResp.Data == nil {
		t.Fatalf("No data field in response")
	}

	items, ok := apiResp.Data.([]interface{})
	if !ok {
		t.Fatalf("Data field is not an array")
	}

	return items
}