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
	property.WebCards.PropertyDetails.AgeOfProperty.Value = input.AgeOfProperty
	property.WebCards.PropertyDetails.FloorNumber.Value = input.FloorNumber
	property.WebCards.PropertyDetails.Facing.Value = input.Facing
	property.WebCards.PropertyDetails.FurnishingType.Value = input.Furnishing
	property.WebCards.PropertyDetails.Balconies.Value = input.BalconyCount
	property.WebCards.PropertyDetails.Bedrooms.Value = input.BedroomsCount
	property.WebCards.PropertyDetails.CoveredParking.Value = input.CoveredParking

	propertyID, err := c.repo.AddProperty(property)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to add property")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to add property", err.Error())
	}
	return &response.AddPropertyResponse{PropertyID: propertyID}, nil
}
