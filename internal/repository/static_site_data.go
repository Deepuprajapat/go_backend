package repository

import (
	"context"
	"time"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/staticsitedata"
	"github.com/VI-IM/im_backend_go/shared/logger"
)

func (r *repository) GetStaticSiteData() (*ent.StaticSiteData, error) {
	// Get the first (and should be only) static site data record
	staticData, err := r.db.StaticSiteData.Query().First(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			logger.Get().Error().Err(err).Msg("Static site data not found")
			return nil, err
		}
		logger.Get().Error().Err(err).Msg("Failed to get static site data")
		return nil, err
	}
	return staticData, nil
}

func (r *repository) UpdateStaticSiteData(data *ent.StaticSiteData) error {
	_, err := r.db.StaticSiteData.UpdateOneID(data.ID).
		SetCategoriesWithAmenities(data.CategoriesWithAmenities).
		SetTestimonials(data.Testimonials).
		SetUpdatedAt(time.Now()).
		Save(context.Background())
	return err
}

func (r *repository) CheckCategoryExists(category string) (bool, error) {
	staticSiteData, err := r.db.StaticSiteData.Query().Where(staticsitedata.IsActive(true)).All(context.Background())
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to check if category exists")
		return false, err
	}

	var isExist bool

	isExist = false
	for existedCategory := range staticSiteData[0].CategoriesWithAmenities.Categories {
		if existedCategory == category {
			isExist = true
		}
	}

	return isExist, nil
}

func (r *repository) AddCategoryWithAmenities(data *ent.StaticSiteData) error {

	return r.UpdateStaticSiteData(data)
}
