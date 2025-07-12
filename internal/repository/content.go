package repository

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/VI-IM/im_backend_go/ent"
)

func (r *repository) GetProjectByCanonicalURL(ctx context.Context, url string) (*ent.Project, error) {
	return r.db.Project.Query().
		Where(func(s *sql.Selector) {
			// Correctly grouped SQL expression
			s.Where(sql.ExprP("(meta_info->>'canonical') = ?", url))
		}).
		Only(ctx)
}

func (r *repository) GetPropertyByName(ctx context.Context, url string) (*ent.Property, error) {
	fmt.Println("name:from repo ", url)
	return r.db.Property.Query().
		Where(func(s *sql.Selector) {
			// Correctly grouped SQL expression
			s.Where(sql.ExprP("(meta_info->>'canonical') = ?", url))
		}).
		Only(ctx)
}
