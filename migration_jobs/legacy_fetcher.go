package migration_jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/VI-IM/im_backend_go/shared/logger"
	"github.com/rs/zerolog/log"
)

var projectIDToAmenitiesMap = map[int64][]int64{}

// fetch all tables from legacy database - now using JSON data loader

func FetchCityByID(ctx context.Context, id int64) (LCity, error) {
	// Use JSON data loader instead of database query
	loader := GetJSONDataLoader()
	if loader == nil {
		return LCity{}, fmt.Errorf("JSON data loader not initialized")
	}

	return loader.GetCityByID(id)
}

func FetchAllLocality(ctx context.Context) ([]LLocality, error) {
	// Use JSON data loader instead of database query
	loader := GetJSONDataLoader()
	if loader == nil {
		return nil, fmt.Errorf("JSON data loader not initialized")
	}

	return loader.GetAllLocalities()
}

func FetchAllDevelopers(ctx context.Context) ([]LDeveloper, error) {
	log.Info().Msg("Fetching all developers from JSON data")

	// Use JSON data loader instead of database query
	loader := GetJSONDataLoader()
	if loader == nil {
		return nil, fmt.Errorf("JSON data loader not initialized")
	}

	return loader.GetAllDevelopers()
}

func FetchPropertyConfigurationByID(ctx context.Context, id int64) (*LPropertyConfiguration, error) {
	// Use JSON data loader instead of database query
	loader := GetJSONDataLoader()
	if loader == nil {
		return nil, fmt.Errorf("JSON data loader not initialized")
	}

	return loader.GetPropertyConfigurationByID(id)
}

func FetchhAllProject(ctx context.Context) ([]LProject, error) {
	fmt.Println("Fetching all projects from JSON data")

	// Use JSON data loader instead of database query
	loader := GetJSONDataLoader()
	if loader == nil {
		return nil, fmt.Errorf("JSON data loader not initialized")
	}

	return loader.GetAllProjects()
}

func fetchAllProperty(ctx context.Context) ([]LProperty, error) {
	// Use JSON data loader instead of database query
	loader := GetJSONDataLoader()
	if loader == nil {
		return nil, fmt.Errorf("JSON data loader not initialized")
	}

	return loader.GetAllProperties()
}

func FetchProjectByID(ctx context.Context, id int64) (*LProject, error) {
	// Use JSON data loader instead of database query
	loader := GetJSONDataLoader()
	if loader == nil {
		return nil, fmt.Errorf("JSON data loader not initialized")
	}

	return loader.GetProjectByID(id)
}

func FetchProjectConfigurationsByID(ctx context.Context, id int64) (*LPropertyConfiguration, error) {
	// Use JSON data loader instead of database query
	loader := GetJSONDataLoader()
	if loader == nil {
		return nil, fmt.Errorf("JSON data loader not initialized")
	}

	return loader.GetPropertyConfigurationByID(id)
}

func FetchLocalityByID(ctx context.Context, id int64) (*LLocality, error) {
	// Use JSON data loader instead of database query
	loader := GetJSONDataLoader()
	if loader == nil {
		return nil, fmt.Errorf("JSON data loader not initialized")
	}

	return loader.GetLocalityByID(id)
}

func FetchDeveloperByID(ctx context.Context, id int64) (*LDeveloper, error) {
	// Use JSON data loader instead of database query
	loader := GetJSONDataLoader()
	if loader == nil {
		return nil, fmt.Errorf("JSON data loader not initialized")
	}

	return loader.GetDeveloperByID(id)
}

func FetchProjectImagesByProjectID(ctx context.Context, id int64) (*[]LProjectImage, error) {
	// Use JSON data loader instead of database query
	loader := GetJSONDataLoader()
	if loader == nil {
		return nil, fmt.Errorf("JSON data loader not initialized")
	}

	return loader.GetProjectImagesByProjectID(id)
}

func FetchFloorPlansByProjectID(ctx context.Context, projectID int64) (*[]LFloorPlan, error) {
	// Use JSON data loader instead of database query
	loader := GetJSONDataLoader()
	if loader == nil {
		return nil, fmt.Errorf("JSON data loader not initialized")
	}

	return loader.GetFloorPlansByProjectID(projectID)
}

func FetchReraByProjectID(ctx context.Context, projectID int64) ([]*LRera, error) {
	// Use JSON data loader instead of database query
	loader := GetJSONDataLoader()
	if loader == nil {
		return nil, fmt.Errorf("JSON data loader not initialized")
	}

	return loader.GetRerasByProjectID(projectID)
}

func FetchFloorPlanByProjectID(ctx context.Context, projectID int64) (*[]LFloorPlan, error) {
	// Use JSON data loader instead of database query
	loader := GetJSONDataLoader()
	if loader == nil {
		return nil, fmt.Errorf("JSON data loader not initialized")
	}

	return loader.GetFloorPlansByProjectID(projectID)
}

type LProjectAmenity struct {
	ProjectID int64 `json:"project_id"`
	AmenityID int64 `json:"amenity_id"`
}

func FetchAmenityByID(ctx context.Context, id int64) (*LAmenity, error) {
	// Use JSON data loader instead of database query
	loader := GetJSONDataLoader()
	if loader == nil {
		return nil, fmt.Errorf("JSON data loader not initialized")
	}

	return loader.GetAmenityByID(id)
}

func FetchProjectAmenitiesByProjectID(ctx context.Context, projectID int64) ([]*LAmenity, error) {
	// Use JSON data loader instead of database query
	loader := GetJSONDataLoader()
	if loader == nil {
		return nil, fmt.Errorf("JSON data loader not initialized")
	}

	return loader.GetProjectAmenitiesByProjectID(projectID)
}

func FetchPaymentPlansByProjectID(ctx context.Context, projectID int64) ([]*LPaymentPlan, error) {
	// Use JSON data loader instead of database query
	loader := GetJSONDataLoader()
	if loader == nil {
		return nil, fmt.Errorf("JSON data loader not initialized")
	}

	return loader.GetPaymentPlansByProjectID(projectID)
}

func FetchFaqsByProjectID(ctx context.Context, projectID int64) ([]*LFAQ, error) {
	// Use JSON data loader instead of database query
	loader := GetJSONDataLoader()
	if loader == nil {
		return nil, fmt.Errorf("JSON data loader not initialized")
	}

	return loader.GetFaqsByProjectID(projectID)
}

func FetchPropertyConfigurationTypeByID(ctx context.Context, id int64) (*LPropertyConfigurationType, error) {
	// Use JSON data loader instead of database query
	loader := GetJSONDataLoader()
	if loader == nil {
		return nil, fmt.Errorf("JSON data loader not initialized")
	}

	return loader.GetPropertyConfigurationTypeByID(id)
}

func FetchProjectConfigurationByID(ctx context.Context, id int64) (*LPropertyConfigurationType, error) {
	// Use JSON data loader instead of database query
	loader := GetJSONDataLoader()
	if loader == nil {
		return nil, fmt.Errorf("JSON data loader not initialized")
	}

	return loader.GetPropertyConfigurationTypeByID(id)
}

func FetchAllBlogs(ctx context.Context) ([]LBlog, error) {
	// Use JSON data loader instead of database query
	loader := GetJSONDataLoader()
	if loader == nil {
		return nil, fmt.Errorf("JSON data loader not initialized")
	}

	return loader.GetAllBlogs()
}

func fetchAllProjectIDs(client *http.Client) (*JavaGetProjectByIDResponse, error) {
	// This function uses external API, so it remains unchanged
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

func FetchAllGenericSearchData(ctx context.Context) ([]LGenericSearchData, error) {
	// Use JSON data loader instead of database query
	loader := GetJSONDataLoader()
	if loader == nil {
		return nil, fmt.Errorf("JSON data loader not initialized")
	}

	return loader.GetAllGenericSearchData()
}
