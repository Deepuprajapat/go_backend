package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VI-IM/im_backend_go/internal/auth"
	"github.com/VI-IM/im_backend_go/internal/middleware"
	"github.com/VI-IM/im_backend_go/internal/testutil"
	"github.com/gorilla/mux"
)

func setupTestEnv(t *testing.T) (*httptest.Server, func()) {
	cleanup := testutil.SetupTestDB(t)
	router := setupRouter()
	ts := httptest.NewServer(router)
	return ts, func() {
		cleanup()
		ts.Close()
	}
}

func setupRouter() *mux.Router {
	r := mux.NewRouter()

	// Add routes
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	r.HandleFunc("/auth/token", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Phone string `json:"phone"`
			OTP   string `json:"otp"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// TODO: Verify OTP here using ent client
		// For now, we'll just generate a token
		token, err := auth.GenerateToken(1, true, request.Phone) // Replace with actual user data
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		// Set the token as a cookie
		middleware.SetAuthCookie(w, token)

		// Return success response
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Authentication successful",
		})
	}).Methods(http.MethodPost)

	return r
}

func TestHealthEndpoint(t *testing.T) {
	ts, cleanup := setupTestEnv(t)
	defer cleanup()

	resp, err := http.Get(ts.URL + "/health")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	body := make([]byte, 2)
	_, err = resp.Body.Read(body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if string(body) != "OK" {
		t.Errorf("Expected response body 'OK', got '%s'", string(body))
	}
}

func TestAuthTokenEndpoint(t *testing.T) {
	ts, cleanup := setupTestEnv(t)
	defer cleanup()

	tests := []struct {
		name       string
		phone      string
		otp        string
		wantStatus int
	}{
		{
			name:       "Valid request",
			phone:      "+1234567890",
			otp:        "123456",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid JSON",
			phone:      "",
			otp:        "",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Missing phone",
			phone:      "",
			otp:        "123456",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Missing OTP",
			phone:      "+1234567890",
			otp:        "",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request body
			reqBody := struct {
				Phone string `json:"phone"`
				OTP   string `json:"otp"`
			}{
				Phone: tt.phone,
				OTP:   tt.otp,
			}

			jsonBody, err := json.Marshal(reqBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			// Make request
			resp, err := http.Post(ts.URL+"/auth/token", "application/json", bytes.NewBuffer(jsonBody))
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			// Check status code
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("Expected status code %d, got %d", tt.wantStatus, resp.StatusCode)
			}

			// If successful, check for cookie
			if resp.StatusCode == http.StatusOK {
				cookies := resp.Cookies()
				found := false
				for _, cookie := range cookies {
					if cookie.Name == "auth_token" {
						found = true
						break
					}
				}
				if !found {
					t.Error("Expected auth_token cookie, not found")
				}
			}
		})
	}
}
