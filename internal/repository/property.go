package repository

import (
	"context"
	"errors"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/project"
	"github.com/VI-IM/im_backend_go/ent/property"
	"github.com/VI-IM/im_backend_go/internal/domain"
	"github.com/VI-IM/im_backend_go/shared/logger"
	"github.com/google/uuid"
)

func (r *repository) GetPropertyByID(id string) (*ent.Property, error) {
	property, err := r.db.Property.Query().
		Where(property.ID(id)).
		WithDeveloper().
		WithProject().
		Only(context.Background())
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

	// property details
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

	// why choose us
	// if len(input.WebCards.WhyChooseUs.ImageUrls) > 0 || len(input.WebCards.WhyChooseUs.USP_List) > 0 {
	// 	if len(input.WebCards.WhyChooseUs.ImageUrls) > 0 {
	// 		newWebCards.WhyChooseUs.ImageUrls = input.WebCards.WhyChooseUs.ImageUrls
	// 	}
	// 	if len(input.WebCards.WhyChooseUs.USP_List) > 0 {
	// 		newWebCards.WhyChooseUs.USP_List = input.WebCards.WhyChooseUs.USP_List
	// 	}
	// 	hasWebCardChanges = true
	// }

	// // know about
	// if input.WebCards.KnowAbout.Description != "" {
	// 	newWebCards.KnowAbout.Description = input.WebCards.KnowAbout.Description
	// 	hasWebCardChanges = true
	// }

	// // video presentation
	// if input.WebCards.VideoPresentation.Title != "" || input.WebCards.VideoPresentation.VideoUrl != "" {
	// 	if input.WebCards.VideoPresentation.Title != "" {
	// 		newWebCards.VideoPresentation.Title = input.WebCards.VideoPresentation.Title
	// 	}
	// 	if input.WebCards.VideoPresentation.VideoUrl != "" {
	// 		newWebCards.VideoPresentation.VideoUrl = input.WebCards.VideoPresentation.VideoUrl
	// 	}
	// 	hasWebCardChanges = true
	// }

	// location map
	if input.WebCards.LocationMap.Description != "" || input.WebCards.LocationMap.GoogleMapLink != "" {
		if input.WebCards.LocationMap.Description != "" {
			newWebCards.LocationMap.Description = input.WebCards.LocationMap.Description
		}
		if input.WebCards.LocationMap.GoogleMapLink != "" {
			newWebCards.LocationMap.GoogleMapLink = input.WebCards.LocationMap.GoogleMapLink
		}
		hasWebCardChanges = true
	}

	// property floor plan
	if input.WebCards.PropertyFloorPlan.Title != "" || len(input.WebCards.PropertyFloorPlan.Plans) > 0 {
		if input.WebCards.PropertyFloorPlan.Title != "" {
			newWebCards.PropertyFloorPlan.Title = input.WebCards.PropertyFloorPlan.Title
		}
		if len(input.WebCards.PropertyFloorPlan.Plans) > 0 {
			newWebCards.PropertyFloorPlan.Plans = input.WebCards.PropertyFloorPlan.Plans
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

func (r *repository) GetPropertiesOfProject(projectID string) ([]*ent.Property, error) {
	if projectID == "" {
		return nil, errors.New("projectID is required")
	}

	properties, err := r.db.Property.Query().
		Where(property.ProjectID(projectID)).
		WithDeveloper().
		WithProject().
		All(context.Background())
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get properties of project")
		return nil, err
	}
	return properties, nil
}

func (r *repository) AddProperty(input domain.Property) (string, error) {
	project, err := r.db.Project.Query().
		Where(project.ID(input.ProjectID)).
		WithDeveloper().
		WithLocation().
		First(context.Background())
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get project")
		return "", err
	}

	propertyID := uuid.New().String()
	property := r.db.Property.Create().
		SetID(propertyID).
		SetProjectID(input.ProjectID).
		SetName(input.Name).
		SetPropertyType(input.PropertyType)
	if project.Edges.Developer.ID != "" {
		property.SetDeveloperID(project.Edges.Developer.ID)
	}
	if project.Edges.Location.ID != "" {
		property.SetLocationID(project.Edges.Location.ID)
	}

	if err = property.Exec(context.Background()); err != nil {
		logger.Get().Error().Err(err).Msg("Failed to add property")
		return "", err
	}
	return propertyID, nil
}

func (r *repository) GetAllProperties(offset, limit int, filters map[string]interface{}) ([]*ent.Property, int, error) {
	ctx := context.Background()

	// Start building the query
	query := r.db.Property.Query().Where(property.IsDeletedEQ(false))

	// Get all properties first
	properties, err := query.
		Order(ent.Desc(property.FieldID)).
		WithDeveloper().
		WithLocation().
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Apply filters in memory
	filteredProperties := make([]*ent.Property, 0)
	for _, p := range properties {
		// Check if property matches all filters
		matches := true

		if propertyType, ok := filters["property_type"].(string); ok && propertyType != "" {
			if p.WebCards.PropertyDetails.PropertyType.Value != propertyType {
				matches = false
			}
		}

		if configuration, ok := filters["configuration"].(string); ok && configuration != "" {
			if p.WebCards.PropertyDetails.Configuration.Value != configuration {
				matches = false
			}
		}

		if city, ok := filters["city"].(string); ok && city != "" {
			if p.Edges.Location.City != city {
				matches = false
			}
		}

		if matches {
			filteredProperties = append(filteredProperties, p)
		}
	}

	// Calculate total after filtering
	total := len(filteredProperties)

	// Apply pagination in memory
	start := offset
	end := offset + limit
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	return filteredProperties[start:end], total, nil
}

func (r *repository) DeleteProperty(id string, hardDelete bool) error {
	if hardDelete {
		// Perform hard delete
		err := r.db.Property.DeleteOneID(id).Exec(context.Background())
		if err != nil {
			if ent.IsNotFound(err) {
				return errors.New("property not found")
			}
			logger.Get().Error().Err(err).Msg("Failed to delete property")
			return err
		}
		return nil
	}

	// Perform soft delete by updating IsDeleted flag
	_, err := r.db.Property.UpdateOneID(id).
		SetIsDeleted(true).
		Save(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.New("property not found")
		}
		logger.Get().Error().Err(err).Msg("Failed to soft delete property")
		return err
	}
	return nil
}

func (r *repository) IsPropertyDeleted(id string) (bool, error) {
	property, err := r.db.Property.Query().
		Where(property.ID(id)).
		Only(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			return false, nil
		}
		logger.Get().Error().Err(err).Msg("Failed to check if property is deleted")
		return false, err
	}
	return property.IsDeleted, nil
}
