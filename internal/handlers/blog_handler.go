package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/VI-IM/im_backend_go/request"
	"github.com/VI-IM/im_backend_go/response"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/gorilla/mux"
)

func (h *Handler) ListBlogs(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	// ✅ Check for is_published query parameter
	isPublishedParam := r.URL.Query().Get("is_published")

	var isPublished *bool
	if isPublishedParam != "" {
		// Parse the boolean value
		switch isPublishedParam {
		case "true":
			value := true
			isPublished = &value
		case "false":
			value := false
			isPublished = &value
		}
		// If it's neither "true" nor "false", isPublished remains nil (no filter)
	}

	var response *response.BlogListResponse
	var err *imhttp.CustomError

	if isPublished != nil {
		// ✅ Use filtered method
		response, err = h.app.ListBlogsWithFilter(isPublished)
	} else {
		// ✅ Use original method (no filter)
		response, err = h.app.ListBlogs(nil)
	}

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
