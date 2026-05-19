package main

import (
	"bytes"
	"fmt"
	"sync"
)

// === sync.Once — run a function exactly once, thread-safe ===

type Database struct {
	Host string
	Port int
}

var (
	dbInstance *Database
	dbOnce     sync.Once
)

// GetDB — singleton pattern with sync.Once
// Thread-safe: if 1000 goroutines call simultaneously, init only runs once
func GetDB() *Database {
	dbOnce.Do(func() {
		fmt.Println("  [Once] Initializing database connection...")
		dbInstance = &Database{Host: "localhost", Port: 5432}
	})
	return dbInstance
}

// === sync.Pool — object pooling to reduce GC pressure ===
//
// Pool allows reusing objects instead of continuously allocating/GC-ing them
// Use for: bytes.Buffer, JSON encoders, temporary buffers
// IMPORTANT: Pool may be cleared by GC at any time — do not use for state

var bufPool = sync.Pool{
	New: func() any {
		// Create a new buffer when pool is empty
		buf := &bytes.Buffer{}
		buf.Grow(256) // pre-allocate 256 bytes
		return buf
	},
}

func processDataWithPool(data string) string {
	// Get buffer from pool (or create new one if pool is empty)
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset() // IMPORTANT: reset before use
	defer bufPool.Put(buf) // return to pool when done

	buf.WriteString("processed: ")
	buf.WriteString(data)
	return buf.String()
}

// For comparison: without pool
func processDataWithoutPool(data string) string {
	var buf bytes.Buffer // allocate new each time
	buf.WriteString("processed: ")
	buf.WriteString(data)
	return buf.String()
}

func demoOncePool() {
	fmt.Println("\n--- sync.Once (singleton) ---")

	var wg sync.WaitGroup
	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			db := GetDB()
			_ = db
		}()
	}
	wg.Wait()
	fmt.Printf("  DB instance: %+v\n", GetDB())
	fmt.Println("  \"Initializing...\" chỉ in 1 lần dù 5 goroutines gọi GetDB()")

	fmt.Println("\n--- sync.Pool ---")
	// Use pool multiple times — buffer is reused
	results := make([]string, 5)
	for i := range 5 {
		results[i] = processDataWithPool(fmt.Sprintf("item %d", i))
	}
	for _, r := range results {
		fmt.Printf("  %s\n", r)
	}

	fmt.Println("\n  Khi nào dùng sync.Pool:")
	fmt.Println("  - Tạo và GC nhiều objects tương tự")
	fmt.Println("  - JSON encoding buffers, HTTP request scratch buffers")
	fmt.Println("  - Đo benchmark trước và sau để confirm improvement")
	fmt.Println()
	fmt.Println("  Không dùng Pool cho:")
	fmt.Println("  - Long-lived objects (Pool có thể bị GC clear)")
	fmt.Println("  - Objects với external state (connection, file handle)")
}
