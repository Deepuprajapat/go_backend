package application

import (
	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/internal/repository"
	"github.com/VI-IM/im_backend_go/request"
	"github.com/VI-IM/im_backend_go/response"
)

type application struct {
	repo repository.AppRepository
}

type ApplicationInterface interface {
	GetAccessToken(username string, password string) (*response.GenerateTokenResponse, error)
	RefreshToken(refreshToken string) (*response.GenerateTokenResponse, error)
	AddProject(input request.AddProjectRequest) (*response.AddProjectResponse, error)
	GetProjectByID(id string) (*ent.Project, error)
}

func NewApplication(repo repository.AppRepository) ApplicationInterface {
	return &application{repo: repo}
}
