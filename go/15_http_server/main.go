// Bài 15: HTTP Server — xây dựng web server với Go 1.22+ mux
// Chạy: go run .
// Test: curl http://localhost:8080/users
//       curl http://localhost:8080/users/1
//       curl -X POST -H "Content-Type: application/json" -d '{"name":"Alice","email":"alice@example.com"}' http://localhost:8080/users
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Khởi tạo in-memory store
	store := NewUserStore()
	store.Create(User{Name: "Alice", Email: "alice@example.com", Age: 30})
	store.Create(User{Name: "Bob", Email: "bob@example.com", Age: 25})

	// Tạo handler
	h := NewUserHandler(store)

	// Go 1.22+: mux hỗ trợ method + path pattern
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users", h.ListUsers)
	mux.HandleFunc("POST /users", h.CreateUser)
	mux.HandleFunc("GET /users/{id}", h.GetUser)       // {id} = path value
	mux.HandleFunc("PUT /users/{id}", h.UpdateUser)
	mux.HandleFunc("DELETE /users/{id}", h.DeleteUser)
	mux.HandleFunc("GET /health", healthHandler)

	// Middleware chain
	handler := Chain(mux,
		RecoveryMiddleware,
		LoggingMiddleware,
		CORSMiddleware,
	)

	// Server với proper timeouts — QUAN TRỌNG cho production
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,   // thời gian đọc request
		WriteTimeout: 10 * time.Second,  // thời gian gửi response
		IdleTimeout:  120 * time.Second, // keep-alive connections
	}

	// Graceful shutdown
	go func() {
		slog.Info("server starting", "addr", srv.Addr)
		slog.Info("try: curl http://localhost:8080/users")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	// Chờ signal (Ctrl+C hoặc SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	sig := <-quit
	slog.Info("shutdown signal received", "signal", sig)

	// Graceful shutdown: cho phép requests đang xử lý hoàn thành
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("forced shutdown", "err", err)
	} else {
		slog.Info("server stopped gracefully")
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"status":"ok","time":"`+time.Now().Format(time.RFC3339)+`"}`)
}
