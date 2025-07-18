package repository

import (
	"context"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/blogs"
	"github.com/VI-IM/im_backend_go/ent/schema"
	"github.com/VI-IM/im_backend_go/shared/logger"
	"github.com/google/uuid"
)

// type BlogRepository interface {
// 	GetAllBlogs() ([]*ent.Blogs, error)
// 	GetBlogByID(id string) (*ent.Blogs, error)
// 	CreateBlog(ctx context.Context, blogURL string, blogContent schema.BlogContent, seoMetaInfo schema.SEOMetaInfo, isPriority bool) (*ent.Blogs, error)
// 	DeleteBlog(ctx context.Context, id string) error
// 	UpdateBlog(ctx context.Context, id string, blogURL *string, blogContent *schema.BlogContent, seoMetaInfo *schema.SEOMetaInfo, isPriority *bool) (*ent.Blogs, error)
// }

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

func (r *repository) CreateBlog(ctx context.Context, slug string, blogContent schema.BlogContent, seoMetaInfo schema.SEOMetaInfo, isPriority bool) (*ent.Blogs, error) {
	blog, err := r.db.Blogs.Create().
		SetID(uuid.New().String()).
		SetSlug(slug).
		SetBlogContent(blogContent).
		SetSeoMetaInfo(seoMetaInfo).
		SetIsPriority(isPriority).
		Save(ctx)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to create blog")
		return nil, err
	}
	return blog, nil
}

func (r *repository) DeleteBlog(ctx context.Context, id string) error {
	// First check if blog exists
	blog, err := r.GetBlogByID(id)
	if err != nil {
		return err
	}
	if blog == nil {
		return nil // Blog doesn't exist, nothing to delete
	}

	// Soft delete by updating is_deleted flag
	_, err = r.db.Blogs.UpdateOneID(id).
		SetIsDeleted(true).
		Save(ctx)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to delete blog")
		return err
	}

	return nil
}

func (r *repository) UpdateBlog(ctx context.Context, id string, blogURL *string, blogContent *schema.BlogContent, seoMetaInfo *schema.SEOMetaInfo, isPriority *bool) (*ent.Blogs, error) {
	// First check if blog exists
	blog, err := r.GetBlogByID(id)
	if err != nil {
		return nil, err
	}
	if blog == nil {
		return nil, nil // Blog doesn't exist
	}

	update := r.db.Blogs.UpdateOneID(id)

	if blogURL != nil {
		update.SetSlug(*blogURL)
	}
	if blogContent != nil {
		update.SetBlogContent(*blogContent)
	}
	if seoMetaInfo != nil {
		update.SetSeoMetaInfo(*seoMetaInfo)
	}
	if isPriority != nil {
		update.SetIsPriority(*isPriority)
	}

	blog, err = update.Save(ctx)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to update blog")
		return nil, err
	}

	return blog, nil
}

