package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/VI-IM/im_backend_go/internal/application"
	"github.com/VI-IM/im_backend_go/request"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/rs/zerolog/log"
)

type AuthHandler struct {
	controller application.ApplicationInterface
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(controller application.ApplicationInterface) *AuthHandler {
	return &AuthHandler{
		controller: controller,
	}
}

func (h *AuthHandler) GenerateToken(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	var req request.GenerateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Error decoding request")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request", err.Error())
	}

	resp, err := h.controller.GetAccessToken(req.Username, req.Password)
	if err != nil {
		log.Error().Err(err).Msg("Error generating token")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request", err.Error())
	}

	return &imhttp.Response{
		Data:       resp,
		StatusCode: http.StatusOK,
	}, nil
}

func RefreshToken(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	var req request.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Error decoding request")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request", err.Error())
	}

	return &imhttp.Response{
		Data: "Token refreshed",
	}, nil
}
