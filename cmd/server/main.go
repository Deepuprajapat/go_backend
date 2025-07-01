package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"strconv"
	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/internal/application"
	s3client "github.com/VI-IM/im_backend_go/internal/client"
	"github.com/VI-IM/im_backend_go/internal/config"
	"github.com/VI-IM/im_backend_go/internal/database"
	"github.com/VI-IM/im_backend_go/internal/repository"
	"github.com/VI-IM/im_backend_go/internal/router"
	"github.com/VI-IM/im_backend_go/migration_jobs"
	"github.com/VI-IM/im_backend_go/shared/logger"
	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
	"github.com/joho/godotenv"
)

var (
	legacyDB *sql.DB
	newDB    *ent.Client
	txn      *ent.Tx
)

func main() {
	// Initialize logger
	logger.Init()
	ctx := context.Background()

	logger.Get().Info().Msg("Starting application...")

	// Load .env file
	if err := godotenv.Load(); err != nil {
		logger.Get().Fatal().Err(err).Msg("Error loading .env file")
	}

	if len(os.Args) > 1 && os.Args[1] == "run-migration" {

		var err error
		legacyDB, err = migration_jobs.NewLegacyDBConnection()
		if err != nil {
			logger.Get().Fatal().Err(err).Msg("Failed to connect to legacy database")
		}
		defer legacyDB.Close()

		newDB, err = migration_jobs.NewNewDBConnection()
		if err != nil {
			logger.Get().Fatal().Err(err).Msg("Failed to connect to new database")
		}
		defer newDB.Close()

		// Start a transaction for the entire migration process
		txn, err = newDB.BeginTx(ctx, nil)
		if err != nil {
			logger.Get().Fatal().Err(err).Msg("Failed to begin transaction")
		}
		defer txn.Rollback()

		logger.Get().Info().Msg("Migrating static site data------------>>>>>>>>>>>>>>>>>>>>")
		if err := migration_jobs.MigrateStaticSiteData(ctx, txn); err != nil {
			logger.Get().Fatal().Err(err).Msg("Failed to migrate static site data")
		}

		// Execute migrations in sequence

		logger.Get().Info().Msg("Migrating localities------------>>>>>>>>>>>>>>>>>>>>")
		if err := migration_jobs.MigrateLocality(ctx, txn); err != nil {
			logger.Get().Fatal().Err(err).Msg("Failed to migrate localities")
		}

		logger.Get().Info().Msg("Migrating developers------------>>>>>>>>>>>>>>>>>>>>")
		if err := migration_jobs.MigrateDeveloper(ctx, txn); err != nil {
			logger.Get().Fatal().Err(err).Msg("Failed to migrate developers")
		}

		logger.Get().Info().Msg("Migrating projects------------>>>>>>>>>>>>>>>>>>>>")
		if err := migration_jobs.MigrateProject(ctx, txn); err != nil {
			logger.Get().Fatal().Err(err).Msg("Failed to migrate projects")
		}

		logger.Get().Info().Msg("Migrating properties------------>>>>>>>>>>>>>>>>>>>>")
		if err := migration_jobs.MigrateProperty(ctx, txn); err != nil {
			logger.Get().Fatal().Err(err).Msg("Failed to migrate properties")
		}

		logger.Get().Info().Msg("Committing transaction------------>>>>>>>>>>>>>>>>>>>>")
		if err := txn.Commit(); err != nil {
			logger.Get().Fatal().Err(err).Msg("Failed to commit transaction")
		}

		logger.Get().Info().Msg("Migration completed successfully")
		return
	}

	// Load configuration
	if err := config.LoadConfig(); err != nil {
		logger.Get().Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Log configuration values
	cfg := config.GetConfig()

	client := database.NewClient(cfg.Database.URL)
	defer client.Close()

	s3Client, err := s3client.NewS3Client(cfg.S3.Bucket)
	if err != nil {
		logger.Get().Fatal().Err(err).Msg("Failed to create S3 client")
	}

	repo := repository.NewRepository(client)
	app := application.NewApplication(repo, s3Client)

	// Initialize router
	router.Init(app)

	// Start server
	logger.Get().Info().Msgf("Server starting on port %d", cfg.Port)
	if err := http.ListenAndServe(":"+strconv.Itoa(cfg.Port), router.Router); err != nil {
		logger.Get().Fatal().Err(err).Msg("Failed to start server")
	}
}
