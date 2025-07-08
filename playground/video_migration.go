package playground

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/schema"
	"github.com/VI-IM/im_backend_go/internal/config"
	"github.com/VI-IM/im_backend_go/internal/database"
	"github.com/VI-IM/im_backend_go/shared/logger"
)

// JavaProject represents the minimal project structure we need from Java API
type JavaProject struct {
	ID        string `json:"id"`
	VideoURL  string `json:"videoUrl"`
	ProjectID string `json:"projectId"`
}

const (
	// Update these with actual Java API endpoints
	javaAPIBaseURL     = "https://api.investmango.com"
	getAllProjectsPath = "/project/get/all"
	getProjectByIDPath = "/project/get/by/id/%s" // %s will be replaced with project ID
)

var err error

func MigrateVideoURLs() {
	// Initialize configuration and database connection
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	cfg := config.GetConfig()
	client := database.NewClient(cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connection succesfull")

	// Create HTTP client
	httpClient := &http.Client{}

	projectIDs, err := fetchAllProjectIDs(httpClient)

	if err != nil {
		log.Fatalf("Failed to fetch project IDs: %v", err)
	}

	projectIDToURL := make(map[string]string)

	for _, projectID := range projectIDs {
		project, err := fetchProjectDetails(httpClient, projectID)
		if err != nil {
			log.Printf("Warning: Failed to fetch details for project %s: %v", projectID, err)
			continue
		}
		if project.VideoURL != "" {
			projectIDToURL[project.ProjectID] = project.VideoURL
		}
	}

	ctx := context.Background()
	for projectID, videoURL := range projectIDToURL {
		err := updateProjectVideoURL(ctx, client, projectID, videoURL)
		if err != nil {
			log.Printf("Warning: Failed to update video URL for project %s: %v", projectID, err)
			continue
		}
		log.Printf("Successfully updated video URL for project %s", projectID)
	}
}

func fetchAllProjectIDs(client *http.Client) ([]string, error) {
	resp, err := client.Get(javaAPIBaseURL + getAllProjectsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch projects: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	logger.Get().Info().Msg(string(body))

	var projects []JavaProject
	if err := json.Unmarshal(body, &projects); err != nil {
		return nil, fmt.Errorf("failed to unmarshal projects: %v", err)
	}

	var projectIDs []string
	for _, project := range projects {
		projectIDs = append(projectIDs, project.ID)
	}
	return projectIDs, nil
}

func fetchProjectDetails(client *http.Client, projectID string) (*JavaProject, error) {
	resp, err := client.Get(fmt.Sprintf(javaAPIBaseURL+getProjectByIDPath, projectID))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch project details: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var project JavaProject
	if err := json.Unmarshal(body, &project); err != nil {
		return nil, fmt.Errorf("failed to unmarshal project details: %v", err)
	}

	return &project, nil
}

func updateProjectVideoURL(ctx context.Context, client *ent.Client, projectID, videoURL string) error {
	project, err := client.Project.Get(ctx, projectID)
	if err != nil {
		return fmt.Errorf("failed to get project: %v", err)
	}

	// Get existing web cards
	webCards := project.WebCards

	// Update video presentation URL
	webCards.VideoPresentation = schema.VideoPresentation{
		URL:         videoURL,
		Description: webCards.VideoPresentation.Description, // Preserve existing description
	}

	// Update the project with new web cards
	_, err = client.Project.UpdateOne(project).
		SetWebCards(webCards).
		Save(ctx)

	if err != nil {
		return fmt.Errorf("failed to update project: %v", err)
	}

	return nil
}
