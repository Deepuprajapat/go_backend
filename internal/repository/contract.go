package repository

import (
	"context"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/schema"
	"github.com/VI-IM/im_backend_go/internal/domain"
)

type repository struct {
	db *ent.Client
}

type AppRepository interface {
	// Auth
	GetUserDetailsByUsername(username string) (*ent.User, error)

	// Project
	GetProjectByID(id string) (*ent.Project, error)
	AddProject(input domain.Project) (string, error)
	UpdateProject(input domain.Project) (*ent.Project, error)
	DeleteProject(id string, hardDelete bool) error
	IsProjectDeleted(id string) (bool, error)
	GetAllProjects(filters map[string]interface{}) ([]*ent.Project, error)
	GetProjectByURL(url string) (*ent.Project, error)

	// Developer
	ExistDeveloperByID(id string) (bool, error)
	GetAllDevelopers() ([]*ent.Developer, error)
	GetDeveloperByID(id string) (*ent.Developer, error)
	SoftDeleteDeveloper(id string) error

	// Location
	ListLocations() ([]*ent.Location, error)
	GetLocationByID(id string) (*ent.Location, error)
	SoftDeleteLocation(id string) error

	// Property
	GetPropertyByID(id string) (*ent.Property, error)
	UpdateProperty(input domain.Property) (*ent.Property, error)
	GetPropertiesOfProject(projectID string) ([]*ent.Property, error)
	AddProperty(input domain.Property) (string, error)
	GetAllProperties(offset, limit int, filters map[string]interface{}) ([]*ent.Property, int, error)
	DeleteProperty(id string, hardDelete bool) error
	IsPropertyDeleted(id string) (bool, error)

	// Static Site Data
	GetStaticSiteData() (*ent.StaticSiteData, error)
	UpdateStaticSiteData(data *ent.StaticSiteData) error
	CheckCategoryExists(category string) (bool, error)

	// Blogs
	GetAllBlogs() ([]*ent.Blogs, error)
	GetBlogByID(id string) (*ent.Blogs, error)
	CreateBlog(ctx context.Context, blogURL string, blogContent schema.BlogContent, seoMetaInfo schema.SEOMetaInfo, isPriority bool) (*ent.Blogs, error)
	DeleteBlog(ctx context.Context, id string) error
	UpdateBlog(ctx context.Context, id string, blogURL *string, blogContent *schema.BlogContent, seoMetaInfo *schema.SEOMetaInfo, isPriority *bool) (*ent.Blogs, error)
}

func NewRepository(db *ent.Client) AppRepository {
	return &repository{db: db}
}
