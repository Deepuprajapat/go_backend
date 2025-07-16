package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/VI-IM/im_backend_go/request"
	"github.com/VI-IM/im_backend_go/response"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/gorilla/mux"
)

func (h *Handler) GetCustomSearchPage(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	slug := mux.Vars(r)["slug"]

	if slug == "" {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Slug is required", "Slug is required")
	}

	customSearchPage, err := h.app.GetCustomSearchPage(r.Context(), slug)
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusNotFound, "Custom search page not found", "Custom search page not found")
	}

	userAgent := r.UserAgent()
	isBot := strings.Contains(strings.ToLower(userAgent), "bot") ||
		strings.Contains(strings.ToLower(userAgent), "crawler") ||
		strings.Contains(strings.ToLower(userAgent), "spider")

	if isBot {
		htmlContent := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
	<title>%s</title>
	<meta name="description" content="%s">
	<link rel="canonical" href="https://investmango.com/s/%s">
</head>
<body>
	<h1>%s</h1>
	<p>%s</p>
	<div class="project-list">
		%s
	</div>
</body>
</html>`,
			customSearchPage.Title,
			customSearchPage.Description,
			slug,
			customSearchPage.Title,
			customSearchPage.Description,
			func() string {
				var projectCards string
				for _, project := range customSearchPage.Projects {
					projectCards += fmt.Sprintf(`
						<div class="project-card">
							<h2>%s</h2>
							<p class="address">%s</p>
							<p class="price">Starting from %s</p>
						</div>`,
						project.ProjectName,
						project.FullDetails.LocationInfo.ShortAddress,
						project.FullDetails.MinPrice,
					)
				}
				return projectCards
			}(),
		)

		return &imhttp.Response{
			Data:       htmlContent,
			StatusCode: http.StatusOK,
		}, nil
	}

	return &imhttp.Response{
		Data:       customSearchPage,
		StatusCode: http.StatusOK,
	}, nil
}

func (h *Handler) GetLinks(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {

	ctx := r.Context()
	links, err := h.app.GetLinks(ctx)
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusNotFound, "Links not found", "Links not found")
	}

	return &imhttp.Response{
		Data:       links,
		StatusCode: http.StatusOK,
	}, nil
}

func (h *Handler) GetAllCustomSearchPages(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {

	ctx := r.Context()
	customSearchPages, err := h.app.GetAllCustomSearchPages(ctx)
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusNotFound, "Custom search pages not found", "Custom search pages not found")
	}

	customSearchPagesResponse := make([]*response.CustomSearchPage, len(customSearchPages))
	for i, customSearchPage := range customSearchPages {
		customSearchPagesResponse[i] = &response.CustomSearchPage{
			ID:          customSearchPage.ID,
			Title:       customSearchPage.Title,
			Description: customSearchPage.Description,
			Filters:     customSearchPage.Filters,
			MetaInfo: &response.MetaInfo{
				Title:       customSearchPage.MetaInfo.Title,
				Description: customSearchPage.MetaInfo.Description,
				Keywords:    customSearchPage.MetaInfo.Keywords,
			},
		}
	}

	return &imhttp.Response{
		Data:       customSearchPagesResponse,
		StatusCode: http.StatusOK,
	}, nil
}

func (h *Handler) AddCustomSearchPage(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {

	ctx := r.Context()

	var customSearchPage *request.CustomSearchPage
	err := json.NewDecoder(r.Body).Decode(&customSearchPage)
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request body", "Invalid request body")
	}

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

	customSearchPageResponse, err := h.app.AddCustomSearchPage(ctx, customSearchPage)

	return &imhttp.Response{
		Data:       customSearchPageResponse,
		StatusCode: http.StatusOK,
	}, nil
}

func (h *Handler) UpdateCustomSearchPage(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {

	ctx := r.Context()

	var customSearchPage *request.CustomSearchPage
	err := json.NewDecoder(r.Body).Decode(&customSearchPage)
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request body", "Invalid request body")
	}
	customSearchPageResponse, err := h.app.UpdateCustomSearchPage(ctx, customSearchPage)
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusNotFound, "Custom search page not found", "Custom search page not found")
	}

	return &imhttp.Response{
		Data:       customSearchPageResponse,
		StatusCode: http.StatusOK,
	}, nil
}

func (h *Handler) DeleteCustomSearchPage(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {

	ctx := r.Context()
	id := r.URL.Query().Get("id")
	if id == "" {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "ID is required", "ID is required")
	}

	err := h.app.DeleteCustomSearchPage(ctx, id)
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusNotFound, "Custom search page not found", "Custom search page not found")
	}

	return &imhttp.Response{
		StatusCode: http.StatusOK,
	}, nil
}
