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
	GetAccessToken(username string, password string) (*response.GenerateTokenResponse, *imhttp.CustomError)
	RefreshToken(refreshToken string) (*response.GenerateTokenResponse, *imhttp.CustomError)
	AddProject(input request.AddProjectRequest) (*response.AddProjectResponse, *imhttp.CustomError)
	GetProjectByID(id string) (*response.Project, *imhttp.CustomError)
	UpdateProject(input request.UpdateProjectRequest) (*response.Project, *imhttp.CustomError)
	DeleteProject(id string) *imhttp.CustomError
	GetPropertyByID(id string) (*response.Property, *imhttp.CustomError)
	UpdateProperty(input request.UpdatePropertyRequest) (*response.Property, *imhttp.CustomError)
	ListProjects() ([]*response.ProjectListResponse, *imhttp.CustomError)
	GetAllLocations() ([]*response.Location, *imhttp.CustomError)
}

func NewApplication(repo repository.AppRepository) ApplicationInterface {
	return &application{repo: repo}
}
