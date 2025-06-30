package handlers

import (
	"net/http"

	"github.com/VI-IM/im_backend_go/request"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/rs/zerolog/log"
)

// UploadFile handles file upload requests
func (h *Handler) UploadFile(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	// Parse multipart form
	err := r.ParseMultipartForm(32 << 20) // 32MB max memory
	if err != nil {
		log.Error().Err(err).Msg("Error parsing multipart form")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid multipart form", err.Error())
	}

	// Get the file from the form
	file, _, err := r.FormFile("file")
	if err != nil {
		log.Error().Err(err).Msg("Error getting file from form")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "File is required", err.Error())
	}
	defer file.Close()

	// Get other form values
	altKeywords := r.FormValue("alt_keywords")
	filePath := r.FormValue("file_path")

	// Create the request struct
	req := request.UploadFileRequest{
		File:        file,
		AltKeywords: altKeywords,
		FilePath:    filePath,
	}

	// Call the application layer
	fileURL, customErr := h.app.UploadFile(file, req)
	if customErr != nil {
		log.Error().Err(customErr).Msg("Error uploading file")
		return nil, customErr
	}

	return &imhttp.Response{
		Data:       map[string]string{"file_url": fileURL},
		StatusCode: http.StatusOK,
	}, nil
}
