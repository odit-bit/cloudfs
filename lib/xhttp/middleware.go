package xhttp

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/cors"
	"golang.org/x/time/rate"
)

func RateLimit(limit, burst int) func(http.Handler) http.Handler {
	var limiter *rate.Limiter
	useLimiter := false
	if limit > 0 {
		limiter = rate.NewLimiter(rate.Limit(limit), burst)
		useLimiter = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if useLimiter {
				if !limiter.Allow() {
					http.Error(w, "limit exceed", http.StatusTooManyRequests)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func LogRequest(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()
			h.ServeHTTP(w, r)
			logger.Info("Request", "et", time.Since(now), "[origin]", r.RemoteAddr, "[path]", r.URL.Path)
		})
	}
}

func CorsDefault() func(http.Handler) http.Handler {
	cors := cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	return func(h http.Handler) http.Handler {
		return cors(h)
	}
}
