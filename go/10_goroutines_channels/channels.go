package main

import (
	"fmt"
	"time"
)

func demoChannels() {
	fmt.Println("\n--- Unbuffered Channel (synchronous) ---")
	// Unbuffered: sender and receiver must both be ready at the same time (rendezvous)
	ch := make(chan string)
	go func() {
		ch <- "hello from goroutine"
	}()
	msg := <-ch
	fmt.Printf("  Received: %q\n", msg)

	fmt.Println("\n--- Buffered Channel (asynchronous) ---")
	// Buffered: send does not block if buffer is not full
	bch := make(chan int, 3)
	bch <- 1
	bch <- 2
	bch <- 3
	// bch <- 4  // ← BLOCKS because buffer is full (no receiver)
	fmt.Printf("  len=%d, cap=%d\n", len(bch), cap(bch))
	fmt.Printf("  Received: %d, %d, %d\n", <-bch, <-bch, <-bch)

	fmt.Println("\n--- close() and range over channel ---")
	jobs := make(chan int, 5)
	go func() {
		for i := range 5 {
			jobs <- i + 1
		}
		close(jobs) // signal that no more values will be sent
	}()

	for j := range jobs { // range stops automatically when channel is closed AND empty
		fmt.Printf("  job: %d\n", j)
	}

	// Check closed channel with 2-value form
	done := make(chan struct{})
	close(done)
	_, ok := <-done
	fmt.Printf("  Receive from closed: ok=%t (false = closed)\n", ok)

	fmt.Println("\n--- Directional Channels ---")
	// chan<- T: send-only
	// <-chan T: receive-only
	// Used to enforce direction at compile time
	ping := make(chan string, 1)
	pong := make(chan string, 1)

	go sender(ping)
	go forwarder(ping, pong)
	result := <-pong
	fmt.Printf("  Received from pipeline: %q\n", result)

	fmt.Println("\n--- select —- multiplex channels ---")
	one := make(chan string)
	two := make(chan string)

	go func() {
		time.Sleep(10 * time.Millisecond)
		one <- "from one"
	}()
	go func() {
		time.Sleep(20 * time.Millisecond)
		two <- "from two"
	}()

	// select picks whichever case is ready first
	for i := range 2 {
		select {
		case v := <-one:
			fmt.Printf("  select [%d]: %s\n", i, v)
		case v := <-two:
			fmt.Printf("  select [%d]: %s\n", i, v)
		}
	}

	fmt.Println("\n--- select with timeout ---")
	slow := make(chan int)
	go func() {
		time.Sleep(200 * time.Millisecond)
		slow <- 42
	}()

	select {
	case v := <-slow:
		fmt.Printf("  Got: %d\n", v)
	case <-time.After(50 * time.Millisecond):
		fmt.Println("  Timeout! (50ms)")
	}

	fmt.Println("\n--- select with default (non-blocking) ---")
	ch2 := make(chan int, 1)
	select {
	case v := <-ch2:
		fmt.Printf("  Got: %d\n", v)
	default:
		fmt.Println("  No value ready (default branch)")
	}
	ch2 <- 42
	select {
	case v := <-ch2:
		fmt.Printf("  Got: %d\n", v)
	default:
		fmt.Println("  No value ready")
	}

	fmt.Println("\n--- Done channel pattern ---")
	// Common pattern: done channel to signal goroutine to stop
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case <-doneCh:
				fmt.Println("  worker: received done signal, stopping")
				return
			default:
				// do work...
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	time.Sleep(50 * time.Millisecond)
	close(doneCh) // broadcast: closing channel wakes up ALL receivers
	time.Sleep(20 * time.Millisecond)
}

// sender can only SEND to ch (chan<-)
func sender(ch chan<- string) {
	ch <- "ping"
}

// forwarder receives from in (<-chan) and sends to out (chan<-)
func forwarder(in <-chan string, out chan<- string) {
	msg := <-in
	out <- msg + " → pong"
}
