package handlers

import (
	"github.com/VI-IM/im_backend_go/internal/application"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	app      application.ApplicationInterface
	validate *validator.Validate
}

func NewHandler(app application.ApplicationInterface) *Handler {
	return &Handler{
		app:      app,
		validate: validator.New(),
	}
}
