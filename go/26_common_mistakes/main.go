// Lesson 26: 14 Common Mistakes in Go
// The most frequent mistakes and how to fix them
// Run: go run .
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

	// BAD: goroutine blocks forever if there is no consumer
	badLeak := func() chan int {
		ch := make(chan int)
		go func() {
			ch <- 42 // if caller ignores channel → goroutine leaks forever
		}()
		return ch
	}

	ch := badLeak()
	val := <-ch // if this line is forgotten → goroutine leak
	fmt.Printf("  BAD pattern (but consumed correctly): %d\n", val)

	// GOOD: use buffered channel or done channel to cancel
	done := make(chan struct{})
	results := make(chan int, 1) // buffered: goroutine won't block

	go func() {
		select {
		case results <- 42:
		case <-done: // goroutine can exit when needed
		}
	}()

	close(done) // signal goroutine to exit
	fmt.Println("  GOOD: goroutine has a done channel to cancel")

	// PRINCIPLE: every goroutine needs a way to exit (done channel / context.Done)
}

// ============================================================
// Mistake 2: Nil Interface Comparison
// ============================================================

type MyError struct{ msg string }

func (e *MyError) Error() string { return e.msg }

// BAD: returning *MyError through error interface — nil check will FAIL!
func badGetError(fail bool) error {
	var err *MyError // nil pointer
	if fail {
		err = &MyError{"something went wrong"}
	}
	// GOTCHA: interface{type=*MyError, value=nil} != nil
	return err
}

// GOOD: return nil directly when there is no error
func goodGetError(fail bool) error {
	if fail {
		return &MyError{"something went wrong"}
	}
	return nil // nil interface — check works correctly
}

func mistake2_nilInterfaceComparison() {
	fmt.Println("\n--- Mistake 2: Nil Interface Comparison ---")

	err := badGetError(false)
	// err has type *MyError but nil value → interface is not nil!
	fmt.Printf("  badGetError(false) == nil: %v (sai! dù không có lỗi)\n", err == nil)

	err2 := goodGetError(false)
	fmt.Printf("  goodGetError(false) == nil: %v (đúng)\n", err2 == nil)

	// PRINCIPLE: do not return a concrete nil pointer through an interface
}

// ============================================================
// Mistake 3: Mutex Copy
// ============================================================

// BAD pattern (go vet will report: "passes lock by value: contains sync.Mutex"):
// func (c BadCounter) BadIncrement() {  // value receiver → copies mutex!
//     c.mu.Lock()
//     defer c.mu.Unlock()
//     c.count++  // changes are on the copy, original is unaffected
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

	// PRINCIPLE: struct with Mutex → always use pointer receiver
	// go vet warns: "assignment copies lock value"
}

// ============================================================
// Mistake 4: Concurrent Map Write
// ============================================================

func mistake4_concurrentMapWrite() {
	fmt.Println("\n--- Mistake 4: Concurrent Map Write ---")

	// BAD: concurrent write to a regular map → panic!
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

	// BAD pattern (before Go 1.22):
	// for i := 0; i < 3; i++ {
	//     go func() { fmt.Println(i) }() // prints 3,3,3 — all capture the same variable i
	// }

	var wg sync.WaitGroup
	results := make([]int, 3)

	// Old correct approach (Go < 1.22): pass through parameter
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx] = idx * idx
		}(i)
	}
	wg.Wait()
	fmt.Printf("  Correct capture (param copy): %v\n", results)

	// Go 1.22+: loop variable per-iteration, range also works correctly
	results2 := make([]int, 3)
	for i := range 3 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results2[i] = i * i // Go 1.22+: i is per-iteration
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

	// BAD: defer in loop runs when the function returns, not at the end of each iteration
	// for _, f := range files {
	//     open(f) → defer close(f) // all closes run at once when func returns!
	// }
	fmt.Println("  BAD: defer trong loop → tất cả resource đóng cùng lúc khi func return")

	// GOOD: wrap in anonymous function
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

	// BAD: function receives slice by value, append inside doesn't affect caller
	badAppend := func(s []int, val int) {
		s = append(s, val) // s is a copy of the header, caller won't see it
		_ = s
	}

	// GOOD option 1: return new slice
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

	// BAD: ignoring error
	val, _ := riskyOperation() // error is silently swallowed
	fmt.Printf("  BAD: val=%d, error silently ignored\n", val)

	// GOOD: always handle error
	val2, err := riskyOperation()
	if err != nil {
		fmt.Printf("  GOOD: handled error: %v\n", err)
	} else {
		fmt.Printf("  GOOD: val=%d\n", val2)
	}

	// PRINCIPLE: only use _ when you are truly certain the error is unneeded
	// golangci-lint errcheck will warn about unhandled errors
}

// ============================================================
// Mistake 9: String Concatenation in Loop
// ============================================================

func mistake9_stringConcatInLoop() {
	fmt.Println("\n--- Mistake 9: String Concat in Loop ---")

	n := 100

	// BAD: O(n²) allocations — each += creates a new string
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

	// PRINCIPLE: use strings.Builder or []byte when concatenating many times
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

	// BAD: use any instead of a typed interface
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

	// BAD: using Python/Java-style format strings doesn't work in Go
	// now.Format("yyyy-MM-dd") → completely wrong!
	// now.Format("YYYY-MM-DD") → wrong!
	fmt.Println("  BAD: Format(\"yyyy-MM-dd\") → Go không dùng pattern letters!")

	// GOOD: Go uses reference time: Mon Jan 2 15:04:05 MST 2006
	fmt.Printf("  GOOD: %s\n", now.Format("2006-01-02 15:04:05"))
	fmt.Printf("  Date only: %s\n", now.Format("2006-01-02"))
	fmt.Printf("  Time only: %s\n", now.Format("15:04:05"))
	fmt.Printf("  RFC3339: %s\n", now.Format(time.RFC3339))
	fmt.Printf("  Custom: %s\n", now.Format("02/01/2006"))

	// PRINCIPLE: reference time = January 2, 15:04:05, 2006, UTC-7
	// 2006=year, 01=month, 02=day, 15=hour(24h), 04=minute, 05=second
}

// ============================================================
// Mistake 12: Map Iteration Order
// ============================================================

func mistake12_rangeMapOrder() {
	fmt.Println("\n--- Mistake 12: Map Iteration Order ---")

	m := map[string]int{"c": 3, "a": 1, "b": 2}

	// BAD: assuming map iteration order is fixed
	fmt.Print("  Map iteration (random, do not rely on order): ")
	for k, v := range m {
		fmt.Printf("%s=%d ", k, v)
	}
	fmt.Println()

	// GOOD: use slice + sort if a deterministic order is needed
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// sort.Strings(keys) — needed to guarantee order
	fmt.Printf("  GOOD: collect keys first: %v\n", keys)
	fmt.Println("  GOOD: then sort.Strings(keys) for deterministic order")

	// PRINCIPLE: Go intentionally randomizes map iteration
	// to prevent code from depending on implementation-defined behavior
}

// ============================================================
// Mistake 13: init() Side Effects
// ============================================================

// BAD: init() with unclear side effects
// func init() {
//     db, _ = sql.Open("postgres", os.Getenv("DATABASE_URL")) // error ignored!
//     globalConfig = loadConfig() // panic if file is missing
// }

// GOOD: explicit initialization, testable and with error handling
func initializeApp() error {
	// setup database, config, etc. with proper error handling
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

	// BAD: main exits before goroutines finish
	// go func() { fmt.Println("may never print!") }()
	// main returns → goroutine is killed
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
	wg.Wait() // ensure all goroutines finish before using results
	fmt.Printf("  GOOD: all goroutines done: %v\n", results)

	// PRINCIPLE: always synchronize with goroutines via WaitGroup, channel, or context
}
