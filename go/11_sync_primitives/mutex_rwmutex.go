package main

import (
	"fmt"
	"sync"
)

// SafeCounter bảo vệ counter với Mutex
// Mutex: mutual exclusion — chỉ 1 goroutine được access cùng lúc
type SafeCounter struct {
	mu    sync.Mutex
	value int
}

// GOTCHA: Không copy Mutex! Luôn dùng pointer receiver khi có Mutex
func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock() // defer đảm bảo luôn Unlock dù có panic
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

// ReadWriteCache dùng RWMutex — tối ưu khi đọc nhiều hơn viết
// RWMutex: nhiều readers CÓ THỂ đọc cùng lúc, nhưng writer độc quyền
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
	c.mu.RLock()          // shared read lock — nhiều goroutines có thể RLock cùng lúc
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

	// 1000 goroutines cùng increment
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
