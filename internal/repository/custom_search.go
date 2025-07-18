package repository

import (
	"context"

	"github.com/VI-IM/im_backend_go/ent/blogs"
	"github.com/VI-IM/im_backend_go/ent/customsearchpage"
	"github.com/VI-IM/im_backend_go/ent/project"
	"github.com/VI-IM/im_backend_go/response"
)

func (r *repository) CheckURLExists(ctx context.Context, url string) (*response.URLExistsResult, error) {
	// Check blogs first
	blog, err := r.db.Blogs.Query().
		Where(blogs.Slug(url)).
		First(ctx)
	if err == nil && blog != nil {
		return &response.URLExistsResult{
			Exists:     true,
			EntityType: "blog",
			EntityID:   blog.ID,
		}, nil
	}

	// Check projects
	project, err := r.db.Project.Query().
		Where(project.Slug(url)).
		First(ctx)
	if err == nil && project != nil {
		return &response.URLExistsResult{
			Exists:     true,
			EntityType: "project",
			EntityID:   project.ID,
		}, nil
	}

	// Check custom search pages
	customSearchPage, err := r.db.CustomSearchPage.Query().
		Where(customsearchpage.Slug(url)).
		First(ctx)
	if err == nil && customSearchPage != nil {
		return &response.URLExistsResult{
			Exists:     true,
			EntityType: "custom_search_page",
			EntityID:   customSearchPage.ID,
		}, nil
	}

	// If we get here, nothing was found
	return &response.URLExistsResult{
		Exists:     false,
		EntityType: "",
		EntityID:   "",
	}, nil
}
