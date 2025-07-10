package application

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/internal/utils"
	"github.com/VI-IM/im_backend_go/request"
	"github.com/VI-IM/im_backend_go/response"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/google/uuid"
)

func (c *application) GetAccessToken(username string, password string) (*response.GenerateTokenResponse, *imhttp.CustomError) {

	if username == "" || password == "" {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "No credentials provided", "No credentials provided")
	}

	user, err := c.repo.GetUserDetailsByUsername(username)
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get user details", err.Error())
	}

	if user == nil {
		return nil, imhttp.NewCustomErr(http.StatusNotFound, "User not found", "User not found")
	}

	if !utils.ComparePassword(user.Password, password) {
		return nil, imhttp.NewCustomErr(http.StatusUnauthorized, "Invalid password", "Invalid password")
	}

	// Generate access token
	accessToken, err := utils.GenerateToken(user.ID, time.Now().Add(time.Hour*24))
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to generate access token", err.Error())
	}

	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(user.ID, time.Now().Add(time.Hour*24*7))
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to generate refresh token", err.Error())
	}

	return &response.GenerateTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (c *application) RefreshToken(refreshToken string) (*response.GenerateTokenResponse, *imhttp.CustomError) {

	if refreshToken == "" {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Refresh token is empty", "Refresh token is empty")
	}

	userID, err := utils.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to verify refresh token", err.Error())
	}

	user, err := c.repo.GetUserDetailsByUsername(strconv.Itoa(userID))
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get user details", err.Error())
	}

	if user == nil {
		return nil, imhttp.NewCustomErr(http.StatusNotFound, "User not found", "User not found")
	}

	// Note: DeletedAt field will be available after regenerating entities
	// if user.DeletedAt != nil {
	// 	return nil, imhttp.NewCustomErr(http.StatusNotFound, "User is deleted", "User is deleted")
	// }

	if !user.IsActive {
		return nil, imhttp.NewCustomErr(http.StatusUnauthorized, "User is not active", "User is not active")
	}

	accessToken, err := utils.GenerateToken(strconv.Itoa(userID), time.Now().Add(time.Hour*24))
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to generate access token", err.Error())
	}

	return &response.GenerateTokenResponse{
		AccessToken: accessToken,
	}, nil
}

func (c *application) Signup(ctx context.Context, req *request.SignupRequest) (*response.GenerateTokenResponse, *imhttp.CustomError) {
	// Check if username already exists
	existingUser, err := c.repo.GetUserDetailsByUsername(req.Username)
	if err != nil && !ent.IsNotFound(err) {
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to check username", err.Error())
	}
	if existingUser != nil {
		return nil, imhttp.NewCustomErr(http.StatusConflict, "Username already exists", "Username already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to hash password", err.Error())
	}

	// Create user
	user := &ent.User{
		ID:          uuid.New().String(),
		Username:    req.Username,
		Password:    hashedPassword,
		Email:       req.Email,
		Name:        req.FirstName + " " + req.LastName,
		PhoneNumber: req.PhoneNumber,
		IsActive:    true,
	}

	createdUser, err := c.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to create user", err.Error())
	}

	// Generate tokens
	accessToken, err := utils.GenerateToken(createdUser.ID, time.Now().Add(time.Hour*24))
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to generate access token", err.Error())
	}

	refreshToken, err := utils.GenerateRefreshToken(createdUser.ID, time.Now().Add(time.Hour*24*7))
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to generate refresh token", err.Error())
	}

	return &response.GenerateTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
