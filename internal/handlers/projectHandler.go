package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/VI-IM/im_backend_go/internal/application"
	"github.com/VI-IM/im_backend_go/request"
	"github.com/VI-IM/im_backend_go/response"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/VI-IM/im_backend_go/shared/logger"
	"github.com/gorilla/mux"
)

type ProjectHandler struct {
	app application.ApplicationInterface
}

func NewProjectHandler(app application.ApplicationInterface) *ProjectHandler {
	return &ProjectHandler{app: app}
}

func (h *ProjectHandler) GetProject(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
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

func (h *ProjectHandler) AddProject(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
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

func (h *ProjectHandler) UpdateProject(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
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

func (h *ProjectHandler) DeleteProject(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
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

func (h *ProjectHandler) ListProjects(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	// Parse pagination parameters from query
	pagination := &request.PaginationRequest{
		Page:     1,
		PageSize: 10,
	}

	if page := r.URL.Query().Get("page"); page != "" {
		if pageNum, err := strconv.Atoi(page); err == nil {
			pagination.Page = pageNum
		}
	}

	if pageSize := r.URL.Query().Get("page_size"); pageSize != "" {
		if pageSizeNum, err := strconv.Atoi(pageSize); err == nil {
			pagination.PageSize = pageSizeNum
		}
	}

	pagination.Validate()

	projects, totalItems, err := h.app.ListProjects(pagination)
	if err != nil {
		return nil, err
	}

	paginatedResponse := response.NewPaginatedResponse(
		projects,
		pagination.Page,
		pagination.PageSize,
		totalItems,
	)

	return &imhttp.Response{
		Data:       paginatedResponse,
		StatusCode: http.StatusOK,
	}, nil
}
