package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/VI-IM/im_backend_go/internal/auth"
	"github.com/VI-IM/im_backend_go/internal/config"
	"github.com/gorilla/mux"
)

type contextKey string

const UserContextKey contextKey = "user"

// AuthMiddleware creates a middleware that validates JWT tokens from cookies
func AuthMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(config.DefaultConfig.JWTCookieName)
			if err != nil {
				// No cookie found, continue without authentication
				next.ServeHTTP(w, r)
				return
			}

			claims, err := auth.ValidateToken(cookie.Value)
			if err != nil {
				// Invalid token, continue without authentication
				next.ServeHTTP(w, r)
				return
			}

			// Add user info to context
			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAuth creates a middleware that requires valid authentication
func RequireAuth() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(UserContextKey).(*auth.Claims)
			if !ok || claims == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// SetAuthCookie sets the JWT token as a cookie
func SetAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     config.DefaultConfig.JWTCookieName,
		Value:    token,
		Path:     "/",
		Domain:   config.DefaultConfig.JWTCookieDomain,
		Expires:  time.Now().Add(time.Duration(config.DefaultConfig.JWTExpirationHours) * time.Hour),
		Secure:   config.DefaultConfig.JWTCookieSecure,
		HttpOnly: config.DefaultConfig.JWTCookieHTTPOnly,
		SameSite: http.SameSiteStrictMode,
	})
}

// ClearAuthCookie removes the auth cookie
func ClearAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     config.DefaultConfig.JWTCookieName,
		Value:    "",
		Path:     "/",
		Domain:   config.DefaultConfig.JWTCookieDomain,
		Expires:  time.Now().Add(-time.Hour),
		Secure:   config.DefaultConfig.JWTCookieSecure,
		HttpOnly: config.DefaultConfig.JWTCookieHTTPOnly,
		SameSite: http.SameSiteStrictMode,
	})
}
