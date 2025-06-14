package controller

import (
	"github.com/VI-IM/im_backend_go/internal/repository"
	"github.com/VI-IM/im_backend_go/response"
)

type Controller struct {
	repo repository.AppRepository
}

type ControllerInterface interface {
	GetAccessToken(username string, password string) (*response.GenerateTokenResponse, error)
	GetProjectByID(id int) (*response.ProjectResponse, error)
}

func NewController(repo repository.AppRepository) *Controller {
	return &Controller{repo: repo}
}
