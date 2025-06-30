package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/VI-IM/im_backend_go/request"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/gorilla/mux"
)

func (h *Handler) GetAmenities(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	amenities, err := h.app.GetAmenities()
	if err != nil {
		return nil, err
	}

	return &imhttp.Response{
		Data:       amenities,
		StatusCode: http.StatusOK,
	}, nil
}

func (h *Handler) GetAmenity(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	vars := mux.Vars(r)
	amenityID := vars["amenity_id"]
	if amenityID == "" {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Amenity ID is required", "Amenity ID is required")
	}

	amenity, err := h.app.GetAmenityByID(amenityID)
	if err != nil {
		return nil, err
	}

	return &imhttp.Response{
		Data:       amenity,
		StatusCode: http.StatusOK,
	}, nil
}

func (h *Handler) CreateAmenity(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	var req request.CreateAmenityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Manual validation
	if req.Category == "" {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Category is required", "Category field cannot be empty")
	}
	if req.Icon == "" {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Icon is required", "Icon field cannot be empty")
	}
	if req.Value == "" {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Value is required", "Value field cannot be empty")
	}

	if err := h.app.CreateAmenity(&req); err != nil {
		return nil, err
	}

	return &imhttp.Response{
		StatusCode: http.StatusCreated,
		Message:    "Amenity created successfully",
	}, nil
}

func (h *Handler) UpdateAmenity(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	vars := mux.Vars(r)
	amenityID := vars["amenity_id"]
	if amenityID == "" {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Amenity ID is required", "Amenity ID is required")
	}

	var req request.UpdateAmenityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate that at least one field is being updated
	if req.Category == "" && req.Icon == "" && req.Value == "" {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "At least one field must be provided for update", "No update fields provided")
	}

	if err := h.app.UpdateAmenity(amenityID, &req); err != nil {
		return nil, err
	}

	return &imhttp.Response{
		StatusCode: http.StatusOK,
		Message:    "Amenity updated successfully",
	}, nil
}
