package application

import (
	"context"
	"net/http"

	"github.com/VI-IM/im_backend_go/ent/schema"
	"github.com/VI-IM/im_backend_go/request"
	"github.com/VI-IM/im_backend_go/response"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/VI-IM/im_backend_go/shared/logger"
)

func (a *application) GetGenericSearchData(ctx context.Context) ([]response.GenericSearchData, *imhttp.CustomError) {
	genericSearchData, err := a.repo.GetGenericSearchData(ctx)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get generic search data")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get generic search data", err.Error())
	}

	responseGenericSearchData := make([]response.GenericSearchData, len(genericSearchData))
	for i, data := range genericSearchData {
		responseGenericSearchData[i] = response.GenericSearchData{
			Index:        i,
			CanonicalURL: data.CanonicalURL,
			SearchTerm:   data.SearchTerm,
			Filters:      data.Filters,
		}
	}
	return responseGenericSearchData, nil
}

func (a *application) AddGenericSearchData(ctx context.Context, input request.GenericSearchData) ([]*response.GenericSearchData, *imhttp.CustomError) {
	if input.CanonicalURL == "" {
		logger.Get().Error().Msg("Canonical URL is required")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Canonical URL is required", "Canonical URL is required")
	}

	if input.SearchTerm == "" {
		logger.Get().Error().Msg("Search term is required")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Search term is required", "Search term is required")
	}
	genericSearchData, err := a.repo.AddGenericSearchData(ctx, &schema.GenericSearchData{
		CanonicalURL: input.CanonicalURL,
		SearchTerm:   input.SearchTerm,
		Filters:      input.Filters,
	})
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to add generic search data")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to add generic search data", err.Error())
	}
	responseGenericSearchData := make([]response.GenericSearchData, len(genericSearchData))
	for i, data := range genericSearchData {
		responseGenericSearchData[i] = response.GenericSearchData{
			Index:        i,
			CanonicalURL: data.CanonicalURL,
			SearchTerm:   data.SearchTerm,
			Filters:      data.Filters,
		}
	}

	responseGenericSearchDataPtr := make([]*response.GenericSearchData, len(responseGenericSearchData))
	for i, data := range responseGenericSearchData {
		responseGenericSearchDataPtr[i] = &data
	}
	return responseGenericSearchDataPtr, nil
}

func (a *application) UpdateGenericSearchData(ctx context.Context, input request.GenericSearchData) (*response.GenericSearchData, *imhttp.CustomError) {
	if input.Index == 0 || input.Index < 0 {
		logger.Get().Error().Msg("Index is required")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Index is required", "Index is required")
	}

	genericSearchData, err := a.repo.UpdateGenericSearchData(ctx, &schema.GenericSearchData{
		CanonicalURL: input.CanonicalURL,
		SearchTerm:   input.SearchTerm,
		Filters:      input.Filters,
	}, input.Index)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to update generic search data")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to update generic search data", err.Error())
	}
	responseGenericSearchData := response.GenericSearchData{
		Index:        input.Index,
		CanonicalURL: genericSearchData.CanonicalURL,
		SearchTerm:   genericSearchData.SearchTerm,
		Filters:      genericSearchData.Filters,
	}
	return &responseGenericSearchData, nil
}

func (a *application) DeleteGenericSearchData(ctx context.Context, index int) *imhttp.CustomError {
	err := a.repo.DeleteGenericSearchData(ctx, index)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to delete generic search data")
		return imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to delete generic search data", err.Error())
	}
	return nil
}
