package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// State của Circuit Breaker
type CBState int

const (
	StateClosed   CBState = iota // bình thường, cho phép request
	StateOpen                    // lỗi nhiều, chặn tất cả request
	StateHalfOpen                // thử phục hồi, cho 1 request qua
)

func (s CBState) String() string {
	return [...]string{"Closed", "Open", "HalfOpen"}[s]
}

var ErrCircuitOpen = errors.New("circuit breaker is open")

// CircuitBreaker ngăn cascade failure bằng cách "ngắt mạch" khi dịch vụ lỗi
type CircuitBreaker struct {
	mu sync.Mutex

	state      CBState
	failures   int
	lastFailed time.Time

	maxFailures int
	timeout     time.Duration // thời gian chờ trước khi thử lại (HalfOpen)
}

func NewCircuitBreaker(maxFailures int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures: maxFailures,
		timeout:     timeout,
	}
}

// Call thực thi fn qua circuit breaker
func (cb *CircuitBreaker) Call(fn func() error) error {
	cb.mu.Lock()

	switch cb.state {
	case StateOpen:
		// Kiểm tra timeout để chuyển sang HalfOpen
		if time.Since(cb.lastFailed) > cb.timeout {
			cb.state = StateHalfOpen
			fmt.Printf("    CB: %s → HalfOpen (retrying)\n", StateOpen)
		} else {
			cb.mu.Unlock()
			return ErrCircuitOpen
		}
	case StateHalfOpen:
		// Chỉ cho 1 request qua để test
	}

	cb.mu.Unlock()

	// Thực thi request
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

	// Thành công → reset
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

	// Simulate một service hay bị lỗi
	unreliableService := func() error {
		failCount++
		if failCount <= 4 {
			return fmt.Errorf("service unavailable (attempt %d)", failCount)
		}
		return nil // phục hồi sau 4 lần thất bại
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

		// Simulate timeout sau request 5 để circuit trở về HalfOpen
		if i == 4 {
			time.Sleep(120 * time.Millisecond)
		}
	}
}
