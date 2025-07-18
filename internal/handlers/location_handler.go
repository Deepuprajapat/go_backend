package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/VI-IM/im_backend_go/request"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/VI-IM/im_backend_go/shared/logger"
	"github.com/gorilla/mux"
)

func (h *Handler) ListLocations(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {

	filters := make(map[string]interface{})
	if city := r.URL.Query().Get("city"); city != "" {
		filters["city"] = city
	}

	response, err := h.app.GetAllLocations(filters)
	if err != nil {
		return nil, err
	}

	return &imhttp.Response{
		Data:       response,
		StatusCode: http.StatusOK,
	}, nil
}

func (h *Handler) GetLocation(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
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

func (h *Handler) AddLocation(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	var input request.AddLocationRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Get().Error().Msg("Invalid request body")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate required fields
	if err := h.validate.Struct(input); err != nil {
		logger.Get().Error().Err(err).Msg("Validation failed")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Validation failed", err.Error())
	}

	location, err := h.app.AddLocation(input)
	if err != nil {
		return nil, err
	}

	return &imhttp.Response{
		Data:       location,
		StatusCode: http.StatusCreated,
	}, nil
}

func (h *Handler) DeleteLocation(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
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
