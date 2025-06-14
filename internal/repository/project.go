package repository

import (
	"context"
	"errors"

	"github.com/VI-IM/im_backend_go/ent"
)

func (r *repository) GetProjectByID(id int) (*ent.Project, error) {
	project, err := r.db.Project.Get(context.Background(), id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("project not found")
		}
		return nil, err
	}

	if project.IsDeleted {
		return nil, errors.New("project is deleted")
	}

	return project, nil
}
