package repository

import (
	"context"
	"errors"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/location"
	"github.com/VI-IM/im_backend_go/shared/logger"
)

func (r *repository) ListLocations(filters map[string]interface{}) ([]*ent.Location, error) {
	ctx := context.Background()
	query := r.db.Location.Query()

	if city, ok := filters["city"].(string); ok && city != "" {
		query = query.Where(location.CityEQ(city))
	}

	locations, err := query.All(ctx)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to list locations")
		return nil, err
	}
	return locations, nil
}

func (r *repository) GetLocationByID(id string) (*ent.Location, error) {
	location, err := r.db.Location.Get(context.Background(), id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("location not found")
		}
		logger.Get().Error().Err(err).Msg("Failed to get location")
		return nil, err
	}
	return location, nil
}

func (r *repository) SoftDeleteLocation(id string) error {
	_, err := r.db.Location.UpdateOneID(id).
		SetIsActive(false).
		Save(context.Background())
	return err
}
