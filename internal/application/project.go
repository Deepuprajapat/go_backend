package application

import (
	"errors"
	"strconv"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/internal/domain"
	"github.com/VI-IM/im_backend_go/request"
	"github.com/VI-IM/im_backend_go/response"
	"github.com/google/uuid"
)

func (c *application) GetProjectByID(id string) (*ent.Project, error) {
	project, err := c.repo.GetProjectByID(id)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (c *application) AddProject(input request.AddProjectRequest) (*response.AddProjectResponse, error) {

	var project domain.AddProjectInput

	// get developer from developer id
	exist, err := c.repo.ExistDeveloperByID(input.DeveloperID)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("developer not found")
	}

	ID := uuid.New().ID()
	project.ProjectID = strconv.FormatUint(uint64(ID), 10)
	projectID, err := c.repo.AddProject(project)
	if err != nil {
		return nil, err
	}

	return &response.AddProjectResponse{
		ProjectID: projectID,
	}, nil
}
