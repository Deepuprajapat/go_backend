package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/internal/application"
	s3client "github.com/VI-IM/im_backend_go/internal/client"
	"github.com/VI-IM/im_backend_go/internal/config"
	"github.com/VI-IM/im_backend_go/internal/database"
	"github.com/VI-IM/im_backend_go/internal/repository"
	"github.com/VI-IM/im_backend_go/internal/router"
	"github.com/VI-IM/im_backend_go/internal/utils"
	"github.com/VI-IM/im_backend_go/migration_jobs"
	"github.com/VI-IM/im_backend_go/shared/logger"
	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
	"github.com/google/uuid"
)

var (
	legacyDB *sql.DB
	newDB    *ent.Client
	txn      *ent.Tx
)

func main() {

	logger.Init()
	ctx := context.Background()

	logger.Get().Info().Msg("Starting application...")

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "run-migration":
			runMigration(ctx)
			return
		case "seed-admin":
			seedAdmin(ctx)
			return
		}
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

func runMigration(ctx context.Context) {
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

	//logger.Get().Info().Msg("Migrating static site data------------>>>>>>>>>>>>>>>>>>>>")
	//if err := migration_jobs.MigrateStaticSiteData(ctx, txn); err != nil {
	//	logger.Get().Fatal().Err(err).Msg("Failed to migrate static site data")
	//}
	//
	//logger.Get().Info().Msg("Migrating localities------------>>>>>>>>>>>>>>>>>>>>")
	//if err := migration_jobs.MigrateLocality(ctx, txn); err != nil {
	//	logger.Get().Fatal().Err(err).Msg("Failed to migrate localities")
	//}
	//
	//logger.Get().Info().Msg("Migrating developers------------>>>>>>>>>>>>>>>>>>>>")
	//if err := migration_jobs.MigrateDeveloper(ctx, txn); err != nil {
	//	logger.Get().Fatal().Err(err).Msg("Failed to migrate developers")
	//}

	//logger.Get().Info().Msg("Migrating projects------------>>>>>>>>>>>>>>>>>>>>")
	//if err := migration_jobs.MigrateProject(ctx, txn); err != nil {
	//	logger.Get().Fatal().Err(err).Msg("Failed to migrate projects")
	//}

	logger.Get().Info().Msg("Migrating properties------------>>>>>>>>>>>>>>>>>>>>")
	if err := migration_jobs.MigrateProperty(ctx, txn); err != nil {
		logger.Get().Fatal().Err(err).Msg("Failed to migrate properties")
	}

	logger.Get().Info().Msg("Migrating blogs------------>>>>>>>>>>>>>>>>>>>>")
	if err = migration_jobs.MigrateBlogs(ctx, txn); err != nil {
		logger.Get().Fatal().Err(err).Msg("Failed to migrate blogs")
	}

	logger.Get().Info().Msg("Committing transaction------------>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	if err := txn.Commit(); err != nil {
		logger.Get().Fatal().Err(err).Msg("Failed to commit transaction")
	}

	logger.Get().Info().Msg("Migration completed successfully")
}

func seedAdmin(ctx context.Context) {
	logger.Get().Info().Msg("Starting admin user seeding...")

	// Load configuration
	if err := config.LoadConfig(); err != nil {
		logger.Get().Fatal().Err(err).Msg("Failed to load configuration")
	}

	cfg := config.GetConfig()
	client := database.NewClient(cfg.Database.URL)
	defer client.Close()

	// Default admin credentials
	adminUsername := "admin"
	adminPassword := "admin123"
	adminEmail := "admin@example.com"
	adminName := "System Administrator"

	// Hash the password
	hashedPassword, err := utils.HashPassword(adminPassword)
	if err != nil {
		logger.Get().Fatal().Err(err).Msg("Failed to hash password")
	}

	// Create admin user
	adminUser, err := client.User.Create().
		SetID(uuid.New().String()).
		SetUsername(adminUsername).
		SetPassword(hashedPassword).
		SetEmail(adminEmail).
		SetName(adminName).
		SetIsActive(true).
		SetIsEmailVerified(true).
		SetIsVerified(true).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		Save(ctx)

	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to create admin user")
	} else {
		logger.Get().Info().Msgf("Admin user created successfully!")
		logger.Get().Info().Msgf("Username: %s", adminUsername)
		logger.Get().Info().Msgf("Password: %s", adminPassword)
		logger.Get().Info().Msgf("Email: %s", adminEmail)
		logger.Get().Info().Msgf("User ID: %s", adminUser.ID)
		fmt.Printf("\n=== ADMIN CREDENTIALS ===\n")
		fmt.Printf("Username: %s\n", adminUsername)
		fmt.Printf("Password: %s\n", adminPassword)
		fmt.Printf("Email: %s\n", adminEmail)
		fmt.Printf("========================\n\n")
	}
}
