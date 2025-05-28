package auth

import (
	"errors"
	"time"
	"github.com/VI-IM/im_backend_go/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID  int    `json:"user_id"`
	IsAdmin bool   `json:"is_admin"`
	Phone   string `json:"phone"`
	jwt.RegisteredClaims
}

func GenerateToken(userID int, isAdmin bool, phone string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(config.DefaultConfig.JWTExpirationHours) * time.Hour)

	claims := &Claims{
		UserID:  userID,
		IsAdmin: isAdmin,
		Phone:   phone,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.DefaultConfig.JWTSecret))
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.DefaultConfig.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
