package application

import (
	"net/http"
	"strings"

	"github.com/VI-IM/im_backend_go/request"
	"github.com/VI-IM/im_backend_go/response"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/VI-IM/im_backend_go/shared/logger"
)

func (c *application) GetAllCategoriesWithAmenities() (*response.AmenityResponse, *imhttp.CustomError) {
	staticData, err := c.repo.GetStaticSiteData()
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get amenities")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get amenities", err.Error())
	}

	cwa := make(map[string][]response.Amenity)
	for category, amenities := range staticData.CategoriesWithAmenities.Categories {
		amenities := make([]response.Amenity, len(amenities))
		for i, amenity := range amenities {
			amenities[i] = response.Amenity{
				Icon:  amenity.Icon,
				Value: amenity.Value,
			}
		}
		cwa[category] = amenities
	}

	return &response.AmenityResponse{
		Categories: cwa,
	}, nil
}

func (c *application) AddCategoryWithAmenities(req *request.CreateAmenityRequest) *imhttp.CustomError {

	// Check if the amenity already exists
	var categoryName string
	for category := range req.Category {
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

func (c *application) UpdateStaticSiteData(req *request.UpdateStaticSiteDataRequest) (*response.StaticSiteDataResponse, *imhttp.CustomError) {
	// Get current static site data
	staticData, err := c.repo.GetStaticSiteData()
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get static site data")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get static site data", err.Error())
	}

	// Check if static site data is active
	if !staticData.IsActive {
		return nil, imhttp.NewCustomErr(http.StatusForbidden, "Cannot update inactive static site data", "Static site data must be active to update")
	}

	// Update fields if provided
	if req.PropertyTypes != nil {
		staticData.PropertyTypes = *req.PropertyTypes
	}
	if req.CategoriesWithAmenities.Categories != nil {
		staticData.CategoriesWithAmenities.Categories = req.CategoriesWithAmenities.Categories
	}
	staticData.IsActive = req.IsActive

	// Save updates
	err = c.repo.UpdateStaticSiteData(staticData)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to update static site data")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to update static site data", err.Error())
	}

	// Get the updated data and return it
	updatedData, err := c.repo.GetStaticSiteData()
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get updated static site data")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get updated static site data", err.Error())
	}

	return response.GetStaticSiteDataFromEnt(updatedData), nil
}

// patch update static site data which is active
