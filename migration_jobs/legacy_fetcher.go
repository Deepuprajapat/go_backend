package migration_jobs

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rs/zerolog/log"
)

// fetch all tables from legacy database

func FetchCityByID(ctx context.Context, db *sql.DB, id int64) (LCity, error) {
	query := `SELECT id, city_name, city_url, created_date, is_active, state_name, updated_date, phone_number FROM city WHERE id = ?`
	rows, err := db.QueryContext(ctx, query, id)
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

func FetchAllLocality(ctx context.Context, db *sql.DB) ([]LLocality, error) {
	query := `SELECT id, created_date, locality_name, locality_url, updated_date, city_id FROM locality`
	rows, err := db.QueryContext(ctx, query)
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
	log.Info().Msgf("Fetched localities %+v", localities)
	return localities, nil
}

func FetchAllDevelopers(ctx context.Context, db *sql.DB) ([]LDeveloper, error) {
	log.Info().Msg("Fetching all developers")
	query := `SELECT id, about, alt_developer_logo, created_date, developer_address, developer_legal_name, developer_logo, developer_name, developer_url, disclaimer, established_year, is_active, is_verified, overview, project_done_no, updated_date, city_name, phone FROM developer`
	rows, err := db.QueryContext(ctx, query)
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

func FetchConfigurationByID(ctx context.Context, db *sql.DB, id int64) (*LPropertyConfiguration, error) {
	query := `SELECT * FROM property_configuration WHERE id = ?`
	rows, err := db.QueryContext(ctx, query, id)
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

func FetchhAllProject(ctx context.Context, db *sql.DB) ([]LProject, error) {
	query := `SELECT * FROM project`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []LProject
	for rows.Next() {
		var project LProject
		if err := rows.Scan(&project.ID, &project.ProjectName, &project.ProjectDescription, &project.Status, &project.ProjectConfigurations, &project.TotalFloor, &project.TotalTowers, &project.ProjectLaunchDate, &project.ProjectPossessionDate, &project.MetaTitle, &project.MetaDescription, &project.MetaKeywords, &project.ProjectURL, &project.ProjectSchema, &project.ProjectLogo, &project.ProjectBrochure, &project.ProjectVideos, &project.ProjectVideoCount, &project.IsFeatured, &project.IsPremium, &project.IsPriority, &project.IsDeleted, &project.DeveloperID, &project.LocalityID, &project.UserID, &project.CreatedDate, &project.UpdatedDate); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func fetchAllProperty(ctx context.Context, db *sql.DB) ([]LProperty, error) {
	query := `SELECT * FROM property`
	rows, err := db.QueryContext(ctx, query)
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
		properties = append(properties, property)
	}
	return properties, nil
}

func FetchProjectByID(ctx context.Context, db *sql.DB, id int64) (*LProject, error) {
	query := `SELECT * FROM project WHERE id = ?`
	rows, err := db.QueryContext(ctx, query, id)
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

func FetchProjectConfigurationsByID(ctx context.Context, db *sql.DB, id int64) (*LPropertyConfiguration, error) {
	query := `SELECT * FROM property_configuration WHERE id = ?`
	rows, err := db.QueryContext(ctx, query, id)
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

func FetchLocalityByID(ctx context.Context, db *sql.DB, id int64) (*LLocality, error) {
	query := `SELECT * FROM locality WHERE id = ?`
	rows, err := db.QueryContext(ctx, query, id)
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

func FetchDeveloperByID(ctx context.Context, db *sql.DB, id int64) (*LDeveloper, error) {
	query := `SELECT * FROM developer WHERE id = ?`
	rows, err := db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var developer LDeveloper
	if err := rows.Scan(&developer.ID, &developer.About, &developer.AltDeveloperLogo, &developer.CreatedDate, &developer.DeveloperAddress, &developer.DeveloperLegalName, &developer.DeveloperLogo, &developer.DeveloperName, &developer.DeveloperURL, &developer.Disclaimer, &developer.EstablishedYear, &developer.IsActive, &developer.IsVerified, &developer.Overview, &developer.ProjectDoneNo, &developer.UpdatedDate, &developer.CityName, &developer.Phone); err != nil {
		return nil, err
	}
	return &developer, nil
}

func FetchProjectImagesByProjectID(ctx context.Context, db *sql.DB, id int64) (*[]LProjectImage, error) {
	query := `SELECT * FROM project_image WHERE project_id = ?`
	rows, err := db.QueryContext(ctx, query, id)
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
