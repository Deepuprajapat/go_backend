package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"github.com/VI-IM/im_backend_go/internal/controller"
	"github.com/VI-IM/im_backend_go/request"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/rs/zerolog/log"
)

type AuthHandler struct {
	controller controller.ControllerInterface
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(controller controller.ControllerInterface) *AuthHandler {
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

func (h *AuthHandler) Signup(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	var req request.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Error decoding signup request")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request", err.Error())
	}

	resp, err := h.controller.Signup(&req)
	if err != nil {
		log.Error().Err(err).Msg("Error during signup")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Signup failed", err.Error())
	}

	return &imhttp.Response{
		Data:       resp,
		StatusCode: http.StatusCreated,
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

func (h *AuthHandler) Signout(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, imhttp.NewCustomErr(http.StatusUnauthorized, "Unauthorized", "No token provided")
	}

	// Check if the header starts with "Bearer "
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, imhttp.NewCustomErr(http.StatusUnauthorized, "Unauthorized", "Invalid token format")
	}

	token := parts[1]
	err := h.controller.InvalidateToken(token)
	if err != nil {
		log.Error().Err(err).Msg("Error invalidating token")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Signout failed", err.Error())
	}

	return &imhttp.Response{
		Data:       "Successfully signed out",
		StatusCode: http.StatusOK,
	}, nil
}
