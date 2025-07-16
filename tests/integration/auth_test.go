package integration

import (
	"testing"

	"github.com/VI-IM/im_backend_go/internal/testutil"
)

func TestAuthEndpoints(t *testing.T) {
	// Setup test server
	server := testutil.NewTestServer(t)
	defer server.Close()

	client := testutil.NewHTTPClient(server.URL())

	t.Run("Health Check", func(t *testing.T) {
		resp, err := client.GET("/health", nil)
		if err != nil {
			t.Fatalf("Failed to make health check request: %v", err)
		}

		testutil.AssertStatus(t, resp, 200)
	})

	t.Run("Generate Token - Valid Credentials", func(t *testing.T) {
		// Create a test user first
		testUser := testutil.CreateBusinessPartner(t, server.Client)

		// Prepare login request
		loginReq := map[string]string{
			"username": testUser.Username,
			"password": testUser.Password,
		}

		var loginResp testutil.LoginResponse
		resp, err := client.POST("/v1/api/auth/generate-token", loginReq, &loginResp)
		if err != nil {
			t.Fatalf("Failed to make login request: %v", err)
		}

		// Assert successful login
		testutil.AssertSuccessResponse(t, resp)
		
		if loginResp.Token == "" {
			t.Error("Expected token in response, got empty string")
		}
	})

	t.Run("Generate Token - Invalid Credentials", func(t *testing.T) {
		loginReq := map[string]string{
			"username": "nonexistent",
			"password": "wrongpassword",
		}

		var loginResp testutil.LoginResponse
		resp, err := client.POST("/v1/api/auth/generate-token", loginReq, &loginResp)
		if err != nil {
			t.Fatalf("Failed to make login request: %v", err)
		}

		// Assert error response
		testutil.AssertErrorResponse(t, resp, 401)
	})

	t.Run("Signup - Valid Data", func(t *testing.T) {
		signupReq := map[string]interface{}{
			"username":   "newuser" + testutil.RandomString(8),
			"email":      testutil.RandomEmail(),
			"password":   "newpassword123",
			"name":       "New User",
			"role":       "business_partner",
		}

		var signupResp testutil.APIResponse
		resp, err := client.POST("/v1/api/auth/signup", signupReq, &signupResp)
		if err != nil {
			t.Fatalf("Failed to make signup request: %v", err)
		}

		// Assert successful signup
		testutil.AssertSuccessResponse(t, resp)
	})

	t.Run("Signup - Duplicate Username", func(t *testing.T) {
		// Create a test user first
		testUser := testutil.CreateBusinessPartner(t, server.Client)

		signupReq := map[string]interface{}{
			"username": testUser.Username, // Use existing username
			"email":    testutil.RandomEmail(),
			"password": "newpassword123",
			"name":     "New User",
			"role":     "business_partner",
		}

		var signupResp testutil.APIResponse
		resp, err := client.POST("/v1/api/auth/signup", signupReq, &signupResp)
		if err != nil {
			t.Fatalf("Failed to make signup request: %v", err)
		}

		// Assert error response
		testutil.AssertBadRequest(t, resp)
	})

	t.Run("Signup - Invalid Data", func(t *testing.T) {
		signupReq := map[string]interface{}{
			"username": "", // Empty username
			"email":    "invalid-email",
			"password": "123", // Too short
		}

		var signupResp testutil.APIResponse
		resp, err := client.POST("/v1/api/auth/signup", signupReq, &signupResp)
		if err != nil {
			t.Fatalf("Failed to make signup request: %v", err)
		}

		// Assert validation error
		testutil.AssertValidationError(t, resp, "username", "email", "password")
	})
}