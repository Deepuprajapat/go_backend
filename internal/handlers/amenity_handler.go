package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/VI-IM/im_backend_go/request"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/gorilla/mux"
)

func (h *Handler) GetAllCategoriesWithAmenities(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {

	cwa, err := h.app.GetAllCategoriesWithAmenities()
	if err != nil {
		return nil, err
	}

	return &imhttp.Response{
		Data:       cwa,
		StatusCode: http.StatusOK,
	}, nil
}

func (h *Handler) AddCategoryWithAmenities(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	var req request.CreateAmenityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Manual validation
	for _, category := range req.Category {
		if category == nil {
			return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Amenities are required", "Amenities field cannot be empty")
		}
		for _, amenity := range category {
			if amenity.Icon == "" {
				return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Icon is required", "Icon field cannot be empty")
			}
			if amenity.Value == "" {
				return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Value is required", "Value field cannot be empty")
			}
		}
	}

	if err := h.app.AddCategoryWithAmenities(&req); err != nil {
		return nil, err
	}

	return &imhttp.Response{
		StatusCode: http.StatusCreated,
		Message:    "Amenity created successfully",
	}, nil
}

func (h *Handler) UpdateStaticSiteData(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	var req request.UpdateStaticSiteDataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := h.validate.Struct(req); err != nil {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request", err.Error())
	}

	staticSiteData, err := h.app.UpdateStaticSiteData(&req)
	if err != nil {
		return nil, err
	}

	return &imhttp.Response{
		Data:       staticSiteData,
		StatusCode: http.StatusOK,
		Message:    "Static site data updated successfully",
	}, nil
}

func (h *Handler) AddCategory(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	var req struct {
		CategoryName string `json:"category_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request body", err.Error())
	}

	if req.CategoryName == "" {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Category name is required", "Category name field cannot be empty")
	}

	if err := h.app.AddCategory(req.CategoryName); err != nil {
		return nil, err
	}

	return &imhttp.Response{
		StatusCode: http.StatusCreated,
		Message:    "Category created successfully",
	}, nil
}

func (h *Handler) AddAmenityToCategory(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {

	var req request.AddAmenityToCategoryRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request body", err.Error())
	}

	categoryName := mux.Vars(r)["category_name"]

	if categoryName == "" {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Category name is required", "Category name field cannot be empty")
	}

	req.CategoryName = categoryName

	if len(req.Amenities) == 0 {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Amenities are required", "Amenities field cannot be empty")
	}

	if err := h.app.AddAmenityToCategory(&req); err != nil {
		return nil, err
	}

	return &imhttp.Response{
		StatusCode: http.StatusCreated,
		Message:    "Amenity added to category successfully",
	}, nil
}

func (h *Handler) DeleteAmenityFromCategory(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	var req request.DeleteAmenityFromCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request body", err.Error())
	}

	categoryName := mux.Vars(r)["category_name"]

	if categoryName == "" {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Category name is required", "Category name field cannot be empty")
	}

	amenityName := mux.Vars(r)["amenity_name"]

	if amenityName == "" {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Amenity name is required", "Amenity name field cannot be empty")
	}

	req.CategoryName = categoryName
	req.AmenityName = amenityName

	if req.AmenityName == "" {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Amenity name is required", "Amenity name field cannot be empty")
	}

	if err := h.app.DeleteAmenityFromCategory(&req); err != nil {
		return nil, err
	}

	return &imhttp.Response{
		StatusCode: http.StatusOK,
		Message:    "Amenity deleted from category successfully",
	}, nil
}

func (h *Handler) DeleteCategoryWithAmenities(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	var req struct {
		CategoryName string `json:"category_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request body", err.Error())
	}

	categoryName := mux.Vars(r)["category_name"]

	if categoryName == "" {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Category name is required", "Category name field cannot be empty")
	}

	if err := h.app.DeleteCategoryWithAmenities(req.CategoryName); err != nil {
		return nil, err
	}

	return &imhttp.Response{
		StatusCode: http.StatusOK,
		Message:    "Category deleted successfully",
	}, nil
}
