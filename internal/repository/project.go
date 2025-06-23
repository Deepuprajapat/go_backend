package repository

import (
	"context"
	"errors"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/developer"
	"github.com/VI-IM/im_backend_go/ent/schema"
	"github.com/VI-IM/im_backend_go/internal/domain"
	"github.com/VI-IM/im_backend_go/internal/domain/enums"
)

func (r *repository) GetProjectByID(id string) (*ent.Project, error) {
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

func (r *repository) AddProject(input domain.AddProjectInput) (string, error) {
	if err := r.db.Project.Create().
		SetID(input.ProjectID).
		SetName(input.ProjectName).
		SetStatus(enums.ProjectStatus("")).
		SetMetaInfo(schema.SEOMeta{
			Canonical: input.ProjectURL,
		}).
		SetDescription("").
		SetDeveloperID(input.DeveloperID).
		Exec(context.Background()); err != nil {
		return "", err
	}
	return input.ProjectID, nil
}

func (r *repository) ExistDeveloperByID(id string) (bool, error) {
	exist, err := r.db.Developer.Query().
		Where(developer.ID(id)).
		Exist(context.Background())
	if err != nil {
		return false, err
	}
	return exist, nil
}
