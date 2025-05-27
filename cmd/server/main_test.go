package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/VI-IM/im_backend_go/internal/router"
	"github.com/VI-IM/im_backend_go/internal/testutil"
)

func TestMain(m *testing.M) {
	// Setup test database
	cleanup := testutil.SetupTestDB(nil)
	defer cleanup()

	// Initialize router for tests
	router.Init()

	// Run tests
	os.Exit(m.Run())
}

func TestHealthEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	router.Router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", response["status"])
	}
}

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name       string
		payload    map[string]interface{}
		wantStatus int
	}{
		{
			name: "valid request",
			payload: map[string]interface{}{
				"phone": "1234567890",
				"code":  "123456",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "invalid phone",
			payload: map[string]interface{}{
				"phone": "invalid",
				"code":  "123456",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "missing code",
			payload: map[string]interface{}{
				"phone": "1234567890",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.payload)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/auth/token", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.Router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status code %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}
