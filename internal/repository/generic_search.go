package repository

import (
	"context"
	"errors"

	"github.com/VI-IM/im_backend_go/ent/schema"
	"github.com/VI-IM/im_backend_go/shared/logger"
)

func (r *repository) GetGenericSearchData(ctx context.Context) ([]schema.GenericSearchData, error) {
	staticSiteData, err := r.db.StaticSiteData.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	if len(staticSiteData) == 0 {
		return nil, nil
	}

	if len(staticSiteData) > 1 {
		logger.Get().Debug().Msg("multiple static site data found")
		return nil, errors.New("multiple static site data found")
	}

	if len(staticSiteData[0].GenericSearchData) == 0 {
		logger.Get().Debug().Msg("No generic search data found")
		return nil, nil
	}

	return staticSiteData[0].GenericSearchData, nil
}

func (r *repository) AddGenericSearchData(ctx context.Context, data *schema.GenericSearchData) ([]schema.GenericSearchData, error) {
	genericSearchData, err := r.db.StaticSiteData.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	if len(genericSearchData) == 0 {
		return nil, nil
	}

	if len(genericSearchData) > 1 {
		return nil, errors.New("multiple static site data found")
	}

	id := genericSearchData[0].ID

	genericSearchData[0].GenericSearchData = append(genericSearchData[0].GenericSearchData, *data)

	updatedStaticSiteData, err := r.db.StaticSiteData.UpdateOneID(id).SetGenericSearchData(genericSearchData[0].GenericSearchData).Save(ctx)
	if err != nil {
		return nil, err
	}

	return updatedStaticSiteData.GenericSearchData, nil
}

func (r *repository) UpdateGenericSearchData(ctx context.Context, data *schema.GenericSearchData, index int) (*schema.GenericSearchData, error) {
	genericSearchData, err := r.db.StaticSiteData.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	if len(genericSearchData) == 0 {
		return nil, errors.New("no static site data found")
	}

	if len(genericSearchData) > 1 {
		return nil, errors.New("multiple static site data found")
	}

	if index < 0 || index >= len(genericSearchData[0].GenericSearchData) {
		return nil, errors.New("index out of range")
	}

	id := genericSearchData[0].ID

	genericSearchData[0].GenericSearchData[index] = *data

	updateStaticSiteData, err := r.db.StaticSiteData.UpdateOneID(id).SetGenericSearchData(genericSearchData[0].GenericSearchData).Save(ctx)
	if err != nil {
		return nil, err
	}

	return &updateStaticSiteData.GenericSearchData[index], nil
}

func (r *repository) DeleteGenericSearchData(ctx context.Context, index int) error {
	genericSearchData, err := r.db.StaticSiteData.Query().All(ctx)
	if err != nil {
		return err
	}

	if len(genericSearchData) == 0 {
		return errors.New("no static site data found")
	}

	if len(genericSearchData) > 1 {
		return errors.New("multiple static site data found")
	}

	if index < 0 || index >= len(genericSearchData[0].GenericSearchData) {
		return errors.New("index out of range")
	}

	genericSearchData[0].GenericSearchData = append(genericSearchData[0].GenericSearchData[:index], genericSearchData[0].GenericSearchData[index+1:]...)

	err = r.db.StaticSiteData.UpdateOneID(genericSearchData[0].ID).SetGenericSearchData(genericSearchData[0].GenericSearchData).Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}
