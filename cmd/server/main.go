package main

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"
	"fmt"
	"github.com/VI-IM/im_backend_go/internal/application"
	"github.com/VI-IM/im_backend_go/internal/config"
	"github.com/VI-IM/im_backend_go/internal/database"
	"github.com/VI-IM/im_backend_go/internal/repository"
	"github.com/VI-IM/im_backend_go/internal/router"
	"github.com/VI-IM/im_backend_go/migration_jobs"
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
	})

	if len(os.Args) > 1 && os.Args[1] == "run-migration" {

		legacyDB, err := migration_jobs.NewLegacyDBConnection()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to connect to legacy database")
		}

		newDB, err := migration_jobs.NewNewDBConnection()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to connect to new database")
		}

		defer newDB.Close()
		defer legacyDB.Close()

		err = migration_jobs.MigrateLocality(context.Background(), legacyDB, newDB)
		if err != nil {
			fmt.Println("Error in migrating localities", err)
			return
		}
		
		err = migration_jobs.MigrateDeveloper(context.Background(), legacyDB,newDB)
		if err != nil {
			fmt.Println("Error in migrating developers", err)
			return
		}

		// ------ follow this sequence ------
		// migrate city
		// migrate locality
		// migrate developer
		// migrate properties
		// migrate project

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
