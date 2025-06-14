package repository

import (
	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/internal/domain"
)

type repository struct {
	db *ent.Client
}

type AppRepository interface {
	GetUserDetailsByUsername(username string) (*ent.User, error)
	GetProjectByID(id int) (*ent.Project, error)
	AddProject(input domain.AddProjectInput) (string, error)
	ExistDeveloperByID(id string) (bool, error)
}

func NewRepository(db *ent.Client) AppRepository {
	return &repository{db: db}
}
