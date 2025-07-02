package database

import (
	"context"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/shared/logger"
	_ "github.com/lib/pq"
)

// NewClient creates a new ent client
func NewClient(dsn string) *ent.Client {

	drv, err := entsql.Open(dialect.Postgres, dsn)
	if err != nil {
		logger.Get().Fatal().Err(err).Msg("failed opening connection to postgres")

	}

	client := ent.NewClient(ent.Driver(drv))

	// Run the auto migration tool
	if err := client.Schema.Create(context.Background()); err != nil {
		logger.Get().Fatal().Err(err).Msg("failed creating schema resources")
	}

	logger.Get().Info().Msg("Connected to PostgreSQL")

	return client
}
