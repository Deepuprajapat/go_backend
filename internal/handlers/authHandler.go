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
	application application.ApplicationInterface
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(application application.ApplicationInterface) *AuthHandler {
	return &AuthHandler{
		application: application,
	}
}

func (h *AuthHandler) GenerateToken(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	var req request.GenerateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Error decoding request")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request", err.Error())
	}

	resp, err := h.application.GetAccessToken(req.Username, req.Password)
	if err != nil {
		log.Error().Err(err).Msg("Error generating token")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request", err.Error())
	}

	return &imhttp.Response{
		Data:       resp,
		StatusCode: http.StatusOK,
	}, nil
}

func (h *AuthHandler) Signup(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	var req request.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Error decoding signup request")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request", err.Error())
	}

	resp, err := h.application.Signup(&req)
	if err != nil {
		log.Error().Err(err).Msg("Error during signup")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Signup failed", err.Error())
	}

	return &imhttp.Response{
		Data:       resp,
		StatusCode: http.StatusCreated,
	}, nil
}

func (h *AuthHandler) RefreshToken(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	var req request.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Error decoding request")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request", err.Error())
	}

	resp, err := h.application.RefreshToken(req.RefreshToken)
	if err != nil {
		log.Error().Err(err).Msg("Error refreshing token")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Refresh token failed", err.Error())
	}

	return &imhttp.Response{
		Data:       resp,
		StatusCode: http.StatusOK,
	}, nil
}
