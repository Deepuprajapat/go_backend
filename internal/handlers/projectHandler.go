package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/VI-IM/im_backend_go/request"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/VI-IM/im_backend_go/shared/logger"
	"github.com/gorilla/mux"
)

func (h *Handler) GetProject(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	vars := mux.Vars(r)
	projectID := vars["project_id"]

	response, err := h.app.GetProjectByID(projectID)
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get project", err.Error())
	}

	return &imhttp.Response{
		Data:       response,
		StatusCode: http.StatusOK,
	}, nil
}

func (h *Handler) AddProject(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	var input request.AddProjectRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Get().Error().Msg("Invalid request body")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request body", err.Error())
	}

	if input.ProjectName == "" || input.ProjectURL == "" || input.ProjectType == "" || input.Locality == "" || input.ProjectCity == "" || input.DeveloperID == "" {
		logger.Get().Error().Msg("Invalid request body")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request body", "Project name, project URL, project type, locality, project city, and developer ID are required")
	}

	response, err := h.app.AddProject(input)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to add project")
		return nil, err
	}

	return &imhttp.Response{
		Data:       response,
		StatusCode: http.StatusOK,
	}, nil
}

func (h *Handler) UpdateProject(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	vars := mux.Vars(r)
	projectID := vars["project_id"]
	if projectID == "" {
		logger.Get().Error().Msg("Project ID is required")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Project ID is required", "Project ID is required")
	}

	var input request.UpdateProjectRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Get().Error().Msg("Invalid request body")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request body", err.Error())
	}

	response, err := h.app.UpdateProject(input)
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to update project", err.Error())
	}

	return &imhttp.Response{
		Data:       response,
		StatusCode: http.StatusOK,
	}, nil
}

func (h *Handler) DeleteProject(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	vars := mux.Vars(r)
	projectID := vars["project_id"]

	if err := h.app.DeleteProject(projectID); err != nil {
		return nil, err
	}

	return &imhttp.Response{
		Data:       nil,
		StatusCode: http.StatusOK,
		Message:    "Project deleted successfully",
	}, nil
}

func (h *Handler) ListProjects(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	// Create filter request
	filters := &request.ProjectFilterRequest{}

	// Parse query parameters
	if configurations := r.URL.Query()["configurations"]; len(configurations) > 0 {
		filters.Configurations = configurations
	}
	if isPremium := r.URL.Query().Get("is_premium"); isPremium == "true" {
		filters.IsPremium = true
	}
	if isPriority := r.URL.Query().Get("is_priority"); isPriority == "true" {
		filters.IsPriority = true
	}
	if isFeatured := r.URL.Query().Get("is_featured"); isFeatured == "true" {
		filters.IsFeatured = true
	}
	if locationID := r.URL.Query().Get("location_id"); locationID != "" {
		filters.LocationID = locationID
	}
	if developerID := r.URL.Query().Get("developer_id"); developerID != "" {
		filters.DeveloperID = developerID
	}
	if name := r.URL.Query().Get("name"); name != "" {
		filters.Name = name
	}
	if projectType := r.URL.Query().Get("type"); projectType != "" {
		filters.Type = projectType
	}

	projects, err := h.app.ListProjects(filters.ToMap())
	if err != nil {
		return nil, err
	}

	return &imhttp.Response{
		Data:       projects,
		StatusCode: http.StatusOK,
	}, nil
}
