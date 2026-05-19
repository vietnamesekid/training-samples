package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// State of the Circuit Breaker
type CBState int

const (
	StateClosed   CBState = iota // normal, allows requests
	StateOpen                    // too many errors, blocks all requests
	StateHalfOpen                // trying to recover, allows 1 request through
)

func (s CBState) String() string {
	return [...]string{"Closed", "Open", "HalfOpen"}[s]
}

var ErrCircuitOpen = errors.New("circuit breaker is open")

// CircuitBreaker prevents cascade failure by "opening the circuit" when a service fails
type CircuitBreaker struct {
	mu sync.Mutex

	state      CBState
	failures   int
	lastFailed time.Time

	maxFailures int
	timeout     time.Duration // wait time before retrying (HalfOpen)
}

func NewCircuitBreaker(maxFailures int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures: maxFailures,
		timeout:     timeout,
	}
}

// Call executes fn through the circuit breaker
func (cb *CircuitBreaker) Call(fn func() error) error {
	cb.mu.Lock()

	switch cb.state {
	case StateOpen:
		// Check timeout to transition to HalfOpen
		if time.Since(cb.lastFailed) > cb.timeout {
			cb.state = StateHalfOpen
			fmt.Printf("    CB: %s → HalfOpen (retrying)\n", StateOpen)
		} else {
			cb.mu.Unlock()
			return ErrCircuitOpen
		}
	case StateHalfOpen:
		// Only let 1 request through to test
	}

	cb.mu.Unlock()

	// Execute request
	err := fn()

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failures++
		cb.lastFailed = time.Now()

		if cb.state == StateHalfOpen || cb.failures >= cb.maxFailures {
			cb.state = StateOpen
			fmt.Printf("    CB: → Open (failures=%d)\n", cb.failures)
		}
		return err
	}

	// Success → reset
	if cb.state == StateHalfOpen {
		fmt.Printf("    CB: HalfOpen → Closed (recovered)\n")
	}
	cb.state = StateClosed
	cb.failures = 0
	return nil
}

func (cb *CircuitBreaker) State() CBState {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

func demoCircuitBreaker() {
	failCount := 0

	// Simulate a service that fails frequently
	unreliableService := func() error {
		failCount++
		if failCount <= 4 {
			return fmt.Errorf("service unavailable (attempt %d)", failCount)
		}
		return nil // recovers after 4 failures
	}

	// maxFailures=3, timeout=100ms
	cb := NewCircuitBreaker(3, 100*time.Millisecond)

	for i := range 8 {
		err := cb.Call(unreliableService)
		state := cb.State()
		if err != nil {
			if errors.Is(err, ErrCircuitOpen) {
				fmt.Printf("  Request %d: BLOCKED (circuit=%s)\n", i+1, state)
			} else {
				fmt.Printf("  Request %d: ERROR: %v (circuit=%s)\n", i+1, err, state)
			}
		} else {
			fmt.Printf("  Request %d: OK (circuit=%s)\n", i+1, state)
		}

		// Simulate timeout after request 5 to let circuit return to HalfOpen
		if i == 4 {
			time.Sleep(120 * time.Millisecond)
		}
	}
}
