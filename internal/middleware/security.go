package middleware

import (
	"net/http"

	"github.com/VI-IM/im_backend_go/internal/config"
	"github.com/gorilla/mux"
)

// SecurityHeadersMiddleware adds security headers to HTTPS responses
func SecurityHeadersMiddleware(cfg config.Config) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only add HSTS header if TLS is enabled
			if cfg.TLS.Enabled {
				// HTTP Strict Transport Security (HSTS) - 1 year
				w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			}

			// Prevent MIME type sniffing
			w.Header().Set("X-Content-Type-Options", "nosniff")

			// Clickjacking protection
			w.Header().Set("X-Frame-Options", "SAMEORIGIN")

			// XSS protection for older browsers
			w.Header().Set("X-XSS-Protection", "1; mode=block")

			// Content Security Policy - basic policy, can be customized based on needs
			w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self' https:;")

			// Referrer Policy
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

			// Permissions Policy (formerly Feature Policy)
			w.Header().Set("Permissions-Policy", "accelerometer=(), camera=(), geolocation=(), gyroscope=(), magnetometer=(), microphone=(), payment=(), usb=()")

			next.ServeHTTP(w, r)
		})
	}
}