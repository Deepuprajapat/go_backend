package handlers

import (
	"net/http"

	imhttp "github.com/VI-IM/im_backend_go/shared"
)


func (h *Handler) CheckURLExists(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	// Get URL from query parameter
	url := r.URL.Query().Get("url")
	if url == "" {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "URL parameter is required", "URL parameter is required")
	}

	result, err := h.app.CheckURLExists(r.Context(), url)
	if err != nil {
		return nil, err
	}

	return &imhttp.Response{
		Data:       result,
		StatusCode: http.StatusOK,
	}, nil
}