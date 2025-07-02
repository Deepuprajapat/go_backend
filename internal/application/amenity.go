package application

import (
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

// add amemities to category
// delete amenity from category
// delete category with its amenities
// patch update static site data which is active
