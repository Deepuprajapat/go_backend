package application

import (
	"context"
	"net/http"

	"github.com/VI-IM/im_backend_go/response"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/VI-IM/im_backend_go/shared/logger"
)

func (c *application) CheckURLExists(ctx context.Context, url string) (*response.CheckURLExistsResponse, *imhttp.CustomError) {
	result, err := c.repo.CheckURLExists(ctx, url)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to check URL existence")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to check URL existence", err.Error())
	}

	return &response.CheckURLExistsResponse{
		Exists:     result.Exists,
		EntityType: result.EntityType,
		EntityID:   result.EntityID,
	}, nil
}
