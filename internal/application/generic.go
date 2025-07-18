package application

import (
	"context"
	"net/http"
	"strings"

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

	if customSearchPage.Title == "" ||
		customSearchPage.Description == "" ||
		customSearchPage.Slug == "" ||
		customSearchPage.Filters == nil ||
		customSearchPage.MetaInfo == nil ||
		customSearchPage.MetaInfo.Title == "" ||
		customSearchPage.MetaInfo.Description == "" ||
		customSearchPage.MetaInfo.Keywords == "" ||
		customSearchPage.SearchTerm == "" {

		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Missing or invalid required fields", "One or more required fields are missing or invalid")
	}

	if customSearchPage.Slug != "" {
		customSearchPage.Slug = strings.ReplaceAll(customSearchPage.Slug, " ", "-")
	}
	// check if slug is already in use
	existingCustomSearchPage, err := a.repo.GetCustomSearchPageFromSlug(ctx, customSearchPage.Slug)
	if err != nil {
		// If it's a "not found" error, that's what we want - slug is available
		if !ent.IsNotFound(err) {
			return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to check if slug is already in use", err.Error())
		}
		// Slug doesn't exist, which is good for creating a new one
	} else if existingCustomSearchPage != nil {
		logger.Get().Info().Msg("Slug already in use")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Slug already in use", "Slug already in use")
	}

	customSearchPageEntity := &ent.CustomSearchPage{
		Title:       customSearchPage.Title,
		Description: customSearchPage.Description,
		Slug:        customSearchPage.Slug,
		Filters:     customSearchPage.Filters,
		SearchTerm:  customSearchPage.SearchTerm,
		MetaInfo: schema.MetaInfo{
			Title:       customSearchPage.MetaInfo.Title,
			Description: customSearchPage.MetaInfo.Description,
			Keywords:    customSearchPage.MetaInfo.Keywords,
		},
	}

	customSearchPageEntity, err = a.repo.AddCustomSearchPage(ctx, customSearchPageEntity)
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

func (a *application) UpdateCustomSearchPage(ctx context.Context, id string, customSearchPage *request.CustomSearchPage) (*response.CustomSearchPage, *imhttp.CustomError) {

	customSearchPageEntity := &ent.CustomSearchPage{
		ID:          id,
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
