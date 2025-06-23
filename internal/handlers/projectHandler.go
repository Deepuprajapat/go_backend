package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/VI-IM/im_backend_go/internal/application"
	"github.com/VI-IM/im_backend_go/request"
	imhttp "github.com/VI-IM/im_backend_go/shared"
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
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request body", err.Error())
	}

	if input.ProjectName == "" || input.ProjectURL == "" || input.ProjectType == "" || input.Locality == "" || input.ProjectCity == "" || input.DeveloperID == "" {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request body", "Project name, project URL, project type, locality, project city, and developer ID are required")
	}

	response, err := h.app.AddProject(input)
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to add project", err.Error())
	}

	return &imhttp.Response{
		Data:       response,
		StatusCode: http.StatusOK,
	}, nil
}
