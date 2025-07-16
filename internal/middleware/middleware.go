package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	zlog "github.com/rs/zerolog/log"
	"github.com/VI-IM/im_backend_go/internal/config"
)

// LoggingMiddleware logs request details including method, URL, status code, and duration
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		wrapper := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(wrapper, r)

		// Log the request details - only for /v1/ API paths
		if !strings.HasPrefix(r.URL.Path, "/v1/") {
			return
		}
		
		duration := time.Since(start)
		cfg := config.GetConfig()
		
		if cfg.Logger.Mode == "simple" {
			// Simple plain text format: Time Method path status code duration
			simpleLogger := log.New(os.Stdout, "", 0)
			simpleLogger.Printf("%s %s %s %d %v", 
				start.Format("15:04:05"), 
				r.Method, 
				r.URL.Path, 
				wrapper.statusCode, 
				duration.Truncate(time.Millisecond))
		} else {
			// Verbose format (original structured logging)
			zlog.Info().
				Str("method", r.Method).
				Str("url", r.URL.String()).
				Str("remote_addr", r.RemoteAddr).
				Str("user_agent", r.UserAgent()).
				Int("status_code", wrapper.statusCode).
				Dur("duration", duration).
				Msg("HTTP Request")
		}
	})
}

// responseWriter wraps http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
