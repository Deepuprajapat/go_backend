package handlers

import (
	"github.com/VI-IM/im_backend_go/internal/application"
)

type Handler struct {
	app application.ApplicationInterface
}

func NewHandler(app application.ApplicationInterface) *Handler {
	return &Handler{app: app}
}
