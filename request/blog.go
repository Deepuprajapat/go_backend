package request

import "github.com/VI-IM/im_backend_go/ent/schema"

type CreateBlogRequest struct {
	BlogURL     string             `json:"blog_url" validate:"required"`
	BlogContent schema.BlogContent `json:"blog_content" validate:"required"`
	SEOMetaInfo schema.SEOMetaInfo `json:"seo_meta_info" validate:"required"`
	IsPriority  bool               `json:"is_priority"`
}

type UpdateBlogRequest struct {
	BlogURL     *string             `json:"blog_url,omitempty"`
	BlogContent *schema.BlogContent `json:"blog_content,omitempty"`
	SEOMetaInfo *schema.SEOMetaInfo `json:"seo_meta_info,omitempty"`
	IsPriority  *bool               `json:"is_priority,omitempty"`
}
