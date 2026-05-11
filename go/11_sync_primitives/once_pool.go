package main

import (
	"bytes"
	"fmt"
	"sync"
)

// === sync.Once — chạy function đúng 1 lần, thread-safe ===

type Database struct {
	Host string
	Port int
}

var (
	dbInstance *Database
	dbOnce     sync.Once
)

// GetDB — singleton pattern với sync.Once
// Thread-safe: nếu 1000 goroutines gọi đồng thời, chỉ 1 lần init
func GetDB() *Database {
	dbOnce.Do(func() {
		fmt.Println("  [Once] Initializing database connection...")
		dbInstance = &Database{Host: "localhost", Port: 5432}
	})
	return dbInstance
}

// === sync.Pool — object pooling để giảm GC pressure ===
//
// Pool cho phép tái sử dụng objects thay vì allocate/GC liên tục
// Dùng cho: bytes.Buffer, JSON encoders, temporary buffers
// QUAN TRỌNG: Pool có thể bị GC clear bất kỳ lúc nào — không dùng cho state

var bufPool = sync.Pool{
	New: func() any {
		// Tạo buffer mới khi pool rỗng
		buf := &bytes.Buffer{}
		buf.Grow(256) // pre-allocate 256 bytes
		return buf
	},
}

func processDataWithPool(data string) string {
	// Lấy buffer từ pool (hoặc tạo mới nếu pool rỗng)
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset() // QUAN TRỌNG: reset trước khi dùng
	defer bufPool.Put(buf) // trả lại pool khi xong

	buf.WriteString("processed: ")
	buf.WriteString(data)
	return buf.String()
}

// So sánh: không dùng pool
func processDataWithoutPool(data string) string {
	var buf bytes.Buffer // allocate mới mỗi lần
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
	// Dùng pool nhiều lần — buffer được tái sử dụng
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
