package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/VI-IM/im_backend_go/request"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/gorilla/mux"
)

func (h *Handler) ListBlogs(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	response, err := h.app.ListBlogs(nil)
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
