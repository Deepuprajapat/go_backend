package repository

import (
	"context"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/blogs"
	"github.com/VI-IM/im_backend_go/shared/logger"
)

func (r *repository) GetAllBlogs() ([]*ent.Blogs, error) {
	ctx := context.Background()

	// Get all blogs
	blogList, err := r.db.Blogs.Query().
		Order(ent.Desc(blogs.FieldID)).
		All(ctx)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get blogs")
		return nil, err
	}

	return blogList, nil
}

func (r *repository) GetBlogByID(id string) (*ent.Blogs, error) {
	blog, err := r.db.Blogs.Query().
		Where(blogs.ID(id)).
		Only(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		logger.Get().Error().Err(err).Msg("Failed to get blog")
		return nil, err
	}
	return blog, nil
}
