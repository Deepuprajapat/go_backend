package application

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/VI-IM/im_backend_go/request"
	"github.com/VI-IM/im_backend_go/response"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/VI-IM/im_backend_go/shared/logger"
)

func (c *application) GetAmenities() (*response.AmenityResponse, *imhttp.CustomError) {
	staticData, err := c.repo.GetStaticSiteData()
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get amenities")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get amenities", err.Error())
	}

	amenities := &response.AmenityResponse{
		Categories: make(map[string][]response.Amenity),
	}

	// Convert from static site data format to response format
	for category, amenityList := range staticData.CategoriesWithAmenities.Categories {
		amenities.Categories[category] = make([]response.Amenity, len(amenityList))
		for i, a := range amenityList {
			amenities.Categories[category][i] = response.Amenity{
				Icon:  a.Icon,
				Value: a.Value,
			}
		}
	}

	return amenities, nil
}

func (c *application) GetAmenityByID(id string) (*response.SingleAmenityResponse, *imhttp.CustomError) {
	staticData, err := c.repo.GetStaticSiteData()
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get amenity")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get amenity", err.Error())
	}

	// Search through all categories for the amenity
	for category, amenities := range staticData.CategoriesWithAmenities.Categories {
		for _, amenity := range amenities {
			if amenity.Value == id {
				return &response.SingleAmenityResponse{
					Category: category,
					Icon:     amenity.Icon,
					Value:    amenity.Value,
				}, nil
			}
		}
	}

	return nil, imhttp.NewCustomErr(http.StatusNotFound, "Amenity not found", "Amenity not found")
}

func (c *application) AddCategoryWithAmenities(req *request.CreateAmenityRequest) *imhttp.CustomError {

	// Check if the amenity already exists
	var categoryName string
	for category, _ := range req.Category {
		categoryName = strings.ToLower(category)
	}

	exist, err := c.repo.CheckCategoryExists(categoryName)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get amenity")
		return imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to create amenity", err.Error())
	}
	if exist {
		return imhttp.NewCustomErr(http.StatusConflict, "Amenity already exists", "Amenity already exists")
	}

	staticData, err := c.repo.GetStaticSiteData()
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get static site data")
		return imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to create amenity", err.Error())
	}
	// Add to existing category or create new category

	if staticData.CategoriesWithAmenities.Categories == nil {
		staticData.CategoriesWithAmenities.Categories = make(map[string][]struct {
			Icon  string `json:"icon"`
			Value string `json:"value"`
		})
	}

	for category, amenities := range req.Category {
		for _, amenity := range amenities {
			staticData.CategoriesWithAmenities.Categories[category] = append(staticData.CategoriesWithAmenities.Categories[category], struct {
				Icon  string `json:"icon"`
				Value string `json:"value"`
			}{
				Icon:  amenity.Icon,
				Value: amenity.Value,
			})
		}
	}

	// Update static site data
	if err := c.repo.UpdateStaticSiteData(staticData); err != nil {
		logger.Get().Error().Err(err).Msg("Failed to update static site data")
		return imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to create amenity", err.Error())
	}

	return nil
}

func (c *application) UpdateAmenity(id string, req *request.UpdateAmenityRequest) *imhttp.CustomError {
	staticData, err := c.repo.GetStaticSiteData()
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get static site data")
		return imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to update amenity", err.Error())
	}

	var foundCategory string
	var foundIndex int
	var found bool

	// Find the amenity to update
	for category, amenities := range staticData.CategoriesWithAmenities.Categories {
		for i, amenity := range amenities {
			if amenity.Value == id {
				foundCategory = category
				foundIndex = i
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		return imhttp.NewCustomErr(http.StatusNotFound, "Amenity not found", "Amenity not found")
	}

	// If updating value, check if new value already exists
	if req.Value != "" && req.Value != id {
		for _, amenities := range staticData.CategoriesWithAmenities.Categories {
			for _, amenity := range amenities {
				if amenity.Value == req.Value {
					return imhttp.NewCustomErr(http.StatusConflict, "Amenity with this value already exists", "Duplicate amenity value")
				}
			}
		}
	}

	// Update the amenity
	amenity := staticData.CategoriesWithAmenities.Categories[foundCategory][foundIndex]
	if req.Icon != "" {
		amenity.Icon = req.Icon
	}
	if req.Value != "" {
		amenity.Value = req.Value
	}

	// If category is being updated, move the amenity to the new category
	if req.Category != "" && req.Category != foundCategory {
		// Remove from old category
		oldCategoryAmenities := staticData.CategoriesWithAmenities.Categories[foundCategory]
		staticData.CategoriesWithAmenities.Categories[foundCategory] = append(
			oldCategoryAmenities[:foundIndex],
			oldCategoryAmenities[foundIndex+1:]...,
		)

		// Add to new category
		staticData.CategoriesWithAmenities.Categories[req.Category] = append(
			staticData.CategoriesWithAmenities.Categories[req.Category],
			amenity,
		)

		// Clean up empty category if needed
		if len(staticData.CategoriesWithAmenities.Categories[foundCategory]) == 0 {
			delete(staticData.CategoriesWithAmenities.Categories, foundCategory)
		}
	} else {
		// Update in place
		staticData.CategoriesWithAmenities.Categories[foundCategory][foundIndex] = amenity
	}

	// Update static site data
	if err := c.repo.UpdateStaticSiteData(staticData); err != nil {
		logger.Get().Error().Err(err).Msg("Failed to update static site data")
		return imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to update amenity", err.Error())
	}

	return nil
}

func (c *application) AddAmenitiesToCategory(req *request.AddAmenitiesToCategoryRequest) *imhttp.CustomError {
	// Get the current static site data
	staticData, err := c.repo.GetStaticSiteData()
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get static site data")
		return imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to add amenities", err.Error())
	}

	// Check if the category exists
	exist, err := c.repo.CheckCategoryExists(req.Category)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to check category existence")
		return imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to add amenities", err.Error())
	}
	if !exist {
		return imhttp.NewCustomErr(http.StatusNotFound, "Category not found", "Category not found")
	}

	// Check for duplicate amenity values in the request
	valueMap := make(map[string]bool)
	for _, item := range req.Items {
		if item.Icon == "" {
			return imhttp.NewCustomErr(http.StatusBadRequest, "Icon is required", "Icon field cannot be empty")
		}
		if item.Value == "" {
			return imhttp.NewCustomErr(http.StatusBadRequest, "Value is required", "Value field cannot be empty")
		}
		if valueMap[item.Value] {
			return imhttp.NewCustomErr(http.StatusConflict, "Duplicate amenity values in request", "Each amenity value must be unique")
		}
		valueMap[item.Value] = true
	}

	// Check for duplicate values in existing amenities
	for _, amenity := range staticData.CategoriesWithAmenities.Categories[req.Category] {
		if valueMap[amenity.Value] {
			return imhttp.NewCustomErr(http.StatusConflict, "Amenity already exists in category", "Amenity value must be unique within category")
		}
	}

	// Add new amenities to the category
	for _, item := range req.Items {
		staticData.CategoriesWithAmenities.Categories[req.Category] = append(
			staticData.CategoriesWithAmenities.Categories[req.Category],
			struct {
				Icon  string `json:"icon"`
				Value string `json:"value"`
			}{
				Icon:  item.Icon,
				Value: item.Value,
			},
		)
	}

	// Update static site data
	if err := c.repo.UpdateStaticSiteData(staticData); err != nil {
		logger.Get().Error().Err(err).Msg("Failed to update static site data")
		return imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to add amenities", err.Error())
	}

	return nil
}

func (c *application) DeleteAmenitiesFromCategory(req *request.DeleteAmenitiesFromCategoryRequest) *imhttp.CustomError {
	// Get the current static site data
	staticData, err := c.repo.GetStaticSiteData()
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get static site data")
		return imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to delete amenities", err.Error())
	}

	// Check if the category exists
	exist, err := c.repo.CheckCategoryExists(req.Category)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to check category existence")
		return imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to delete amenities", err.Error())
	}
	if !exist {
		return imhttp.NewCustomErr(http.StatusNotFound, "Category not found", "Category not found")
	}

	// Create a map of values to delete for O(1) lookup
	valuesToDelete := make(map[string]bool)
	for _, value := range req.Values {
		valuesToDelete[value] = true
	}

	// Create a new slice to hold the amenities we want to keep
	var newAmenities []struct {
		Icon  string `json:"icon"`
		Value string `json:"value"`
	}

	// Track which values were actually found and deleted
	deletedValues := make(map[string]bool)

	// Filter out the amenities that should be deleted
	for _, amenity := range staticData.CategoriesWithAmenities.Categories[req.Category] {
		if !valuesToDelete[amenity.Value] {
			newAmenities = append(newAmenities, amenity)
		} else {
			deletedValues[amenity.Value] = true
		}
	}

	// Check if any requested values were not found
	var notFoundValues []string
	for _, value := range req.Values {
		if !deletedValues[value] {
			notFoundValues = append(notFoundValues, value)
		}
	}
	if len(notFoundValues) > 0 {
		return imhttp.NewCustomErr(http.StatusNotFound, "Some amenities not found", fmt.Sprintf("The following amenities were not found: %v", notFoundValues))
	}

	// Update the category with the filtered amenities
	staticData.CategoriesWithAmenities.Categories[req.Category] = newAmenities

	// If the category is now empty, optionally remove it
	if len(newAmenities) == 0 {
		delete(staticData.CategoriesWithAmenities.Categories, req.Category)
	}

	// Update static site data
	if err := c.repo.UpdateStaticSiteData(staticData); err != nil {
		logger.Get().Error().Err(err).Msg("Failed to update static site data")
		return imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to delete amenities", err.Error())
	}

	return nil
}

func (c *application) DeleteCategory(req *request.DeleteCategoryRequest) *imhttp.CustomError {
	// Get the current static site data
	staticData, err := c.repo.GetStaticSiteData()
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get static site data")
		return imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to delete category", err.Error())
	}

	// Check if the category exists
	exist, err := c.repo.CheckCategoryExists(req.Category)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to check category existence")
		return imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to delete category", err.Error())
	}
	if !exist {
		return imhttp.NewCustomErr(http.StatusNotFound, "Category not found", "Category not found")
	}

	// Delete the category and all its amenities
	delete(staticData.CategoriesWithAmenities.Categories, req.Category)

	// Update static site data
	if err := c.repo.UpdateStaticSiteData(staticData); err != nil {
		logger.Get().Error().Err(err).Msg("Failed to update static site data")
		return imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to delete category", err.Error())
	}

	return nil
}

// patch update static site data which is active
