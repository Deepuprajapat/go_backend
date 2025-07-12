package application

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/VI-IM/im_backend_go/request"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/VI-IM/im_backend_go/shared/logger"
)

func (a *application) UploadFile(request request.UploadFileRequest) (string, string, *imhttp.CustomError) {
	presignedURL, err := a.s3Client.GeneratePresignedURL(context.Background(), request.FileName, "PUT", 1*time.Hour)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to generate presigned URL")
		return "", "", imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to generate presigned URL", err.Error())
	}

	imageURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", os.Getenv("AWS_S3_BUCKET"), request.FileName)

	return presignedURL, imageURL, nil
}
