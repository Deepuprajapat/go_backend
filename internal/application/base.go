package application

import (
	"github.com/VI-IM/im_backend_go/internal/repository"
	"github.com/VI-IM/im_backend_go/request"
	"github.com/VI-IM/im_backend_go/response"
)

type Application struct {
	repo repository.AppRepository
}

type ApplicationInterface interface {
	GetAccessToken(username string, password string) (*response.GenerateTokenResponse, error)
	RefreshToken(refreshToken string) (*response.GenerateTokenResponse, error)
	Signup(req *request.SignupRequest) (*response.SignupResponse, error)
}

func NewApplication(repo repository.AppRepository) *Application {
	return &Application{repo: repo}
}
