package repository

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/user"
)

func (r *repository) GetUserDetailsByEmail(ctx context.Context, email string) (*ent.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	user, err := r.db.User.Query().Where(user.Email(email)).First(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repository) CreateUser(ctx context.Context, input *ent.User) (*ent.User, error) {
	if input == nil {
		return nil, errors.New("user input is required")
	}

	id := fmt.Sprintf("%x", sha256.Sum256([]byte(input.Email)))[:16]
	// Create the user
	createdUser, err := r.db.User.Create().
		SetID(id).
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

func (r *repository) CheckIfUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	exists, err := r.db.User.Query().Where(user.Email(email)).Exist(ctx)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *repository) CheckIfUserExistsByID(ctx context.Context, userID string) (bool, error) {
	if userID == "" {
		return false, errors.New("user ID is required")
	}
	exists, err := r.db.User.Query().Where(user.ID(userID)).Exist(ctx)
	if err != nil {
		return false, err
	}
	return exists, nil
}
