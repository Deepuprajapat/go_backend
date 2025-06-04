package repository

import (
	"context"

	"github.com/VI-IM/im_backend_go/ent"
)

type repository struct {
	db *ent.Client
}

type AppRepository interface {
	GetUserDetailsByUsername(username string) (*ent.User, error)
	CreateUser(ctx context.Context, input *ent.User) (*ent.User, error)
	BlacklistToken(ctx context.Context, token string, userID int) error
}

func NewRepository(db *ent.Client) AppRepository {
	return &repository{db: db}
}
