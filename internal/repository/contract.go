package repository

import (
	"context"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/schema"
	"github.com/VI-IM/im_backend_go/internal/domain"
	"github.com/VI-IM/im_backend_go/response"
)

type repository struct {
	db *ent.Client
}

type AppRepository interface {
	// Auth
	GetUserDetailsByEmail(ctx context.Context, email string) (*ent.User, error)
	CreateUser(ctx context.Context, user *ent.User) (*ent.User, error)
	CheckIfUserExistsByEmail(ctx context.Context, email string) (bool, error)
	CheckIfUserExistsByID(ctx context.Context, userID string) (bool, error)

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
	ListLocations(filters map[string]interface{}) ([]*ent.Location, error)
	GetLocationByID(id string) (*ent.Location, error)
	AddLocation(localityName, city, state, phoneNumber, country, pincode string) (*ent.Location, error)
	SoftDeleteLocation(id string) error
	GetAllUniqueCities() ([]string, error)
	GetAllUniqueLocations() ([]string, error)

	// Property
	GetPropertyByID(id string) (*ent.Property, error)
	UpdateProperty(input domain.Property) (*ent.Property, error)
	GetPropertiesOfProject(projectID string) ([]*ent.Property, error)
	AddProperty(input domain.Property) (*PropertyResult, error)
	GetAllProperties(offset, limit int, filters map[string]interface{}) ([]*ent.Property, int, error)
	DeleteProperty(id string, hardDelete bool) error
	IsPropertyDeleted(id string) (bool, error)

	// Static Site Data
	GetStaticSiteData() (*ent.StaticSiteData, error)
	UpdateStaticSiteData(data *ent.StaticSiteData) error
	CheckCategoryExists(category string) (bool, error)

	// Custom Search
	CheckURLExists(ctx context.Context, url string) (*response.URLExistsResult, error)

	// Blogs
	GetAllBlogs() ([]*ent.Blogs, error)
	GetAllBlogsWithFilter(isPublished *bool) ([]*ent.Blogs, error)  // ✅ Add this
	GetBlogByID(id string) (*ent.Blogs, error)
	CreateBlog(ctx context.Context, slug string, blogContent schema.BlogContent, seoMetaInfo schema.SEOMetaInfo, isPriority bool, isPublished bool) (*ent.Blogs, error)
	DeleteBlog(ctx context.Context, id string) error
	UpdateBlog(ctx context.Context, id string, blogURL *string, blogContent *schema.BlogContent, seoMetaInfo *schema.SEOMetaInfo, isPriority *bool) (*ent.Blogs, error)

	//content

	GetProjectByCanonicalURL(ctx context.Context, canonicalURL string) (*ent.Project, error)
	GetPropertyByCanonicalURL(ctx context.Context, canonicalURL string) (*ent.Property, error)
	GetBlogByCanonicalURL(ctx context.Context, canonicalURL string) (*ent.Blogs, error)

	// Generic Search
	GetCustomSearchPageFromSlug(ctx context.Context, slug string) (*ent.CustomSearchPage, error)
	GetAllCustomSearchPages(ctx context.Context) ([]*ent.CustomSearchPage, error)
	AddCustomSearchPage(ctx context.Context, customSearchPage *ent.CustomSearchPage) (*ent.CustomSearchPage, error)
	UpdateCustomSearchPage(ctx context.Context, customSearchPage *ent.CustomSearchPage) (*ent.CustomSearchPage, error)
	DeleteCustomSearchPage(ctx context.Context, id string) error

	// Leads
	CreateLead(ctx context.Context, lead *ent.Leads) (*ent.Leads, error)
	GetLeadByID(ctx context.Context, id int) (*ent.Leads, error)
	GetLeadByPhone(ctx context.Context, phone string) (*ent.Leads, error)
	GetLeadByPhoneAndOTP(ctx context.Context, phone, otp string) (*ent.Leads, error)
	UpdateLead(ctx context.Context, lead *ent.Leads) (*ent.Leads, error)
	GetAllLeads(ctx context.Context, filters map[string]interface{}) ([]*ent.Leads, error)
	GetLeadsByDate(ctx context.Context, date string) ([]*ent.Leads, error)
}

func NewRepository(db *ent.Client) AppRepository {
	return &repository{db: db}
}
