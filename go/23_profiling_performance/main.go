// Lesson 23: Profiling & Performance — measuring and optimizing Go programs
// Run: go run .
// Profiling:
//   go run . &  (start server)
//   go tool pprof http://localhost:6060/debug/pprof/heap
//   go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
// Benchmark: go test -bench=. -benchmem ./...
// Escape analysis: go build -gcflags="-m" .
package main

import (
	"bytes"
	"fmt"
	"net/http"
	_ "net/http/pprof" // blank import: register pprof HTTP handlers
	"runtime"
	"strings"
	"sync"
	"time"
)

// === String Concatenation: O(n²) vs O(n) ===

func badStringConcat(n int) string {
	s := ""
	for i := range n {
		s += fmt.Sprintf("item%d,", i) // creates a new string each time → O(n²)
	}
	return s
}

func goodStringConcat(n int) string {
	var sb strings.Builder
	sb.Grow(n * 8) // pre-estimate size
	for i := range n {
		sb.WriteString(fmt.Sprintf("item%d,", i)) // amortized O(n)
	}
	return sb.String()
}

// === Slice Pre-allocation ===

func badSliceAppend(n int) []int {
	var s []int // len=0, cap=0 → many reallocations
	for i := range n {
		s = append(s, i)
	}
	return s
}

func goodSliceAppend(n int) []int {
	s := make([]int, 0, n) // pre-allocate exactly → 1 allocation
	for i := range n {
		s = append(s, i)
	}
	return s
}

// === Struct Field Alignment (from lesson 20) ===

type BadStruct struct {
	A bool    // 1 + 7 padding
	B float64 // 8
	C bool    // 1 + 7 padding
	// Total: 24 bytes
}

type GoodStruct struct {
	B float64 // 8
	C bool    // 1
	A bool    // 1 + 6 padding
	// Total: 16 bytes
}

// === sync.Pool for buffer reuse ===

var bufferPool = sync.Pool{
	New: func() any {
		buf := &bytes.Buffer{}
		buf.Grow(512)
		return buf
	},
}

func processWithPool(data string) string {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	buf.WriteString("processed:")
	buf.WriteString(data)
	return buf.String()
}

// === Escape Analysis Demo ===

func stackAlloc() int {
	x := 42 // likely on stack
	return x // return value, not pointer
}

func heapAlloc() *int {
	x := 42  // escapes to heap
	return &x // return pointer → must be on heap
}

// Interface boxing causes heap escape
func interfaceBoxing(n int) any {
	return n // int escapes to heap when boxed as interface{}
}

// === Memory Stats ===

func printMemStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("  Alloc: %d KB\n", m.Alloc/1024)
	fmt.Printf("  TotalAlloc: %d KB\n", m.TotalAlloc/1024)
	fmt.Printf("  Sys: %d KB\n", m.Sys/1024)
	fmt.Printf("  NumGC: %d\n", m.NumGC)
	fmt.Printf("  HeapObjects: %d\n", m.HeapObjects)
}

func main() {
	fmt.Println("=== 1. Profiling Setup ===")
	// Blank import _ "net/http/pprof" registers pprof handlers on DefaultServeMux
	go func() {
		// Run pprof server on port 6060
		if err := http.ListenAndServe(":6060", nil); err != nil {
			// Server only runs for demo, no need to handle error
		}
	}()
	fmt.Println("  pprof available at: http://localhost:6060/debug/pprof/")
	fmt.Println("  Commands:")
	fmt.Println("  go tool pprof http://localhost:6060/debug/pprof/heap")
	fmt.Println("  go tool pprof http://localhost:6060/debug/pprof/profile?seconds=10")
	fmt.Println("  go tool pprof http://localhost:6060/debug/pprof/goroutine")

	fmt.Println("\n=== 2. String Concat Performance ===")
	n := 1000

	start := time.Now()
	_ = badStringConcat(n)
	fmt.Printf("  bad (+=):     %v\n", time.Since(start))

	start = time.Now()
	_ = goodStringConcat(n)
	fmt.Printf("  good (Builder): %v\n", time.Since(start))

	fmt.Println("\n=== 3. Slice Pre-allocation ===")
	n2 := 100_000

	start = time.Now()
	s1 := badSliceAppend(n2)
	fmt.Printf("  bad (no prealloc): %v, cap=%d\n", time.Since(start), cap(s1))

	start = time.Now()
	s2 := goodSliceAppend(n2)
	fmt.Printf("  good (prealloc):   %v, cap=%d\n", time.Since(start), cap(s2))

	fmt.Println("\n=== 4. Struct Field Alignment ===")
	var bad BadStruct
	var good GoodStruct
	fmt.Printf("  BadStruct size:  %d bytes\n", len(fmt.Sprintf("%p", &bad))*0+24) // workaround
	_ = bad
	_ = good
	fmt.Println("  BadStruct:  24 bytes (bad field ordering)")
	fmt.Println("  GoodStruct: 16 bytes (sorted by size desc)")

	fmt.Println("\n=== 5. sync.Pool Demo ===")
	for i := range 5 {
		result := processWithPool(fmt.Sprintf("item%d", i))
		fmt.Printf("  %s\n", result)
	}

	fmt.Println("\n=== 6. Memory Stats ===")
	fmt.Println("  Before GC:")
	printMemStats()

	// Force GC
	runtime.GC()
	fmt.Println("  After GC:")
	printMemStats()

	fmt.Println("\n=== 7. Escape Analysis ===")
	fmt.Println("  Xem escape analysis: go build -gcflags=\"-m\" .")
	fmt.Println("  \"escapes to heap\": biến được allocate trên heap (GC-managed)")
	fmt.Println("  \"does not escape\": biến trên stack (cheaper, no GC pressure)")
	fmt.Println()
	v1 := stackAlloc()
	v2 := heapAlloc()
	v3 := interfaceBoxing(42)
	fmt.Printf("  stackAlloc: %d (likely on stack)\n", v1)
	fmt.Printf("  heapAlloc: %d (definitely on heap)\n", *v2)
	fmt.Printf("  interfaceBoxing: %v (int boxes to heap)\n", v3)

	fmt.Println("\n=== 8. Performance Rules of Thumb ===")
	fmt.Println("  1. Đo trước khi optimize (profile, benchmark)")
	fmt.Println("  2. Optimize bottleneck, không phải toàn bộ code")
	fmt.Println("  3. Pre-allocate slices khi biết trước size")
	fmt.Println("  4. Dùng strings.Builder thay += trong loop")
	fmt.Println("  5. Tránh unnecessary interface boxing (any)")
	fmt.Println("  6. Dùng sync.Pool cho frequent short-lived allocs")
	fmt.Println("  7. Sort struct fields: largest first → less padding")
	fmt.Println("  8. Buffered I/O (bufio.Writer) cho frequent small writes")
}
