package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

// HTTPClient provides utilities for making HTTP requests in tests
type HTTPClient struct {
	Client      *http.Client
	BaseURL     string
	AuthToken   string
}

// NewHTTPClient creates a new HTTP client for testing
func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		Client:  &http.Client{},
		BaseURL: baseURL,
	}
}

// SetAuthToken sets the JWT token for authenticated requests
func (c *HTTPClient) SetAuthToken(token string) {
	c.AuthToken = token
}

// NewHTTPRequest creates a new HTTP request with optional JSON body
func NewHTTPRequest(method, url string, body interface{}) (*http.Request, error) {
	var bodyReader io.Reader

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

// DoRequest performs an HTTP request and returns the response
func (c *HTTPClient) DoRequest(req *http.Request) (*http.Response, error) {
	if c.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}
	return c.Client.Do(req)
}

// DoJSONRequest performs a request and unmarshals JSON response
func (c *HTTPClient) DoJSONRequest(req *http.Request, result interface{}) (*http.Response, error) {
	resp, err := c.DoRequest(req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return resp, fmt.Errorf("failed to decode JSON response: %w", err)
		}
	}

	return resp, nil
}

// GET performs a GET request
func (c *HTTPClient) GET(path string, result interface{}) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.BaseURL+path, nil)
	if err != nil {
		return nil, err
	}
	return c.DoJSONRequest(req, result)
}

// POST performs a POST request with JSON body
func (c *HTTPClient) POST(path string, body interface{}, result interface{}) (*http.Response, error) {
	req, err := NewHTTPRequest("POST", c.BaseURL+path, body)
	if err != nil {
		return nil, err
	}
	return c.DoJSONRequest(req, result)
}

// PUT performs a PUT request with JSON body
func (c *HTTPClient) PUT(path string, body interface{}, result interface{}) (*http.Response, error) {
	req, err := NewHTTPRequest("PUT", c.BaseURL+path, body)
	if err != nil {
		return nil, err
	}
	return c.DoJSONRequest(req, result)
}

// DELETE performs a DELETE request
func (c *HTTPClient) DELETE(path string, result interface{}) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", c.BaseURL+path, nil)
	if err != nil {
		return nil, err
	}
	return c.DoJSONRequest(req, result)
}

// AssertStatus checks if response status matches expected
func AssertStatus(t *testing.T, resp *http.Response, expectedStatus int) {
	if resp.StatusCode != expectedStatus {
		t.Errorf("Expected status %d, got %d", expectedStatus, resp.StatusCode)
	}
}

// AssertJSON compares response JSON with expected structure
func AssertJSON(t *testing.T, resp *http.Response, expected interface{}) {
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("Failed to marshal expected JSON: %v", err)
	}

	var actualData, expectedData interface{}
	
	if err := json.Unmarshal(body, &actualData); err != nil {
		t.Fatalf("Failed to unmarshal actual JSON: %v", err)
	}
	
	if err := json.Unmarshal(expectedJSON, &expectedData); err != nil {
		t.Fatalf("Failed to unmarshal expected JSON: %v", err)
	}

	if !jsonEqual(actualData, expectedData) {
		t.Errorf("JSON mismatch.\nExpected: %s\nActual: %s", expectedJSON, body)
	}
}

// jsonEqual compares two JSON structures
func jsonEqual(a, b interface{}) bool {
	aJSON, _ := json.Marshal(a)
	bJSON, _ := json.Marshal(b)
	return string(aJSON) == string(bJSON)
}