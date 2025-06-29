package repository

import (
	"context"
	"errors"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/developer"
	"github.com/VI-IM/im_backend_go/ent/project"
	"github.com/VI-IM/im_backend_go/ent/schema"
	"github.com/VI-IM/im_backend_go/internal/domain"
	"github.com/VI-IM/im_backend_go/internal/domain/enums"
	"github.com/VI-IM/im_backend_go/shared/logger"
)

func (r *repository) GetProjectByID(id string) (*ent.Project, error) {
	project, err := r.db.Project.Get(context.Background(), id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("project not found")
		}
		return nil, err
	}

	if project.IsDeleted {
		return nil, errors.New("project is deleted")
	}
	return project, nil
}

func (r *repository) AddProject(input domain.Project) (string, error) {
	if err := r.db.Project.Create().
		SetID(input.ProjectID).
		SetName(input.ProjectName).
		SetStatus(enums.ProjectStatusNEWLAUNCH).
		SetMetaInfo(schema.SEOMeta{
			Canonical: input.ProjectURL,
		}).
		SetDescription("").
		SetDeveloperID(input.DeveloperID).
		Exec(context.Background()); err != nil {
		logger.Get().Error().Err(err).Msg("Failed to add project")
		return "", err
	}
	return input.ProjectID, nil
}

func (r *repository) ExistDeveloperByID(id string) (bool, error) {
	exist, err := r.db.Developer.Query().
		Where(developer.ID(id)).
		Exist(context.Background())
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to check if developer exists")
		return false, err
	}
	return exist, nil
}

func (r *repository) IsProjectDeleted(id string) (bool, error) {
	project, err := r.db.Project.Query().
		Where(project.ID(id)).
		Only(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			return false, nil
		}
		logger.Get().Error().Err(err).Msg("Failed to check if project is deleted")
		return false, err
	}
	return project.IsDeleted, nil
}

func (r *repository) UpdateProject(input domain.Project) (*ent.Project, error) {
	// get the project
	oldProject, err := r.GetProjectByID(input.ProjectID)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get project")
		return nil, err
	}

	project := r.db.Project.UpdateOneID(input.ProjectID)

	if input.ProjectName != "" {
		project.SetName(input.ProjectName)
	}
	if input.Status != "" {
		project.SetStatus(input.Status)
	}

	// Handle TimelineInfo updates
	if input.TimelineInfo != (schema.TimelineInfo{}) {
		newTimelineInfo := oldProject.TimelineInfo
		if input.TimelineInfo.ProjectLaunchDate != "" {
			newTimelineInfo.ProjectLaunchDate = input.TimelineInfo.ProjectLaunchDate
		}
		if input.TimelineInfo.ProjectPossessionDate != "" {
			newTimelineInfo.ProjectPossessionDate = input.TimelineInfo.ProjectPossessionDate
		}
		project.SetTimelineInfo(newTimelineInfo)
	}

	// Handle MetaInfo updates
	if input.MetaInfo != (schema.SEOMeta{}) {
		newMetaInfo := oldProject.MetaInfo
		if input.MetaInfo.Title != "" {
			newMetaInfo.Title = input.MetaInfo.Title
		}
		if input.MetaInfo.Description != "" {
			newMetaInfo.Description = input.MetaInfo.Description
		}
		if input.MetaInfo.Keywords != "" {
			newMetaInfo.Keywords = input.MetaInfo.Keywords
		}
		if input.MetaInfo.Canonical != "" {
			newMetaInfo.Canonical = input.MetaInfo.Canonical
		}
		if input.MetaInfo.ProjectSchema != "" {
			newMetaInfo.ProjectSchema = input.MetaInfo.ProjectSchema
		}
		project.SetMetaInfo(newMetaInfo)
	}

	// Handle WebCards updates
	newWebCards := oldProject.WebCards
	hasWebCardChanges := false

	// Update Images if provided
	if len(input.WebCards.Images) > 0 {
		newWebCards.Images = input.WebCards.Images
		hasWebCardChanges = true
	}

	// Update ReraInfo if provided
	if input.WebCards.ReraInfo.WebsiteLink != "" || len(input.WebCards.ReraInfo.ReraList) > 0 {
		if input.WebCards.ReraInfo.WebsiteLink != "" {
			newWebCards.ReraInfo.WebsiteLink = input.WebCards.ReraInfo.WebsiteLink
		}
		if len(input.WebCards.ReraInfo.ReraList) > 0 {
			newWebCards.ReraInfo.ReraList = input.WebCards.ReraInfo.ReraList
		}
		hasWebCardChanges = true
	}

	// Update Details if provided
	if input.WebCards.Details != (schema.ProjectDetails{}) {
		if input.WebCards.Details.Area.Value != "" {
			newWebCards.Details.Area.Value = input.WebCards.Details.Area.Value
		}
		if input.WebCards.Details.Sizes.Value != "" {
			newWebCards.Details.Sizes.Value = input.WebCards.Details.Sizes.Value
		}
		if input.WebCards.Details.Units.Value != "" {
			newWebCards.Details.Units.Value = input.WebCards.Details.Units.Value
		}
		if input.WebCards.Details.Configuration.Value != "" {
			newWebCards.Details.Configuration.Value = input.WebCards.Details.Configuration.Value
		}
		if input.WebCards.Details.TotalFloor.Value != "" {
			newWebCards.Details.TotalFloor.Value = input.WebCards.Details.TotalFloor.Value
		}
		if input.WebCards.Details.TotalTowers.Value != "" {
			newWebCards.Details.TotalTowers.Value = input.WebCards.Details.TotalTowers.Value
		}
		if input.WebCards.Details.LaunchDate.Value != "" {
			newWebCards.Details.LaunchDate.Value = input.WebCards.Details.LaunchDate.Value
		}
		if input.WebCards.Details.PossessionDate.Value != "" {
			newWebCards.Details.PossessionDate.Value = input.WebCards.Details.PossessionDate.Value
		}
		if input.WebCards.Details.Type.Value != "" {
			newWebCards.Details.Type.Value = input.WebCards.Details.Type.Value
		}
		hasWebCardChanges = true
	}

	// Update WhyToChoose if provided
	if len(input.WebCards.WhyToChoose.ImageUrls) > 0 {
		newWebCards.WhyToChoose.ImageUrls = input.WebCards.WhyToChoose.ImageUrls
		hasWebCardChanges = true
	}

	if len(input.WebCards.WhyToChoose.USP_List) > 0 {
		if len(newWebCards.WhyToChoose.USP_List) == 0 {
			newWebCards.WhyToChoose.USP_List = make([]string, len(input.WebCards.WhyToChoose.USP_List))
		}

		for i, usp := range input.WebCards.WhyToChoose.USP_List {
			if i >= len(newWebCards.WhyToChoose.USP_List) || newWebCards.WhyToChoose.USP_List[i] != usp {
				if i >= len(newWebCards.WhyToChoose.USP_List) {
					newWebCards.WhyToChoose.USP_List = append(newWebCards.WhyToChoose.USP_List, "")
				}
				newWebCards.WhyToChoose.USP_List[i] = usp
				hasWebCardChanges = true
			}
		}
	}

	// Update KnowAbout if provided
	if input.WebCards.KnowAbout.Description != "" || input.WebCards.KnowAbout.DownloadLink != "" {
		if input.WebCards.KnowAbout.Description != "" {
			newWebCards.KnowAbout.Description = input.WebCards.KnowAbout.Description
		}
		if input.WebCards.KnowAbout.DownloadLink != "" {
			newWebCards.KnowAbout.DownloadLink = input.WebCards.KnowAbout.DownloadLink
		}
		hasWebCardChanges = true
	}

	// Update FloorPlan if provided
	if len(input.WebCards.FloorPlan.Products) > 0 || input.WebCards.FloorPlan.Description != "" {
		if input.WebCards.FloorPlan.Description != "" {
			newWebCards.FloorPlan.Description = input.WebCards.FloorPlan.Description
		}
		if len(input.WebCards.FloorPlan.Products) > 0 {
			newWebCards.FloorPlan.Products = input.WebCards.FloorPlan.Products
		}
		hasWebCardChanges = true
	}

	// Update PriceList if provided
	if len(input.WebCards.PriceList.BHKOptionsWithPrices) > 0 || input.WebCards.PriceList.Description != "" {
		if input.WebCards.PriceList.Description != "" {
			newWebCards.PriceList.Description = input.WebCards.PriceList.Description
		}
		if len(input.WebCards.PriceList.BHKOptionsWithPrices) > 0 {
			newWebCards.PriceList.BHKOptionsWithPrices = input.WebCards.PriceList.BHKOptionsWithPrices
		}
		hasWebCardChanges = true
	}

	// Update Amenities if provided
	if len(input.WebCards.Amenities.CategoriesWithAmenities) > 0 || input.WebCards.Amenities.Description != "" {
		if input.WebCards.Amenities.Description != "" {
			newWebCards.Amenities.Description = input.WebCards.Amenities.Description
		}
		if len(input.WebCards.Amenities.CategoriesWithAmenities) > 0 {
			newWebCards.Amenities.CategoriesWithAmenities = input.WebCards.Amenities.CategoriesWithAmenities
		}
		hasWebCardChanges = true
	}

	// Update VideoPresentation if provided
	if input.WebCards.VideoPresentation.Description != "" || len(input.WebCards.VideoPresentation.URL) > 0 {
		if input.WebCards.VideoPresentation.Description != "" {
			newWebCards.VideoPresentation.Description = input.WebCards.VideoPresentation.Description
		}
		if len(input.WebCards.VideoPresentation.URL) > 0 {
			newWebCards.VideoPresentation.URL = input.WebCards.VideoPresentation.URL
		}
		hasWebCardChanges = true
	}

	// Update PaymentPlans if provided
	if len(input.WebCards.PaymentPlans.Plans) > 0 || input.WebCards.PaymentPlans.Description != "" {
		if input.WebCards.PaymentPlans.Description != "" {
			newWebCards.PaymentPlans.Description = input.WebCards.PaymentPlans.Description
		}
		if len(input.WebCards.PaymentPlans.Plans) > 0 {
			newWebCards.PaymentPlans.Plans = input.WebCards.PaymentPlans.Plans
		}
		hasWebCardChanges = true
	}

	if hasWebCardChanges {
		project.SetWebCards(newWebCards)
	}

	// Handle LocationInfo updates
	if input.LocationInfo != (schema.LocationInfo{}) {
		newLocationInfo := oldProject.LocationInfo
		if input.LocationInfo.ShortAddress != "" {
			newLocationInfo.ShortAddress = input.LocationInfo.ShortAddress
		}
		if input.LocationInfo.Longitude != "" {
			newLocationInfo.Longitude = input.LocationInfo.Longitude
		}
		if input.LocationInfo.Latitude != "" {
			newLocationInfo.Latitude = input.LocationInfo.Latitude
		}
		if input.LocationInfo.GoogleMapLink != "" {
			newLocationInfo.GoogleMapLink = input.LocationInfo.GoogleMapLink
		}
		project.SetLocationInfo(newLocationInfo)
	}

	// Handle boolean flags
	if input.IsFeatured {
		project.SetIsFeatured(input.IsFeatured)
	}
	if input.IsPremium {
		project.SetIsPremium(input.IsPremium)
	}
	if input.IsPriority {
		project.SetIsPriority(input.IsPriority)
	}

	// Handle other fields
	if input.Description != "" {
		project.SetDescription(input.Description)
	}
	if input.IsDeleted {
		project.SetIsDeleted(input.IsDeleted)
	}

	updatedProject, err := project.Save(context.Background())
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to update project")
		return nil, err
	}

	return updatedProject, nil
}

func (r *repository) DeleteProject(id string, hardDelete bool) error {
	if hardDelete {
		// Perform hard delete
		err := r.db.Project.DeleteOneID(id).Exec(context.Background())
		if err != nil {
			if ent.IsNotFound(err) {
				return errors.New("project not found")
			}
			logger.Get().Error().Err(err).Msg("Failed to delete project")
			return err
		}
		return nil
	}

	// Perform soft delete by updating IsDeleted flag
	_, err := r.db.Project.UpdateOneID(id).
		SetIsDeleted(true).
		Save(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.New("project not found")
		}
		logger.Get().Error().Err(err).Msg("Failed to soft delete project")
		return err
	}
	return nil
}

func (r *repository) GetAllProjects(offset, limit int) ([]*ent.Project, int, error) {
	ctx := context.Background()

	// Get total count
	total, err := r.db.Project.Query().
		Where(project.IsDeletedEQ(false)).
		Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	projects, err := r.db.Project.Query().
		Where(project.IsDeletedEQ(false)).
		Order(ent.Desc(project.FieldID)). // Order by ID as a fallback since CreatedAt is not available
		Offset(offset).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}
