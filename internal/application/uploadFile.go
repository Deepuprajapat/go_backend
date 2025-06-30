package application

import (
	"context"
	"io"
	"net/http"

	"github.com/VI-IM/im_backend_go/request"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/VI-IM/im_backend_go/shared/logger"
	"github.com/rs/zerolog/log"
)

func (a *application) UploadFile(file io.Reader, request request.UploadFileRequest) (string, *imhttp.CustomError) {
	if file == nil {
		log.Error().Msg("File is required")
		return "", imhttp.NewCustomErr(http.StatusBadRequest, "File is required", "File is required")
	}

	imageURL, err := a.s3Client.UploadFile(context.Background(), request.AltKeywords, file)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to upload file")
		return "", imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to upload file", err.Error())
	}

	if imageURL == "" {
		log.Error().Msg("Failed to upload file")
		return "", imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to upload file", "Failed to upload file")
	}

	return imageURL, nil
}
