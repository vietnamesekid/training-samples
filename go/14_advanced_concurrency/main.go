// Bài 14: Advanced Concurrency Patterns
// Chạy: go run .
// go get golang.org/x/sync (nếu chưa có)
package main

import "fmt"

func main() {
	fmt.Println("=== WORKER POOL ===")
	demoWorkerPool()

	fmt.Println("\n=== ERRGROUP ===")
	demoErrGroup()

	fmt.Println("\n=== SEMAPHORE ===")
	demoSemaphore()

	fmt.Println("\n=== FAN-OUT / FAN-IN ===")
	demoFanOutFanIn()

	fmt.Println("\n=== RATE LIMITER ===")
	demoRateLimiter()
}
