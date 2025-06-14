package handlers

import (
	"net/http"
	"strconv"

	"github.com/VI-IM/im_backend_go/internal/application"
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
	projectID, err := strconv.Atoi(vars["project_id"])
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid project ID", err.Error())
	}

	response, err := h.app.GetProjectByID(projectID)
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get project", err.Error())
	}

	return &imhttp.Response{
		Data:       response,
		StatusCode: http.StatusOK,
	}, nil
}
