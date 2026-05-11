package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func demoMapAtomic() {
	fmt.Println("\n--- sync.Map ---")
	// sync.Map: concurrent-safe map
	// Phù hợp khi: read nhiều hơn write, hoặc keys ít khi thay đổi
	// KHÔNG phù hợp cho high-write workloads (dùng map + Mutex tốt hơn)
	var sm sync.Map

	// Store
	sm.Store("key1", "value1")
	sm.Store("key2", "value2")
	sm.Store("key3", "value3")

	// Load
	if v, ok := sm.Load("key1"); ok {
		fmt.Printf("  Load key1: %v\n", v)
	}

	// LoadOrStore — atomic check-and-set
	actual, loaded := sm.LoadOrStore("key1", "new-value1")
	fmt.Printf("  LoadOrStore key1: actual=%v, loaded=%t\n", actual, loaded)

	actual2, loaded2 := sm.LoadOrStore("key4", "value4")
	fmt.Printf("  LoadOrStore key4 (new): actual=%v, loaded=%t\n", actual2, loaded2)

	// Delete
	sm.Delete("key2")

	// Range — iterate (thứ tự không đảm bảo)
	fmt.Println("  Range:")
	sm.Range(func(k, v any) bool {
		fmt.Printf("    %v = %v\n", k, v)
		return true // return false để dừng iteration
	})

	// LoadAndDelete — atomic load + delete
	if v, ok := sm.LoadAndDelete("key3"); ok {
		fmt.Printf("  LoadAndDelete key3: %v\n", v)
	}

	fmt.Println("\n--- sync/atomic — lock-free operations ---")
	// atomic: operations không cần lock (dùng CPU instructions)
	// Nhanh hơn Mutex cho single value operations

	// atomic.Int64 (Go 1.19+ — typed atomics)
	var counter atomic.Int64
	var wg sync.WaitGroup

	for range 1000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Add(1)
		}()
	}
	wg.Wait()
	fmt.Printf("  atomic.Int64 counter: %d\n", counter.Load())

	// atomic.Bool (Go 1.19+)
	var running atomic.Bool
	running.Store(true)
	fmt.Printf("  atomic.Bool: %t\n", running.Load())
	running.Swap(false) // trả về giá trị cũ
	fmt.Printf("  after Swap(false): %t\n", running.Load())

	// CompareAndSwap (CAS) — conditional atomic update
	var state atomic.Int32
	state.Store(1)

	// CAS: chỉ set nếu giá trị hiện tại = old
	ok := state.CompareAndSwap(1, 2) // if state == 1, set to 2
	fmt.Printf("  CAS(1→2): ok=%t, state=%d\n", ok, state.Load())

	ok2 := state.CompareAndSwap(1, 3) // if state == 1 (but it's 2 now)
	fmt.Printf("  CAS(1→3): ok=%t, state=%d\n", ok2, state.Load())

	// atomic.Value — store/load arbitrary value atomically
	var config atomic.Value
	type AppConfig struct{ Debug bool; LogLevel string }

	config.Store(AppConfig{Debug: false, LogLevel: "info"})
	if v := config.Load(); v != nil {
		cfg := v.(AppConfig)
		fmt.Printf("  atomic.Value config: %+v\n", cfg)
	}

	// Hot swap config (không cần lock)
	config.Store(AppConfig{Debug: true, LogLevel: "debug"})
	if v := config.Load(); v != nil {
		cfg := v.(AppConfig)
		fmt.Printf("  config after hot-swap: %+v\n", cfg)
	}

	fmt.Println("\n  Khi nào dùng atomic vs Mutex:")
	fmt.Println("  atomic: counter, flag, single value hot-swap")
	fmt.Println("  Mutex: complex state, multiple fields, non-atomic operations")
}
