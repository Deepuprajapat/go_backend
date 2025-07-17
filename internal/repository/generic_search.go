package repository

import (
	"context"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/customsearchpage"
	"github.com/VI-IM/im_backend_go/ent/schema"
	"github.com/VI-IM/im_backend_go/shared/logger"
	"github.com/google/uuid"
)

func (r *repository) GetCustomSearchPageFromSlug(ctx context.Context, slug string) (*ent.CustomSearchPage, error) {
	customSearchPage, err := r.db.CustomSearchPage.Query().Where(customsearchpage.Slug(slug)).First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			logger.Get().Debug().Str("slug", slug).Msg("Custom search page not found")
			return nil, err
		}
		logger.Get().Error().Err(err).Msg("Failed to get custom search page by slug")
		return nil, err
	}
	return customSearchPage, nil
}

func (r *repository) GetAllCustomSearchPages(ctx context.Context) ([]*ent.CustomSearchPage, error) {
	customSearchPages, err := r.db.CustomSearchPage.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	return customSearchPages, nil
}

func (r *repository) AddCustomSearchPage(ctx context.Context, customSearchPage *ent.CustomSearchPage) (*ent.CustomSearchPage, error) {
	logger.Get().Info().Msg("Adding custom search page from repository")
	customSearchPage, err := r.db.CustomSearchPage.Create().
		SetID(uuid.New().String()).
		SetTitle(customSearchPage.Title).
		SetDescription(customSearchPage.Description).
		SetFilters(customSearchPage.Filters).
		SetMetaInfo(customSearchPage.MetaInfo).
		SetSearchTerm(customSearchPage.SearchTerm).
		SetSlug(customSearchPage.Slug).
		Save(ctx)
	if err != nil {
		logger.Get().Error().Err(err).Msg(err.Error())
		return nil, err
	}
	return customSearchPage, nil
}

func (r *repository) UpdateCustomSearchPage(ctx context.Context, customSearchPage *ent.CustomSearchPage) (*ent.CustomSearchPage, error) {
	update := r.db.CustomSearchPage.UpdateOneID(customSearchPage.ID)

	if customSearchPage.Title != "" {
		update.SetTitle(customSearchPage.Title)
	}
	if customSearchPage.Description != "" {
		update.SetDescription(customSearchPage.Description)
	}

	if customSearchPage.MetaInfo != (schema.MetaInfo{}) {
		update.SetMetaInfo(customSearchPage.MetaInfo)
	}

	if customSearchPage.SearchTerm != "" {
		update.SetSearchTerm(customSearchPage.SearchTerm)
	}

	if customSearchPage.Filters != nil {
		update.SetFilters(customSearchPage.Filters)
	}

	return update.Save(ctx)
}

func (r *repository) DeleteCustomSearchPage(ctx context.Context, id string) error {
	err := r.db.CustomSearchPage.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
