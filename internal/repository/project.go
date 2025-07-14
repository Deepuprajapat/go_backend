package repository

import (
	"context"
	"errors"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/VI-IM/im_backend_go/ent"
	developerEnt "github.com/VI-IM/im_backend_go/ent/developer"
	locationEnt "github.com/VI-IM/im_backend_go/ent/location"
	predicateEnt "github.com/VI-IM/im_backend_go/ent/predicate"
	projectEnt "github.com/VI-IM/im_backend_go/ent/project"
	"github.com/VI-IM/im_backend_go/ent/schema"
	"github.com/VI-IM/im_backend_go/internal/domain"
	"github.com/VI-IM/im_backend_go/internal/domain/enums"
	"github.com/VI-IM/im_backend_go/shared/logger"
)

func (r *repository) GetProjectByID(id string) (*ent.Project, error) {
	project, err := r.db.Project.Query().
		Where(projectEnt.ID(id)).
		WithDeveloper().
		WithLocation().
		Only(context.Background())
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
		SetProjectType(projectEnt.ProjectType(input.ProjectType)).
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
		Where(developerEnt.ID(id)).
		Exist(context.Background())
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to check if developer exists")
		return false, err
	}
	return exist, nil
}

func (r *repository) IsProjectDeleted(id string) (bool, error) {
	project, err := r.db.Project.Query().
		Where(projectEnt.ID(id)).
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
	tx, err := r.db.Tx(context.Background())
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to create transaction")
		return nil, err
	}
	defer tx.Rollback()

	if input.MaxPrice != ""{
		project.SetMaxPrice(input.MaxPrice)
	}

	if input.MinPrice != ""{
		project.SetMaxPrice(input.MinPrice)
	}
	
	if input.ProjectName != "" {
		project.SetName(input.ProjectName)
	}
	if input.Status != "" {
		project.SetStatus(input.Status)
	}

	if input.Description != "" {
		project.SetDescription(input.Description)
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
	if input.MetaInfo.Title != "" || input.MetaInfo.Description != "" || input.MetaInfo.Keywords != "" || input.MetaInfo.Canonical != "" || len(input.MetaInfo.ProjectSchema) > 0 {
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
		if len(input.MetaInfo.ProjectSchema) > 0 {
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
		newWebCards.WhyToChoose.USP_List = input.WebCards.WhyToChoose.USP_List
		hasWebCardChanges = true
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
			if len(newWebCards.FloorPlan.Products) == 0 {
				newWebCards.FloorPlan.Products = make([]schema.FloorPlanItem, len(input.WebCards.FloorPlan.Products))
			}
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
			if len(newWebCards.PriceList.BHKOptionsWithPrices) == 0 {
				newWebCards.PriceList.BHKOptionsWithPrices = make([]schema.ProductConfiguration, len(input.WebCards.PriceList.BHKOptionsWithPrices))
			}
			newWebCards.PriceList.BHKOptionsWithPrices = input.WebCards.PriceList.BHKOptionsWithPrices
		}
		hasWebCardChanges = true
	}

	// Update VideoPresentation if provided
	if input.WebCards.VideoPresentation.Description != "" || len(input.WebCards.VideoPresentation.URLs) > 0 {
		if input.WebCards.VideoPresentation.Description != "" {
			newWebCards.VideoPresentation.Description = input.WebCards.VideoPresentation.Description
		}
		if len(input.WebCards.VideoPresentation.URLs) > 0 {
			newWebCards.VideoPresentation.URLs = input.WebCards.VideoPresentation.URLs
		}
		hasWebCardChanges = true
	}

	// Update PaymentPlans if provided
	if len(input.WebCards.PaymentPlans.Plans) > 0 || input.WebCards.PaymentPlans.Description != "" {
		if input.WebCards.PaymentPlans.Description != "" {
			newWebCards.PaymentPlans.Description = input.WebCards.PaymentPlans.Description
		}
		if len(input.WebCards.PaymentPlans.Plans) > 0 {
			if len(newWebCards.PaymentPlans.Plans) == 0 {
				newWebCards.PaymentPlans.Plans = make([]schema.Plan, len(input.WebCards.PaymentPlans.Plans))
			}
			newWebCards.PaymentPlans.Plans = input.WebCards.PaymentPlans.Plans
		}
		hasWebCardChanges = true
	}

	// Update SitePlan if provided
	if input.WebCards.SitePlan.Description != "" || input.WebCards.SitePlan.Image != "" {
		if input.WebCards.SitePlan.Description != "" {
			newWebCards.SitePlan.Description = input.WebCards.SitePlan.Description
		}
		if input.WebCards.SitePlan.Image != "" {
			newWebCards.SitePlan.Image = input.WebCards.SitePlan.Image
		}
		hasWebCardChanges = true
	}

	// Update About if provided
	if input.WebCards.About.Description != "" || input.WebCards.About.LogoURL != "" || input.WebCards.About.EstablishmentYear != "" || input.WebCards.About.TotalProjects != "" || input.WebCards.About.ContactDetails.Name != "" || input.WebCards.About.ContactDetails.ProjectAddress != "" || input.WebCards.About.ContactDetails.Phone != "" || input.WebCards.About.ContactDetails.BookingLink != "" {
		if input.WebCards.About.Description != "" {
			newWebCards.About.Description = input.WebCards.About.Description
		}
		if input.WebCards.About.LogoURL != "" {
			newWebCards.About.LogoURL = input.WebCards.About.LogoURL
		}
		if input.WebCards.About.EstablishmentYear != "" {
			newWebCards.About.EstablishmentYear = input.WebCards.About.EstablishmentYear
		}
		if input.WebCards.About.TotalProjects != "" {
			newWebCards.About.TotalProjects = input.WebCards.About.TotalProjects
		}
		if input.WebCards.About.ContactDetails.Name != "" {
			newWebCards.About.ContactDetails.Name = input.WebCards.About.ContactDetails.Name
		}
		if input.WebCards.About.ContactDetails.ProjectAddress != "" {
			newWebCards.About.ContactDetails.ProjectAddress = input.WebCards.About.ContactDetails.ProjectAddress
		}
		if input.WebCards.About.ContactDetails.Phone != "" {
			newWebCards.About.ContactDetails.Phone = input.WebCards.About.ContactDetails.Phone
		}
		if input.WebCards.About.ContactDetails.BookingLink != "" {
			newWebCards.About.ContactDetails.BookingLink = input.WebCards.About.ContactDetails.BookingLink
		}
		hasWebCardChanges = true
	}

	// Update Faqs if provided
	if len(input.WebCards.Faqs) > 0 {
		if len(newWebCards.Faqs) == 0 {
			newWebCards.Faqs = make([]schema.FAQ, len(input.WebCards.Faqs))
		}
		newWebCards.Faqs = input.WebCards.Faqs
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

	if _, err := project.Save(context.Background()); err != nil {
		logger.Get().Error().Err(err).Msg("Failed to update project")
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		logger.Get().Error().Err(err).Msg("Failed to commit transaction")
		return nil, err
	}

	// Fetch the updated project with edges loaded
	projectWithEdges, err := r.db.Project.Query().
		Where(projectEnt.ID(input.ProjectID)).
		WithDeveloper().
		WithLocation().
		Only(context.Background())
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to fetch updated project with edges")
		return nil, err
	}
	return projectWithEdges, nil
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

func (r *repository) GetAllProjects(filters map[string]interface{}) ([]*ent.Project, error) {
	ctx := context.Background()

	// Start building the query
	query := r.db.Project.Query().Where(projectEnt.IsDeletedEQ(false))

	// Apply filters if any are set
	if len(filters) > 0 {
		predicates := []predicateEnt.Project{}

		// Apply boolean filters
		if isPremium, ok := filters["is_premium"].(bool); ok && isPremium {
			predicates = append(predicates, projectEnt.IsPremiumEQ(true))
		}
		if isPriority, ok := filters["is_priority"].(bool); ok && isPriority {
			predicates = append(predicates, projectEnt.IsPriorityEQ(true))
		}
		if isFeatured, ok := filters["is_featured"].(bool); ok && isFeatured {
			predicates = append(predicates, projectEnt.IsFeaturedEQ(true))
		}

		// Apply location filter
		if locationID, ok := filters["location_id"].(string); ok && locationID != "" {
			query = query.Where(projectEnt.HasLocationWith(locationEnt.ID(locationID)))
		}

		// Apply developer filter
		if developerID, ok := filters["developer_id"].(string); ok && developerID != "" {
			query = query.Where(projectEnt.HasDeveloperWith(developerEnt.ID(developerID)))
		}

		// Apply name filter
		if name, ok := filters["name"].(string); ok && name != "" {
			query = query.Where(projectEnt.NameContainsFold(name))
		}

		if city, ok := filters["city"].(string); ok && city != "" {
			// Remove quotes if present
			city = strings.Trim(city, "\"")
			// Filter projects that have a location with matching city
			query = query.Where(projectEnt.HasLocationWith(locationEnt.CityEQ(city)))
		}

		// Apply type filter
		if projectType, ok := filters["type"].(string); ok && projectType != "" {
			query = query.Where(projectEnt.ProjectTypeEQ(projectEnt.ProjectType(projectType)))
		}

		if len(predicates) > 0 {
			query = query.Where(projectEnt.Or(predicates...))
		}
	}

	// Execute the query with eager loading of related entities
	projects, err := query.
		Order(ent.Desc(projectEnt.FieldID)).
		WithDeveloper().
		WithLocation().
		All(ctx)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *repository) GetProjectByURL(url string) (*ent.Project, error) {

	project, err := r.db.Project.Query().
		Where(
			projectEnt.IsDeletedEQ(false),
			func(s *sql.Selector) {
				s.Where(sql.ExprP("meta_info->>'canonical' = ?", url))
			},
		).
		First(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		logger.Get().Error().Err(err).Msg("Failed to get project by URL")
		return nil, err
	}

	return project, nil
}
