package testutil

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	testDBName = "im_test_db"
	maxRetries = 5
)

// SetupTestDB creates a test database and returns a cleanup function
func SetupTestDB(t *testing.T) func() {
	// Connect to MySQL without specifying a database
	dsn := fmt.Sprintf("root:password@tcp(localhost:3306)/")

	// Try to connect with retries
	var db *sql.DB
	var err error
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			// Test the connection
			err = db.Ping()
			if err == nil {
				break
			}
		}
		if i < maxRetries-1 {
			time.Sleep(time.Second * 2)
		}
	}
	if err != nil {
		t.Fatalf("Failed to connect to MySQL after %d retries: %v", maxRetries, err)
	}

	// Create test database
	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", testDBName))
	if err != nil {
		t.Fatalf("Failed to drop existing test database: %v", err)
	}

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", testDBName))
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Set environment variable for the test
	os.Setenv("DB_DSN", fmt.Sprintf("root:password@tcp(localhost:3306)/%s", testDBName))

	// Return cleanup function
	return func() {
		_, err := db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", testDBName))
		if err != nil {
			t.Errorf("Failed to drop test database: %v", err)
		}
		db.Close()
	}
}
