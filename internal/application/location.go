package application

import (
	"net/http"

	"github.com/VI-IM/im_backend_go/response"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/VI-IM/im_backend_go/shared/logger"
)

func (c *application) GetAllLocations(filters map[string]interface{}) ([]*response.Location, *imhttp.CustomError) {
	locations, err := c.repo.ListLocations(filters)
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

func (c *application) GetLocationByID(id string) (*response.Location, *imhttp.CustomError) {
	location, err := c.repo.GetLocationByID(id)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get location")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get location", err.Error())
	}

	return response.GetLocationFromEnt(location), nil
}

func (c *application) DeleteLocation(id string) *imhttp.CustomError {
	if err := c.repo.SoftDeleteLocation(id); err != nil {
		logger.Get().Error().Err(err).Msg("Failed to delete location")
		return imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to delete location", err.Error())
	}
	return nil
}
