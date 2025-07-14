package repository

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/project"
)

func (r *repository) GetProjectByCanonicalURL(ctx context.Context, url string) (*ent.Project, error) {
	return r.db.Project.Query().
		Where(
			project.IsDeletedEQ(false),
			func(s *sql.Selector) {
				s.Where(sql.ExprP("meta_info->>'canonical' = $1", url)).Limit(1)
			},
		).
		Only(ctx)
}

func (r *repository) GetPropertyByCanonicalURL(ctx context.Context, url string) (*ent.Property, error) {
	return r.db.Property.Query().
		Where(func(s *sql.Selector) {
			// Correctly grouped SQL expression
			s.Where(sql.ExprP("(meta_info->>'canonical') = $1", url))
		}).
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
