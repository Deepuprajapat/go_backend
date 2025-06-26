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
	GetProjectByID(id string) (*ent.Project, error)
	AddProject(input domain.Project) (string, error)
	UpdateProject(input domain.Project) (*ent.Project, error)
	DeleteProject(id string, hardDelete bool) error
	ExistDeveloperByID(id string) (bool, error)
	IsProjectDeleted(id string) (bool, error)
	GetPropertyByID(id string) (*ent.Property, error)
	UpdateProperty(input domain.Property) (*ent.Property, error)
	GetAllProjects() ([]*ent.Project, error)
	ListLocations() ([]*ent.Location, error)
}

func NewRepository(db *ent.Client) AppRepository {
	return &repository{db: db}
}
