package repository

import (
	"context"
	"errors"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/developer"
	"github.com/VI-IM/im_backend_go/shared/logger"
)

func (r *repository) GetAllDevelopers(offset, limit int) ([]*ent.Developer, int, error) {
	ctx := context.Background()

	// Get total count
	total, err := r.db.Developer.Query().Count(ctx)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to count developers")
		return nil, 0, err
	}

	// Get paginated results
	developers, err := r.db.Developer.Query().
		Order(ent.Desc(developer.FieldID)).
		Offset(offset).
		Limit(limit).
		All(ctx)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get developers")
		return nil, 0, err
	}

	return developers, total, nil
}

func (r *repository) GetDeveloperByID(id string) (*ent.Developer, error) {
	developer, err := r.db.Developer.Get(context.Background(), id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("developer not found")
		}
		logger.Get().Error().Err(err).Msg("Failed to get developer")
		return nil, err
	}
	return developer, nil
}
