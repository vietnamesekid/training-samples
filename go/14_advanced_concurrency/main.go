// Lesson 14: Advanced Concurrency Patterns
// Run: go run .
// go get golang.org/x/sync (if not already installed)
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
