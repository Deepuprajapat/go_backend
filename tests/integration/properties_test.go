package integration

import (
	"fmt"
	"testing"

	"github.com/VI-IM/im_backend_go/internal/testutil"
)

func TestPropertyEndpoints(t *testing.T) {
	// Setup test server
	server := testutil.NewTestServer(t)
	defer server.Close()

	// Clear database before tests
	server.ClearDatabase(t)

	t.Run("List Properties - Public Endpoint", func(t *testing.T) {
		// Create test property
		propertyFactory := testutil.NewPropertyFactory(server.Client)
		propertyFactory.Create(t)

		// Make request to list properties
		client := testutil.NewHTTPClient(server.URL())
		resp, err := client.GET("/v1/api/properties", nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		// Assert successful response with at least one property
		testutil.AssertListResponse(t, resp, 1)
	})

	t.Run("Get Property - Valid ID", func(t *testing.T) {
		// Create test property
		propertyFactory := testutil.NewPropertyFactory(server.Client)
		property := propertyFactory.Create(t)

		// Make request to get specific property
		client := testutil.NewHTTPClient(server.URL())
		path := fmt.Sprintf("/v1/api/properties/%s", property.ID)
		resp, err := client.GET(path, nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		// Assert successful response with property data
		apiResp := testutil.AssertContainsData(t, resp, "id", "name", "description", "price")
		
		// Verify property ID matches
		propertyID := testutil.ExtractDataField(t, apiResp, "id")
		if propertyID != property.ID {
			t.Errorf("Expected property ID %s, got %v", property.ID, propertyID)
		}
	})

	t.Run("Get Property - Invalid ID", func(t *testing.T) {
		client := testutil.NewHTTPClient(server.URL())
		resp, err := client.GET("/v1/api/properties/nonexistent-id", nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		// Assert not found response
		testutil.AssertNotFound(t, resp)
	})

	t.Run("Add Property - Business Partner", func(t *testing.T) {
		// Create test user and authenticate
		testUser := testutil.CreateBusinessPartner(t, server.Client)
		client := testUser.AuthorizedClient(server.URL())

		// Create test project
		projectFactory := testutil.NewProjectFactory(server.Client)
		project := projectFactory.Create(t)

		// Prepare property data
		propertyReq := map[string]interface{}{
			"name":          "New Test Property",
			"description":   "A new test property description",
			"property_type": "apartment",
			"price":         1500000,
			"project_id":    project.ID,
		}

		var propertyResp testutil.APIResponse
		resp, err := client.POST("/v1/api/properties", propertyReq, &propertyResp)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		// Assert successful creation
		testutil.AssertSuccessResponse(t, resp)
	})

	t.Run("Add Property - Unauthenticated", func(t *testing.T) {
		client := testutil.NewHTTPClient(server.URL())

		propertyReq := map[string]interface{}{
			"name":        "New Test Property",
			"description": "A new test property description",
			"price":       1500000,
		}

		var propertyResp testutil.APIResponse
		resp, err := client.POST("/v1/api/properties", propertyReq, &propertyResp)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		// Assert unauthorized response (business partner required)
		testutil.AssertUnauthorized(t, resp)
	})

	t.Run("Update Property - Business Partner", func(t *testing.T) {
		// Create test user and authenticate
		testUser := testutil.CreateBusinessPartner(t, server.Client)
		client := testUser.AuthorizedClient(server.URL())

		// Create test property
		propertyFactory := testutil.NewPropertyFactory(server.Client)
		property := propertyFactory.Create(t)

		// Prepare update data
		updateReq := map[string]interface{}{
			"name":        "Updated Property Name",
			"description": "Updated property description",
			"price":       2000000,
		}

		var updateResp testutil.APIResponse
		path := fmt.Sprintf("/v1/api/properties/%s", property.ID)
		resp, err := client.PUT(path, updateReq, &updateResp)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		// Assert successful update
		testutil.AssertSuccessResponse(t, resp)
	})

	t.Run("Delete Property - Business Partner", func(t *testing.T) {
		// Create test user and authenticate
		testUser := testutil.CreateBusinessPartner(t, server.Client)
		client := testUser.AuthorizedClient(server.URL())

		// Create test property
		propertyFactory := testutil.NewPropertyFactory(server.Client)
		property := propertyFactory.Create(t)

		// Make delete request
		var deleteResp testutil.APIResponse
		path := fmt.Sprintf("/v1/api/properties/%s", property.ID)
		resp, err := client.DELETE(path, &deleteResp)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		// Assert successful deletion
		testutil.AssertSuccessResponse(t, resp)

		// Verify property no longer exists
		resp2, err := client.GET(path, nil)
		if err != nil {
			t.Fatalf("Failed to make verification request: %v", err)
		}
		testutil.AssertNotFound(t, resp2)
	})

	t.Run("Get Properties of Project", func(t *testing.T) {
		// Create test project with properties
		projectFactory := testutil.NewProjectFactory(server.Client)
		project := projectFactory.Create(t)

		propertyFactory := testutil.NewPropertyFactory(server.Client)
		propertyFactory.CreateWithProject(t, project)
		propertyFactory.CreateWithProject(t, project)

		// Make request to get properties of project
		client := testutil.NewHTTPClient(server.URL())
		path := fmt.Sprintf("/v1/api/projects/%s/properties", project.ID)
		resp, err := client.GET(path, nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		// Assert successful response with at least 2 properties
		testutil.AssertListResponse(t, resp, 2)
	})

	t.Run("Admin List Properties - Business Partner", func(t *testing.T) {
		// Create test user and authenticate
		testUser := testutil.CreateBusinessPartner(t, server.Client)
		client := testUser.AuthorizedClient(server.URL())

		// Create test property
		propertyFactory := testutil.NewPropertyFactory(server.Client)
		propertyFactory.Create(t)

		// Make request to admin list properties
		resp, err := client.GET("/v1/api/admin/dashboard/properties", nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		// Assert successful response (business partner can access)
		testutil.AssertSuccessResponse(t, resp)
	})

	t.Run("Admin List Properties - Unauthenticated", func(t *testing.T) {
		client := testutil.NewHTTPClient(server.URL())

		// Make request to admin list properties without auth
		resp, err := client.GET("/v1/api/admin/dashboard/properties", nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		// Assert unauthorized response (business partner required)
		testutil.AssertUnauthorized(t, resp)
	})
}