package application

import (
	"context"
	"net/http"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/schema"
	"github.com/VI-IM/im_backend_go/internal/domain"
	"github.com/VI-IM/im_backend_go/request"
	"github.com/VI-IM/im_backend_go/response"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/VI-IM/im_backend_go/shared/logger"
)

func (c *application) GetPropertyByID(id string) (*response.Property, *imhttp.CustomError) {
	property, err := c.repo.GetPropertyByID(id)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get property")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get property", err.Error())
	}

	return response.GetPropertyFromEnt(property), nil
}

func (c *application) GetPropertyBySlug(ctx context.Context, slug string) (*response.Property, *imhttp.CustomError) {
	property, err := c.repo.GetPropertyBySlug(ctx, slug)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get property by slug")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get property", err.Error())
	}

	return response.GetPropertyFromEnt(property), nil
}

func (c *application) UpdateProperty(input request.UpdatePropertyRequest) (*response.Property, *imhttp.CustomError) {
	existingProperty, err := c.repo.GetPropertyByID(input.PropertyID)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get property")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get property", err.Error())
	}
	if existingProperty == nil {
		return nil, imhttp.NewCustomErr(http.StatusNotFound, "Property not found", "Property not found")
	}

	if input.DeveloperID != "" {
		exists, err := c.repo.ExistDeveloperByID(input.DeveloperID)
		if err != nil {
			logger.Get().Error().Err(err).Msg("Failed to check developer")
			return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to check developer", err.Error())
		}
		if !exists {
			return nil, imhttp.NewCustomErr(http.StatusNotFound, "Developer not found", "Developer not found")
		}
	}

	// Start with existing property data to preserve current values
	property := domain.Property{
		PropertyID:       existingProperty.ID,
		Name:             existingProperty.Name,
		PropertyType:     existingProperty.PropertyType,
		PropertyImages:   existingProperty.PropertyImages,
		WebCards:         existingProperty.WebCards,
		PricingInfo:      existingProperty.PricingInfo,
		PropertyReraInfo: existingProperty.PropertyReraInfo,
		MetaInfo:         existingProperty.MetaInfo,
		IsFeatured:       existingProperty.IsFeatured,
		IsDeleted:        existingProperty.IsDeleted,
		DeveloperID:      existingProperty.DeveloperID,
		LocationID:       existingProperty.LocationID,
		ProjectID:        existingProperty.ProjectID,
	}

	// Selectively update only the fields that are provided and non-empty
	if input.Name != "" {
		property.Name = input.Name
	}
	if len(input.PropertyImages) > 0 {
		property.PropertyImages = input.PropertyImages
	}
	// Only update PricingInfo if it contains actual data (non-empty Price field)
	if input.PricingInfo.Price != "" {
		property.PricingInfo = input.PricingInfo
	}
	// Only update WebCards if it's not empty (check if any field is provided)
	if !isWebCardsEmpty(input.WebCards) {
		property.WebCards = input.WebCards
	}
	// Only update PropertyReraInfo if it contains actual data
	if input.PropertyReraInfo.ReraNumber != "" {
		property.PropertyReraInfo = input.PropertyReraInfo
	}
	// Only update MetaInfo if it contains actual data
	if !isMetaInfoEmpty(input.MetaInfo) {
		property.MetaInfo = input.MetaInfo
	}
	if input.DeveloperID != "" {
		property.DeveloperID = input.DeveloperID
	}
	if input.LocationID != "" {
		property.LocationID = input.LocationID
	}
	if input.ProjectID != "" {
		property.ProjectID = input.ProjectID
	}
	// Note: IsFeatured and IsDeleted are boolean fields that can be explicitly set to false,
	// so we need a different approach. For now, we'll always update them as they might be intentional changes.
	property.IsFeatured = input.IsFeatured
	property.IsDeleted = input.IsDeleted

	updatedProperty, err := c.repo.UpdateProperty(property)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to update property")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to update property", err.Error())
	}

	return response.GetPropertyFromEnt(updatedProperty), nil
}

// Helper function to check if WebCards contains any meaningful data
func isWebCardsEmpty(webCards schema.WebCards) bool {
	return webCards.PropertyDetails.BuiltUpArea.Value == "" &&
		webCards.PropertyDetails.Sizes.Value == "" &&
		webCards.PropertyDetails.FloorNumber.Value == "" &&
		webCards.PropertyDetails.Configuration.Value == "" &&
		webCards.PropertyDetails.PossessionStatus.Value == "" &&
		webCards.PropertyDetails.Balconies.Value == "" &&
		webCards.PropertyDetails.CoveredParking.Value == "" &&
		webCards.PropertyDetails.Bedrooms.Value == "" &&
		webCards.PropertyDetails.PropertyType.Value == "" &&
		webCards.PropertyDetails.AgeOfProperty.Value == "" &&
		webCards.PropertyDetails.FurnishingType.Value == "" &&
		webCards.PropertyDetails.Facing.Value == "" &&
		webCards.PropertyDetails.ReraNumber.Value == "" &&
		webCards.PropertyDetails.Bathrooms.Value == "" &&
		webCards.PropertyFloorPlan.Title == "" &&
		len(webCards.PropertyFloorPlan.Plans) == 0 &&
		webCards.KnowAbout.Description == "" &&
		len(webCards.WhyToChoose.UspList) == 0 &&
		len(webCards.WhyToChoose.ImageUrls) == 0 &&
		len(webCards.VideoPresentation.Urls) == 0 &&
		webCards.LocationMap.Description == "" &&
		webCards.LocationMap.GoogleMapLink == ""
}

// Helper function to check if MetaInfo contains any meaningful data
func isMetaInfoEmpty(metaInfo schema.PropertyMetaInfo) bool {
	return metaInfo.Title == "" &&
		metaInfo.Description == "" &&
		metaInfo.Keywords == "" &&
		metaInfo.Canonical == ""
}

// Helper function to prefill property data from project information
func prefillPropertyFromProject(project *ent.Project, property *domain.Property) {
	if project == nil {
		return
	}

	// Initialize WebCards if empty - we don't need to check since we're going to set fields directly

	// Property Address: Inherit from project location info
	if project.LocationInfo.GoogleMapLink != "" || project.LocationInfo.ShortAddress != "" {
		property.WebCards.LocationMap.GoogleMapLink = project.LocationInfo.GoogleMapLink
		property.WebCards.LocationMap.Description = project.LocationInfo.ShortAddress
	}

	// RERA Details: Copy from project's RERA info
	if len(project.WebCards.ReraInfo.ReraList) > 0 && project.WebCards.ReraInfo.ReraList[0].ReraNumber != "" {
		// Use first RERA entry for property RERA number
		property.PropertyReraInfo.ReraNumber = project.WebCards.ReraInfo.ReraList[0].ReraNumber
		property.WebCards.PropertyDetails.ReraNumber.Value = project.WebCards.ReraInfo.ReraList[0].ReraNumber
	}

	// Built Up Area: Copy from project details
	if project.WebCards.Details.Area.Value != "" {
		property.WebCards.PropertyDetails.BuiltUpArea.Value = project.WebCards.Details.Area.Value
	}
	if project.WebCards.Details.Sizes.Value != "" {
		property.WebCards.PropertyDetails.Sizes.Value = project.WebCards.Details.Sizes.Value
	}

	// Possession Status: Copy from project possession date
	if project.WebCards.Details.PossessionDate.Value != "" {
		property.WebCards.PropertyDetails.PossessionStatus.Value = project.WebCards.Details.PossessionDate.Value
	}

	// Property Type: Map project type to property type
	if project.ProjectType.String() != "" {
		property.WebCards.PropertyDetails.PropertyType.Value = project.ProjectType.String()
	}
	if project.WebCards.Details.Type.Value != "" {
		property.WebCards.PropertyDetails.PropertyType.Value = project.WebCards.Details.Type.Value
	}

	// Configuration: Copy from project
	if project.WebCards.Details.Configuration.Value != "" {
		property.WebCards.PropertyDetails.Configuration.Value = project.WebCards.Details.Configuration.Value
	}

	// Why To Choose (USP): Copy from project
	if len(project.WebCards.WhyToChoose.USP_List) > 0 {
		property.WebCards.WhyToChoose.UspList = make([]string, len(project.WebCards.WhyToChoose.USP_List))
		copy(property.WebCards.WhyToChoose.UspList, project.WebCards.WhyToChoose.USP_List)
	}
	if len(project.WebCards.WhyToChoose.ImageUrls) > 0 {
		property.WebCards.WhyToChoose.ImageUrls = make([]string, len(project.WebCards.WhyToChoose.ImageUrls))
		copy(property.WebCards.WhyToChoose.ImageUrls, project.WebCards.WhyToChoose.ImageUrls)
	}

	// Know About: Copy from project
	if project.WebCards.KnowAbout.Description != "" {
		property.WebCards.KnowAbout.Description = project.WebCards.KnowAbout.Description
	}

	// Video Presentation: Copy from project
	if len(project.WebCards.VideoPresentation.URLs) > 0 {
		property.WebCards.VideoPresentation.Urls = make([]string, len(project.WebCards.VideoPresentation.URLs))
		copy(property.WebCards.VideoPresentation.Urls, project.WebCards.VideoPresentation.URLs)
	}

	// Note: Amenities and Builder About are not directly copied to property schema
	// as they are typically accessed via the project relationship
	// However, they can be accessed through the project relationship when needed
}

func (c *application) GetPropertiesOfProject(projectID string) ([]*response.Property, *imhttp.CustomError) {
	properties, err := c.repo.GetPropertiesOfProject(projectID)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get properties of project")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get properties of project", err.Error())
	}

	var propertyResponses []*response.Property
	for _, property := range properties {
		propertyResponses = append(propertyResponses, response.GetPropertyFromEnt(property))
	}

	return propertyResponses, nil
}

func (c *application) AddProperty(input request.AddPropertyRequest) (*response.AddPropertyResponse, *imhttp.CustomError) {
	var property domain.Property

	// Set basic property fields from request
	property.ProjectID = input.ProjectID
	property.Name = input.Name
	property.PropertyType = input.PropertyType
	property.CreatedByUserID = input.CreatedByUserID

	// Fetch project data to prefill property information
	if input.ProjectID != "" {
		project, err := c.repo.GetProjectByID(input.ProjectID)
		if err != nil {
			logger.Get().Error().Err(err).Str("projectID", input.ProjectID).Msg("Failed to get project for prefilling")
			return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get project data", err.Error())
		}

		if project == nil {
			return nil, imhttp.NewCustomErr(http.StatusNotFound, "Project not found", "Project not found")
		}

		// Prefill property data from project information
		prefillPropertyFromProject(project, &property)
	}

	result, err := c.repo.AddProperty(property)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to add property")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to add property", err.Error())
	}
	return &response.AddPropertyResponse{
		PropertyID: result.PropertyID,
		Slug:       result.Slug,
	}, nil
}

func (c *application) ListProperties(pagination *request.GetAllAPIRequest) ([]*response.PropertyListResponse, int, *imhttp.CustomError) {
	properties, totalItems, err := c.repo.GetAllProperties(pagination.GetOffset(), pagination.GetLimit(), pagination.Filters)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to list properties")
		return nil, 0, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to list properties", err.Error())
	}

	var propertyResponses []*response.PropertyListResponse
	for _, property := range properties {
		var developerName, location string
		if property.Edges.Developer != nil {
			developerName = property.Edges.Developer.Name
		}
		if property.Edges.Location != nil {
			location = property.Edges.Location.LocalityName
		}
		propertyResponses = append(propertyResponses, response.GetPropertyListResponse(property, developerName, location))
	}

	return propertyResponses, totalItems, nil
}

func (c *application) DeleteProperty(id string) *imhttp.CustomError {
	isDeleted, err := c.repo.IsPropertyDeleted(id)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to check if property is deleted")
		return imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to check if property is deleted", err.Error())
	}
	if isDeleted {
		return imhttp.NewCustomErr(http.StatusBadRequest, "Property is already deleted", "Property is already deleted")
	}

	if err := c.repo.DeleteProperty(id, false); err != nil {
		logger.Get().Error().Err(err).Msg("Failed to delete property")
		return imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to delete property", err.Error())
	}

	return nil
}
