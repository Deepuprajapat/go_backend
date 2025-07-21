package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/VI-IM/im_backend_go/request"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/gorilla/mux"
)

func (h *Handler) ListBlogs(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	// ✅ Check for is_published query parameter
	var req request.GetAllAPIRequest
	req.Filters = make(map[string]interface{})

	isPublishedParam := r.URL.Query().Get("is_published")
	isRecentBlogsParam := r.URL.Query().Get("recent_blogs")
	page := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("page_size")

	if isPublishedParam != "" {
		req.Filters["is_published"] = isPublishedParam == "true"
	}

	if isRecentBlogsParam != "" {
		req.Filters["recent_blogs"] = isRecentBlogsParam == "true"
	}

	if isPublishedParam == "" && isRecentBlogsParam == "" {
		req.Filters["recent_blogs"] = true
	}

	if page != "" {
		req.Page, _ = strconv.Atoi(page)
	} else {
		req.Page = 1
	}
	if pageSize != "" {
		req.PageSize, _ = strconv.Atoi(pageSize)
	} else {
		req.PageSize = 10000
	}

	response, err := h.app.ListBlogs(&req)
	if err != nil {
		return nil, err
	}

	return &imhttp.Response{
		Data:       response,
		StatusCode: http.StatusOK,
	}, nil
}

func (h *Handler) GetBlog(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	vars := mux.Vars(r)
	blogID := vars["blog_id"]

	blog, err := h.app.GetBlogByID(blogID)
	if err != nil {
		return nil, err
	}

	return &imhttp.Response{
		Data:       blog,
		StatusCode: http.StatusOK,
	}, nil
}

func (h *Handler) CreateBlog(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	var req request.CreateBlogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := h.validate.Struct(req); err != nil {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request", err.Error())
	}

	blog, err := h.app.CreateBlog(r.Context(), &req)
	if err != nil {
		return nil, err
	}

	return &imhttp.Response{
		Data:       blog,
		StatusCode: http.StatusCreated,
	}, nil
}

func (h *Handler) DeleteBlog(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	vars := mux.Vars(r)
	blogID := vars["blog_id"]

	if err := h.app.DeleteBlog(r.Context(), blogID); err != nil {
		return nil, err
	}

	return &imhttp.Response{
		Data:       map[string]string{"message": "Blog deleted successfully"},
		StatusCode: http.StatusOK,
	}, nil
}

func (h *Handler) UpdateBlog(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	vars := mux.Vars(r)
	blogID := vars["blog_id"]

	var req request.UpdateBlogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := h.validate.Struct(req); err != nil {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request", err.Error())
	}

	blog, err := h.app.UpdateBlog(r.Context(), blogID, &req)
	if err != nil {
		return nil, err
	}

	return &imhttp.Response{
		Data:       blog,
		StatusCode: http.StatusOK,
	}, nil
}
