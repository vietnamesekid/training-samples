package main

import (
	"fmt"
	"sync"
)

// SafeCounter protects a counter with a Mutex
// Mutex: mutual exclusion — only 1 goroutine can access at a time
type SafeCounter struct {
	mu    sync.Mutex
	value int
}

// GOTCHA: Never copy a Mutex! Always use pointer receiver when holding a Mutex
func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock() // defer ensures Unlock is always called even on panic
	c.value++
}

func (c *SafeCounter) Add(n int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value += n
}

func (c *SafeCounter) Get() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// ReadWriteCache uses RWMutex — optimized for read-heavy workloads
// RWMutex: multiple readers CAN read simultaneously, but writer gets exclusive access
type ReadWriteCache struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewReadWriteCache() *ReadWriteCache {
	return &ReadWriteCache{
		data: make(map[string]string),
	}
}

func (c *ReadWriteCache) Set(key, value string) {
	c.mu.Lock()          // exclusive write lock
	defer c.mu.Unlock()
	c.data[key] = value
}

func (c *ReadWriteCache) Get(key string) (string, bool) {
	c.mu.RLock()          // shared read lock — multiple goroutines can RLock simultaneously
	defer c.mu.RUnlock()
	v, ok := c.data[key]
	return v, ok
}

func (c *ReadWriteCache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := make([]string, 0, len(c.data))
	for k := range c.data {
		keys = append(keys, k)
	}
	return keys
}

func demoMutex() {
	fmt.Println("\n--- sync.Mutex ---")
	counter := &SafeCounter{}
	var wg sync.WaitGroup

	// 1000 goroutines incrementing concurrently
	for range 1000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}
	wg.Wait()
	fmt.Printf("  Counter (1000 goroutines): %d\n", counter.Get())

	fmt.Println("\n--- sync.RWMutex ---")
	cache := NewReadWriteCache()

	// Writers
	var wg2 sync.WaitGroup
	for i := range 5 {
		wg2.Add(1)
		go func(n int) {
			defer wg2.Done()
			key := fmt.Sprintf("key%d", n)
			cache.Set(key, fmt.Sprintf("value%d", n))
		}(i)
	}
	wg2.Wait()

	// Many concurrent readers
	for i := range 10 {
		wg2.Add(1)
		go func(n int) {
			defer wg2.Done()
			key := fmt.Sprintf("key%d", n%5)
			v, ok := cache.Get(key)
			if ok {
				fmt.Printf("  Read %s = %s\n", key, v)
			}
		}(i)
	}
	wg2.Wait()

	fmt.Println("\n  NGUYÊN TẮC Mutex:")
	fmt.Println("  - Dùng sync.Mutex khi cần exclusive access")
	fmt.Println("  - Dùng sync.RWMutex khi read >> write")
	fmt.Println("  - Không bao giờ COPY mutex (dùng pointer)")
	fmt.Println("  - Lock và Unlock trong cùng 1 goroutine")
	fmt.Println("  - Dùng defer Unlock() để tránh deadlock")
}
