package application

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/VI-IM/im_backend_go/internal/domain"
	"github.com/VI-IM/im_backend_go/request"
	"github.com/VI-IM/im_backend_go/response"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/VI-IM/im_backend_go/shared/logger"
)

func (c *application) GetProjectByID(id string) (*response.GetProjectResponse, *imhttp.CustomError) {

	// check if the project is deleted
	isDeleted, err := c.repo.IsProjectDeleted(id)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to check if project is deleted")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to check if project is deleted", err.Error())
	}
	if isDeleted {
		return nil, imhttp.NewCustomErr(http.StatusNotFound, "Project not found or deleted", "Project not found or deleted")
	}

	project, err := c.repo.GetProjectByID(id)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get project")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get project", err.Error())
	}

	return response.GetProjectFromEnt(project), nil
}

func (c *application) AddProject(input request.AddProjectRequest) (*response.AddProjectResponse, *imhttp.CustomError) {

	var project domain.AddProjectInput

	// get developer from developer id
	exist, err := c.repo.ExistDeveloperByID(input.DeveloperID)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to add project")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to add project", err.Error())
	}
	if !exist {
		return nil, imhttp.NewCustomErr(http.StatusNotFound, "Developer not found", "Developer not found")
	}

	project.ProjectID = fmt.Sprintf("%x", sha256.Sum256([]byte(strconv.FormatInt(time.Now().Unix(), 10))))[:16]
	project.ProjectName = input.ProjectName
	project.ProjectURL = input.ProjectURL
	project.DeveloperID = input.DeveloperID

	projectID, err := c.repo.AddProject(project)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to add project")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to add project", err.Error())
	}

	return &response.AddProjectResponse{
		ProjectID: projectID,
	}, nil
}

func (c *application) UpdateProject(input request.UpdateProjectRequest) (*response.UpdateProjectResponse, *imhttp.CustomError) {
	return nil, nil
}
