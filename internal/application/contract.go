package application

import (
	"context"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/internal/client"
	"github.com/VI-IM/im_backend_go/internal/repository"
	"github.com/VI-IM/im_backend_go/request"
	"github.com/VI-IM/im_backend_go/response"
	imhttp "github.com/VI-IM/im_backend_go/shared"
)

type application struct {
	repo      repository.AppRepository
	s3Client  client.S3ClientInterface
	smsClient client.SMSClientInterface
	crmClient client.CRMClientInterface
}

type ApplicationInterface interface {
	// Auth
	GetAccessToken(username, password string) (*response.GenerateTokenResponse, *imhttp.CustomError)
	RefreshToken(refreshToken string) (*response.GenerateTokenResponse, *imhttp.CustomError)
	Signup(ctx context.Context, req *request.SignupRequest) (*response.GenerateTokenResponse, *imhttp.CustomError)

	// Project
	GetProjectByID(id string) (*response.Project, *imhttp.CustomError)
	AddProject(input request.AddProjectRequest) (*response.AddProjectResponse, *imhttp.CustomError)
	UpdateProject(input request.UpdateProjectRequest) (*response.Project, *imhttp.CustomError)
	DeleteProject(id string) *imhttp.CustomError
	ListProjects(request *request.GetAllAPIRequest) ([]*response.ProjectListResponse, *imhttp.CustomError)
	CompareProjects(projectIDs []string) (*response.ProjectComparisonResponse, *imhttp.CustomError)
	GetProjectByURL(url string) (*ent.Project, *imhttp.CustomError)

	// Developer
	ListDevelopers(pagination *request.GetAllAPIRequest) ([]*response.Developer, *imhttp.CustomError)
	GetDeveloperByID(id string) (*response.Developer, *imhttp.CustomError)
	DeleteDeveloper(id string) *imhttp.CustomError

	// Location
	GetAllLocations(filters map[string]interface{}) ([]*response.Location, *imhttp.CustomError)
	GetLocationByID(id string) (*response.Location, *imhttp.CustomError)
	DeleteLocation(id string) *imhttp.CustomError

	// Property
	GetPropertyByID(id string) (*response.Property, *imhttp.CustomError)
	UpdateProperty(input request.UpdatePropertyRequest) (*response.Property, *imhttp.CustomError)
	GetPropertiesOfProject(projectID string) ([]*response.Property, *imhttp.CustomError)
	AddProperty(input request.AddPropertyRequest) (*response.AddPropertyResponse, *imhttp.CustomError)
	ListProperties(pagination *request.GetAllAPIRequest) ([]*response.PropertyListResponse, int, *imhttp.CustomError)
	DeleteProperty(id string) *imhttp.CustomError

	// Amenity
	GetAllCategoriesWithAmenities() (*response.AmenityResponse, *imhttp.CustomError)
	// GetAmenities() (*response.AmenityResponse, *imhttp.CustomError)
	// GetAmenityByName(name string) (*response.SingleAmenityResponse, *imhttp.CustomError)
	AddCategoryWithAmenities(req *request.CreateAmenityRequest) *imhttp.CustomError
	// UpdateAmenity(id string, req *request.UpdateAmenityRequest) *imhttp.CustomError
	// AddAmenitiesToCategory(req *request.AddAmenitiesToCategoryRequest) *imhttp.CustomError
	// DeleteAmenitiesFromCategory(req *request.DeleteAmenitiesFromCategoryRequest) *imhttp.CustomError
	// DeleteCategory(req *request.DeleteCategoryRequest) *imhttp.CustomError
	// UpdateStaticSiteData(req *request.UpdateStaticSiteDataRequest) *imhttp.CustomError
	UpdateStaticSiteData(req *request.UpdateStaticSiteDataRequest) *imhttp.CustomError

	// Upload File
	UploadFile(request request.UploadFileRequest) (string, string, *imhttp.CustomError)

	// Blogs
	ListBlogs(pagination *request.GetAllAPIRequest) (*response.BlogListResponse, *imhttp.CustomError)
	GetBlogByID(id string) (*response.BlogResponse, *imhttp.CustomError)
	CreateBlog(ctx context.Context, req *request.CreateBlogRequest) (*response.BlogResponse, *imhttp.CustomError)
	DeleteBlog(ctx context.Context, id string) *imhttp.CustomError
	UpdateBlog(ctx context.Context, id string, req *request.UpdateBlogRequest) (*response.BlogResponse, *imhttp.CustomError)

	// Content
	GetProjectByCanonicalURL(ctx context.Context, url string) (*ent.Project, *imhttp.CustomError)
	GetPropertyByCanonicalURL(ctx context.Context, url string) (*ent.Property, *imhttp.CustomError)
	GetBlogByCanonicalURL(ctx context.Context, url string) (*ent.Blogs, *imhttp.CustomError)

	// Generic Search
	GetCustomSearchPage(ctx context.Context, slug string) (*response.CustomSearchPage, *imhttp.CustomError)
	GetLinks(ctx context.Context) ([]*response.Link, *imhttp.CustomError)
	GetAllCustomSearchPages(ctx context.Context) ([]*response.CustomSearchPage, *imhttp.CustomError)
	AddCustomSearchPage(ctx context.Context, customSearchPage *request.CustomSearchPage) (*response.CustomSearchPage, *imhttp.CustomError)
	UpdateCustomSearchPage(ctx context.Context, customSearchPage *request.CustomSearchPage) (*response.CustomSearchPage, *imhttp.CustomError)
	DeleteCustomSearchPage(ctx context.Context, id string) *imhttp.CustomError

	// Leads
	CreateLeadWithOTP(ctx context.Context, req *request.CreateLeadRequest) (*response.CreateLeadResponse, *imhttp.CustomError)
	CreateLead(ctx context.Context, req *request.CreateLeadRequest) (*response.CreateLeadResponse, *imhttp.CustomError)
	GetLeadByID(ctx context.Context, id int) (*response.Lead, *imhttp.CustomError)
	GetAllLeads(ctx context.Context, req *request.GetLeadsRequest) (*response.LeadListResponse, *imhttp.CustomError)
	ValidateOTP(ctx context.Context, req *request.ValidateOTPRequest) (*response.ValidateOTPResponse, *imhttp.CustomError)
	ResendOTP(ctx context.Context, req *request.ResendOTPRequest) (*response.ResendOTPResponse, *imhttp.CustomError)
}

func NewApplication(repo repository.AppRepository, s3Client client.S3ClientInterface, smsClient client.SMSClientInterface, crmClient client.CRMClientInterface) ApplicationInterface {
	return &application{repo: repo, s3Client: s3Client, smsClient: smsClient, crmClient: crmClient}
}
