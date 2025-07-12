package migration_jobs

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/VI-IM/im_backend_go/shared/logger"
	"github.com/rs/zerolog/log"
)

var projectIDToAmenitiesMap = map[int64][]int64{}

// fetch all tables from legacy database

func FetchCityByID(ctx context.Context, id int64) (LCity, error) {
	query := `SELECT id, city_name, city_url, created_date, is_active, state_name, updated_date, phone_number FROM city WHERE id = ?`
	rows, err := legacyDB.QueryContext(ctx, query, id)
	if err != nil {
		return LCity{}, err
	}
	defer rows.Close()

	var city LCity

	// 👇 This is the key fix
	if rows.Next() {
		if err := rows.Scan(&city.ID, &city.Name, &city.URL, &city.CreatedDate, &city.IsActive, &city.StateName, &city.UpdatedDate, &city.Phone); err != nil {
			return LCity{}, err
		}
		return city, nil
	}

	// If no rows found
	return LCity{}, fmt.Errorf("no city found with ID %d", id)
}

func FetchAllLocality(ctx context.Context) ([]LLocality, error) {
	query := `SELECT id, created_date, locality_name, locality_url, updated_date, city_id FROM locality`
	rows, err := legacyDB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var localities []LLocality
	for rows.Next() {
		var locality LLocality
		if err := rows.Scan(
			&locality.ID,
			&locality.CreatedDate,
			&locality.Name,
			&locality.URL,
			&locality.UpdatedDate,
			&locality.CityID,
		); err != nil {
			return nil, err
		}
		localities = append(localities, locality)
	}
	// Check for errors after loop
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return localities, nil
}

func FetchAllDevelopers(ctx context.Context) ([]LDeveloper, error) {
	log.Info().Msg("Fetching all developers")
	query := `SELECT id, about, alt_developer_logo, created_date, developer_address, developer_legal_name, developer_logo, developer_name, developer_url, disclaimer, established_year, is_active, is_verified, overview, project_done_no, updated_date, city_name, phone FROM developer`
	rows, err := legacyDB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var developers []LDeveloper
	for rows.Next() {
		var developer LDeveloper
		if err := rows.Scan(&developer.ID, &developer.About, &developer.AltDeveloperLogo, &developer.CreatedDate, &developer.DeveloperAddress, &developer.DeveloperLegalName, &developer.DeveloperLogo, &developer.DeveloperName, &developer.DeveloperURL, &developer.Disclaimer, &developer.EstablishedYear, &developer.IsActive, &developer.IsVerified, &developer.Overview, &developer.ProjectDoneNo, &developer.UpdatedDate, &developer.CityName, &developer.Phone); err != nil {
			return nil, err
		}
		developers = append(developers, developer)
	}
	return developers, nil
}

func FetchPropertyConfigurationByID(ctx context.Context, id int64) (*LPropertyConfiguration, error) {
	query := `SELECT * FROM project_configuration WHERE id = ?`
	rows, err := legacyDB.QueryContext(ctx, query, id)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to fetch property configuration by ID %d", id)
		return nil, err
	}
	defer rows.Close()

	var configuration LPropertyConfiguration

	for rows.Next() {
		if err := rows.Scan(&configuration.ID, &configuration.CreatedDate, &configuration.ProjectConfigurationName, &configuration.UpdatedDate, &configuration.ConfigurationTypeID); err != nil {
			log.Error().Err(err).Msgf("Failed to scan property configuration by ID %d", id)
			return nil, err
		}
	}

	return &configuration, nil
}

func FetchhAllProject(ctx context.Context) ([]LProject, error) {
	fmt.Println("Fetching all projects")
	query := `SELECT * FROM project`
	rows, err := legacyDB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []LProject
	for rows.Next() {
		var project LProject
		if err := rows.Scan(
			&project.ID,
			&project.AltProjectLogo,
			&project.AltSitePlanImg,
			&project.AmenitiesPara,
			&project.AvailableUnit,
			&project.CoverPhoto,
			&project.CreatedDate,
			&project.FloorPara,
			&project.IsDeleted,
			&project.IsFeatured,
			&project.IsPremium,
			&project.IsPriority,
			&project.LocationMap,
			&project.LocationPara,
			&project.MetaDescription,
			&project.MetaTitle,
			&project.OverviewPara,
			&project.PaymentPara,
			&project.PriceListPara,
			&project.ProjectAbout,
			&project.ProjectAddress,
			&project.ProjectArea,
			&project.ProjectBrochure,
			&project.ProjectConfigurations,
			&project.ProjectDescription,
			&project.ProjectLaunchDate,
			&project.ProjectLocationURL,
			&project.ProjectLogo,
			&project.ProjectName,
			&project.ProjectPossessionDate,
			&project.ProjectRERA,
			&project.ProjectSchema,
			&project.ProjectUnits,
			&project.ProjectURL,
			&project.ProjectVideoCount,
			&project.ProjectVideos,
			&project.ReraLink,
			&project.ShortAddress,
			&project.SitePlanImg,
			&project.SitePlanPara,
			&project.Status,
			&project.TotalFloor,
			&project.TotalTowers,
			&project.UpdatedDate,
			&project.USP,
			&project.VideoPara,
			&project.WhyPara,
			&project.PropertyConfigTypeID,
			&project.DeveloperID,
			&project.LocalityID,
			&project.UserID,
			&project.MetaKeywords,
		); err != nil {
			return nil, err
		}
		// fmt.Println("Project Videos", project.ProjectVideos)
		projects = append(projects, project)
	}
	return projects, nil
}

func fetchAllProperty(ctx context.Context) ([]LProperty, error) {
	query := `SELECT * FROM property`
	rows, err := legacyDB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var properties []LProperty
	for rows.Next() {
		var property LProperty
		if err := rows.Scan(
			&property.ID,
			&property.About,
			&property.AgeOfProperty,
			&property.AmenitiesPara,
			&property.Balcony,
			&property.Bathrooms,
			&property.Bedrooms,
			&property.BuiltupArea,
			&property.CoverPhoto,
			&property.CoveredParking,
			&property.CreatedDate,
			&property.Facing,
			&property.FloorImage2D,
			&property.FloorImage3D,
			&property.FloorPara,
			&property.Floors,
			&property.FurnishingType,
			&property.Images,
			&property.IsDeleted,
			&property.IsFeatured,
			&property.Latlong,
			&property.ListingType,
			&property.LocationMap,
			&property.LocationAdvantage,
			&property.LocationPara,
			&property.LogoImage,
			&property.MetaDescription,
			&property.MetaKeywords,
			&property.MetaTitle,
			&property.OpenParking,
			&property.OverviewPara,
			&property.PossessionDate,
			&property.PossessionStatus,
			&property.Price,
			&property.ProductSchema,
			&property.PropertyAddress,
			&property.PropertyName,
			&property.PropertyURL,
			&property.PropertyVideo,
			&property.Rera,
			&property.Size,
			&property.UpdatedByID,
			&property.UpdatedDate,
			&property.USP,
			&property.VideoPara,
			&property.ConfigurationID,
			&property.DeveloperID,
			&property.LocalityID,
			&property.ProjectID,
			&property.UserID,
			&property.Highlights,
			&property.LocaionMap,
		); err != nil {
			return nil, err
		}
		log.Info().Msgf("Product Schema: %+v", property.ProductSchema)
		properties = append(properties, property)
	}
	return properties, nil
}

func FetchProjectByID(ctx context.Context, id int64) (*LProject, error) {
	query := `SELECT * FROM project WHERE id = ?`
	rows, err := legacyDB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var project LProject
	if err := rows.Scan(&project.ID, &project.ProjectName, &project.ProjectDescription, &project.Status, &project.ProjectConfigurations, &project.TotalFloor, &project.TotalTowers, &project.ProjectLaunchDate, &project.ProjectPossessionDate, &project.MetaTitle, &project.MetaDescription, &project.MetaKeywords, &project.ProjectURL, &project.ProjectSchema, &project.ProjectLogo, &project.ProjectBrochure, &project.ProjectVideos, &project.ProjectVideoCount, &project.IsFeatured, &project.IsPremium, &project.IsPriority, &project.IsDeleted, &project.DeveloperID, &project.LocalityID, &project.UserID, &project.CreatedDate, &project.UpdatedDate); err != nil {
		return nil, err
	}
	return &project, nil
}

func FetchProjectConfigurationsByID(ctx context.Context, id int64) (*LPropertyConfiguration, error) {
	query := `SELECT * FROM property_configuration WHERE id = ?`
	rows, err := legacyDB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configuration LPropertyConfiguration
	if err := rows.Scan(&configuration.ID, &configuration.CreatedDate, &configuration.ProjectConfigurationName, &configuration.UpdatedDate, &configuration.ConfigurationTypeID); err != nil {
		return nil, err
	}
	return &configuration, nil
}

func FetchLocalityByID(ctx context.Context, id int64) (*LLocality, error) {
	query := `SELECT * FROM locality WHERE id = ?`
	rows, err := legacyDB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locality LLocality
	if err := rows.Scan(&locality.ID, &locality.CreatedDate, &locality.Name, &locality.URL, &locality.UpdatedDate, &locality.CityID); err != nil {
		return nil, err
	}
	return &locality, nil
}

func FetchDeveloperByID(ctx context.Context, id int64) (*LDeveloper, error) {
	query := `SELECT * FROM developer WHERE id = ?`
	rows, err := legacyDB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var developer LDeveloper
	if rows.Next() {
		if err := rows.Scan(
			&developer.ID,
			&developer.About,
			&developer.AltDeveloperLogo,
			&developer.CreatedDate,
			&developer.DeveloperAddress,
			&developer.DeveloperLegalName,
			&developer.DeveloperLogo,
			&developer.DeveloperName,
			&developer.DeveloperURL,
			&developer.Disclaimer,
			&developer.EstablishedYear,
			&developer.IsActive,
			&developer.IsVerified,
			&developer.Overview,
			&developer.ProjectDoneNo,
			&developer.UpdatedDate,
			&developer.CityName,
			&developer.Phone,
		); err != nil {
			return nil, err
		}
		return &developer, nil
	}

	// No rows found
	return nil, sql.ErrNoRows
}

func FetchProjectImagesByProjectID(ctx context.Context, id int64) (*[]LProjectImage, error) {
	query := `SELECT * FROM project_image WHERE project_id = ?`
	rows, err := legacyDB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var propertyImages []LProjectImage
	for rows.Next() {
		var propertyImage LProjectImage
		if err := rows.Scan(&propertyImage.ProjectID, &propertyImage.ImageAltName, &propertyImage.ImageURL); err != nil {
			return nil, err
		}
		propertyImages = append(propertyImages, propertyImage)
	}
	return &propertyImages, nil
}

func FetchFloorPlansByProjectID(ctx context.Context, projectID int64) (*[]LFloorPlan, error) {
	query := `SELECT id, created_date, img_url, is_sold_out, price, size, title, updated_date, configuration_id, project_id, user_id FROM floorplan WHERE project_id = ?`
	rows, err := legacyDB.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var floorPlans []LFloorPlan
	for rows.Next() {
		var floorPlan LFloorPlan
		if err := rows.Scan(
			&floorPlan.ID,
			&floorPlan.CreatedDate,
			&floorPlan.ImgURL,
			&floorPlan.IsSoldOut,
			&floorPlan.Price,
			&floorPlan.Size,
			&floorPlan.Title,
			&floorPlan.UpdatedDate,
			&floorPlan.ConfigurationID,
			&floorPlan.ProjectID,
			&floorPlan.UserID,
		); err != nil {
			return nil, err
		}
		floorPlans = append(floorPlans, floorPlan)
	}
	return &floorPlans, nil
}

func FetchReraByProjectID(ctx context.Context, projectID int64) ([]*LRera, error) {
	query := `SELECT * FROM rera_info WHERE project_id = ?`
	rows, err := legacyDB.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reras []*LRera
	for rows.Next() {
		var rera LRera
		if err := rows.Scan(&rera.ID, &rera.CreatedDate, &rera.Phase, &rera.ProjectReraName, &rera.QRImages, &rera.ReraNumber, &rera.Status, &rera.UpdatedDate, &rera.ProjectID, &rera.UserID); err != nil {
			return nil, err
		}
		reras = append(reras, &rera)
	}
	return reras, nil
}

func FetchFloorPlanByProjectID(ctx context.Context, projectID int64) (*[]LFloorPlan, error) {
	query := `SELECT * FROM floorplan WHERE project_id = ?`
	rows, err := legacyDB.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var floorPlans []LFloorPlan
	for rows.Next() {
		var floorPlan LFloorPlan
		if err := rows.Scan(&floorPlan.ID, &floorPlan.CreatedDate, &floorPlan.ImgURL, &floorPlan.IsSoldOut, &floorPlan.Price, &floorPlan.Size, &floorPlan.Title, &floorPlan.UpdatedDate, &floorPlan.ConfigurationID, &floorPlan.ProjectID, &floorPlan.UserID); err != nil {
			return nil, err
		}
		floorPlans = append(floorPlans, floorPlan)
	}
	return &floorPlans, nil
}

type LProjectAmenity struct {
	ProjectID int64 `json:"project_id"`
	AmenityID int64 `json:"amenity_id"`
}

func FetchAmenityByID(ctx context.Context, id int64) (*LAmenity, error) {
	query := `SELECT * FROM amenities WHERE id = ?`
	rows, err := legacyDB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var amenity LAmenity
	if rows.Next() {
		if err := rows.Scan(
			&amenity.ID,
			&amenity.AmenitiesCategory,
			&amenity.AmenitiesName,
			&amenity.AmenitiesURL,
			&amenity.CreatedDate,
			&amenity.UpdatedDate,
		); err != nil {
			return nil, err
		}
		return &amenity, nil
	}
	// No rows found
	return nil, sql.ErrNoRows
}

func FetchProjectAmenitiesByProjectID(ctx context.Context, projectID int64) ([]*LAmenity, error) {
	query := `SELECT * FROM project_amenities WHERE project_id = ?`
	rows, err := legacyDB.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var amenityIDs []int64
	for rows.Next() {
		var projectAmenity LProjectAmenity
		if err := rows.Scan(&projectAmenity.ProjectID, &projectAmenity.AmenityID); err != nil {
			return nil, err
		}
		amenityIDs = append(amenityIDs, projectAmenity.AmenityID)
	}

	var amenities []*LAmenity
	for _, amenityID := range amenityIDs {
		amenity, err := FetchAmenityByID(ctx, amenityID)
		if err != nil {
			return nil, err
		}
		amenities = append(amenities, amenity)
	}

	return amenities, nil
}

func FetchPaymentPlansByProjectID(ctx context.Context, projectID int64) ([]*LPaymentPlan, error) {
	query := `SELECT * FROM payment_plan WHERE project_id = ?`
	rows, err := legacyDB.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var paymentPlans []*LPaymentPlan
	for rows.Next() {
		var paymentPlan LPaymentPlan
		if err := rows.Scan(&paymentPlan.ID, &paymentPlan.PaymentPlanName, &paymentPlan.PaymentPlanValue, &paymentPlan.ProjectID); err != nil {
			return nil, err
		}
		paymentPlans = append(paymentPlans, &paymentPlan)
	}
	return paymentPlans, nil
}

func FetchFaqsByProjectID(ctx context.Context, projectID int64) ([]*LFAQ, error) {
	query := `SELECT * FROM faq WHERE project_id = ?`
	rows, err := legacyDB.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var faqs []*LFAQ
	for rows.Next() {
		var faq LFAQ
		if err := rows.Scan(&faq.ID, &faq.Question, &faq.Answer, &faq.ProjectID); err != nil {
			return nil, err
		}
		faqs = append(faqs, &faq)
	}
	return faqs, nil
}

func FetchPropertyConfigurationTypeByID(ctx context.Context, id int64) (*LPropertyConfigurationType, error) {
	query := `SELECT * FROM project_configuration_type WHERE id = ?`
	rows, err := legacyDB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var propertyConfigurationType LPropertyConfigurationType
	for rows.Next() {
		if err := rows.Scan(&propertyConfigurationType.ID, &propertyConfigurationType.ConfigurationTypeName, &propertyConfigurationType.CreatedDate, &propertyConfigurationType.PropertyType, &propertyConfigurationType.UpdatedDate); err != nil {
			return nil, err
		}
	}

	return &propertyConfigurationType, nil
}
func FetchProjectConfigurationByID(ctx context.Context, id int64) (*LPropertyConfigurationType, error) {
	query := `SELECT * FROM project_configuration_type WHERE id = ?`
	rows, err := legacyDB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projectConfiguration LPropertyConfigurationType
	for rows.Next() {
		if err := rows.Scan(&projectConfiguration.ID, &projectConfiguration.ConfigurationTypeName, &projectConfiguration.CreatedDate, &projectConfiguration.PropertyType, &projectConfiguration.UpdatedDate); err != nil {
			return nil, err
		}
	}
	return &projectConfiguration, nil
}

func FetchAllBlogs(ctx context.Context) ([]LBlog, error) {
	query := `SELECT id, alt, coalesce(blog_schema, '[]'), blog_url, canonical, created_date, description, headings, images, 
					 is_priority, sub_headings, updated_date, user_id, meta_keywords, meta_title 
			  FROM blogs`
	rows, err := legacyDB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []LBlog
	for rows.Next() {
		var blog LBlog
		if err := rows.Scan(
			&blog.ID,
			&blog.Alt,
			&blog.BlogSchema,
			&blog.BlogURL,
			&blog.Canonical,
			&blog.CreatedDate,
			&blog.Description,
			&blog.Headings,
			&blog.Images,
			&blog.IsPriority,
			&blog.SubHeadings,
			&blog.UpdatedDate,
			&blog.UserID,
			&blog.MetaKeywords,
			&blog.MetaTitle,
		); err != nil {
			return nil, err
		}
		blogs = append(blogs, blog)
	}
	return blogs, nil
}

func fetchAllProjectIDs(client *http.Client) (*JavaGetProjectByIDResponse, error) {

	resp, err := client.Get(javaAPIBaseURL + getAllProjectsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch projects: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	logger.Get().Info().Msg(string(body))

	var projects JavaGetProjectByIDResponse
	if err := json.Unmarshal(body, &projects); err != nil {
		return nil, fmt.Errorf("failed to unmarshal projects: %v", err)
	}

	return &projects, nil
}
