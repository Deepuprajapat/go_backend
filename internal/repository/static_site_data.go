package repository

import (
	"context"
	"errors"
	"time"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/staticsitedata"
	"github.com/VI-IM/im_backend_go/request"
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

func (r *repository) AddCategory(categoryName string) error {
	ssd, err := r.GetStaticSiteData()
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get static site data")
		return err
	}

	for category, _ := range ssd.CategoriesWithAmenities.Categories {
		if category == categoryName {
			return errors.New("category already exists")
		}
	}

	ssd.CategoriesWithAmenities.Categories[categoryName] = []struct {
		Icon  string `json:"icon"`
		Value string `json:"value"`
	}{}

	if err := r.UpdateStaticSiteData(ssd); err != nil {
		logger.Get().Error().Err(err).Msg("Failed to update static site data")
		return err
	}

	return nil
}

func (r *repository) AddAmenityToCategory(req *request.AddAmenityToCategoryRequest) error {

	ssd, err := r.GetStaticSiteData()
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get static site data")
		return err
	}

	ssd.CategoriesWithAmenities.Categories[req.CategoryName] = append(ssd.CategoriesWithAmenities.Categories[req.CategoryName], req.Amenities...)

	if err := r.UpdateStaticSiteData(ssd); err != nil {
		logger.Get().Error().Err(err).Msg("Failed to update static site data")
		return err
	}

	return nil
}

func (r *repository) DeleteAmenityFromCategory(req *request.DeleteAmenityFromCategoryRequest) error {

	ssd, err := r.GetStaticSiteData()
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get static site data")
		return err
	}
	var amenities []struct {
		Icon  string `json:"icon"`
		Value string `json:"value"`
	}
	for _, amenity := range ssd.CategoriesWithAmenities.Categories[req.CategoryName] {
		if amenity.Value != req.AmenityName {
			amenities = append(amenities, amenity)
		}
	}

	ssd.CategoriesWithAmenities.Categories[req.CategoryName] = amenities

	if err := r.UpdateStaticSiteData(ssd); err != nil {
		logger.Get().Error().Err(err).Msg("Failed to update static site data")
		return err
	}

	return nil
}

func (r *repository) DeleteCategoryWithAmenities(categoryName string) error {

	ssd, err := r.GetStaticSiteData()
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get static site data")
		return err
	}

	delete(ssd.CategoriesWithAmenities.Categories, categoryName)

	if err := r.UpdateStaticSiteData(ssd); err != nil {
		logger.Get().Error().Err(err).Msg("Failed to update static site data")
		return err
	}

	return nil
}
