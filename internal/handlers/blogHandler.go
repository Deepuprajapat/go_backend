package handlers

import (
	"net/http"

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
