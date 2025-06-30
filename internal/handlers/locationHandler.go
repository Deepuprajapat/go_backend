package handlers

import (
	"net/http"

	"github.com/VI-IM/im_backend_go/internal/application"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/gorilla/mux"
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

func (h *LocationHandler) GetLocation(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	vars := mux.Vars(r)
	locationID := vars["location_id"]
	if locationID == "" {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Location ID is required", "Location ID is required")
	}

	location, err := h.app.GetLocationByID(locationID)
	if err != nil {
		return nil, err
	}

	return &imhttp.Response{
		Data:       location,
		StatusCode: http.StatusOK,
	}, nil
}

func (h *LocationHandler) DeleteLocation(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	vars := mux.Vars(r)
	locationID := vars["location_id"]
	if locationID == "" {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Location ID is required", "Location ID is required")
	}

	if err := h.app.DeleteLocation(locationID); err != nil {
		return nil, err
	}

	return &imhttp.Response{
		StatusCode: http.StatusOK,
		Message:    "Location deleted successfully",
	}, nil
}
