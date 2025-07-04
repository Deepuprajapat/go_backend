package migration_jobs

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/internal/config"
	"github.com/VI-IM/im_backend_go/internal/database"
	"github.com/VI-IM/im_backend_go/shared/logger"
)

var (
	newDB    *ent.Client
	legacyDB *sql.DB
)

func NewNewDBConnection() (*ent.Client, error) {

	if err := config.LoadConfig(); err != nil {
		logger.Get().Fatal().Err(err).Msg("Failed to load config")
	}

	logger.Get().Info().Msgf("Database URL: %s", config.GetConfig().Database.URL)
	logger.Get().Info().Msgf("Database URL: %s", config.GetConfig().Database.DB_HOST)

	client := database.NewClient(config.GetConfig().Database.URL)

	// Run the auto migration tool
	if err := client.Schema.Create(context.Background()); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed creating schema resources: %v", err)
	}

	return client, nil
}

// LegacyDBConfig holds the configuration for the legacy database connection
type LegacyDBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// NewLegacyDBConnection creates a new connection to the legacy MySQL database
func NewLegacyDBConnection() (*sql.DB, error) {
	config := LegacyDBConfig{
		Host:     "invest.c74oiwy6c0gc.ap-south-1.rds.amazonaws.com",
		Port:     "3306",
		User:     "admin",
		Password: "Abhishek202408",
		DBName:   "investmango",
	}

	// Format the connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)

	// Open the database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %v", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	// Set the global variable
	legacyDB = db

	log.Println("Successfully connected to legacy database")
	return db, nil
}
