package repository

import (
	"context"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/blogs"
	"github.com/VI-IM/im_backend_go/ent/project"
	"github.com/VI-IM/im_backend_go/ent/property"
)

func (r *repository) GetProjectByCanonicalURL(ctx context.Context, url string) (*ent.Project, error) {
	return r.db.Project.Query().
		Where(
			project.IsDeletedEQ(false),
			project.SlugEQ(url),
		).
		WithLocation().
		WithDeveloper().
		Only(ctx)
}

func (r *repository) GetPropertyByCanonicalURL(ctx context.Context, url string) (*ent.Property, error) {
	return r.db.Property.Query().
		Where(
			// Correctly grouped SQL expression
			property.SlugEQ(url),
		).
		Only(ctx)
}

func (r *repository) GetBlogByCanonicalURL(ctx context.Context, url string) (*ent.Blogs, error) {
	return r.db.Blogs.Query().
		Where(
			// Search by slug field instead of blog_url
			blogs.Slug(url),
		).
		Only(ctx)
}
