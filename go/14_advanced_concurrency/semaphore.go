package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Semaphore = buffered channel — limits the number of concurrent operations
// Use when: limiting concurrent DB connections, API calls, file I/O

type Semaphore chan struct{}

func NewSemaphore(n int) Semaphore {
	return make(Semaphore, n)
}

func (s Semaphore) Acquire() {
	s <- struct{}{} // block if full
}

func (s Semaphore) Release() {
	<-s
}

func demoSemaphore() {
	fmt.Println("\n--- Semaphore (buffered channel) ---")
	sem := NewSemaphore(3) // at most 3 concurrent operations
	var wg sync.WaitGroup
	var active atomic.Int32

	for i := range 10 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sem.Acquire()
			defer sem.Release()

			n := active.Add(1)
			fmt.Printf("  Task %d started (concurrent: %d)\n", id, n)
			time.Sleep(50 * time.Millisecond)
			active.Add(-1)
		}(i)
	}

	wg.Wait()
	fmt.Println("  All tasks done, max concurrency was 3")
}

func demoRateLimiter() {
	fmt.Println("\n--- Rate Limiter (token bucket) ---")
	// Token bucket: ticker fills bucket at fixed rate
	// Each request consumes 1 token

	rate := time.NewTicker(100 * time.Millisecond) // 10 req/sec
	defer rate.Stop()

	for i := range 5 {
		<-rate.C // wait for token
		fmt.Printf("  Request %d processed at %v\n", i+1,
			time.Now().Format("15:04:05.000"))
	}

	fmt.Println("\n  Dùng golang.org/x/time/rate cho production rate limiter:")
	fmt.Println("  limiter := rate.NewLimiter(rate.Every(time.Second/10), 5)")
	fmt.Println("  if err := limiter.Wait(ctx); err != nil { ... }")
}
