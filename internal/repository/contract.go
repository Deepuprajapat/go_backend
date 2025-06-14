package repository

import "github.com/VI-IM/im_backend_go/ent"

type repository struct {
	db *ent.Client
}

type AppRepository interface {
	GetUserDetailsByUsername(username string) (*ent.User, error)
	GetProjectByID(id int) (*ent.Project, error)
}

func NewRepository(db *ent.Client) AppRepository {
	return &repository{db: db}
}
