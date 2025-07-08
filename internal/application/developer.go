package application

import (
	"net/http"

	"github.com/VI-IM/im_backend_go/request"
	"github.com/VI-IM/im_backend_go/response"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/VI-IM/im_backend_go/shared/logger"
)

func (c *application) ListDevelopers(input *request.GetAllAPIRequest) ([]*response.Developer, int, *imhttp.CustomError) {
	developers, totalItems, err := c.repo.GetAllDevelopers(input.GetOffset(), input.GetLimit())
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to list developers")
		return nil, 0, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to list developers", err.Error())
	}

	var developerResponses []*response.Developer
	for _, developer := range developers {
		developerResponses = append(developerResponses, response.GetDeveloperFromEnt(developer))
	}

	return developerResponses, totalItems, nil
}

func (c *application) GetDeveloperByID(id string) (*response.Developer, *imhttp.CustomError) {
	developer, err := c.repo.GetDeveloperByID(id)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get developer")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get developer", err.Error())
	}

	return response.GetDeveloperFromEnt(developer), nil
}

func (c *application) DeleteDeveloper(id string) *imhttp.CustomError {
	if err := c.repo.SoftDeleteDeveloper(id); err != nil {
		logger.Get().Error().Err(err).Msg("Failed to delete developer")
		return imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to delete developer", err.Error())
	}
	return nil
}
