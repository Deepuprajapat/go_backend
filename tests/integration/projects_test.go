package integration

import (
	"fmt"
	"testing"

	"github.com/VI-IM/im_backend_go/internal/testutil"
)

func TestProjectEndpoints(t *testing.T) {
	// Setup test server
	server := testutil.NewTestServer(t)
	defer server.Close()

	// Clear database before tests
	server.ClearDatabase(t)

	t.Run("List Projects - Public Endpoint", func(t *testing.T) {
		// Create test data
		locationFactory := testutil.NewLocationFactory(server.Client)
		location := locationFactory.Create(t)

		developerFactory := testutil.NewDeveloperFactory(server.Client)
		developer := developerFactory.Create(t)

		projectFactory := testutil.NewProjectFactory(server.Client)
		projectFactory.CreateWithDependencies(t, location, developer)

		// Make request to list projects
		client := testutil.NewHTTPClient(server.URL())
		resp, err := client.GET("/v1/api/projects", nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		// Assert successful response with at least one project
		testutil.AssertListResponse(t, resp, 1)
	})

	t.Run("Get Project - Valid ID", func(t *testing.T) {
		// Create test project
		projectFactory := testutil.NewProjectFactory(server.Client)
		project := projectFactory.Create(t)

		// Make request to get specific project
		client := testutil.NewHTTPClient(server.URL())
		path := fmt.Sprintf("/v1/api/projects/%s", project.ID)
		resp, err := client.GET(path, nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		// Assert successful response with project data
		apiResp := testutil.AssertContainsData(t, resp, "id", "name", "description")
		
		// Verify project ID matches
		projectID := testutil.ExtractDataField(t, apiResp, "id")
		if projectID != project.ID {
			t.Errorf("Expected project ID %s, got %v", project.ID, projectID)
		}
	})

	t.Run("Get Project - Invalid ID", func(t *testing.T) {
		client := testutil.NewHTTPClient(server.URL())
		resp, err := client.GET("/v1/api/projects/nonexistent-id", nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		// Assert not found response
		testutil.AssertNotFound(t, resp)
	})

	t.Run("Add Project - Authenticated User", func(t *testing.T) {
		// Create test user and authenticate
		testUser := testutil.CreateBusinessPartner(t, server.Client)
		client := testUser.AuthorizedClient(server.URL())

		// Create dependencies
		locationFactory := testutil.NewLocationFactory(server.Client)
		location := locationFactory.Create(t)

		developerFactory := testutil.NewDeveloperFactory(server.Client)
		developer := developerFactory.Create(t)

		// Prepare project data
		projectReq := map[string]interface{}{
			"name":           "New Test Project",
			"description":    "A new test project description",
			"project_status": "ongoing",
			"location_id":    location.ID,
			"developer_id":   developer.ID,
		}

		var projectResp testutil.APIResponse
		resp, err := client.POST("/v1/api/projects", projectReq, &projectResp)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		// Assert successful creation
		testutil.AssertSuccessResponse(t, resp)
	})

	t.Run("Add Project - Unauthenticated", func(t *testing.T) {
		client := testutil.NewHTTPClient(server.URL())

		projectReq := map[string]interface{}{
			"name":        "New Test Project",
			"description": "A new test project description",
		}

		var projectResp testutil.APIResponse
		resp, err := client.POST("/v1/api/projects", projectReq, &projectResp)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		// Assert unauthorized response
		testutil.AssertUnauthorized(t, resp)
	})

	t.Run("Update Project - Authenticated User", func(t *testing.T) {
		// Create test user and authenticate
		testUser := testutil.CreateBusinessPartner(t, server.Client)
		client := testUser.AuthorizedClient(server.URL())

		// Create test project
		projectFactory := testutil.NewProjectFactory(server.Client)
		project := projectFactory.Create(t)

		// Prepare update data
		updateReq := map[string]interface{}{
			"name":        "Updated Project Name",
			"description": "Updated project description",
		}

		var updateResp testutil.APIResponse
		path := fmt.Sprintf("/v1/api/projects/%s", project.ID)
		resp, err := client.PUT(path, updateReq, &updateResp)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		// Assert successful update
		testutil.AssertSuccessResponse(t, resp)
	})

	t.Run("Delete Project - Authenticated User", func(t *testing.T) {
		// Create test user and authenticate
		testUser := testutil.CreateBusinessPartner(t, server.Client)
		client := testUser.AuthorizedClient(server.URL())

		// Create test project
		projectFactory := testutil.NewProjectFactory(server.Client)
		project := projectFactory.Create(t)

		// Make delete request
		var deleteResp testutil.APIResponse
		path := fmt.Sprintf("/v1/api/projects/%s", project.ID)
		resp, err := client.DELETE(path, &deleteResp)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		// Assert successful deletion
		testutil.AssertSuccessResponse(t, resp)

		// Verify project no longer exists
		resp2, err := client.GET(path, nil)
		if err != nil {
			t.Fatalf("Failed to make verification request: %v", err)
		}
		testutil.AssertNotFound(t, resp2)
	})

	t.Run("Compare Projects", func(t *testing.T) {
		// Create test projects
		projectFactory := testutil.NewProjectFactory(server.Client)
		project1 := projectFactory.Create(t)
		project2 := projectFactory.Create(t)

		// Prepare compare request
		compareReq := map[string]interface{}{
			"project_ids": []string{project1.ID, project2.ID},
		}

		client := testutil.NewHTTPClient(server.URL())
		var compareResp testutil.APIResponse
		resp, err := client.POST("/v1/api/projects/compare", compareReq, &compareResp)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		// Assert successful comparison
		testutil.AssertSuccessResponse(t, resp)
	})
}