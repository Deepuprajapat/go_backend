package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/VI-IM/im_backend_go/request"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/rs/zerolog/log"
)

// UploadFile handles file upload requests
// UploadFile handles file upload requests
func (h *Handler) UploadFile(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	// Parse multipart form
	var req request.UploadFileRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Error().Err(err).Msg("Error decoding request")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request", err.Error())
	}

	if req.FileName == "" || req.FilePath == "" || req.AltKeywords == "" {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "File name, file path and alt keywords are required", "File name, file path and alt keywords are required")
	}
	presignedURL, imageURL, err := h.app.UploadFile(req)
	
	return &imhttp.Response{
		Data: struct {
			PresignedURL string `json:"presigned_url"`
			ImageURL     string `json:"image_url"`
		}{
			PresignedURL: presignedURL,
			ImageURL:     imageURL,
		},
		StatusCode: http.StatusOK,
	}, nil
}
