package http

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

type contextKey string

const RequestIDKey contextKey = "dink_request_id"

// RequestID generates a unique request ID and attaches it to request context and header.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-Id")
		if reqID == "" {
			b := make([]byte, 16)
			_, _ = rand.Read(b)
			reqID = "req_" + hex.EncodeToString(b)
		}
		w.Header().Set("X-Request-Id", reqID)
		next.ServeHTTP(w, r)
	})
}

// Recoverer catches panics and returns an RFC 9457 internal server error.
func Recoverer(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					if logger != nil {
						logger.Error("panic recovered in HTTP handler",
							slog.Any("recover", rvr),
							slog.String("stack", string(debug.Stack())),
						)
					}
					http.Error(w, `{"success":false,"error":{"status":500,"title":"Internal Server Error","code":"INTERNAL_SERVER_ERROR"}}`, http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// RequestLogger logs requests using standard library log/slog.
func RequestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			if logger != nil {
				logger.Info("http request",
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.Duration("duration", time.Since(start)),
				)
			}
		})
	}
}
