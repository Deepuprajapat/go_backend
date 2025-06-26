package handlers

import (
	"net/http"

	"github.com/VI-IM/im_backend_go/internal/application"
	imhttp "github.com/VI-IM/im_backend_go/shared"
)

type LocationHandler struct {
	app application.ApplicationInterface
}

func NewLocationHandler(app application.ApplicationInterface) *LocationHandler {
	return &LocationHandler{app: app}
}

func (h *LocationHandler) ListLocations(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	response, err := h.app.GetAllLocations()
	if err != nil {
		return nil, err
	}

	return &imhttp.Response{
		Data:       response,
		StatusCode: http.StatusOK,
	}, nil
}
