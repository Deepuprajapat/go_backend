package handlers

import (
	"net/http"
	"strconv"

	"github.com/VI-IM/im_backend_go/internal/application"
	"github.com/VI-IM/im_backend_go/request"
	"github.com/VI-IM/im_backend_go/response"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/gorilla/mux"
)

type DeveloperHandler struct {
	app application.ApplicationInterface
}

func NewDeveloperHandler(app application.ApplicationInterface) *DeveloperHandler {
	return &DeveloperHandler{app: app}
}

func (h *DeveloperHandler) ListDevelopers(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
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

	developers, totalItems, err := h.app.ListDevelopers(pagination)
	if err != nil {
		return nil, err
	}

	paginatedResponse := response.NewPaginatedResponse(
		developers,
		pagination.Page,
		pagination.PageSize,
		totalItems,
	)

	return &imhttp.Response{
		Data:       paginatedResponse,
		StatusCode: http.StatusOK,
	}, nil
}

func (h *DeveloperHandler) GetDeveloper(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
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
