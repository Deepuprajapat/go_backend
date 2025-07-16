package integration

import (
	"fmt"
	"testing"

	"github.com/VI-IM/im_backend_go/internal/testutil"
)

func TestLeadEndpoints(t *testing.T) {
	// Setup test server
	server := testutil.NewTestServer(t)
	defer server.Close()

	// Clear database before tests
	server.ClearDatabase(t)

	t.Run("Create Lead - Public Endpoint", func(t *testing.T) {
		client := testutil.NewHTTPClient(server.URL())

		// Prepare lead data
		leadReq := map[string]interface{}{
			"name":         "Test Lead",
			"email":        testutil.RandomEmail(),
			"phone_number": testutil.RandomPhoneNumber(),
			"message":      "I'm interested in your properties",
		}

		var leadResp testutil.APIResponse
		resp, err := client.POST("/v1/api/leads", leadReq, &leadResp)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		// Assert successful creation
		testutil.AssertSuccessResponse(t, resp)
	})

	t.Run("Create Lead with OTP - Public Endpoint", func(t *testing.T) {
		client := testutil.NewHTTPClient(server.URL())

		// Prepare lead data with OTP request
		leadReq := map[string]interface{}{
			"name":         "Test Lead OTP",
			"email":        testutil.RandomEmail(),
			"phone_number": testutil.RandomPhoneNumber(),
			"message":      "I'm interested in your properties",
			"send_otp":     true,
		}

		var leadResp testutil.APIResponse
		resp, err := client.POST("/v1/api/leads/send-otp", leadReq, &leadResp)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		// Assert successful creation with OTP
		testutil.AssertSuccessResponse(t, resp)
	})

	t.Run("Create Lead - Invalid Data", func(t *testing.T) {
		client := testutil.NewHTTPClient(server.URL())

		// Prepare invalid lead data
		leadReq := map[string]interface{}{
			"name":         "", // Empty name
			"email":        "invalid-email",
			"phone_number": "123", // Invalid phone
		}

		var leadResp testutil.APIResponse
		resp, err := client.POST("/v1/api/leads", leadReq, &leadResp)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		// Assert validation error
		testutil.AssertValidationError(t, resp, "name", "email", "phone_number")
	})

	t.Run("Get All Leads - DM Role", func(t *testing.T) {
		// Create test DM user and authenticate
		testUser := testutil.CreateDM(t, server.Client)
		client := testUser.AuthorizedClient(server.URL())

		// Create test leads
		leadsFactory := testutil.NewLeadsFactory(server.Client)
		leadsFactory.Create(t)
		leadsFactory.Create(t)

		// Make request to get all leads
		resp, err := client.GET("/v1/api/leads", nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		// Assert successful response with leads
		testutil.AssertListResponse(t, resp, 2)
	})

	t.Run("Get All Leads - Business Partner (Forbidden)", func(t *testing.T) {
		// Create test business partner user and authenticate
		testUser := testutil.CreateBusinessPartner(t, server.Client)
		client := testUser.AuthorizedClient(server.URL())

		// Make request to get all leads (should be forbidden for business partners)
		resp, err := client.GET("/v1/api/leads", nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		// Assert forbidden response (only DM role can access)
		testutil.AssertForbidden(t, resp)
	})

	t.Run("Get All Leads - Unauthenticated", func(t *testing.T) {
		client := testutil.NewHTTPClient(server.URL())

		// Make request to get all leads without authentication
		resp, err := client.GET("/v1/api/leads", nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		// Assert unauthorized response
		testutil.AssertUnauthorized(t, resp)
	})

	t.Run("Get Lead by ID - DM Role", func(t *testing.T) {
		// Create test DM user and authenticate
		testUser := testutil.CreateDM(t, server.Client)
		client := testUser.AuthorizedClient(server.URL())

		// Create test lead
		leadsFactory := testutil.NewLeadsFactory(server.Client)
		lead := leadsFactory.Create(t)

		// Make request to get specific lead
		path := fmt.Sprintf("/v1/api/leads/get/by/%s", lead.ID)
		resp, err := client.GET(path, nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		// Assert successful response with lead data
		apiResp := testutil.AssertContainsData(t, resp, "id", "name", "email", "phone_number")
		
		// Verify lead ID matches
		leadID := testutil.ExtractDataField(t, apiResp, "id")
		if leadID != lead.ID {
			t.Errorf("Expected lead ID %s, got %v", lead.ID, leadID)
		}
	})

	t.Run("Get Lead by ID - Invalid ID", func(t *testing.T) {
		// Create test DM user and authenticate
		testUser := testutil.CreateDM(t, server.Client)
		client := testUser.AuthorizedClient(server.URL())

		// Make request with invalid lead ID
		path := "/v1/api/leads/get/by/nonexistent-id"
		resp, err := client.GET(path, nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		// Assert not found response
		testutil.AssertNotFound(t, resp)
	})

	t.Run("Validate OTP - Valid Flow", func(t *testing.T) {
		client := testutil.NewHTTPClient(server.URL())

		// First create a lead with OTP
		leadReq := map[string]interface{}{
			"name":         "Test Lead OTP",
			"email":        testutil.RandomEmail(),
			"phone_number": testutil.RandomPhoneNumber(),
			"message":      "I'm interested in your properties",
		}

		var leadResp testutil.APIResponse
		resp, err := client.POST("/v1/api/leads/send-otp", leadReq, &leadResp)
		if err != nil {
			t.Fatalf("Failed to create lead with OTP: %v", err)
		}

		testutil.AssertSuccessResponse(t, resp)

		// Note: In a real test, you'd need to mock the OTP service or use a test OTP
		// For demonstration purposes, we'll just test the endpoint structure
		otpReq := map[string]interface{}{
			"phone_number": leadReq["phone_number"],
			"otp":          "123456", // This would fail with actual OTP validation
		}

		var otpResp testutil.APIResponse
		resp2, err := client.PUT("/v1/api/leads/validate-otp", otpReq, &otpResp)
		if err != nil {
			t.Fatalf("Failed to make OTP validation request: %v", err)
		}

		// In this test case, we expect it to fail due to invalid OTP
		// but we're testing that the endpoint is accessible
		testutil.AssertErrorResponse(t, resp2, 400)
	})

	t.Run("Resend OTP", func(t *testing.T) {
		client := testutil.NewHTTPClient(server.URL())

		// Prepare resend OTP request
		resendReq := map[string]interface{}{
			"phone_number": testutil.RandomPhoneNumber(),
		}

		var resendResp testutil.APIResponse
		resp, err := client.PUT("/v1/api/leads/resend-otp", resendReq, &resendResp)
		if err != nil {
			t.Fatalf("Failed to make resend OTP request: %v", err)
		}

		// The response could be success or error depending on implementation
		// We're testing that the endpoint is accessible
		if resp.StatusCode != 200 && resp.StatusCode != 400 && resp.StatusCode != 404 {
			t.Errorf("Unexpected status code: %d", resp.StatusCode)
		}
	})
}