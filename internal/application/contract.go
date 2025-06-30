package application

import (
	"github.com/VI-IM/im_backend_go/internal/repository"
	"github.com/VI-IM/im_backend_go/request"
	"github.com/VI-IM/im_backend_go/response"
	imhttp "github.com/VI-IM/im_backend_go/shared"
)

type application struct {
	repo repository.AppRepository
}

type ApplicationInterface interface {
	// Auth
	GetAccessToken(username string, password string) (*response.GenerateTokenResponse, *imhttp.CustomError)
	RefreshToken(refreshToken string) (*response.GenerateTokenResponse, *imhttp.CustomError)

	// Project
	AddProject(input request.AddProjectRequest) (*response.AddProjectResponse, *imhttp.CustomError)
	GetProjectByID(id string) (*response.Project, *imhttp.CustomError)
	UpdateProject(input request.UpdateProjectRequest) (*response.Project, *imhttp.CustomError)
	DeleteProject(id string) *imhttp.CustomError

	// Property
	GetPropertyByID(id string) (*response.Property, *imhttp.CustomError)
	UpdateProperty(input request.UpdatePropertyRequest) (*response.Property, *imhttp.CustomError)
	GetPropertiesOfProject(projectID string) ([]*response.Property, *imhttp.CustomError)
	AddProperty(input request.AddPropertyRequest) (*response.AddPropertyResponse, *imhttp.CustomError)
	ListProperties(pagination *request.PaginationRequest) ([]*response.PropertyListResponse, int, *imhttp.CustomError)
	DeleteProperty(id string) *imhttp.CustomError

	// Developer
	ListDevelopers(pagination *request.PaginationRequest) ([]*response.Developer, int, *imhttp.CustomError)
	GetDeveloperByID(id string) (*response.Developer, *imhttp.CustomError)
	DeleteDeveloper(id string) *imhttp.CustomError

	// Location
	ListProjects(pagination *request.PaginationRequest) ([]*response.ProjectListResponse, int, *imhttp.CustomError)
	GetAllLocations() ([]*response.Location, *imhttp.CustomError)
	GetLocationByID(id string) (*response.Location, *imhttp.CustomError)
	DeleteLocation(id string) *imhttp.CustomError

	// Amenity
	GetAmenities() (*response.AmenityResponse, *imhttp.CustomError)
	GetAmenityByID(id string) (*response.SingleAmenityResponse, *imhttp.CustomError)
	CreateAmenity(req *request.CreateAmenityRequest) *imhttp.CustomError
	UpdateAmenity(id string, req *request.UpdateAmenityRequest) *imhttp.CustomError
}

func NewApplication(repo repository.AppRepository) ApplicationInterface {
	return &application{repo: repo}
}
