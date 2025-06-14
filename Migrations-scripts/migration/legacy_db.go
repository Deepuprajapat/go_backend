package migration

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

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

	log.Println("Successfully connected to legacy database")
	return db, nil
}
