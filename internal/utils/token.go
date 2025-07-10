package utils

import (
	"errors"
	"time"

	"github.com/VI-IM/im_backend_go/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID string, expiresAt time.Time) (string, error) {
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

func GenerateRefreshToken(userID string, expiresAt time.Time) (string, error) {
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
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}
	val, ok := claims["user_id"]
	if !ok || val == nil {
		return 0, errors.New("user_id not found in token")
	}
	userIDFloat, ok := val.(float64)
	if !ok {
		return 0, errors.New("user_id in token is not a float64")
	}
	return int(userIDFloat), nil
}

func VerifyRefreshToken(tokenString string) (userID int, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig().JWTConfig.AuthSecret), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}
	if claims["token_type"] != "refresh" {
		return 0, errors.New("invalid token type")
	}
	val, ok := claims["user_id"]
	if !ok || val == nil {
		return 0, errors.New("user_id not found in token")
	}
	userIDFloat, ok := val.(float64)
	if !ok {
		return 0, errors.New("user_id in token is not a float64")
	}
	return int(userIDFloat), nil
}
