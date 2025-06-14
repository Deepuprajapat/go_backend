package application

import (
	"errors"
	"strconv"
	"time"

	"github.com/VI-IM/im_backend_go/internal/utils"
	"github.com/VI-IM/im_backend_go/response"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserDeleted       = errors.New("user is deleted")
	ErrInvalidToken      = errors.New("invalid token")
	ErrUserNotActive     = errors.New("user is not active")
	ErrRefreshTokenEmpty = errors.New("refresh token is empty")
	ErrNoCredentials     = errors.New("no credentials provided")
)

func (c *application) GetAccessToken(username string, password string) (*response.GenerateTokenResponse, error) {

	if username == "" || password == "" {
		return nil, ErrNoCredentials
	}

	user, err := c.repo.GetUserDetailsByUsername(username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	if err := utils.ComparePassword(user.Password, password); err != nil {
		return nil, err
	}

	// Generate access token
	accessToken, err := utils.GenerateToken(user.ID, time.Now().Add(time.Hour*24))
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(user.ID, time.Now().Add(time.Hour*24*7))
	if err != nil {
		return nil, err
	}

	return &response.GenerateTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (c *application) RefreshToken(refreshToken string) (*response.GenerateTokenResponse, error) {

	if refreshToken == "" {
		return nil, ErrRefreshTokenEmpty
	}

	userID, err := utils.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	user, err := c.repo.GetUserDetailsByUsername(strconv.Itoa(userID))
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	if user.IsDeleted {
		return nil, ErrUserDeleted
	}

	if !user.IsActive {
		return nil, ErrUserNotActive
	}

	accessToken, err := utils.GenerateToken(userID, time.Now().Add(time.Hour*24))
	if err != nil {
		return nil, err
	}

	return &response.GenerateTokenResponse{
		AccessToken: accessToken,
	}, nil
}
