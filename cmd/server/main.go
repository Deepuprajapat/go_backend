package main

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/VI-IM/im_backend_go/Migrations-scripts/migration"
	"github.com/VI-IM/im_backend_go/internal/application"
	"github.com/VI-IM/im_backend_go/internal/config"
	"github.com/VI-IM/im_backend_go/internal/database"
	"github.com/VI-IM/im_backend_go/internal/repository"
	"github.com/VI-IM/im_backend_go/internal/router"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
)

func main() {
	// Initialize zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    true,
	})

	if os.Args[1] == "run-migration" {
		legacyDB, err := migration.NewLegacyDBConnection()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to connect to legacy database")
		}

		newDB, err := migration.NewNewDBConnection()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to connect to new database")
		}
		defer newDB.Close()
		defer legacyDB.Close()

		projects, err := migration.FetchLegacyProjectData(context.Background(), legacyDB)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to fetch legacy project data")
		}

		err = migration.MigrateCommonFields(context.Background(), newDB, projects)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to migrate common fields")
		}

		log.Info().Msg("Migration completed successfully")
		return

	}

	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	client := database.NewClient("im_db_dev:password@tcp(im_mysql_db:3306)/mydb?parseTime=true")

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
