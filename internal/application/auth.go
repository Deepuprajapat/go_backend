package application

import (
	"net/http"
	"strconv"
	"time"

	"github.com/VI-IM/im_backend_go/internal/utils"
	"github.com/VI-IM/im_backend_go/response"
	imhttp "github.com/VI-IM/im_backend_go/shared"
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

	if err := utils.ComparePassword(user.Password, password); err != nil {
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

	if user.IsDeleted {
		return nil, imhttp.NewCustomErr(http.StatusNotFound, "User is deleted", "User is deleted")
	}

	if !user.IsActive {
		return nil, imhttp.NewCustomErr(http.StatusUnauthorized, "User is not active", "User is not active")
	}

	accessToken, err := utils.GenerateToken(userID, time.Now().Add(time.Hour*24))
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to generate access token", err.Error())
	}

	return &response.GenerateTokenResponse{
		AccessToken: accessToken,
	}, nil
}
