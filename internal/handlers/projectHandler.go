package handlers

import (
	"net/http"
	"strconv"

	"github.com/VI-IM/im_backend_go/internal/controller"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/gorilla/mux"
)

type ProjectHandler struct {
	ctrl controller.ControllerInterface
}

func NewProjectHandler(ctrl controller.ControllerInterface) *ProjectHandler {
	return &ProjectHandler{ctrl: ctrl}
}

func (h *ProjectHandler) GetProject(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["project_id"])
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid project ID", err.Error())
	}

	response, err := h.ctrl.GetProjectByID(projectID)
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get project", err.Error())
	}

	return &imhttp.Response{
		Data:       response,
		StatusCode: http.StatusOK,
	}, nil
}
