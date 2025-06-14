package application

import "github.com/VI-IM/im_backend_go/ent"

func (c *application) GetProjectByID(id int) (*ent.Project, error) {
	project, err := c.repo.GetProjectByID(id)
	if err != nil {
		return nil, err
	}

	return project, nil
}
