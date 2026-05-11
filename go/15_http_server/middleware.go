package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// Middleware type — wrap http.Handler
type Middleware func(http.Handler) http.Handler

// Chain áp dụng middlewares theo thứ tự: first → last wraps handler
// Request đi theo thứ tự: m1 → m2 → m3 → handler → m3 → m2 → m1
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	// Apply reverse order so first middleware is outermost
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// responseWriter wraps http.ResponseWriter để capture status code
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, status: http.StatusOK}
}

func (rw *responseWriter) WriteHeader(status int) {
	if !rw.wroteHeader {
		rw.status = status
		rw.wroteHeader = true
		rw.ResponseWriter.WriteHeader(status)
	}
}

// LoggingMiddleware log mỗi request với method, path, status, duration
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := newResponseWriter(w)

		next.ServeHTTP(rw, r)

		slog.Info("http request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.status,
			"duration", time.Since(start),
			"remote", r.RemoteAddr,
		)
	})
}

// RecoveryMiddleware bắt panic, trả về 500 thay vì crash server
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("handler panic",
					"panic", err,
					"method", r.Method,
					"path", r.URL.Path,
				)
				writeError(w, http.StatusInternalServerError, "internal server error",
					fmt.Sprintf("panic: %v", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// CORSMiddleware thêm CORS headers cho cross-origin requests
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
