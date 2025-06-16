package migration_jobs

import (
	"context"
	"database/sql"
)

// Use ent schema modal to migrate data from legacy database to new database

var (
	legacyToNewProjectIDMAP       = make(map[string]string)
	legacyToNewDeveloperIDMAP     = make(map[string]string)
	legacyToNewLocalityIDMAP      = make(map[string]string)
	legacyToNewConfigurationIDMAP = make(map[string]string)
	legacyToNewConfigTypeIDMAP    = make(map[string]string)
)

func migrateProject(ctx context.Context, db *sql.DB) error {
	// fetch all projects from legacy database

	// create project in iteration
	// map legacy project id to new project id

	return nil
}

func migrateDeveloper(ctx context.Context, db *sql.DB) error {

	return nil
}

func migrateLocality(ctx context.Context, db *sql.DB) error {
	return nil
}

func migrateCity(ctx context.Context, db *sql.DB) error {
	return nil
}

func migrateConfiguration(ctx context.Context, db *sql.DB) error {
	return nil
}
