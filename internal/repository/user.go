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

func (r *repository) CreateUser(ctx context.Context, input *ent.User) (*ent.User, error) {
	if input == nil {
		return nil, errors.New("user input is required")
	}

	// Check if username already exists
	exists, err := r.db.User.Query().Where(user.Username(input.Username)).Exist(ctx)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	exists, err = r.db.User.Query().Where(user.Email(input.Email)).Exist(ctx)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// Create the user
	createdUser, err := r.db.User.Create().
		SetUsername(input.Username).
		SetPassword(input.Password).
		SetEmail(input.Email).
		SetName(input.Name).
		SetPhoneNumber(input.PhoneNumber).
		SetIsActive(true).
		Save(ctx)

	if err != nil {
		return nil, err
	}

	return createdUser, nil
}
