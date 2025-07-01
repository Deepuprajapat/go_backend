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
	GetAllProjects(filters map[string]interface{}) ([]*ent.Project, error)

	// Developer
	ExistDeveloperByID(id string) (bool, error)
	GetAllDevelopers(offset, limit int) ([]*ent.Developer, int, error)
	GetDeveloperByID(id string) (*ent.Developer, error)
	SoftDeleteDeveloper(id string) error

	// Location
	ListLocations() ([]*ent.Location, error)
	GetLocationByID(id string) (*ent.Location, error)
	SoftDeleteLocation(id string) error

	// Property
	GetPropertyByID(id string) (*ent.Property, error)
	UpdateProperty(input domain.Property) (*ent.Property, error)
	GetPropertiesOfProject(projectID string) ([]*ent.Property, error)
	AddProperty(input domain.Property) (string, error)
	GetAllProperties(offset, limit int) ([]*ent.Property, int, error)
	DeleteProperty(id string, hardDelete bool) error
	IsPropertyDeleted(id string) (bool, error)

	// Static Site Data
	GetStaticSiteData() (*ent.StaticSiteData, error)
	UpdateStaticSiteData(data *ent.StaticSiteData) error
}

func NewRepository(db *ent.Client) AppRepository {
	return &repository{db: db}
}
