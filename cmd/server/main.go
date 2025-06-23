package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/internal/application"
	"github.com/VI-IM/im_backend_go/internal/config"
	"github.com/VI-IM/im_backend_go/internal/database"
	"github.com/VI-IM/im_backend_go/internal/repository"
	"github.com/VI-IM/im_backend_go/internal/router"
	"github.com/VI-IM/im_backend_go/migration_jobs"
	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	legacyDB *sql.DB
	newDB    *ent.Client
	txn      *ent.Tx
)

func main() {
	// Initialize zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	})
	ctx := context.Background()

	if len(os.Args) > 1 && os.Args[1] == "run-migration" {
		var err error
		legacyDB, err = migration_jobs.NewLegacyDBConnection()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to connect to legacy database")
		}
		defer legacyDB.Close()

		newDB, err = migration_jobs.NewNewDBConnection()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to connect to new database")
		}
		defer newDB.Close()

		// Start a transaction for the entire migration process
		txn, err = newDB.BeginTx(ctx, nil)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to begin transaction")
		}
		defer txn.Rollback()

		// Execute migrations in sequence

		log.Info().Msg("Migrating localities------------>>>>>>>>>>>>>>>>>>>>")
		if err := migration_jobs.MigrateLocality(ctx, txn); err != nil {
			log.Fatal().Err(err).Msg("Failed to migrate localities")
		}

		log.Info().Msg("Migrating developers------------>>>>>>>>>>>>>>>>>>>>")
		if err := migration_jobs.MigrateDeveloper(ctx, txn); err != nil {
			log.Fatal().Err(err).Msg("Failed to migrate developers")
		}

		log.Info().Msg("Migrating properties------------>>>>>>>>>>>>>>>>>>>>")
		if err := migration_jobs.MigrateProperty(ctx, txn); err != nil {
			log.Fatal().Err(err).Msg("Failed to migrate properties")
		}

		log.Info().Msg("Migrating projects------------>>>>>>>>>>>>>>>>>>>>")
		if err := migration_jobs.MigrateProject(ctx, txn); err != nil {
			log.Fatal().Err(err).Msg("Failed to migrate projects")
		}

		log.Info().Msg("Committing transaction------------>>>>>>>>>>>>>>>>>>>>")
		if err := txn.Commit(); err != nil {
			log.Fatal().Err(err).Msg("Failed to commit transaction")
		}

		log.Info().Msg("Migration completed successfully")
		return
	}

	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	client := database.NewClient("postgres://im_db_dev:password@localhost:5434/mydb?sslmode=disable")

	defer client.Close()

	repo := repository.NewRepository(client)
	app := application.NewApplication(repo)

	// Initialize router
	router.Init(app)

	// Start server
	log.Info().Msgf("Server starting on port %d", config.GetConfig().Port)
	if err := http.ListenAndServe(":"+strconv.Itoa(config.GetConfig().Port), router.Router); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
