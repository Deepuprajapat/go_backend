package repository

import (
	"context"
	"errors"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/internal/domain"
	"github.com/VI-IM/im_backend_go/shared/logger"
)

func (r *repository) GetPropertyByID(id string) (*ent.Property, error) {
	property, err := r.db.Property.Get(context.Background(), id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("property not found")
		}
		logger.Get().Error().Err(err).Msg("Failed to get property")
		return nil, err
	}
	return property, nil
}

func (r *repository) UpdateProperty(input domain.Property) (*ent.Property, error) {
	
	oldProperty, err := r.GetPropertyByID(input.PropertyID)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get property")
		return nil, err
	}

	property := r.db.Property.UpdateOneID(input.PropertyID)

	if input.Name != "" {
		property.SetName(input.Name)
	}
	if len(input.PropertyImages) > 0 {
		property.SetPropertyImages(input.PropertyImages)
	}

	
	newWebCards := oldProperty.WebCards
	hasWebCardChanges := false

	if input.WebCards.PropertyDetails != (oldProperty.WebCards.PropertyDetails) {
		if input.WebCards.PropertyDetails.BuiltUpArea.Value != "" {
			newWebCards.PropertyDetails.BuiltUpArea.Value = input.WebCards.PropertyDetails.BuiltUpArea.Value
			hasWebCardChanges = true
		}
		if input.WebCards.PropertyDetails.Sizes.Value != "" {
			newWebCards.PropertyDetails.Sizes.Value = input.WebCards.PropertyDetails.Sizes.Value
			hasWebCardChanges = true
		}
		if input.WebCards.PropertyDetails.FloorNumber.Value != "" {
			newWebCards.PropertyDetails.FloorNumber.Value = input.WebCards.PropertyDetails.FloorNumber.Value
			hasWebCardChanges = true
		}
		if input.WebCards.PropertyDetails.Configuration.Value != "" {
			newWebCards.PropertyDetails.Configuration.Value = input.WebCards.PropertyDetails.Configuration.Value
			hasWebCardChanges = true
		}
		if input.WebCards.PropertyDetails.PossessionStatus.Value != "" {
			newWebCards.PropertyDetails.PossessionStatus.Value = input.WebCards.PropertyDetails.PossessionStatus.Value
			hasWebCardChanges = true
		}
		if input.WebCards.PropertyDetails.Balconies.Value != "" {
			newWebCards.PropertyDetails.Balconies.Value = input.WebCards.PropertyDetails.Balconies.Value
			hasWebCardChanges = true
		}
		if input.WebCards.PropertyDetails.CoveredParking.Value != "" {
			newWebCards.PropertyDetails.CoveredParking.Value = input.WebCards.PropertyDetails.CoveredParking.Value
			hasWebCardChanges = true
		}
		if input.WebCards.PropertyDetails.Bedrooms.Value != "" {
			newWebCards.PropertyDetails.Bedrooms.Value = input.WebCards.PropertyDetails.Bedrooms.Value
			hasWebCardChanges = true
		}
		if input.WebCards.PropertyDetails.PropertyType.Value != "" {
			newWebCards.PropertyDetails.PropertyType.Value = input.WebCards.PropertyDetails.PropertyType.Value
			hasWebCardChanges = true
		}
		if input.WebCards.PropertyDetails.AgeOfProperty.Value != "" {
			newWebCards.PropertyDetails.AgeOfProperty.Value = input.WebCards.PropertyDetails.AgeOfProperty.Value
			hasWebCardChanges = true
		}
		if input.WebCards.PropertyDetails.FurnishingType.Value != "" {
			newWebCards.PropertyDetails.FurnishingType.Value = input.WebCards.PropertyDetails.FurnishingType.Value
			hasWebCardChanges = true
		}
		if input.WebCards.PropertyDetails.ReraNumber.Value != "" {
			newWebCards.PropertyDetails.ReraNumber.Value = input.WebCards.PropertyDetails.ReraNumber.Value
			hasWebCardChanges = true
		}
		if input.WebCards.PropertyDetails.Facing.Value != "" {
			newWebCards.PropertyDetails.Facing.Value = input.WebCards.PropertyDetails.Facing.Value
			hasWebCardChanges = true
		}
		if input.WebCards.PropertyDetails.Bathrooms.Value != "" {
			newWebCards.PropertyDetails.Bathrooms.Value = input.WebCards.PropertyDetails.Bathrooms.Value
			hasWebCardChanges = true
		}
	}

	
	if len(input.WebCards.WhyChooseUs.ImageUrls) > 0 || len(input.WebCards.WhyChooseUs.USP_List) > 0 {
		if len(input.WebCards.WhyChooseUs.ImageUrls) > 0 {
			newWebCards.WhyChooseUs.ImageUrls = input.WebCards.WhyChooseUs.ImageUrls
		}
		if len(input.WebCards.WhyChooseUs.USP_List) > 0 {
			newWebCards.WhyChooseUs.USP_List = input.WebCards.WhyChooseUs.USP_List
		}
		hasWebCardChanges = true
	}

	if input.WebCards.KnowAbout.Description != "" {
		newWebCards.KnowAbout.Description = input.WebCards.KnowAbout.Description
		hasWebCardChanges = true
	}

	if input.WebCards.VideoPresentation.Title != "" || input.WebCards.VideoPresentation.VideoUrl != "" {
		if input.WebCards.VideoPresentation.Title != "" {
			newWebCards.VideoPresentation.Title = input.WebCards.VideoPresentation.Title
		}
		if input.WebCards.VideoPresentation.VideoUrl != "" {
			newWebCards.VideoPresentation.VideoUrl = input.WebCards.VideoPresentation.VideoUrl
		}
		hasWebCardChanges = true
	}

	if input.WebCards.LocationMap.Description != "" || input.WebCards.LocationMap.GoogleMapLink != "" {
		if input.WebCards.LocationMap.Description != "" {
			newWebCards.LocationMap.Description = input.WebCards.LocationMap.Description
		}
		if input.WebCards.LocationMap.GoogleMapLink != "" {
			newWebCards.LocationMap.GoogleMapLink = input.WebCards.LocationMap.GoogleMapLink
		}
		hasWebCardChanges = true
	}

	if hasWebCardChanges {
		property.SetWebCards(newWebCards)
	}

	if input.PricingInfo != (oldProperty.PricingInfo) {
		property.SetPricingInfo(input.PricingInfo)
	}
	if input.PropertyReraInfo != (oldProperty.PropertyReraInfo) {
		property.SetPropertyReraInfo(input.PropertyReraInfo)
	}
	if input.MetaInfo != (oldProperty.MetaInfo) {
		property.SetMetaInfo(input.MetaInfo)
	}
	if input.DeveloperID != "" {
		property.SetDeveloperID(input.DeveloperID)
	}
	if input.LocationID != "" {
		property.SetLocationID(input.LocationID)
	}
	if input.ProjectID != "" {
		property.SetProjectID(input.ProjectID)
	}

	updatedProperty, err := property.Save(context.Background())
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to update property")
		return nil, err
	}
	return updatedProperty, nil
}
