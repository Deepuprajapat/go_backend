package migration

import (
	"context"
	"fmt"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/internal/database"
	_ "github.com/lib/pq"
)

// NewNewDBConnection creates a new connection using ent client
func NewNewDBConnection() (*ent.Client, error) {
	client := database.NewClient("postgres://im_db_dev:password@localhost:5434/mydb?sslmode=disable")

	// Run the auto migration tool
	if err := client.Schema.Create(context.Background()); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed creating schema resources: %v", err)
	}

	return client, nil
}
