package application

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/internal/utils"
	"github.com/VI-IM/im_backend_go/request"
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

func (a *Application) GetAccessToken(username string, password string) (*response.GenerateTokenResponse, error) {

	if username == "" || password == "" {
		return nil, ErrNoCredentials
	}

	user, err := a.repo.GetUserDetailsByUsername(username)
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

func (a *Application) RefreshToken(refreshToken string) (*response.GenerateTokenResponse, error) {

	if refreshToken == "" {
		return nil, ErrRefreshTokenEmpty
	}

	userID, err := utils.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	user, err := a.repo.GetUserDetailsByUsername(strconv.Itoa(userID))
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

	newRefreshToken, err := utils.GenerateRefreshToken(userID, time.Now().Add(time.Hour*24*7))
	if err != nil {
		return nil, err
	}

	return &response.GenerateTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (a *Application) Signup(req *request.SignupRequest) (*response.SignupResponse, error) {
	// Validate input
	if req.Username == "" || req.Password == "" || req.Email == "" {
		return nil, errors.New("username, password and email are required")
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Split full name into first and last name
	names := strings.Split(req.FullName, " ")
	firstName := names[0]
	lastName := ""
	if len(names) > 1 {
		lastName = strings.Join(names[1:], " ")
	}

	// Create user input
	userInput := &ent.User{
		Username:    req.Username,
		Password:    hashedPassword,
		Email:       req.Email,
		FirstName:   firstName,
		LastName:    lastName,
		PhoneNumber: req.Phone,
	}

	// Create user in database
	user, err := a.repo.CreateUser(context.Background(), userInput)
	if err != nil {
		return nil, err
	}

	// Return response
	return &response.SignupResponse{
		ID:        strconv.Itoa(user.ID),
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FirstName + " " + user.LastName,
		Phone:     user.PhoneNumber,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}, nil
}
