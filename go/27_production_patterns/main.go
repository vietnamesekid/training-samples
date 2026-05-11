// Bài 27: Production Patterns trong Go
// Circuit Breaker, Retry+Backoff, Graceful Shutdown, Health Checks, Config
// Chạy: go run .
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("=== Production Patterns ===")

	fmt.Println("\n=== 1. Circuit Breaker ===")
	demoCircuitBreaker()

	fmt.Println("\n=== 2. Retry with Exponential Backoff ===")
	demoRetry()

	fmt.Println("\n=== 3. Config from Environment ===")
	demoConfig()

	fmt.Println("\n=== 4. Health Check Server ===")
	demoHealthCheck()

	fmt.Println("\n=== 5. Graceful Shutdown ===")
	demoGracefulShutdown()
}

// demoGracefulShutdown minh họa graceful shutdown pattern với signal handling
func demoGracefulShutdown() {
	// Pattern graceful shutdown:
	// 1. Listen SIGINT/SIGTERM
	// 2. Stop accepting new requests
	// 3. Wait for in-flight requests to complete
	// 4. Close resources (DB, cache, etc.)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("  Server running... (send SIGINT to stop)")
	fmt.Println("  Pattern:")
	fmt.Println("  quit := make(chan os.Signal, 1)")
	fmt.Println("  signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)")
	fmt.Println("  <-quit")
	fmt.Println("  server.Shutdown(ctx)")

	// Trong demo này, không block — chỉ show pattern
	signal.Stop(quit)
	fmt.Println("  Graceful shutdown complete")
}
