package application

import (
	"net/http"

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

func (c *application) UpdateProperty(input request.UpdatePropertyRequest) (*response.Property, *imhttp.CustomError) {
	var property domain.Property

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

	property.PropertyID = input.PropertyID
	property.Name = input.Name
	property.PropertyImages = input.PropertyImages
	property.WebCards = input.WebCards
	property.PricingInfo = input.PricingInfo
	property.PropertyReraInfo = input.PropertyReraInfo
	property.MetaInfo = input.MetaInfo
	property.IsFeatured = input.IsFeatured
	property.IsDeleted = input.IsDeleted
	property.DeveloperID = input.DeveloperID
	property.LocationID = input.LocationID
	property.ProjectID = input.ProjectID
	if input.Slug != "" {
		property.Slug = input.Slug
	}

	updatedProperty, err := c.repo.UpdateProperty(property)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to update property")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to update property", err.Error())
	}

	return response.GetPropertyFromEnt(updatedProperty), nil
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

	property.ProjectID = input.ProjectID
	property.Name = input.Name
	property.PropertyType = input.PropertyType
	property.CreatedByUserID = input.CreatedByUserID
	property.Slug = input.Slug

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
