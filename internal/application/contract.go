package application

import (
	"io"

	"github.com/VI-IM/im_backend_go/internal/client"
	"github.com/VI-IM/im_backend_go/internal/repository"
	"github.com/VI-IM/im_backend_go/request"
	"github.com/VI-IM/im_backend_go/response"
	imhttp "github.com/VI-IM/im_backend_go/shared"
)

type application struct {
	repo     repository.AppRepository
	s3Client client.S3ClientInterface
}

type ApplicationInterface interface {
	// Auth
	GetAccessToken(username, password string) (*response.GenerateTokenResponse, *imhttp.CustomError)
	RefreshToken(refreshToken string) (*response.GenerateTokenResponse, *imhttp.CustomError)

	// Project
	GetProjectByID(id string) (*response.Project, *imhttp.CustomError)
	AddProject(input request.AddProjectRequest) (*response.AddProjectResponse, *imhttp.CustomError)
	UpdateProject(input request.UpdateProjectRequest) (*response.Project, *imhttp.CustomError)
	DeleteProject(id string) *imhttp.CustomError
	ListProjects(filters map[string]interface{}) ([]*response.ProjectListResponse, *imhttp.CustomError)

	// Developer
	ListDevelopers(pagination *request.PaginationRequest) ([]*response.Developer, int, *imhttp.CustomError)
	GetDeveloperByID(id string) (*response.Developer, *imhttp.CustomError)
	DeleteDeveloper(id string) *imhttp.CustomError

	// Location
	GetAllLocations() ([]*response.Location, *imhttp.CustomError)
	GetLocationByID(id string) (*response.Location, *imhttp.CustomError)
	DeleteLocation(id string) *imhttp.CustomError

	// Property
	GetPropertyByID(id string) (*response.Property, *imhttp.CustomError)
	UpdateProperty(input request.UpdatePropertyRequest) (*response.Property, *imhttp.CustomError)
	GetPropertiesOfProject(projectID string) ([]*response.Property, *imhttp.CustomError)
	AddProperty(input request.AddPropertyRequest) (*response.AddPropertyResponse, *imhttp.CustomError)
	ListProperties(pagination *request.PaginationRequest) ([]*response.PropertyListResponse, int, *imhttp.CustomError)
	DeleteProperty(id string) *imhttp.CustomError

	// Amenity
	GetAmenities() (*response.AmenityResponse, *imhttp.CustomError)
	GetAmenityByID(id string) (*response.SingleAmenityResponse, *imhttp.CustomError)
	CreateAmenity(req *request.CreateAmenityRequest) *imhttp.CustomError
	UpdateAmenity(id string, req *request.UpdateAmenityRequest) *imhttp.CustomError

	// Upload File
	UploadFile(file io.Reader, request request.UploadFileRequest) (string, *imhttp.CustomError)

	// Blogs
	ListBlogs(pagination *request.PaginationRequest) (*response.BlogListResponse, *imhttp.CustomError)
	GetBlogByID(id string) (*response.BlogResponse, *imhttp.CustomError)
}

func NewApplication(repo repository.AppRepository, s3Client client.S3ClientInterface) ApplicationInterface {
	return &application{repo: repo, s3Client: s3Client}
}
