// Bài 12: Context — quản lý cancellation, deadline, và values
// Chạy: go run .
package main

import (
	"context"
	"fmt"
	"time"
)

// contextKey là kiểu riêng để tránh key collision giữa packages
// NGUYÊN TẮC: Không bao giờ dùng built-in types (string, int) làm context key
type contextKey string

const (
	userIDKey    contextKey = "userID"
	requestIDKey contextKey = "requestID"
	traceIDKey   contextKey = "traceID"
)

// === Simulated service calls ===

func fetchUserFromDB(ctx context.Context, userID int) (string, error) {
	// Simulate DB query
	select {
	case <-time.After(50 * time.Millisecond):
		return fmt.Sprintf("User%d", userID), nil
	case <-ctx.Done():
		return "", fmt.Errorf("fetchUserFromDB: %w", ctx.Err())
	}
}

func callExternalAPI(ctx context.Context, endpoint string) (string, error) {
	select {
	case <-time.After(200 * time.Millisecond): // slow API
		return "API response", nil
	case <-ctx.Done():
		return "", fmt.Errorf("callExternalAPI %s: %w", endpoint, ctx.Err())
	}
}

// handler mô phỏng HTTP request handler
func handler(ctx context.Context, userID int) error {
	// Extract values từ context
	reqID, _ := ctx.Value(requestIDKey).(string)
	fmt.Printf("  [%s] Handling request for userID=%d\n", reqID, userID)

	// Truyền context xuống các layers
	user, err := fetchUserFromDB(ctx, userID)
	if err != nil {
		return fmt.Errorf("handler: %w", err)
	}

	result, err := callExternalAPI(ctx, "/profile")
	if err != nil {
		return fmt.Errorf("handler: %w", err)
	}

	fmt.Printf("  [%s] Done: user=%s, api=%s\n", reqID, user, result)
	return nil
}

func main() {
	fmt.Println("=== 1. context.Background() và context.TODO() ===")
	// Background: root context, không bao giờ cancel — dùng trong main, init, tests
	bg := context.Background()
	// TODO: placeholder khi chưa biết dùng context nào — để sau refactor
	todo := context.TODO()
	fmt.Printf("  Background: %v\n", bg)
	fmt.Printf("  TODO: %v\n", todo)

	fmt.Println("\n=== 2. WithCancel — manual cancellation ===")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // NGUYÊN TẮC: luôn defer cancel() để tránh goroutine leak

	go func() {
		select {
		case <-ctx.Done():
			fmt.Printf("  goroutine cancelled: %v\n", ctx.Err())
		case <-time.After(1 * time.Second):
			fmt.Println("  goroutine done normally")
		}
	}()

	time.Sleep(50 * time.Millisecond)
	cancel() // signal goroutine to stop
	time.Sleep(20 * time.Millisecond)

	fmt.Println("\n=== 3. WithTimeout — auto cancel sau thời gian ===")
	// Timeout: hủy sau duration
	ctx2, cancel2 := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel2()

	err := handler(ctx2, 42)
	if err != nil {
		fmt.Printf("  handler error: %v\n", err)
	}

	fmt.Println("\n=== 4. WithTimeout — bị timeout ===")
	// Context timeout ngắn hơn API call time
	ctx3, cancel3 := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel3()

	// Thêm value vào context trước khi truyền
	ctx3 = context.WithValue(ctx3, requestIDKey, "req-001")
	err = handler(ctx3, 1)
	if err != nil {
		fmt.Printf("  handler timed out: %v\n", err)
	}

	fmt.Println("\n=== 5. WithDeadline — cancel tại thời điểm cụ thể ===")
	deadline := time.Now().Add(200 * time.Millisecond)
	ctx4, cancel4 := context.WithDeadline(context.Background(), deadline)
	defer cancel4()

	fmt.Printf("  Deadline: %v (in ~200ms)\n", deadline.Format("15:04:05.000"))
	select {
	case <-ctx4.Done():
		fmt.Printf("  Context done: %v\n", ctx4.Err())
	case <-time.After(300 * time.Millisecond):
		fmt.Println("  Timer expired first")
	}

	fmt.Println("\n=== 6. WithValue — truyền values qua context chain ===")
	ctx5 := context.WithValue(context.Background(), userIDKey, 123)
	ctx5 = context.WithValue(ctx5, requestIDKey, "req-abc-123")
	ctx5 = context.WithValue(ctx5, traceIDKey, "trace-xyz-789")

	// Extract values — phải type assert
	if uid, ok := ctx5.Value(userIDKey).(int); ok {
		fmt.Printf("  userID: %d\n", uid)
	}
	if reqID, ok := ctx5.Value(requestIDKey).(string); ok {
		fmt.Printf("  requestID: %s\n", reqID)
	}
	if traceID, ok := ctx5.Value(traceIDKey).(string); ok {
		fmt.Printf("  traceID: %s\n", traceID)
	}

	// Value không tồn tại → nil
	missing := ctx5.Value("nonexistent")
	fmt.Printf("  missing key: %v\n", missing)

	fmt.Println("\n=== 7. Errors của Context ===")
	fmt.Printf("  context.Canceled: %v\n", context.Canceled)
	fmt.Printf("  context.DeadlineExceeded: %v\n", context.DeadlineExceeded)

	cancelledCtx, c := context.WithCancel(context.Background())
	c()
	fmt.Printf("  Is Canceled: %t\n", cancelledCtx.Err() == context.Canceled)

	fmt.Println("\n=== 8. Context Best Practices ===")
	fmt.Println("  ✓ Luôn defer cancel() ngay sau WithCancel/WithTimeout/WithDeadline")
	fmt.Println("  ✓ Truyền context làm argument đầu tiên (ctx context.Context)")
	fmt.Println("  ✓ Không store context trong struct — truyền qua function call")
	fmt.Println("  ✓ Dùng typed key (type contextKey string) cho context values")
	fmt.Println("  ✓ Context values: request-scoped data (userID, traceID, authToken)")
	fmt.Println("  ✓ Không dùng context để truyền optional parameters")
	fmt.Println("  ✗ Không return context từ function")
	fmt.Println("  ✗ Không dùng nil context — dùng context.Background()")
}
