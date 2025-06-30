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
		log.Error().Err(err).Msg("Error decoding request")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request", err.Error())
	}
	resp, err := h.app.GetAccessToken(req.Username, req.Password)
	if err != nil {
		log.Error().Err(err).Msg("Error generating token")
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

	return &imhttp.Response{
		Data:       "Token refreshed",
		StatusCode: http.StatusOK,
	}, nil
}
