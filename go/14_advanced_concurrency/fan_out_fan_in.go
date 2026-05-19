package main

import (
	"fmt"
	"sync"
	"time"
)

// Fan-out: 1 input → multiple workers process in parallel
// Fan-in: multiple outputs → 1 channel

// fanOut distributes jobs from input across n worker channels
func fanOut(input <-chan int, n int) []<-chan int {
	outputs := make([]<-chan int, n)
	for i := range n {
		ch := make(chan int)
		outputs[i] = ch
		go func(out chan<- int) {
			defer close(out)
			for v := range input {
				// Simulate processing
				time.Sleep(10 * time.Millisecond)
				out <- v * v // square
			}
		}(ch)
	}
	return outputs
}

// fanIn merges multiple channels into one
func fanIn(inputs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	for _, in := range inputs {
		wg.Add(1)
		go func(ch <-chan int) {
			defer wg.Done()
			for v := range ch {
				out <- v
			}
		}(in)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func demoFanOutFanIn() {
	fmt.Println("\n--- Fan-out / Fan-in ---")

	// Input: 1, 2, 3, 4, 5
	input := make(chan int, 5)
	for i := range 5 {
		input <- i + 1
	}
	close(input)

	// Fan-out to 3 workers
	workers := fanOut(input, 3)
	fmt.Printf("  Distributed across %d workers\n", len(workers))

	// Fan-in results
	results := fanIn(workers...)

	var all []int
	for v := range results {
		all = append(all, v)
	}
	fmt.Printf("  Results (unordered): %v\n", all)
}
