package handlers

import (
	"net/http"

	"github.com/VI-IM/im_backend_go/request"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/gorilla/mux"
)

func (h *Handler) ListDevelopers(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	developers, err := h.app.ListDevelopers(&request.GetAllAPIRequest{})
	if err != nil {
		return nil, err
	}

	return &imhttp.Response{
		Data:       developers,
		StatusCode: http.StatusOK,
	}, nil
}

func (h *Handler) GetDeveloper(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	vars := mux.Vars(r)
	developerID := vars["developer_id"]
	if developerID == "" {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Developer ID is required", "Developer ID is required")
	}

	developer, err := h.app.GetDeveloperByID(developerID)
	if err != nil {
		return nil, err
	}

	return &imhttp.Response{
		Data:       developer,
		StatusCode: http.StatusOK,
	}, nil
}

func (h *Handler) DeleteDeveloper(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	vars := mux.Vars(r)
	developerID := vars["developer_id"]
	if developerID == "" {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Developer ID is required", "Developer ID is required")
	}

	if err := h.app.DeleteDeveloper(developerID); err != nil {
		return nil, err
	}

	return &imhttp.Response{
		StatusCode: http.StatusOK,
		Message:    "Developer deleted successfully",
	}, nil
}
