package response

import (
	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/schema"
)

type BlogResponse struct {
	ID          string             `json:"id"`
	Slug        string             `json:"slug"`
	BlogContent schema.BlogContent `json:"blog_content"`
	SEOMetaInfo schema.SEOMetaInfo `json:"seo_meta_info"`
	IsPriority  bool               `json:"is_priority"`
	CreatedAt   int64              `json:"created_at"`
	UpdatedAt   int64              `json:"updated_at"`
	IsPublished bool               `json:"is_published"`
}

type BlogListItem struct {
	ID          string  `json:"id"`
	Image       string  `json:"image"`
	Title       string  `json:"title"`
	Slug        *string `json:"slug"`
	Description string  `json:"description"`
	IsPriority  bool    `json:"is_priority"`
	CreatedAt   int64   `json:"created_at"`
	UpdatedAt   int64   `json:"updated_at"`
	IsPublished bool    `json:"is_published"`
}

func GetBlogFromEnt(blog *ent.Blogs) *BlogResponse {
	return &BlogResponse{
		ID:          blog.ID,
		Slug:        blog.Slug,
		BlogContent: blog.BlogContent,
		SEOMetaInfo: blog.SeoMetaInfo,
		IsPriority:  blog.IsPriority,
		CreatedAt:   blog.CreatedAt.Unix(),
		UpdatedAt:   blog.UpdatedAt.Unix(),
		IsPublished: blog.IsPublished,
	}
}

func GetBlogListItemFromEnt(blog *ent.Blogs) *BlogListItem {
	return &BlogListItem{
		ID:          blog.ID,
		Image:       blog.BlogContent.Image,
		Slug:        &blog.Slug,
		Title:       blog.BlogContent.Title,
		Description: blog.BlogContent.Description,
		IsPriority:  blog.IsPriority,
		CreatedAt:   blog.CreatedAt.Unix(),
		UpdatedAt:   blog.UpdatedAt.Unix(),
		IsPublished: blog.IsPublished,
	}
}

type BlogListResponse struct {
	Blogs []*BlogListItem `json:"blogs"`
}
