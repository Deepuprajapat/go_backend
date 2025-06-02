package repository

import (
	"context"
	"errors"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/user"
)

func (r *repository) GetUserDetailsByUsername(username string) (*ent.User, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}

	user, err := r.db.User.Query().Where(user.Username(username)).First(context.Background())
	if err != nil {
		return nil, err
	}

	return user, nil
}
