package repository

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/VI-IM/im_backend_go/ent"
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
			func(s *sql.Selector) {
				s.Where(sql.ExprP("blog_url = $1", url))
			},
		).
		Only(ctx)
}
