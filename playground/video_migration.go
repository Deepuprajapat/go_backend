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

type JavaGetProjectByIDResponse struct {
	Content []struct {
		ID     string   `json:"id"`
		Videos []string `json:"videos"`
	} `json:"content"`
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

	projects, err := fetchAllProjectIDs(httpClient)
	if err != nil {
		log.Fatalf("Failed to fetch project IDs: %v", err)
	}

	projectIDToURL := make(map[string][]string)

	for _, project := range projects.Content {
		projectIDToURL[project.ID] = project.Videos
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

func fetchAllProjectIDs(client *http.Client) (*JavaGetProjectByIDResponse, error) {
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

	var projects JavaGetProjectByIDResponse
	if err := json.Unmarshal(body, &projects); err != nil {
		return nil, fmt.Errorf("failed to unmarshal projects: %v", err)
	}

	return &projects, nil
}

func fetchProjectDetails(client *http.Client, projectID string) (*JavaGetProjectByIDResponse, error) {
	resp, err := client.Get(fmt.Sprintf(javaAPIBaseURL+getProjectByIDPath, projectID))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch project details: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var project JavaGetProjectByIDResponse
	if err := json.Unmarshal(body, &project); err != nil {
		return nil, fmt.Errorf("failed to unmarshal project details: %v", err)
	}

	return &project, nil
}

func updateProjectVideoURL(ctx context.Context, client *ent.Client, projectID string, videoURLs []string) error {

	if err := client.Project.UpdateOneID(projectID).
		SetWebCards(schema.ProjectWebCards{
			VideoPresentation: schema.VideoPresentation{
				URLs:        videoURLs,
				Description: "",
			},
		}).
		Exec(ctx); err != nil {
		return fmt.Errorf("failed to update project: %v", err)
	}

	return nil
}
