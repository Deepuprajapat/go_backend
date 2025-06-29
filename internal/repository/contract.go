package repository

import (
	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/internal/domain"
)

type repository struct {
	db *ent.Client
}

type AppRepository interface {
	// Auth
	GetUserDetailsByUsername(username string) (*ent.User, error)

	// Project
	GetProjectByID(id string) (*ent.Project, error)
	AddProject(input domain.Project) (string, error)
	UpdateProject(input domain.Project) (*ent.Project, error)
	DeleteProject(id string, hardDelete bool) error
	IsProjectDeleted(id string) (bool, error)
	GetAllProjects(offset, limit int) ([]*ent.Project, int, error)

	// Developer
	ExistDeveloperByID(id string) (bool, error)

	// Location
	ListLocations() ([]*ent.Location, error)

	// Property
	GetPropertyByID(id string) (*ent.Property, error)
	UpdateProperty(input domain.Property) (*ent.Property, error)
	GetPropertiesOfProject(projectID string) ([]*ent.Property, error)
	AddProperty(input domain.Property) (string, error)
}

func NewRepository(db *ent.Client) AppRepository {
	return &repository{db: db}
}
