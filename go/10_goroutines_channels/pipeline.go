package main

import (
	"context"
	"fmt"
)

// Pipeline pattern: a chain of goroutines processing data in stages
// generate → square → print
// Each stage receives from an upstream channel, processes, and sends to downstream

// Stage 1: generate numbers
func generate(ctx context.Context, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n:
			case <-ctx.Done(): // stop if context is cancelled
				return
			}
		}
	}()
	return out
}

// Stage 2: square each number
func square(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// Stage 3: add 1 to each number
func addOne(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n + 1:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func demoPipeline() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("\n--- Simple Pipeline: generate → square → addOne ---")
	nums := generate(ctx, 1, 2, 3, 4, 5)
	squares := square(ctx, nums)
	results := addOne(ctx, squares)

	for v := range results {
		fmt.Printf("  result: %d\n", v)
	}
	// Output: 2, 5, 10, 17, 26  (n² + 1)

	fmt.Println("\n--- Pipeline with early cancel ---")
	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2() // ensure cancel is always called

	nums2 := generate(ctx2, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	squares2 := square(ctx2, nums2)

	count := 0
	for v := range squares2 {
		fmt.Printf("  value: %d\n", v)
		count++
		if count >= 3 {
			cancel2() // cancel pipeline after 3 values
			break
		}
	}
	fmt.Println("  Pipeline cancelled after 3 values")
}
