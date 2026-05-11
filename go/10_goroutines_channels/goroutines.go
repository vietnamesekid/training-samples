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

	// go keyword: launch goroutine — rất nhẹ (~2KB stack ban đầu, tự grow)
	// OS thread: ~1-8MB stack — goroutine rẻ hơn hàng trăm lần
	go func() {
		fmt.Println("  Goroutine 1: running concurrently")
	}()

	// GOTCHA: main goroutine có thể exit trước khi goroutine khác chạy xong
	time.Sleep(10 * time.Millisecond) // KHÔNG làm thế này trong production!

	fmt.Println("\n--- sync.WaitGroup — chờ nhiều goroutines ===")

	var wg sync.WaitGroup
	results := make([]int, 5)

	for i := range 5 {
		wg.Add(1) // tăng counter TRƯỚC khi launch goroutine
		go func(idx int) {
			defer wg.Done() // giảm counter khi xong (chạy qua defer)
			results[idx] = idx * idx
			fmt.Printf("  goroutine %d: computed %d\n", idx, idx*idx)
		}(i) // truyền i vào để tránh closure capture bug
	}

	wg.Wait() // block cho đến khi counter = 0
	fmt.Printf("  Results: %v\n", results)

	// Go 1.25+: sync.WaitGroup.Go — tiện hơn, tự Add(1) + launch
	fmt.Println("\n--- sync.WaitGroup.Go (Go 1.25+) ---")
	var wg2 sync.WaitGroup
	for i := range 3 {
		wg2.Go(func() { // tự Add(1) và launch goroutine
			fmt.Printf("  wg.Go goroutine %d\n", i)
		})
	}
	wg2.Wait()

	fmt.Println("\n--- Goroutine Runtime Info ---")
	fmt.Printf("  NumCPU: %d\n", runtime.NumCPU())
	fmt.Printf("  GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0)) // 0 = query, không set
	fmt.Printf("  NumGoroutine: %d\n", runtime.NumGoroutine())

	fmt.Println("\n--- Goroutine Leak Prevention ---")
	demoGoroutineLeakFix()
}

// GOTCHA: Goroutine leak — goroutine blocked mãi mãi
// BAD example (chỉ minh họa bằng comment):
//   ch := make(chan int)
//   go func() {
//       val := <-ch  // block mãi mãi nếu không ai gửi
//       process(val)
//   }()
//   // nếu không gửi vào ch, goroutine bị leak!

// GOOD: Dùng context để cancel goroutine
func demoGoroutineLeakFix() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel() // đảm bảo cancel được gọi

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
