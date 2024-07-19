package main

import (
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type midFunc func(http.Handler) http.Handler

func RateLimit(limit, burst int) midFunc {
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

func LogRequest(logger *slog.Logger) midFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()
			h.ServeHTTP(w, r)
			logger.Info("Request", "et", time.Since(now), "[origin]", r.RemoteAddr, "[path]", r.URL.Path)
		})
	}
}
