package request

import "github.com/VI-IM/im_backend_go/ent/schema"

type CreateBlogRequest struct {
	Slug     string             `json:"slug" validate:"required"`
	BlogContent schema.BlogContent `json:"blog_content" validate:"required"`
	SEOMetaInfo schema.SEOMetaInfo `json:"seo_meta_info" validate:"required"`
	IsPriority  bool               `json:"is_priority"`
	IsPublished bool               `json:"is_published"`
}

type UpdateBlogRequest struct {
	BlogURL     *string             `json:"blog_url,omitempty"`
	BlogContent *schema.BlogContent `json:"blog_content,omitempty"`
	SEOMetaInfo *schema.SEOMetaInfo `json:"seo_meta_info,omitempty"`
	IsPriority  *bool               `json:"is_priority,omitempty"`
}
