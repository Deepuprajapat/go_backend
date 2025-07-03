package application

import (
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

	return response.GetBlogFromEnt(blog), nil
}
