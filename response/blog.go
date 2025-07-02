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
}

func GetBlogFromEnt(blog *ent.Blogs) *BlogResponse {
	return &BlogResponse{
		ID:          blog.ID,
		BlogURL:     blog.BlogURL,
		BlogContent: blog.BlogContent,
		SEOMetaInfo: blog.SeoMetaInfo,
		IsPriority:  blog.IsPriority,
	}
}

type BlogListResponse struct {
	Blogs      []*BlogResponse `json:"blogs"`
	TotalCount int             `json:"total_count"`
}
