package utils

import (
	"errors"
	"time"

	"github.com/VI-IM/im_backend_go/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID int, expiresAt time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": expiresAt.Unix(),
	})
	secret := config.GetConfig().JWTConfig.AuthSecret
	if secret == "" {
		return "", errors.New("JWT_SECRET is not set")
	}
	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(userID int, expiresAt time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":        userID,
		"token_type": "refresh",
		"exp":        expiresAt.Unix(),
	})

	secret := config.GetConfig().JWTConfig.AuthSecret
	if secret == "" {
		return "", errors.New("JWT_SECRET is not set")
	}
	return token.SignedString([]byte(secret))
}

func VerifyToken(tokenString string) (userID int, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig().JWTConfig.AuthSecret), nil
	})
	if err != nil {
		return 0, err
	}
	return int(token.Claims.(jwt.MapClaims)["user_id"].(float64)), nil
}

func VerifyRefreshToken(tokenString string) (userID int, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig().JWTConfig.AuthSecret), nil
	})
	if err != nil {
		return 0, err
	}

	if token.Claims.(jwt.MapClaims)["token_type"] != "refresh" {
		return 0, errors.New("invalid token type")
	}

	return int(token.Claims.(jwt.MapClaims)["user_id"].(float64)), nil
}
