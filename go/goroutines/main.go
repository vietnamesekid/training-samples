package main

import (
	"errors"
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

func main() {
}
