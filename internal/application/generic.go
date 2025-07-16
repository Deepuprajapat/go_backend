package application

import (
	"context"
	"net/http"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/schema"
	"github.com/VI-IM/im_backend_go/request"
	"github.com/VI-IM/im_backend_go/response"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/VI-IM/im_backend_go/shared/logger"
)

func (a *application) GetCustomSearchPage(ctx context.Context, slug string) (*response.CustomSearchPage, *imhttp.CustomError) {
	var customSearchPage *response.CustomSearchPage
	csp, err := a.repo.GetCustomSearchPageFromSlug(ctx, slug)
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusNotFound, "Custom search page not found", "Custom search page not found")
	}

	request := &request.GetAllAPIRequest{
		Filters: csp.Filters,
	}

	fprojects, customErr := a.ListProjects(request)
	if customErr != nil {
		return nil, customErr
	}

	customSearchPage = &response.CustomSearchPage{
		ID:          csp.ID,
		Title:       csp.Title,
		Description: csp.Description,
		Projects:    fprojects,
		Slug:        csp.Slug,
		MetaInfo: &response.MetaInfo{
			Title:       csp.MetaInfo.Title,
			Description: csp.MetaInfo.Description,
			Keywords:    csp.MetaInfo.Keywords,
		},
	}

	return customSearchPage, nil
}

func (a *application) GetLinks(ctx context.Context) ([]*response.Link, *imhttp.CustomError) {

	allCustomSearchPages, err := a.repo.GetAllCustomSearchPages(ctx)
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to list links", err.Error())
	}

	var links []*response.Link
	for _, customSearchPage := range allCustomSearchPages {
		links = append(links, &response.Link{
			Title: customSearchPage.Title,
			Slug:  customSearchPage.Slug,
		})
	}

	return links, nil
}

func (a *application) GetAllCustomSearchPages(ctx context.Context) ([]*response.CustomSearchPage, *imhttp.CustomError) {

	allCustomSearchPages, err := a.repo.GetAllCustomSearchPages(ctx)
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to list custom search pages", err.Error())
	}

	customSearchPages := make([]*response.CustomSearchPage, len(allCustomSearchPages))
	for i, customSearchPage := range allCustomSearchPages {
		customSearchPages[i] = &response.CustomSearchPage{
			ID:          customSearchPage.ID,
			Title:       customSearchPage.Title,
			Description: customSearchPage.Description,
			Filters:     customSearchPage.Filters,
			Slug:        customSearchPage.Slug,
			MetaInfo: &response.MetaInfo{
				Title:       customSearchPage.MetaInfo.Title,
				Description: customSearchPage.MetaInfo.Description,
				Keywords:    customSearchPage.MetaInfo.Keywords,
			},
		}
	}

	return customSearchPages, nil
}

func (a *application) AddCustomSearchPage(ctx context.Context, customSearchPage *request.CustomSearchPage) (*response.CustomSearchPage, *imhttp.CustomError) {

	customSearchPageEntity := &ent.CustomSearchPage{
		Title:       customSearchPage.Title,
		Description: customSearchPage.Description,
		Filters:     customSearchPage.Filters,
		SearchTerm:  customSearchPage.SearchTerm,
		MetaInfo: schema.MetaInfo{
			Title:       customSearchPage.MetaInfo.Title,
			Description: customSearchPage.MetaInfo.Description,
			Keywords:    customSearchPage.MetaInfo.Keywords,
		},
	}
	logger.Get().Info().Msg("Adding custom search page from application")
	logger.Get().Info().Msg(customSearchPageEntity.SearchTerm)
	logger.Get().Info().Msg(customSearchPageEntity.Title)
	logger.Get().Info().Msg(customSearchPageEntity.Description)
	logger.Get().Info().Interface("filters", customSearchPageEntity.Filters).Msg("CustomSearchPageEntity Filters")
	logger.Get().Info().Msg(customSearchPageEntity.MetaInfo.Title)
	logger.Get().Info().Msg(customSearchPageEntity.MetaInfo.Description)
	logger.Get().Info().Msg(customSearchPageEntity.MetaInfo.Keywords)

	customSearchPageEntity, err := a.repo.AddCustomSearchPage(ctx, customSearchPageEntity)
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to add custom search page", err.Error())
	}

	response := &response.CustomSearchPage{
		ID:          customSearchPageEntity.ID,
		Title:       customSearchPageEntity.Title,
		Description: customSearchPageEntity.Description,
		Filters:     customSearchPageEntity.Filters,
		SearchTerm:  customSearchPageEntity.SearchTerm,
		MetaInfo: &response.MetaInfo{
			Title:       customSearchPageEntity.MetaInfo.Title,
			Description: customSearchPageEntity.MetaInfo.Description,
			Keywords:    customSearchPageEntity.MetaInfo.Keywords,
		},
	}

	return response, nil
}

func (a *application) UpdateCustomSearchPage(ctx context.Context, customSearchPage *request.CustomSearchPage) (*response.CustomSearchPage, *imhttp.CustomError) {

	customSearchPageEntity := &ent.CustomSearchPage{
		ID:          customSearchPage.ID,
		Title:       customSearchPage.Title,
		Description: customSearchPage.Description,
		Filters:     customSearchPage.Filters,
		MetaInfo: schema.MetaInfo{
			Title:       customSearchPage.MetaInfo.Title,
			Description: customSearchPage.MetaInfo.Description,
			Keywords:    customSearchPage.MetaInfo.Keywords,
		},
	}

	customSearchPageEntity, err := a.repo.UpdateCustomSearchPage(ctx, customSearchPageEntity)
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to update custom search page", err.Error())
	}

	response := &response.CustomSearchPage{

		Title:       customSearchPageEntity.Title,
		Description: customSearchPageEntity.Description,
		Filters:     customSearchPageEntity.Filters,
		MetaInfo: &response.MetaInfo{
			Title:       customSearchPageEntity.MetaInfo.Title,
			Description: customSearchPageEntity.MetaInfo.Description,
			Keywords:    customSearchPageEntity.MetaInfo.Keywords,
		},
	}

	return response, nil
}

func (a *application) DeleteCustomSearchPage(ctx context.Context, id string) *imhttp.CustomError {

	err := a.repo.DeleteCustomSearchPage(ctx, id)
	if err != nil {
		return imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to delete custom search page", err.Error())
	}

	return nil
}
