package database

import (
	"context"
	"fmt"
	"log"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/VI-IM/im_backend_go/ent"
	_ "github.com/lib/pq"
)

// NewClient creates a new ent client
func NewClient(dsn string) *ent.Client {
	drv, err := entsql.Open(dialect.Postgres, dsn)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}

	client := ent.NewClient(ent.Driver(drv))

	// Run the auto migration tool
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	fmt.Println("Connected to PostgreSQL")

	return client
}
