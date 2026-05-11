package main

import (
	"fmt"
	"sync"
	"time"
)

// === sync.Cond — conditional variable ===
// Dùng khi: goroutine cần chờ một điều kiện trở nên đúng
// Phổ biến hơn channel khi: nhiều waiters, cần broadcast

// BlockingQueue: queue thread-safe với blocking get
type BlockingQueue struct {
	mu    sync.Mutex
	cond  *sync.Cond
	items []any
	cap   int
}

func NewBlockingQueue(capacity int) *BlockingQueue {
	q := &BlockingQueue{
		items: make([]any, 0, capacity),
		cap:   capacity,
	}
	q.cond = sync.NewCond(&q.mu)
	return q
}

// Put: block nếu queue đầy
func (q *BlockingQueue) Put(item any) {
	q.mu.Lock()
	for len(q.items) >= q.cap {
		q.cond.Wait() // atomically: release lock + wait + reacquire lock
	}
	q.items = append(q.items, item)
	q.cond.Signal() // wake up one waiter
	q.mu.Unlock()
}

// Get: block nếu queue rỗng
func (q *BlockingQueue) Get() any {
	q.mu.Lock()
	for len(q.items) == 0 {
		q.cond.Wait()
	}
	item := q.items[0]
	q.items = q.items[1:]
	q.cond.Signal() // wake up one waiter (có thể là producer)
	q.mu.Unlock()
	return item
}

func (q *BlockingQueue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}

// Broadcast example: worker pool start signal
type WorkerGroup struct {
	mu      sync.Mutex
	cond    *sync.Cond
	started bool
}

func NewWorkerGroup() *WorkerGroup {
	wg := &WorkerGroup{}
	wg.cond = sync.NewCond(&wg.mu)
	return wg
}

func (wg *WorkerGroup) Start() {
	wg.mu.Lock()
	wg.started = true
	wg.cond.Broadcast() // wake up ALL waiters (khác Signal chỉ wake 1)
	wg.mu.Unlock()
}

func (wg *WorkerGroup) WaitForStart() {
	wg.mu.Lock()
	for !wg.started {
		wg.cond.Wait()
	}
	wg.mu.Unlock()
}

func demoCond() {
	fmt.Println("\n--- BlockingQueue với sync.Cond ---")

	q := NewBlockingQueue(3)
	var wg sync.WaitGroup

	// Consumer goroutines
	for i := range 5 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			item := q.Get()
			fmt.Printf("  Consumer %d received: %v\n", id, item)
		}(i)
	}

	time.Sleep(10 * time.Millisecond)

	// Producer
	for i := range 5 {
		fmt.Printf("  Producing item %d\n", i)
		q.Put(fmt.Sprintf("item-%d", i))
	}

	wg.Wait()
	fmt.Printf("  Queue empty: len=%d\n", q.Len())

	fmt.Println("\n--- Broadcast: workers chờ signal start ---")
	wg2 := NewWorkerGroup()
	var mu sync.Mutex
	started := make([]int, 0)

	for i := range 5 {
		go func(id int) {
			wg2.WaitForStart()
			mu.Lock()
			started = append(started, id)
			mu.Unlock()
		}(i)
	}

	time.Sleep(10 * time.Millisecond)
	fmt.Println("  Broadcasting start signal...")
	wg2.Start()
	time.Sleep(20 * time.Millisecond)

	mu.Lock()
	fmt.Printf("  Workers started: %v\n", started)
	mu.Unlock()

	fmt.Println("\n  Cond.Signal() vs Cond.Broadcast():")
	fmt.Println("  Signal: wake up ONE waiter (biết chính xác ai cần wake)")
	fmt.Println("  Broadcast: wake up ALL waiters (e.g., config changed)")
	fmt.Println()
	fmt.Println("  Lưu ý: Trong thực tế, channel thường đủ dùng hơn sync.Cond")
	fmt.Println("  Cond hữu ích khi cần: multiple conditions, expensive re-check")
}
