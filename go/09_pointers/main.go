// Lesson 9: Pointers — pointers in Go
// Run: go run .
// View escape analysis: go build -gcflags="-m" .
package main

import "fmt"

// === Basics: & (address-of) and * (dereference) ===

func zeroval(n int) {
	n = 0 // only changes the local copy
}

func zeroptr(n *int) {
	*n = 0 // changes the value at the address n points to
}

// === Pointers and Structs ===

type Counter struct {
	value int
}

// Value receiver — receives a COPY, cannot mutate
func (c Counter) Get() int {
	return c.value
}

// Pointer receiver — receives a POINTER, can mutate
func (c *Counter) Increment() {
	c.value++
}

func (c *Counter) Add(n int) {
	c.value += n
}

func (c *Counter) Reset() {
	c.value = 0
}

// === When to use pointers ===

type SmallStruct struct {
	X, Y int
}

type LargeStruct struct {
	Data [1024]byte // 1KB — should use pointer to avoid copying
}

// 1. Need to mutate: use pointer receiver
func (s *SmallStruct) MoveBy(dx, dy int) {
	s.X += dx
	s.Y += dy
}

// 2. Large struct: use pointer to avoid copying
func processLarge(s *LargeStruct) {
	s.Data[0] = 42
}

// 3. Optional / nullable value: *T can be nil
type Config struct {
	Port    int
	Timeout *int // nil = "not set, use default"
}

// 4. Shared state: multiple goroutines share the same object
type SharedCache struct {
	data map[string]string
}

func newSharedCache() *SharedCache {
	return &SharedCache{data: make(map[string]string)}
}

// === new() ===

func useNew() *int {
	p := new(int)  // allocate int, zero-initialized, returns *int
	*p = 42
	return p       // safe to return — escapes to heap
}

// === Stack vs Heap (escape analysis) ===
// Go compiler decides stack vs heap automatically
// Use: go build -gcflags="-m" . to inspect

func stackAlloc() int {
	x := 42    // ← Go may allocate on stack
	return x   // copy value back
}

func heapAlloc() *int {
	x := 42     // ← x MUST escape to heap because we return a pointer
	return &x   // Go guarantees x lives long enough
}

// Returning pointer to local variable — SAFE in Go (unlike C!)
func newCounter(initial int) *Counter {
	c := Counter{value: initial} // c escapes to heap
	return &c                    // safe — GC manages lifetime
}

func main() {
	fmt.Println("=== 1. Pointer Cơ Bản ===")

	i := 10
	fmt.Printf("i = %d, &i = %p\n", i, &i)

	zeroval(i)
	fmt.Printf("sau zeroval: i = %d (không đổi)\n", i)

	zeroptr(&i)
	fmt.Printf("sau zeroptr: i = %d (đã đổi)\n", i)

	// Create a pointer with &
	p := &i
	fmt.Printf("p = %p, *p = %d\n", p, *p)
	*p = 100
	fmt.Printf("sau *p=100: i = %d\n", i)

	fmt.Println("\n=== 2. Nil Pointer ===")
	var nilPtr *int
	fmt.Printf("nil pointer: %v\n", nilPtr)
	fmt.Printf("nil pointer == nil: %t\n", nilPtr == nil)
	// *nilPtr = 1  // ← PANIC: nil pointer dereference
	// Always check nil before dereferencing
	if nilPtr != nil {
		fmt.Println("*nilPtr =", *nilPtr)
	} else {
		fmt.Println("nilPtr is nil — không dereference")
	}

	fmt.Println("\n=== 3. Pointer Receiver ===")
	c := Counter{}
	c.Increment()
	c.Increment()
	c.Add(8)
	fmt.Printf("Counter: %d\n", c.Get())

	// Auto-deref: Go automatically dereferences pointer when calling methods
	cp := &Counter{value: 5}
	cp.Increment() // Go automatically does (*cp).Increment()
	fmt.Printf("Counter via pointer: %d\n", cp.Get())

	fmt.Println("\n=== 4. new() ===")
	ptr := useNew()
	fmt.Printf("new(int) = %d\n", *ptr)

	// new vs make:
	// new(T): allocate, zero-init, return *T — works for any type
	// make(T, ...): only for slice, map, channel — returns T (not pointer)
	s := make([]int, 5)
	m := make(map[string]int)
	fmt.Printf("make slice: %v\n", s)
	fmt.Printf("make map: %v\n", m)

	fmt.Println("\n=== 5. Return Pointer to Local Variable (safe!) ===")
	counter := newCounter(10)
	counter.Increment()
	fmt.Printf("counter from newCounter: %d\n", counter.Get())

	n := heapAlloc()
	fmt.Printf("heapAlloc: %d (lives on heap, not stack)\n", *n)

	sv := stackAlloc()
	fmt.Printf("stackAlloc: %d (value copy)\n", sv)

	fmt.Println("\n=== 6. Khi Nào Dùng Pointer ===")
	fmt.Println("  ✓ Cần mutate: pointer receiver, *T parameter")
	fmt.Println("  ✓ Struct lớn: tránh copy overhead")
	fmt.Println("  ✓ Optional/nullable: *T có thể nil")
	fmt.Println("  ✓ Shared state: nhiều nơi cùng modify")
	fmt.Println()
	fmt.Println("  ✗ Không cần khi: primitive types nhỏ (int, bool)")
	fmt.Println("  ✗ Không cần khi: immutable data")
	fmt.Println()
	fmt.Println("  Go KHÔNG có pointer arithmetic (khác C!)")
	fmt.Println("  → GC tự quản lý lifetime, không cần free manually")
	fmt.Println()
	fmt.Println("  Xem escape analysis: go build -gcflags=\"-m\" .")
	fmt.Println("  \"escapes to heap\" = allocate on heap (GC-managed)")
	fmt.Println("  \"does not escape\" = allocate on stack (cheaper)")
}
