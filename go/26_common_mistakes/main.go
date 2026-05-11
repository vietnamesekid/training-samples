// Bài 26: 14 Common Mistakes trong Go
// Các lỗi phổ biến nhất và cách sửa
// Chạy: go run .
package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

func main() {
	fmt.Println("=== 14 Common Go Mistakes ===")

	mistake1_goroutineLeak()
	mistake2_nilInterfaceComparison()
	mistake3_mutexCopy()
	mistake4_concurrentMapWrite()
	mistake5_loopVariableCapture()
	mistake6_deferInLoop()
	mistake7_sliceHeaderCopy()
	mistake8_ignoreError()
	mistake9_stringConcatInLoop()
	mistake10_emptyInterfaceType()
	mistake11_timeFormat()
	mistake12_rangeMapOrder()
	mistake13_initFunctionSideEffects()
	mistake14_goroutineWithoutWait()
}

// ============================================================
// Mistake 1: Goroutine Leak
// ============================================================

func mistake1_goroutineLeak() {
	fmt.Println("\n--- Mistake 1: Goroutine Leak ---")

	// BAD: goroutine bị block mãi mãi nếu không có consumer
	badLeak := func() chan int {
		ch := make(chan int)
		go func() {
			ch <- 42 // nếu caller bỏ qua channel → goroutine leak mãi mãi
		}()
		return ch
	}

	ch := badLeak()
	val := <-ch // nếu quên dòng này → goroutine leak
	fmt.Printf("  BAD pattern (but consumed correctly): %d\n", val)

	// GOOD: dùng buffered channel hoặc done channel để cancel
	done := make(chan struct{})
	results := make(chan int, 1) // buffered: goroutine không bị block

	go func() {
		select {
		case results <- 42:
		case <-done: // goroutine có thể thoát khi cần
		}
	}()

	close(done) // signal goroutine thoát
	fmt.Println("  GOOD: goroutine có done channel để cancel")

	// NGUYÊN TẮC: mọi goroutine cần có cách thoát (done channel / context.Done)
}

// ============================================================
// Mistake 2: Nil Interface Comparison
// ============================================================

type MyError struct{ msg string }

func (e *MyError) Error() string { return e.msg }

// BAD: trả về *MyError qua interface error — nil check sẽ FAIL!
func badGetError(fail bool) error {
	var err *MyError // nil pointer
	if fail {
		err = &MyError{"something went wrong"}
	}
	// GOTCHA: interface{type=*MyError, value=nil} != nil
	return err
}

// GOOD: trả về nil trực tiếp khi không có lỗi
func goodGetError(fail bool) error {
	if fail {
		return &MyError{"something went wrong"}
	}
	return nil // nil interface — check đúng
}

func mistake2_nilInterfaceComparison() {
	fmt.Println("\n--- Mistake 2: Nil Interface Comparison ---")

	err := badGetError(false)
	// err có type *MyError nhưng value nil → interface không nil!
	fmt.Printf("  badGetError(false) == nil: %v (sai! dù không có lỗi)\n", err == nil)

	err2 := goodGetError(false)
	fmt.Printf("  goodGetError(false) == nil: %v (đúng)\n", err2 == nil)

	// NGUYÊN TẮC: không return concrete nil pointer qua interface
}

// ============================================================
// Mistake 3: Mutex Copy
// ============================================================

// BAD pattern (go vet sẽ báo: "passes lock by value: contains sync.Mutex"):
// func (c BadCounter) BadIncrement() {  // value receiver → copy mutex!
//     c.mu.Lock()
//     defer c.mu.Unlock()
//     c.count++  // thay đổi trên bản copy, original không bị ảnh hưởng
// }

type GoodCounter struct {
	mu    sync.Mutex
	count int
}

// GOOD: pointer receiver
func (c *GoodCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *GoodCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

func mistake3_mutexCopy() {
	fmt.Println("\n--- Mistake 3: Mutex Copy ---")
	fmt.Println("  BAD: value receiver copies the Mutex → count stays 0, lock is broken")
	fmt.Println("  go vet flags: \"passes lock by value: contains sync.Mutex\"")
	fmt.Println("  See commented-out BadCounter above for the problematic pattern")

	good := &GoodCounter{}
	good.Increment()
	good.Increment()
	fmt.Printf("  GOOD: count = %d\n", good.Value())

	// NGUYÊN TẮC: struct có Mutex → luôn dùng pointer receiver
	// go vet cảnh báo: "assignment copies lock value"
}

// ============================================================
// Mistake 4: Concurrent Map Write
// ============================================================

func mistake4_concurrentMapWrite() {
	fmt.Println("\n--- Mistake 4: Concurrent Map Write ---")

	// BAD: concurrent write vào map thông thường → panic!
	// var m = map[string]int{}
	// for i := 0; i < 100; i++ {
	//     go func(i int) { m[fmt.Sprintf("key%d", i)] = i }(i) // PANIC!
	// }
	fmt.Println("  BAD: concurrent map write → fatal: concurrent map writes")

	// GOOD option 1: sync.Map
	var sm sync.Map
	var wg sync.WaitGroup
	for i := range 10 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sm.Store(fmt.Sprintf("key%d", i), i)
		}(i)
	}
	wg.Wait()

	count := 0
	sm.Range(func(_, _ any) bool { count++; return true })
	fmt.Printf("  GOOD sync.Map: %d entries\n", count)

	// GOOD option 2: regular map + Mutex
	var mu sync.Mutex
	m2 := make(map[string]int)
	for i := range 10 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			mu.Lock()
			m2[fmt.Sprintf("key%d", i)] = i
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	fmt.Printf("  GOOD map+Mutex: %d entries\n", len(m2))
}

// ============================================================
// Mistake 5: Loop Variable Capture (pre Go 1.22)
// ============================================================

func mistake5_loopVariableCapture() {
	fmt.Println("\n--- Mistake 5: Loop Variable Capture ---")

	// BAD pattern (trước Go 1.22):
	// for i := 0; i < 3; i++ {
	//     go func() { fmt.Println(i) }() // in 3,3,3 — tất cả capture cùng biến i
	// }

	var wg sync.WaitGroup
	results := make([]int, 3)

	// Cách đúng cũ (Go < 1.22): truyền qua parameter
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx] = idx * idx
		}(i)
	}
	wg.Wait()
	fmt.Printf("  Correct capture (param copy): %v\n", results)

	// Go 1.22+: loop variable per-iteration, range cũng đúng
	results2 := make([]int, 3)
	for i := range 3 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results2[i] = i * i // Go 1.22+: i là per-iteration
		}()
	}
	wg.Wait()
	fmt.Printf("  Go 1.22+ per-iteration: %v\n", results2)
}

// ============================================================
// Mistake 6: Defer in Loop
// ============================================================

func mistake6_deferInLoop() {
	fmt.Println("\n--- Mistake 6: Defer in Loop ---")

	files := []string{"a.txt", "b.txt", "c.txt"}

	// BAD: defer trong loop chạy khi function return, không phải cuối iteration
	// for _, f := range files {
	//     open(f) → defer close(f) // tất cả close chạy cùng lúc khi func return!
	// }
	fmt.Println("  BAD: defer trong loop → tất cả resource đóng cùng lúc khi func return")

	// GOOD: wrap trong anonymous function
	for _, f := range files {
		func(name string) {
			// open(name)
			defer fmt.Printf("  GOOD: close %s (immediate)\n", name)
			// process(name)...
		}(f)
	}
}

// ============================================================
// Mistake 7: Slice Header Copy
// ============================================================

func mistake7_sliceHeaderCopy() {
	fmt.Println("\n--- Mistake 7: Slice Header Copy ---")

	// BAD: function nhận slice by value, append bên trong không thay đổi caller
	badAppend := func(s []int, val int) {
		s = append(s, val) // s là bản copy header, caller không thấy
		_ = s
	}

	// GOOD option 1: trả về slice mới
	goodAppend := func(s []int, val int) []int {
		return append(s, val)
	}

	// GOOD option 2: pointer to slice
	goodAppendPtr := func(s *[]int, val int) {
		*s = append(*s, val)
	}

	original := []int{1, 2, 3}
	badAppend(original, 99)
	fmt.Printf("  BAD: after badAppend: %v (99 không xuất hiện)\n", original)

	result := goodAppend(original, 99)
	fmt.Printf("  GOOD return: %v\n", result)

	goodAppendPtr(&original, 88)
	fmt.Printf("  GOOD pointer: %v\n", original)
}

// ============================================================
// Mistake 8: Ignoring Errors
// ============================================================

func riskyOperation() (int, error) {
	return 0, fmt.Errorf("operation failed")
}

func mistake8_ignoreError() {
	fmt.Println("\n--- Mistake 8: Ignoring Errors ---")

	// BAD: bỏ qua error
	val, _ := riskyOperation() // lỗi bị nuốt mất
	fmt.Printf("  BAD: val=%d, error silently ignored\n", val)

	// GOOD: luôn handle error
	val2, err := riskyOperation()
	if err != nil {
		fmt.Printf("  GOOD: handled error: %v\n", err)
	} else {
		fmt.Printf("  GOOD: val=%d\n", val2)
	}

	// NGUYÊN TẮC: chỉ dùng _ khi thực sự chắc chắn không cần error
	// golangci-lint errcheck sẽ cảnh báo unhandled errors
}

// ============================================================
// Mistake 9: String Concatenation in Loop
// ============================================================

func mistake9_stringConcatInLoop() {
	fmt.Println("\n--- Mistake 9: String Concat in Loop ---")

	n := 100

	// BAD: O(n²) allocations — mỗi += tạo string mới
	badResult := ""
	for i := range n {
		badResult += fmt.Sprintf("%d,", i)
	}

	// GOOD: strings.Builder — O(n)
	var sb strings.Builder
	sb.Grow(n * 4) // pre-allocate
	for i := range n {
		fmt.Fprintf(&sb, "%d,", i)
	}
	goodResult := sb.String()

	fmt.Printf("  BAD concat len: %d (O(n²) allocs)\n", len(badResult))
	fmt.Printf("  GOOD Builder len: %d (O(n) allocs)\n", len(goodResult))

	// NGUYÊN TẮC: dùng strings.Builder hoặc []byte khi concat nhiều lần
}

// ============================================================
// Mistake 10: Empty Interface Overuse
// ============================================================

// GOOD: define explicit interface
type Processor interface {
	Process() string
}

type IntData struct{ val int }

func (d IntData) Process() string { return fmt.Sprintf("int: %d", d.val) }

type StrData struct{ val string }

func (d StrData) Process() string { return fmt.Sprintf("str: %s", d.val) }

func mistake10_emptyInterfaceType() {
	fmt.Println("\n--- Mistake 10: Empty Interface Overuse ---")

	// BAD: dùng any thay vì typed interface
	badProcess := func(data any) string {
		switch v := data.(type) {
		case int:
			return fmt.Sprintf("int: %d", v)
		case string:
			return fmt.Sprintf("str: %s", v)
		default:
			return fmt.Sprintf("unknown: %T", v)
		}
	}

	// GOOD: typed interface — compile-time safety
	goodProcess := func(p Processor) string {
		return p.Process()
	}

	fmt.Printf("  BAD any: %s\n", badProcess(42))
	fmt.Printf("  BAD any: %s\n", badProcess("hello"))
	fmt.Printf("  GOOD interface: %s\n", goodProcess(IntData{99}))
	fmt.Printf("  GOOD interface: %s\n", goodProcess(StrData{"world"}))
}

// ============================================================
// Mistake 11: Wrong Time Format
// ============================================================

func mistake11_timeFormat() {
	fmt.Println("\n--- Mistake 11: Wrong Time Format ---")

	now := time.Now()

	// BAD: dùng format string kiểu Python/Java không hoạt động trong Go
	// now.Format("yyyy-MM-dd") → sai hoàn toàn!
	// now.Format("YYYY-MM-DD") → sai!
	fmt.Println("  BAD: Format(\"yyyy-MM-dd\") → Go không dùng pattern letters!")

	// GOOD: Go dùng reference time: Mon Jan 2 15:04:05 MST 2006
	fmt.Printf("  GOOD: %s\n", now.Format("2006-01-02 15:04:05"))
	fmt.Printf("  Date only: %s\n", now.Format("2006-01-02"))
	fmt.Printf("  Time only: %s\n", now.Format("15:04:05"))
	fmt.Printf("  RFC3339: %s\n", now.Format(time.RFC3339))
	fmt.Printf("  Custom: %s\n", now.Format("02/01/2006"))

	// NGUYÊN TẮC: reference time = January 2, 15:04:05, 2006, UTC-7
	// 2006=year, 01=month, 02=day, 15=hour(24h), 04=minute, 05=second
}

// ============================================================
// Mistake 12: Map Iteration Order
// ============================================================

func mistake12_rangeMapOrder() {
	fmt.Println("\n--- Mistake 12: Map Iteration Order ---")

	m := map[string]int{"c": 3, "a": 1, "b": 2}

	// BAD: assume map iteration order là fixed
	fmt.Print("  Map iteration (random, do not rely on order): ")
	for k, v := range m {
		fmt.Printf("%s=%d ", k, v)
	}
	fmt.Println()

	// GOOD: dùng slice + sort nếu cần thứ tự xác định
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// sort.Strings(keys) — cần để đảm bảo thứ tự
	fmt.Printf("  GOOD: collect keys first: %v\n", keys)
	fmt.Println("  GOOD: then sort.Strings(keys) for deterministic order")

	// NGUYÊN TẮC: Go randomize map iteration intentionally
	// để tránh code depend on implementation-defined behavior
}

// ============================================================
// Mistake 13: init() Side Effects
// ============================================================

// BAD: init() với side effects không rõ ràng
// func init() {
//     db, _ = sql.Open("postgres", os.Getenv("DATABASE_URL")) // error ignored!
//     globalConfig = loadConfig() // panic nếu file thiếu
// }

// GOOD: explicit initialization, testable và có error handling
func initializeApp() error {
	// setup database, config, etc. với proper error handling
	return nil
}

func mistake13_initFunctionSideEffects() {
	fmt.Println("\n--- Mistake 13: init() Side Effects ---")
	fmt.Println("  BAD: init() với DB connections, global state, file I/O, ignored errors")
	fmt.Println("  GOOD: explicit setup trong main() với proper error handling")
	fmt.Println("  OK to use init(): register drivers (sql.Register), validate constants")

	if err := initializeApp(); err != nil {
		fmt.Printf("  init failed: %v\n", err)
		return
	}
	fmt.Println("  explicit initialization: OK")
}

// ============================================================
// Mistake 14: Goroutine without Synchronization
// ============================================================

func mistake14_goroutineWithoutWait() {
	fmt.Println("\n--- Mistake 14: Goroutine without Synchronization ---")

	// BAD: main thoát trước khi goroutines hoàn thành
	// go func() { fmt.Println("may never print!") }()
	// main returns → goroutine bị kill
	fmt.Println("  BAD: goroutine started without WaitGroup/channel → may not complete")

	// GOOD: WaitGroup
	var wg sync.WaitGroup
	results := make([]string, 3)

	for i := range 3 {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx] = fmt.Sprintf("result-%d", idx)
		}(i)
	}
	wg.Wait() // đảm bảo tất cả goroutines xong trước khi dùng results
	fmt.Printf("  GOOD: all goroutines done: %v\n", results)

	// NGUYÊN TẮC: luôn sync với goroutines qua WaitGroup, channel, hoặc context
}
