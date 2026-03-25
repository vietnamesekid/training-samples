package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type CircuitBreaker struct {
	mu            sync.Mutex
	isOpen        bool
	threshold     int
	failureCount  int
	lastErrorTime time.Time
}

func (cb *CircuitBreaker) Call(f func() (string, error)) (string, error) {
	cb.mu.Lock()

	if cb.isOpen {
		if time.Since(cb.lastErrorTime) > 5*time.Second {
			cb.isOpen = false
		} else {
			cb.mu.Unlock()
			return "", errors.New("Circuit Breaker is open (service down)")
		}
	}

	cb.mu.Unlock()

	res, err := f()

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failureCount++
		if cb.failureCount >= cb.threshold {
			cb.isOpen = true
			cb.lastErrorTime = time.Now()
		}

		return "", err
	}

	cb.failureCount = 0
	return res, nil
}

func worker(done chan bool) {
	fmt.Println("Worker: Starting work...")
	time.Sleep(2 * time.Second) // Simulate work
	fmt.Println("Worker: Work completed.")
	done <- true
}

func ping(message string, pings chan<- string) {
	pings <- message
}

func pong(pingss <-chan string, pongs chan<- string) {
	msg := <-pingss
	pongs <- msg
}

func main() {
	// message := make(chan string, 2)

	// message <- "ping"
	// message <- "pong"

	// println(<-message)
	// println(<-message)

	// go func() {
	// 	message <- "ping"
	// }()

	// println(<-message)

	// go func() {
	// 	println(<-message)
	// }()

	// message <- "pong"

	// println(<-message)

	// done := make(chan bool)

	// go worker(done)

	// <-done
	// fmt.Println("Main: Worker has completed.")

	// pings := make(chan string, 1)
	// pongs := make(chan string, 1)

	// ping("Send message", pings)
	// pong(pings, pongs)

	// fmt.Println(<-pongs)

	// one := make(chan string)
	// two := make(chan string)

	// go func() {
	// 	defer close(one)

	// 	time.Sleep(1 * time.Second)
	// 	one <- "Message from one"
	// }()

	// go func() {
	// 	defer close(two)

	// 	time.Sleep(2 * time.Second)
	// 	two <- "Message from two"
	// }()

	// for one != nil || two != nil {
	// 	select {
	// 	case val, closed := <-one:
	// 		if !closed {
	// 			one = nil // Đóng channel one để không nhận thêm
	// 		} else {
	// 			fmt.Println("Received from one:", val)
	// 		}
	// 	case val, closed := <-two:
	// 		if !closed {
	// 			two = nil // Đóng channel two để không nhận thêm
	// 		} else {
	// 			fmt.Println("Received from two:", val)
	// 		}

	// 	}
	// }

	// select {
	// case val := <-one:
	// 	fmt.Println("Received from one:", val)
	// case val := <-two:
	// 	fmt.Println("Received from two:", val)
	// case <-time.After(3 * time.Second):
	// 	fmt.Println("Timeout: No messages received")
	// 	break
	// }

	// jobs := make(chan int, 5)
	// done := make(chan struct{})

	// // Worker Goroutine
	// go func() {
	// 	for job := range jobs {
	// 		fmt.Printf("Processing job %d\n", job)
	// 		time.Sleep(1 * time.Second) // Simulate time-consuming work
	// 	}

	// 	close(done)
	// }()

	// // Producer
	// for job := 1; job <= cap(jobs); job++ {
	// 	fmt.Printf("Sending job %d\n", job)
	// 	jobs <- job
	// }

	// // Đóng channel để báo hiệu không còn job nào nữa
	// close(jobs)

	// <-done
	// fmt.Println("All jobs processed, exiting.")

	queue := make(chan string, 3)

	queue <- "Task 1"
	queue <- "Task 2"

	close(queue)

	for task := range queue {
		fmt.Println("Processing:", task)
	}

}
