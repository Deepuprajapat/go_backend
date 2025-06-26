package application

import (
	"net/http"

	"github.com/VI-IM/im_backend_go/response"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/VI-IM/im_backend_go/shared/logger"
)

func (c *application) ListLocations() ([]*response.Location, *imhttp.CustomError) {
	locations, err := c.repo.ListLocations()
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to list locations")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to list locations", err.Error())
	}

	var locationResponses []*response.Location
	for _, location := range locations {
		locationResponses = append(locationResponses, response.GetLocationFromEnt(location))
	}

	return locationResponses, nil
}
