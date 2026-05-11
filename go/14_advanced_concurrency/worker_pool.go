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

// Worker Pool — sản phẩm production-grade
// Giải quyết vấn đề: giới hạn số goroutines đồng thời, tái sử dụng workers

var (
	ErrPoolClosed      = errors.New("worker pool: closed")
	ErrQueueFull       = errors.New("worker pool: queue full")
	ErrShutdownTimeout = errors.New("worker pool: shutdown timed out")
)

type Job func(ctx context.Context)

type Pool struct {
	queue  chan Job
	stop   chan struct{}     // closed on Shutdown — broadcast to all workers
	once   sync.Once        // đảm bảo close(stop) chỉ 1 lần
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
	log    *slog.Logger
}

func NewPool(workers, capacity int, log *slog.Logger) *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	p := &Pool{
		queue:  make(chan Job, capacity),
		stop:   make(chan struct{}),
		ctx:    ctx,
		cancel: cancel,
		log:    log,
	}
	// Launch worker goroutines
	for range workers {
		p.wg.Add(1)
		go p.runWorker()
	}
	return p
}

// Submit block nếu queue đầy, trả về ErrPoolClosed nếu đã shutdown
func (p *Pool) Submit(job Job) error {
	select {
	case p.queue <- job:
		return nil
	case <-p.stop:
		return ErrPoolClosed
	}
}

// TrySubmit non-blocking — trả về ErrQueueFull hoặc ErrPoolClosed ngay lập tức
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

// SubmitWithContext tôn trọng deadline của caller
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

// Shutdown dừng pool, chờ jobs đang chạy, timeout sau duration
func (p *Pool) Shutdown(timeout time.Duration) error {
	p.once.Do(func() {
		close(p.stop) // broadcast: không nhận jobs mới
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
		p.cancel() // interrupt running jobs qua context
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
			// Drain remaining jobs trước khi exit
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

// safeExecute bắt panic, không làm crash worker
func (p *Pool) safeExecute(job Job) {
	defer func() {
		if r := recover(); r != nil {
			p.log.Error("worker panic recovered", "panic", r)
		}
	}()
	job(p.ctx)
}

func demoWorkerPool() {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))
	pool := NewPool(3, 10, log) // 3 workers, queue capacity 10

	var wg sync.WaitGroup

	// Submit 8 jobs
	for i := range 8 {
		task := i + 1
		wg.Add(1)
		if err := pool.Submit(func(ctx context.Context) {
			defer wg.Done()
			fmt.Printf("  Task %d: processing...\n", task)
			time.Sleep(20 * time.Millisecond)
			fmt.Printf("  Task %d: done\n", task)
		}); err != nil {
			wg.Done()
			fmt.Printf("  Submit error: %v\n", err)
		}
	}

	// Job có panic — pool phải survive
	pool.TrySubmit(func(_ context.Context) {
		panic("intentional panic in job")
	})

	wg.Wait()
	fmt.Println("  All jobs done, shutting down...")

	if err := pool.Shutdown(2 * time.Second); err != nil {
		fmt.Printf("  Shutdown error: %v\n", err)
	} else {
		fmt.Println("  Pool shutdown cleanly")
	}
}
