package response

import (
	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/schema"
)

type BlogResponse struct {
	ID          string             `json:"id"`
	BlogURL     string             `json:"blog_url"`
	BlogContent schema.BlogContent `json:"blog_content"`
	SEOMetaInfo schema.SEOMetaInfo `json:"seo_meta_info"`
	IsPriority  bool               `json:"is_priority"`
	CreatedAt   int64              `json:"created_at"`
	UpdatedAt   int64              `json:"updated_at"`
}

type BlogListItem struct {
	ID          string             `json:"id"`
	Image       string             `json:"image"`
	Title       string             `json:"title"`
	BlogURL     string             `json:"blog_url"`
	Description string             `json:"description"`
	IsPriority  bool               `json:"is_priority"`
	CreatedAt   int64              `json:"created_at"`
	UpdatedAt   int64              `json:"updated_at"`
}

func GetBlogFromEnt(blog *ent.Blogs) *BlogResponse {
	return &BlogResponse{
		ID:          blog.ID,
		BlogURL:     blog.BlogURL,
		BlogContent: blog.BlogContent,
		SEOMetaInfo: blog.SeoMetaInfo,
		IsPriority:  blog.IsPriority,
		CreatedAt:   blog.CreatedAt,
		UpdatedAt:   blog.UpdatedAt,
	}
}

func GetBlogListItemFromEnt(blog *ent.Blogs) *BlogListItem {
	return &BlogListItem{
		ID:          blog.ID,
		Image:       blog.BlogContent.Image,
		BlogURL:     blog.BlogURL,
		Title:       blog.BlogContent.Title,
		Description: blog.BlogContent.Description,
		IsPriority:  blog.IsPriority,
		CreatedAt:   blog.CreatedAt,
		UpdatedAt:   blog.UpdatedAt,
	}
}

type BlogListResponse struct {
	Blogs []*BlogListItem `json:"blogs"`
}
