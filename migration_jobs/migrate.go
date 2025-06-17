package migration_jobs

import (
	"context"
	"database/sql"

	"github.com/VI-IM/im_backend_go/ent/schema"
	"github.com/google/uuid"
)

// Use ent schema modal to migrate data from legacy database to new database

var (
	legacyToNewProjectIDMAP       = make(map[string]string)
	legacyToNewDeveloperIDMAP     = make(map[int64]string)
	legacyToNewLocalityIDMAP      = make(map[int64]string)
	legacyToNewConfigurationIDMAP = make(map[string]string)
	legacyToNewConfigTypeIDMAP    = make(map[string]string)
)

func migrateProject(ctx context.Context, db *sql.DB) error {
	// fetch all projects from legacy database

	// create project in iteration
	// map legacy project id to new project id
	// setDeveloperID(legacyToNewDeveloperIDMAP[lproject.DeveloperID])

	return nil
}

func migrateDeveloper(ctx context.Context, db *sql.DB) error {
	ldeveloper, err := FetchAllDevelopers(ctx, db)
	if err != nil {
		return err
	}

	for _, developer := range ldeveloper {
		id := uuid.New().String()
		legacyToNewDeveloperIDMAP[developer.ID] = id
		if err := newDB.Developer.Create().
			SetID(id).
			SetName(*developer.DeveloperName).
			SetLegalName(*developer.DeveloperLegalName).
			SetIdentifier(*developer.DeveloperName).
			SetEstablishedYear(int(*developer.EstablishedYear)).
			SetMediaContent(schema.DeveloperMediaContent{
				DeveloperAddress: *developer.DeveloperAddress,
				Phone:            *developer.Phone,
				DeveloperLogo:    *developer.DeveloperLogo,
				AltDeveloperLogo: *developer.AltDeveloperLogo,
				About:            *developer.About,
				Overview:         *developer.Overview,
				Disclaimer:       *developer.Disclaimer,
			}).
			SetIsVerified(developer.IsVerified).
			Exec(ctx); err != nil {
			return err
		}
	}

	return nil
}

func migrateLocality(ctx context.Context, db *sql.DB) error {
	//new location id will be generated
	llocality, err := FetchAllLocality(ctx, db)
	if err != nil {
		return err
	}

	for _, locality := range llocality {
		city, err := FetchCityByID(ctx, db, *locality.CityID)
		if err != nil {
			return err
		}
		id := uuid.New().String()
		legacyToNewLocalityIDMAP[locality.ID] = id
		if err := newDB.Location.Create().
			SetID(id).
			SetLocalityName(*locality.Name).
			SetCity(*city.Name).
			SetState(*city.StateName).
			SetPhoneNumber(*city.Phone).
			SetCountry("India").
			SetPincode(*locality.URL).
			SetIsActive(true).
			Exec(ctx); err != nil {
			return err
		}
	}

	return nil
}
