package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"
)

var (
	ErrPoolClosed      = errors.New("worker pool: closed")
	ErrQueueFull       = errors.New("worker pool: queue full")
	ErrShutdownTimeout = errors.New("worker pool: shutdown timed out")
)

type Job func(ctx context.Context)

// Pool

type Pool struct {
	queue  chan Job
	stop   chan struct{} // closed on Shutdown — broadcast to all goroutines
	once   sync.Once     // ensures stop is closed exactly once
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
	log    *slog.Logger
}

func New(workers, capacity int, log *slog.Logger) *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	p := &Pool{
		queue:  make(chan Job, capacity),
		stop:   make(chan struct{}),
		ctx:    ctx,
		cancel: cancel,
		log:    log,
	}

	for range workers {
		p.wg.Add(1)
		go p.runWorker()
	}

	return p
}

// Submit blocks if queue is full. Returns ErrPoolClosed if shut down.
func (p *Pool) Submit(job Job) error {
	select {
	case p.queue <- job:
		return nil
	case <-p.stop:
		return ErrPoolClosed
	}
}

// TrySubmit is non-blocking. Returns ErrQueueFull or ErrPoolClosed.
func (p *Pool) TrySubmit(job Job) error {
	select {
	case p.queue <- job:
		return nil
	case <-p.stop:
		return ErrPoolClosed
	default:
		return ErrQueueFull
	}
}

// SubmitWithContext respects the caller's deadline.
func (p *Pool) SubmitWithContext(ctx context.Context, job Job) error {
	select {
	case p.queue <- job:
		return nil
	case <-p.stop:
		return ErrPoolClosed
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Shutdown drains queued jobs then waits for workers to exit.
// On timeout, cancels the pool context to interrupt running jobs.
func (p *Pool) Shutdown(timeout time.Duration) error {
	p.once.Do(func() {
		close(p.stop) // broadcast: no new submissions accepted
	})

	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		p.cancel()
		return ErrShutdownTimeout
	}
}

func (p *Pool) runWorker() {
	defer p.wg.Done()
	for {
		select {
		case job, ok := <-p.queue:
			if !ok {
				return
			}
			p.safeExecute(job)
		case <-p.stop:
			// drain remaining jobs before exiting
			for {
				select {
				case job := <-p.queue:
					p.safeExecute(job)
				default:
					return
				}
			}
		}
	}
}

func (p *Pool) safeExecute(job Job) {
	defer func() {
		if r := recover(); r != nil {
			p.log.Error("worker panic", "recovered", r)
		}
	}()
	job(p.ctx)
}

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	workers := 5    // number of worker goroutines
	capacity := 100 // max queued jobs

	pool := New(workers, capacity, log)

	var wg sync.WaitGroup

	for i := range 10 {
		task := i + 1 // capture loop variable
		wg.Add(1)

		pool.Submit(func(ctx context.Context) {
			defer wg.Done()
			log.Info("processing", "task", task)
			select {
			case <-time.After(300 * time.Millisecond):
			case <-ctx.Done():
				log.Warn("cancelled", "task", task)
			}
		})
	}

	// Job that panics — pool must survive
	pool.TrySubmit(func(_ context.Context) {
		panic("intentional panic")
	})

	wg.Wait()

	fmt.Println("all jobs done, shutting down...")

	if err := pool.Shutdown(5 * time.Second); err != nil {
		log.Error("shutdown", "err", err)
		os.Exit(1)
	}

	fmt.Println("pool closed cleanly")
}
