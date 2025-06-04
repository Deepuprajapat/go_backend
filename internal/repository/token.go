package repository

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (r *repository) BlacklistToken(ctx context.Context, token string, userID int) error {
	// Parse the token to get its expiration time
	parsedToken, _, err := jwt.NewParser().ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return jwt.ErrInvalidKey
	}

	// Get expiration time from token
	exp, ok := claims["exp"].(float64)
	if !ok {
		return jwt.ErrInvalidKey
	}

	expiresAt := time.Unix(int64(exp), 0)

	// Create blacklisted token entry
	_, err = r.db.BlacklistedToken.Create().
		SetToken(token).
		SetUserID(userID).
		SetExpiresAt(expiresAt).
		Save(ctx)
	return err
}
