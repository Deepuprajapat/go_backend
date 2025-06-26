package repository

import (
	"context"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/shared/logger"
)

func (r *repository) ListLocations() ([]*ent.Location, error) {
	locations, err := r.db.Location.Query().
		Where().
		All(context.Background())
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to list locations")
		return nil, err
	}
	return locations, nil
}
