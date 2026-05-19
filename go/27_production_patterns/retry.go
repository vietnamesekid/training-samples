package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// RetryConfig configures retry behavior
type RetryConfig struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
	Multiplier  float64
}

var DefaultRetryConfig = RetryConfig{
	MaxAttempts: 5,
	BaseDelay:   10 * time.Millisecond,
	MaxDelay:    500 * time.Millisecond,
	Multiplier:  2.0,
}

// IsRetryable determines whether to retry
type IsRetryable func(err error) bool

// RetryWithBackoff executes fn with exponential backoff
// delay: base, base*2, base*4, ... capped at MaxDelay
func RetryWithBackoff(ctx context.Context, cfg RetryConfig, retryable IsRetryable, fn func() error) error {
	delay := cfg.BaseDelay

	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}

		if retryable != nil && !retryable(err) {
			return fmt.Errorf("non-retryable error: %w", err)
		}

		if attempt == cfg.MaxAttempts {
			return fmt.Errorf("max attempts (%d) exceeded: %w", cfg.MaxAttempts, err)
		}

		fmt.Printf("  Attempt %d failed: %v, retrying in %v\n", attempt, err, delay)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}

		// Exponential backoff with cap
		delay = time.Duration(float64(delay) * cfg.Multiplier)
		if delay > cfg.MaxDelay {
			delay = cfg.MaxDelay
		}
	}

	return errors.New("retry exhausted")
}

var ErrNotRetryable = errors.New("permanent error")

func demoRetry() {
	attempt := 0

	// Simulate a service that recovers after 3 failures
	operation := func() error {
		attempt++
		if attempt < 3 {
			return fmt.Errorf("temporary error (attempt %d)", attempt)
		}
		fmt.Printf("  Attempt %d: SUCCESS\n", attempt)
		return nil
	}

	ctx := context.Background()
	cfg := RetryConfig{
		MaxAttempts: 5,
		BaseDelay:   5 * time.Millisecond,
		MaxDelay:    50 * time.Millisecond,
		Multiplier:  2.0,
	}

	err := RetryWithBackoff(ctx, cfg, nil, operation)
	if err != nil {
		fmt.Printf("  Final error: %v\n", err)
	} else {
		fmt.Println("  Operation completed successfully")
	}

	// Non-retryable error
	fmt.Println()
	err2 := RetryWithBackoff(ctx, cfg,
		func(err error) bool { return !errors.Is(err, ErrNotRetryable) },
		func() error { return ErrNotRetryable },
	)
	fmt.Printf("  Non-retryable: %v\n", err2)
}
