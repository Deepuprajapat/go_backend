package application

import (
	"context"
	"net/http"

	"github.com/VI-IM/im_backend_go/request"
	"github.com/VI-IM/im_backend_go/response"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/VI-IM/im_backend_go/shared/logger"
)

func (c *application) ListBlogs(pagination *request.PaginationRequest) (*response.BlogListResponse, *imhttp.CustomError) {
	// Get blogs from repository
	blogs, err := c.repo.GetAllBlogs()
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get blogs")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get blogs", err.Error())
	}

	// Convert to response type
	blogResponses := make([]*response.BlogListItem, len(blogs))
	for i, blog := range blogs {
		blogResponses[i] = response.GetBlogListItemFromEnt(blog)
	}

	return &response.BlogListResponse{
		Blogs: blogResponses,
	}, nil
}

func (c *application) GetBlogByID(id string) (*response.BlogResponse, *imhttp.CustomError) {
	blog, err := c.repo.GetBlogByID(id)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get blog")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get blog", err.Error())
	}

	if blog == nil {
		return nil, imhttp.NewCustomErr(http.StatusNotFound, "Blog not found", "Blog not found")
	}

	if blog.IsDeleted {
		return nil, imhttp.NewCustomErr(http.StatusNotFound, "Blog has been deleted", "Blog has been deleted")
	}

	return response.GetBlogFromEnt(blog), nil
}

func (c *application) CreateBlog(ctx context.Context, req *request.CreateBlogRequest) (*response.BlogResponse, *imhttp.CustomError) {
	blog, err := c.repo.CreateBlog(ctx, req.BlogURL, req.BlogContent, req.SEOMetaInfo, req.IsPriority)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to create blog")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to create blog", err.Error())
	}

	return response.GetBlogFromEnt(blog), nil
}

func (c *application) DeleteBlog(ctx context.Context, id string) *imhttp.CustomError {
	// Check if blog exists
	blog, err := c.repo.GetBlogByID(id)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get blog")
		return imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get blog", err.Error())
	}
	if blog == nil {
		return imhttp.NewCustomErr(http.StatusNotFound, "Blog not found", "Blog not found")
	}

	// Delete the blog
	err = c.repo.DeleteBlog(ctx, id)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to delete blog")
		return imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to delete blog", err.Error())
	}

	return nil
}

func (c *application) UpdateBlog(ctx context.Context, id string, req *request.UpdateBlogRequest) (*response.BlogResponse, *imhttp.CustomError) {
	// Update blog in repository
	blog, err := c.repo.UpdateBlog(ctx, id, req.BlogURL, req.BlogContent, req.SEOMetaInfo, req.IsPriority)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to update blog")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to update blog", err.Error())
	}

	if blog == nil {
		return nil, imhttp.NewCustomErr(http.StatusNotFound, "Blog not found", "Blog not found")
	}

	return response.GetBlogFromEnt(blog), nil
}
