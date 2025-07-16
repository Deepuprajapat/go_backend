package testutil

import (
	"context"
	"fmt"
	"os"
	"testing"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/VI-IM/im_backend_go/ent"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// SetupTestClient creates an Ent client using testcontainers for complete isolation
func SetupTestClient(t *testing.T) (*ent.Client, func()) {
	ctx := context.Background()

	// Start PostgreSQL container with testcontainers
	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2)),
	)
	if err != nil {
		t.Fatalf("Failed to start PostgreSQL container: %v", err)
	}

	// Get connection string from container
	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		postgresContainer.Terminate(ctx)
		t.Fatalf("Failed to get connection string: %v", err)
	}

	// Set environment variable for the test (some parts of the app might read this)
	os.Setenv("DATABASE_URL", connStr)

	// Create Ent client
	drv, err := entsql.Open(dialect.Postgres, connStr)
	if err != nil {
		postgresContainer.Terminate(ctx)
		t.Fatalf("Failed opening connection to test postgres: %v", err)
	}

	client := ent.NewClient(ent.Driver(drv))

	// Run the auto migration tool
	if err := client.Schema.Create(context.Background()); err != nil {
		client.Close()
		postgresContainer.Terminate(ctx)
		t.Fatalf("Failed creating schema resources: %v", err)
	}

	// Return client and cleanup function
	return client, func() {
		client.Close()
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Errorf("Failed to terminate PostgreSQL container: %v", err)
		}
	}
}

// SetupTestDB creates a test database using testcontainers (legacy function, use SetupTestClient instead)
func SetupTestDB(t *testing.T) func() {
	client, cleanup := SetupTestClient(t)
	client.Close() // We only want the database setup, not the client
	return cleanup
}

// GetTestConnectionString returns a connection string to a test PostgreSQL instance
func GetTestConnectionString(t *testing.T) (string, func()) {
	ctx := context.Background()

	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2)),
	)
	if err != nil {
		t.Fatalf("Failed to start PostgreSQL container: %v", err)
	}

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		postgresContainer.Terminate(ctx)
		t.Fatalf("Failed to get connection string: %v", err)
	}

	return connStr, func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Errorf("Failed to terminate PostgreSQL container: %v", err)
		}
	}
}
