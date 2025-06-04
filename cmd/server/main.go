package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/VI-IM/im_backend_go/internal/application"
	"github.com/VI-IM/im_backend_go/internal/config"
	"github.com/VI-IM/im_backend_go/internal/database"
	"github.com/VI-IM/im_backend_go/internal/repository"
	"github.com/VI-IM/im_backend_go/internal/router"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339, NoColor: true})

	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Initialize database
	client := database.NewClient("im_user:password@tcp(127.0.0.1:3306)/mydb?parseTime=true")
	defer client.Close()

	repo := repository.NewRepository(client)
	application := application.NewApplication(repo)

	// Initialize router
	router.Init(application)

	// Start server
	log.Info().Msgf("Server starting on port %d", config.GetConfig().Port)
	if err := http.ListenAndServe(":"+strconv.Itoa(config.GetConfig().Port), router.Router); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
