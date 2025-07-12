package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/VI-IM/im_backend_go/request"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/rs/zerolog/log"
)

func (h *Handler) GenerateToken(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	var req request.GenerateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request", err.Error())
	}
	resp, err := h.app.GetAccessToken(req.Email, req.Password)
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request", err.Error())
	}

	return &imhttp.Response{
		Data:       resp,
		StatusCode: http.StatusOK,
	}, nil
}

// RefreshToken refreshes the access token
func (h *Handler) RefreshToken(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	var req request.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Error decoding request")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request", err.Error())
	}

	resp, err := h.app.RefreshToken(req.RefreshToken)
	if err != nil {
		log.Error().Err(err).Msg("Error refreshing token")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Refresh token failed", err.Error())
	}

	return &imhttp.Response{
		Data:       resp,
		StatusCode: http.StatusOK,
	}, nil
}

func (h *Handler) Signup(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	var req request.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Error decoding request")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request", err.Error())
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request", err.Error())
	}

	// Call application layer to create user
	resp, err := h.app.Signup(r.Context(), &req)
	if err != nil {
		return nil, err
	}

	return &imhttp.Response{
		Data:       resp,
		StatusCode: http.StatusCreated,
		Message:    "User created successfully",
	}, nil
}
