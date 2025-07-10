package repository

import (
	"context"
	"errors"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/developer"
	"github.com/VI-IM/im_backend_go/shared/logger"
)

func (r *repository) GetAllDevelopers() ([]*ent.Developer, error) {
	ctx := context.Background()

	// Get all developers
	developers, err := r.db.Developer.Query().
		Order(ent.Desc(developer.FieldID)).
		All(ctx)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get developers")
		return nil, err
	}

	return developers, nil
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

func (r *repository) SoftDeleteDeveloper(id string) error {
	_, err := r.db.Developer.UpdateOneID(id).
		SetIsActive(false).
		Save(context.Background())
	return err
}
