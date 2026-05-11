// Bài 17: Logging với slog — structured logging trong Go 1.21+
// Chạy: go run .
package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"
)

// === Custom Handler — gửi logs tới nhiều destinations ===

type MultiHandler struct {
	handlers []slog.Handler
}

func NewMultiHandler(handlers ...slog.Handler) *MultiHandler {
	return &MultiHandler{handlers: handlers}
}

func (h *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, r.Level) {
			if err := handler.Handle(ctx, r); err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, h := range h.handlers {
		handlers[i] = h.WithAttrs(attrs)
	}
	return &MultiHandler{handlers: handlers}
}

func (h *MultiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, h := range h.handlers {
		handlers[i] = h.WithGroup(name)
	}
	return &MultiHandler{handlers: handlers}
}

func main() {
	fmt.Println("=== 1. Default Text Logger ===")
	// slog.Default() dùng TextHandler ghi vào os.Stderr
	slog.Info("application starting", "version", "1.0.0", "env", "development")
	slog.Warn("disk usage high", "percent", 85, "path", "/var/log")
	slog.Error("connection failed", "host", "db.example.com", "port", 5432)

	fmt.Println("\n=== 2. Custom JSON Handler (production) ===")
	jsonLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug, // log semua levels
		AddSource: true,            // thêm file:line vào log
	}))
	jsonLogger.Debug("debug message", "key", "value")
	jsonLogger.Info("user created", "userID", 42, "email", "alice@example.com")

	fmt.Println("\n=== 3. Text Handler (development) ===")
	textLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	textLogger.Info("server started", "addr", ":8080", "tls", false)

	fmt.Println("\n=== 4. slog.SetDefault ===")
	// Đặt global default logger — toàn bộ slog.* calls dùng logger này
	slog.SetDefault(textLogger)
	slog.Info("now using text logger globally")

	fmt.Println("\n=== 5. Structured Attributes ===")
	// Các Attr helper functions — type-safe hơn interface{}
	jsonLogger.Info("order placed",
		slog.Int("orderID", 12345),
		slog.String("customerID", "C001"),
		slog.Float64("total", 99.99),
		slog.Bool("isPremium", true),
		slog.Duration("processingTime", 123*time.Millisecond),
		slog.Time("timestamp", time.Now()),
		slog.Any("tags", []string{"electronics", "sale"}),
	)

	fmt.Println("\n=== 6. logger.With() — persistent attributes ===")
	// With(): tạo logger mới với attrs thêm vào mọi log record
	requestLogger := jsonLogger.With(
		slog.String("requestID", "req-abc-123"),
		slog.String("userID", "user-456"),
	)
	requestLogger.Info("request received", "method", "GET", "path", "/api/users")
	requestLogger.Info("query executed", "table", "users", "rows", 42)
	requestLogger.Warn("slow query", "duration", 2500*time.Millisecond)

	fmt.Println("\n=== 7. slog.Group — nhóm attributes ===")
	jsonLogger.Info("database connected",
		slog.Group("db",
			slog.String("host", "localhost"),
			slog.Int("port", 5432),
			slog.String("name", "myapp"),
		),
		slog.Group("pool",
			slog.Int("maxConns", 25),
			slog.Int("minConns", 5),
		),
	)

	fmt.Println("\n=== 8. Context-aware logging ===")
	type reqKey struct{}
	ctx := context.WithValue(context.Background(), reqKey{}, "req-789")

	// logger.InfoContext: dùng context cho potential log filtering
	jsonLogger.InfoContext(ctx, "processing request",
		slog.String("handler", "getUserProfile"),
	)

	fmt.Println("\n=== 9. Log Levels ===")
	// Levels: Debug(-4) < Info(0) < Warn(4) < Error(8)
	// Custom level:
	const LevelTrace = slog.Level(-8)
	const LevelFatal = slog.Level(12)

	for _, level := range []slog.Level{LevelTrace, slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError, LevelFatal} {
		fmt.Printf("  Level %d = %v\n", int(level), level)
	}

	fmt.Println("\n=== 10. Multi-destination Logger ===")
	// Gửi logs tới cả stdout (JSON) và stderr (text) cùng lúc
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	textHandler := slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError}) // chỉ errors tới stderr
	multiLogger := slog.New(NewMultiHandler(jsonHandler, textHandler))
	multiLogger.Info("multi-destination log") // chỉ json handler in
	multiLogger.Error("error to both destinations") // cả 2 handlers in

	fmt.Println("\n=== Best Practices ===")
	fmt.Println("  - Dùng slog thay vì log.Printf (type-safe, structured)")
	fmt.Println("  - JSON handler cho production (machine-readable)")
	fmt.Println("  - Text handler cho development (human-readable)")
	fmt.Println("  - logger.With() cho request-scoped fields (requestID, userID)")
	fmt.Println("  - Không log sensitive data (passwords, tokens, PII)")
	fmt.Println("  - Dùng slog.Group() để nhóm related fields")
}
