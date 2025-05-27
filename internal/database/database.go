package database

import (
	"context"
	"log"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/VI-IM/im_backend_go/ent"
)

// NewClient creates a new ent client
func NewClient(dsn string) *ent.Client {
	drv, err := entsql.Open(dialect.MySQL, dsn)
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}

	client := ent.NewClient(ent.Driver(drv))

	// Run the auto migration tool
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	return client
}
