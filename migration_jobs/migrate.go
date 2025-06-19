package migration_jobs

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/schema"
	"github.com/rs/zerolog/log"
)

// Use ent schema modal to migrate data from legacy database to new database

var (
	legacyToNewProjectIDMAP       = make(map[int64]string)
	legacyToNewDeveloperIDMAP     = make(map[int64]string)
	legacyToNewLocalityIDMAP      = make(map[int64]string)
	legacyToNewConfigurationIDMAP = make(map[string]string)
	legacyToNewConfigTypeIDMAP    = make(map[string]string)
	legacyToNewPropertyIDMAP      = make(map[int64]string)
)

// func migrateProject(ctx context.Context, db *sql.DB) error {
// 	projects, err := FetchhAllProject(ctx, db)
// 	if err != nil {
// 		return err
// 	}

// 	for _, project := range projects {
// 		id := uuid.New().String()
// 		legacyToNewProjectIDMAP[project.ID] = id
// 		if err := newDB.Project.Create().
// 			SetID(id).
// 			SetName(*project.ProjectName).
// 			SetDescription(*project.ProjectDescription).
// 			SetStatus(*project.Status).
// 			SetProjectConfigurations(*project.ProjectConfigurations).
// 			SetTotalFloor(int(*project.TotalFloor)).
// 			SetTotalTowers(int(*project.TotalTowers)).
// 			SetTimelineInfo(schema.TimelineInfo{
// 				ProjectLaunchDate:     *project.ProjectLaunchDate,
// 				ProjectPossessionDate: *project.ProjectPossessionDate,
// 			}).
// 			SetMetaInfo(schema.SEOMeta{
// 				Title:         *project.MetaTitle,
// 				Description:   *project.MetaDescription,
// 				Keywords:      *project.MetaKeywords,
// 				Canonical:     *project.ProjectURL,
// 				ProjectSchema: *project.ProjectSchema,
// 			}).
// 			SetWebCards(schema.ProjectWebCards{

// 			}).
// 			SetIsFeatured(project.IsFeatured).
// 			SetIsPremium(project.IsPremium).
// 			SetIsPriority(project.IsPriority).
// 			SetIsDeleted(project.IsDeleted).
// 			SetDeveloperID(legacyToNewDeveloperIDMAP[*project.DeveloperID]).
// 			SetLocalityID(legacyToNewLocalityIDMAP[*project.LocalityID]).
// 			Exec(ctx); err != nil {
// 			return err
// 		}

// 	}

// 	// fetch all projects from legacy database

// 	// create project in iteration
// 	// map legacy project id to new project id
// 	// setDeveloperID(legacyToNewDeveloperIDMAP[lproject.DeveloperID])

// 	return nil
// }

func MigrateDeveloper(ctx context.Context, db *sql.DB, newDB *ent.Client) error {
	log.Info().Msg("Migrating developers")
	ldeveloper, err := FetchAllDevelopers(ctx, db)
	if err != nil {
		return err
	}

	for _, developer := range ldeveloper {
		id := fmt.Sprintf("%x", sha256.Sum256([]byte(strconv.FormatInt(developer.ID, 10))))[:16]
		legacyToNewDeveloperIDMAP[developer.ID] = id
		if err := newDB.Developer.Create().
			SetID(id).
			SetName(safeStr(developer.DeveloperName)).
			SetLegalName(safeStr(developer.DeveloperLegalName)).
			SetIdentifier(safeStr(developer.DeveloperName)).
			SetEstablishedYear(safeInt(developer.EstablishedYear)).
			SetMediaContent(schema.DeveloperMediaContent{
				DeveloperAddress: safeStr(developer.DeveloperAddress),
				Phone:            safeStr(developer.Phone),
				DeveloperLogo:    safeStr(developer.DeveloperLogo),
				AltDeveloperLogo: safeStr(developer.AltDeveloperLogo),
				About:            safeStr(developer.About),
				Overview:         safeStr(developer.Overview),
				Disclaimer:       safeStr(developer.Disclaimer),
			}).
			SetIsVerified(developer.IsVerified != nil && *developer.IsVerified).
			Exec(ctx); err != nil {
			return err
		}

	}
	return nil
}

func MigrateLocality(ctx context.Context, db *sql.DB, newDB *ent.Client) error {
	//new location id will be generated
	llocality, err := FetchAllLocality(ctx, db)
	if err != nil {
		return err
	}

	for _, locality := range llocality {

		city, err := FetchCityByID(ctx, db, *locality.CityID)

		log.Info().Msgf("Fetched city %+v", city)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to fetch city for locality ID %d", locality.ID)
			continue
		}

		if newDB == nil {
			return fmt.Errorf("newDB is nil — database connection not initialized")
		}

		phoneInt, err := parsePhoneJSONToString(city.Phone)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to convert phone for locality ID %d", locality.ID)
		}

		id := fmt.Sprintf("%x", sha256.Sum256([]byte(strconv.FormatInt(locality.ID, 10))))[:16]
		legacyToNewLocalityIDMAP[locality.ID] = id
		if err := newDB.Location.Create().
			SetID(id).
			SetLocalityName(safeStr(locality.Name)).
			SetCity(safeStr(city.Name)).
			SetState(safeStr(city.StateName)).
			SetPhoneNumber(*phoneInt).
			SetCountry("India").
			SetPincode("112222").
			SetIsActive(true).
			Exec(ctx); err != nil {
			log.Error().Err(err).Msgf("Failed to insert locality ID %d", locality.ID)
			continue
		}
	}
	return nil
}

func migrateProperty(ctx context.Context, db *sql.DB) error {
	properties, err := fetchAllProperty(ctx, db)
	if err != nil {
		return err
	}
	for _, property := range properties {
		id := fmt.Sprintf("%x", sha256.Sum256([]byte(strconv.FormatInt(property.ID, 10))))[:16]
		legacyToNewPropertyIDMAP[property.ID] = id
		project, err := FetchProjectByID(ctx, db, *property.ProjectID)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to fetch project for property ID %d", property.ID)
			continue
		}
		projectConfigurations, err := FetchProjectConfigurationsByID(ctx, db, *property.ConfigurationID)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to fetch project configurations for property ID %d", property.ID)
			continue
		}
		locality, err := FetchLocalityByID(ctx, db, *property.LocalityID)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to fetch locality for property ID %d", property.ID)
			continue
		}
		developer, err := FetchDeveloperByID(ctx, db, *property.DeveloperID)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to fetch developer for property ID %d", property.ID)
		}

		projectImages, err := FetchProjectImagesByProjectID(ctx, db, *property.ProjectID)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to fetch project images for property ID %d", property.ID)
			continue
		}

		propertyImages, err := parsePropertyImagesFromProjectImages(projectImages)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to parse property images for property ID %d", property.ID)
			continue
		}
		floorPlans, err := FetchFloorPlansByProjectID(ctx, db, *property.ProjectID)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to fetch floor plans for property ID %d", property.ID)
			continue
		}
		webCards, err := parseWebCardsFromProject(project)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to parse web cards for property ID %d", property.ID)
			continue
		}

		// property Type and configurationtype(Ground Floor, Apartment, etc.) name from projectConfigurations table

		if err := newDB.Property.Create().
			SetID(id).
			SetName(*property.PropertyName).
			SetPropertyImages(*propertyImages).
			SetProjectID(legacyToNewProjectIDMAP[*property.ProjectID]).
			SetWebCards(*webCards).
			Exec(ctx); err != nil {
			return err
		}
	}
	return nil
}
