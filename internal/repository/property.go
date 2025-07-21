package repository

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/location"
	"github.com/VI-IM/im_backend_go/ent/project"
	"github.com/VI-IM/im_backend_go/ent/property"
	"github.com/VI-IM/im_backend_go/ent/schema"
	"github.com/VI-IM/im_backend_go/internal/domain"
	"github.com/VI-IM/im_backend_go/shared/logger"
	"github.com/google/uuid"
)

// PropertyResult contains the result of adding a property
type PropertyResult struct {
	PropertyID string
	Slug       string
}

// generateSlug creates a URL-friendly slug from property name and property ID
func generateSlug(name, propertyID string) string {
	// Normalize the name: lowercase, replace spaces and special characters with hyphens
	normalized := strings.ToLower(name)
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	normalized = reg.ReplaceAllString(normalized, "-")
	normalized = strings.Trim(normalized, "-")

	// Get last 6 characters of property ID
	suffix := propertyID
	if len(propertyID) > 6 {
		suffix = propertyID[len(propertyID)-6:]
	}

	return normalized + "-" + suffix
}

// createDefaultWebCards creates default empty web cards structure
func createDefaultWebCards() schema.WebCards {
	return schema.WebCards{
		PropertyDetails: schema.PropertyDetails{},
		PropertyFloorPlan: schema.PropertyFloorPlan{
			Title: "",
			Plans: []map[string]string{},
		},
		KnowAbout: struct {
			Description string `json:"description,omitempty"`
		}{
			Description: "",
		},
		WhyToChoose: struct {
			UspList   []string `json:"usp_list,omitempty"`
			ImageUrls []string `json:"image_urls,omitempty"`
		}{
			UspList:   []string{},
			ImageUrls: []string{},
		},
		VideoPresentation: struct {
			Urls []string `json:"urls,omitempty"`
		}{
			Urls: []string{},
		},
		LocationMap: struct {
			Description   string `json:"description,omitempty"`
			GoogleMapLink string `json:"google_map_link,omitempty"`
		}{
			Description:   "",
			GoogleMapLink: "",
		},
	}
}

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

func (r *repository) GetPropertyBySlug(ctx context.Context, slug string) (*ent.Property, error) {
	property, err := r.db.Property.Query().
		Where(property.Slug(slug)).
		WithDeveloper().
		WithProject().
		Only(ctx)
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
		// Regenerate slug when property name changes
		if input.Name != oldProperty.Name {
			newSlug := generateSlug(input.Name, input.PropertyID)
			property.SetSlug(newSlug)
		}
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

	// Only update PricingInfo if it contains actual data (non-empty Price field)
	if input.PricingInfo.Price != "" && input.PricingInfo != (oldProperty.PricingInfo) {
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

func (r *repository) AddProperty(input domain.Property) (*PropertyResult, error) {
	project, err := r.db.Project.Query().
		Where(project.ID(input.ProjectID)).
		WithDeveloper().
		WithLocation().
		First(context.Background())
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get project")
		return nil, err
	}

	// Validate user exists if CreatedByUserID is provided
	if input.CreatedByUserID != nil {
		userExists, err := r.CheckIfUserExistsByID(context.Background(), *input.CreatedByUserID)
		if err != nil {
			logger.Get().Error().Err(err).Msg("Failed to check if user exists")
			return nil, err
		}
		if !userExists {
			logger.Get().Error().Str("user_id", *input.CreatedByUserID).Msg("User does not exist")
			return nil, errors.New("user does not exist")
		}
	}

	propertyID := uuid.New().String()
	slug := generateSlug(input.Name, propertyID)

	// Create default values for required JSON fields
	defaultPricingInfo := schema.PropertyPricingInfo{Price: ""}

	property := r.db.Property.Create().
		SetID(propertyID).
		SetProjectID(input.ProjectID).
		SetName(input.Name).
		SetSlug(slug).
		SetPropertyType(input.PropertyType).
		SetWebCards(input.WebCards).
		SetPricingInfo(defaultPricingInfo).
		SetPropertyReraInfo(input.PropertyReraInfo)
	if project.Edges.Developer != nil && project.Edges.Developer.ID != "" {
		property.SetDeveloperID(project.Edges.Developer.ID)
	}
	if project.Edges.Location != nil && project.Edges.Location.ID != "" {
		property.SetLocationID(project.Edges.Location.ID)
	}
	if input.CreatedByUserID != nil {
		property.SetCreatedByUserID(*input.CreatedByUserID)
	}

	if err = property.Exec(context.Background()); err != nil {
		logger.Get().Error().Err(err).Msg("Failed to add property")
		return nil, err
	}
	return &PropertyResult{
		PropertyID: propertyID,
		Slug:       slug,
	}, nil
}

func (r *repository) GetAllProperties(offset, limit int, filters map[string]interface{}) ([]*ent.Property, int, error) {
	ctx := context.Background()

	// Start building the base query
	baseQuery := r.db.Property.Query().Where(property.IsDeletedEQ(false))

	// Apply filters at query level
	query := r.applyPropertyFilters(baseQuery, filters)

	// Get total count with filters applied (but without pagination)
	total, err := query.Count(ctx)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to count filtered properties")
		return nil, 0, err
	}

	// Apply pagination and fetch results
	properties, err := query.
		Order(ent.Desc(property.FieldCreatedAt)).
		WithDeveloper().
		WithLocation().
		WithProject().
		Offset(offset).
		Limit(limit).
		All(ctx)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to fetch properties")
		return nil, 0, err
	}

	return properties, total, nil
}

// applyPropertyFilters applies filters to the property query
func (r *repository) applyPropertyFilters(query *ent.PropertyQuery, filters map[string]interface{}) *ent.PropertyQuery {
	// Filter by created_by_user_id
	if createdByUserID, ok := filters["created_by_user_id"].(string); ok && createdByUserID != "" {
		query = query.Where(property.CreatedByUserIDEQ(createdByUserID))
	}

	// Filter by property_type - this requires JSON field filtering
	if propertyType, ok := filters["property_type"].(string); ok && propertyType != "" {
		// Note: This is a simplified approach. For complex JSON filtering,
		// you might need raw SQL or restructure the schema
		query = query.Where(property.PropertyTypeContains(propertyType))
	}

	// Filter by city through location relationship
	if city, ok := filters["city"].(string); ok && city != "" {
		query = query.Where(property.HasLocationWith(location.CityEQ(city)))
	}

	// Filter by developer_id
	if developerID, ok := filters["developer_id"].(string); ok && developerID != "" {
		query = query.Where(property.DeveloperIDEQ(developerID))
	}

	// Filter by location_id
	if locationID, ok := filters["location_id"].(string); ok && locationID != "" {
		query = query.Where(property.LocationIDEQ(locationID))
	}

	// Filter by project_id
	if projectID, ok := filters["project_id"].(string); ok && projectID != "" {
		query = query.Where(property.ProjectIDEQ(projectID))
	}

	// Filter by is_featured
	if isFeatured, ok := filters["is_featured"].(bool); ok {
		query = query.Where(property.IsFeaturedEQ(isFeatured))
	}

	// Filter by name (partial match)
	if name, ok := filters["name"].(string); ok && name != "" {
		query = query.Where(property.NameContainsFold(name))
	}

	return query
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

func (r *repository) GetPropertyBySlug(slug string) (*ent.Property, error) {
	property, err := r.db.Property.Query().
		Where(property.Slug(slug)).
		Only(context.Background())
	if err != nil {
		return nil, err
	}
	return property, nil
}
