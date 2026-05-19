package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

func demoGoroutines() {
	fmt.Println("\n--- Goroutine Cơ Bản ---")

	// go keyword: launch goroutine — very lightweight (~2KB initial stack, grows automatically)
	// OS thread: ~1-8MB stack — goroutines are hundreds of times cheaper
	go func() {
		fmt.Println("  Goroutine 1: running concurrently")
	}()

	// GOTCHA: main goroutine may exit before other goroutines finish
	time.Sleep(10 * time.Millisecond) // DO NOT do this in production!

	fmt.Println("\n--- sync.WaitGroup — wait for multiple goroutines ===")

	var wg sync.WaitGroup
	results := make([]int, 5)

	for i := range 5 {
		wg.Add(1) // increment counter BEFORE launching goroutine
		go func(idx int) {
			defer wg.Done() // decrement counter when done (runs via defer)
			results[idx] = idx * idx
			fmt.Printf("  goroutine %d: computed %d\n", idx, idx*idx)
		}(i) // pass i in to avoid closure capture bug
	}

	wg.Wait() // block until counter = 0
	fmt.Printf("  Results: %v\n", results)

	// Go 1.25+: sync.WaitGroup.Go — more convenient, auto Add(1) + launch
	fmt.Println("\n--- sync.WaitGroup.Go (Go 1.25+) ---")
	var wg2 sync.WaitGroup
	for i := range 3 {
		wg2.Go(func() { // auto Add(1) and launch goroutine
			fmt.Printf("  wg.Go goroutine %d\n", i)
		})
	}
	wg2.Wait()

	fmt.Println("\n--- Goroutine Runtime Info ---")
	fmt.Printf("  NumCPU: %d\n", runtime.NumCPU())
	fmt.Printf("  GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0)) // 0 = query, don't set
	fmt.Printf("  NumGoroutine: %d\n", runtime.NumGoroutine())

	fmt.Println("\n--- Goroutine Leak Prevention ---")
	demoGoroutineLeakFix()
}

// GOTCHA: Goroutine leak — goroutine blocked forever
// BAD example (illustrated with comment only):
//   ch := make(chan int)
//   go func() {
//       val := <-ch  // blocks forever if nobody sends
//       process(val)
//   }()
//   // if nothing is sent to ch, goroutine is leaked!

// GOOD: Use context to cancel goroutine
func demoGoroutineLeakFix() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel() // ensure cancel is called

	ch := make(chan int, 1)

	go func() {
		select {
		case ch <- 42:
			fmt.Println("  goroutine: sent value")
		case <-ctx.Done():
			fmt.Printf("  goroutine: cancelled (%v)\n", ctx.Err())
		}
	}()

	select {
	case v := <-ch:
		fmt.Printf("  main: received %d\n", v)
	case <-ctx.Done():
		fmt.Printf("  main: timeout (%v)\n", ctx.Err())
	}
}
