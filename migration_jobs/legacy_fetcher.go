package migration_jobs

import (
	"context"
	"database/sql"
)

// fetch all tables from legacy database

func FetchCityByID(ctx context.Context, db *sql.DB, id int64) (LCity, error) {
	query := `SELECT * FROM city WHERE id = ?`
	rows, err := db.QueryContext(ctx, query, id)
	if err != nil {
		return LCity{}, err
	}
	defer rows.Close()

	var city LCity
	if err := rows.Scan(&city.ID, &city.Name, &city.URL, &city.CreatedDate, &city.IsActive, &city.StateName, &city.UpdatedDate, &city.Phone); err != nil {
		return LCity{}, err
	}
	return city, nil
}

func FetchAllLocality(ctx context.Context, db *sql.DB) ([]LLocality, error) {
	query := `SELECT * FROM locality`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var localities []LLocality
	for rows.Next() {
		var locality LLocality
		if err := rows.Scan(&locality.ID, &locality.Name, &locality.URL, &locality.CreatedDate, &locality.UpdatedDate, &locality.CityID); err != nil {
			return nil, err
		}
		localities = append(localities, locality)
	}
	return localities, nil
}

func FetchAllDevelopers(ctx context.Context, db *sql.DB) ([]LDeveloper, error) {
	query := `SELECT * FROM developer`
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
