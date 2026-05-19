package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

// errgroup: run multiple goroutines, collect errors, cancel all if one fails
// Better than WaitGroup when error handling is needed

type UserProfile struct {
	Name   string
	Orders []string
	Points int
}

func fetchUserName(ctx context.Context, userID int) (string, error) {
	select {
	case <-time.After(30 * time.Millisecond):
		return fmt.Sprintf("User%d", userID), nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func fetchUserOrders(ctx context.Context, userID int) ([]string, error) {
	select {
	case <-time.After(50 * time.Millisecond):
		return []string{"order-1", "order-2"}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func fetchUserPoints(ctx context.Context, userID int) (int, error) {
	select {
	case <-time.After(20 * time.Millisecond):
		return 1500, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

func demoErrGroup() {
	ctx := context.Background()
	userID := 42

	fmt.Println("\n--- Parallel fetch with errgroup ---")

	var profile UserProfile
	g, ctx := errgroup.WithContext(ctx)

	// Launch 3 goroutines simultaneously
	g.Go(func() error {
		name, err := fetchUserName(ctx, userID)
		if err != nil {
			return fmt.Errorf("fetch name: %w", err)
		}
		profile.Name = name
		return nil
	})

	g.Go(func() error {
		orders, err := fetchUserOrders(ctx, userID)
		if err != nil {
			return fmt.Errorf("fetch orders: %w", err)
		}
		profile.Orders = orders
		return nil
	})

	g.Go(func() error {
		points, err := fetchUserPoints(ctx, userID)
		if err != nil {
			return fmt.Errorf("fetch points: %w", err)
		}
		profile.Points = points
		return nil
	})

	// Wait for all to finish, returns the first error (if any)
	if err := g.Wait(); err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Profile: name=%s, orders=%v, points=%d\n",
			profile.Name, profile.Orders, profile.Points)
	}

	fmt.Println("\n--- errgroup with SetLimit ---")
	// SetLimit: limit the number of concurrent goroutines
	g2, _ := errgroup.WithContext(context.Background())
	g2.SetLimit(3) // at most 3 goroutines at a time

	for i := range 10 {
		task := i
		g2.Go(func() error {
			time.Sleep(10 * time.Millisecond)
			fmt.Printf("  Task %d done\n", task)
			return nil
		})
	}

	if err := g2.Wait(); err != nil {
		fmt.Printf("  Error: %v\n", err)
	}
}
