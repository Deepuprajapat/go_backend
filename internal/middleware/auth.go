package middleware

import (
	"context"
	"github.com/VI-IM/im_backend_go/internal/auth"
	"net/http"
	"strings"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := authHeader
		if strings.Contains(authHeader, "Bearer") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add claims to request context for use in handlers
		ctx := context.WithValue(r.Context(), "user_claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireBusinessPartner ensures only business_partner or superadmin can access
func RequireBusinessPartner(next http.Handler) http.Handler {
	return Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("user_claims").(*auth.Claims)
		if !ok {
			http.Error(w, "Invalid user context", http.StatusUnauthorized)
			return
		}

		if claims.Role != "business_partner" && claims.Role != "superadmin" {
			http.Error(w, "Access denied: requires business_partner or superadmin role", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}))
}

// RequireSuperAdmin ensures only superadmin can access
func RequireSuperAdmin(next http.Handler) http.Handler {
	return Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("user_claims").(*auth.Claims)
		if !ok {
			http.Error(w, "Invalid user context", http.StatusUnauthorized)
			return
		}

		if claims.Role != "superadmin" {
			http.Error(w, "Access denied: requires superadmin role", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}))
}
