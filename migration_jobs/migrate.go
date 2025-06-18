package migration_jobs

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/schema"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// Use ent schema modal to migrate data from legacy database to new database

var (
	legacyToNewProjectIDMAP       = make(map[string]string)
	legacyToNewDeveloperIDMAP     = make(map[int64]string)
	legacyToNewLocalityIDMAP      = make(map[int64]string)
	legacyToNewConfigurationIDMAP = make(map[string]string)
	legacyToNewConfigTypeIDMAP    = make(map[string]string)
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
		id := uuid.New().String()
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

		// phoneInt, err := parsePhoneJSONToInt32(city.Phone)
		// if err != nil {
		// 	log.Error().Err(err).Msgf("Failed to convert phone for locality ID %d", locality.ID)
		// 	continue
		// }

		id := uuid.New().String()
		legacyToNewLocalityIDMAP[locality.ID] = id
		if err := newDB.Location.Create().
			SetID(id).
			SetLocalityName(safeStr(locality.Name)).
			SetCity(safeStr(city.Name)).
			SetState(safeStr(city.StateName)).
			SetPhoneNumber(extractNumericPhone(*city.Phone)).
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

// func migrateProperty(ctx context.Context, db *sql.DB) error {
// 	properties, err := fetchAllProperty(ctx, db)
// 	if err != nil {
// 		return err
// 	}
// 	for _, property := range properties {
// 		id := uuid.New().String()
// 		legacyToNewPropertyIDMAP[property.ID] = id
// 		if err := newDB.Property.Create().
// 			SetID(id).
// 			SetName(*property.PropertyName).
// 			SetDescription(*property.PropertyAddress).
// 			SetIsFeatured(property.IsFeatured).
// 			SetIsPremium(property.IsPremium).
// 			SetIsPriority(property.IsPriority).
// 			SetIsDeleted(property.IsDeleted).
// 			SetDeveloperID(legacyToNewDeveloperIDMAP[*property.DeveloperID]).
// 			SetLocalityID(legacyToNewLocalityIDMAP[*property.LocalityID]).
// 			SetProjectID(legacyToNewProjectIDMAP[*property.ProjectID]).
// 			SetStatus(enums.PropertyStatus(*property.Status)).
// 			Exec(ctx); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
