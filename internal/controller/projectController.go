package controller

import (
	"github.com/VI-IM/im_backend_go/response"
)

func (c *Controller) GetProjectByID(id int) (*response.ProjectResponse, error) {
	project, err := c.repo.GetProjectByID(id)
	if err != nil {
		return nil, err
	}

	return response.NewProjectResponse(project), nil
}
