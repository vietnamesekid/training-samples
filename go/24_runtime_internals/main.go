// Lesson 24: Go Runtime Internals — GMP Model & Memory Management
// Run: go run .
// View GMP scheduler trace: GODEBUG=schedtrace=1000 go run .
// View GC trace: GODEBUG=gctrace=1 go run .
package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
	"time"
)

// === GMP Model ===
// G = Goroutine (user-level thread, ~2KB initial stack)
// M = Machine (OS thread, managed by runtime)
// P = Processor (logical CPU, controls goroutine scheduling)
//
// Number of P = GOMAXPROCS (default = number of CPU cores)
// Each P has a local run queue (LRQ) holding goroutines
// Work stealing: idle P takes goroutines from other Ps
//
// Goroutine states:
//   Running: currently executing on M
//   Runnable: ready, waiting for P
//   Waiting: blocked (channel, syscall, sleep)

func printRuntimeInfo() {
	fmt.Printf("  NumCPU: %d\n", runtime.NumCPU())
	fmt.Printf("  GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))
	fmt.Printf("  NumGoroutine: %d\n", runtime.NumGoroutine())
	fmt.Printf("  GOARCH: %s\n", runtime.GOARCH)
	fmt.Printf("  GOOS: %s\n", runtime.GOOS)
	fmt.Printf("  Version: %s\n", runtime.Version())
}

// === GC Deep Dive ===

func printGCStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("  Heap:\n")
	fmt.Printf("    HeapAlloc:    %d KB (in-use)\n", m.HeapAlloc/1024)
	fmt.Printf("    HeapSys:      %d KB (from OS)\n", m.HeapSys/1024)
	fmt.Printf("    HeapIdle:     %d KB (idle spans)\n", m.HeapIdle/1024)
	fmt.Printf("    HeapInuse:    %d KB (in-use spans)\n", m.HeapInuse/1024)
	fmt.Printf("    HeapReleased: %d KB (returned to OS)\n", m.HeapReleased/1024)
	fmt.Printf("    HeapObjects:  %d\n", m.HeapObjects)

	fmt.Printf("  GC:\n")
	fmt.Printf("    NumGC: %d\n", m.NumGC)
	fmt.Printf("    GCCPUFraction: %.4f (%.2f%% of CPU)\n", m.GCCPUFraction, m.GCCPUFraction*100)
	if m.NumGC > 0 {
		fmt.Printf("    LastGC: %v ago\n", time.Since(time.Unix(0, int64(m.LastGC))))
	}
	fmt.Printf("    PauseTotalNs: %d ms\n", m.PauseTotalNs/1_000_000)

	fmt.Printf("  Stack:\n")
	fmt.Printf("    StackInuse: %d KB\n", m.StackInuse/1024)
	fmt.Printf("    StackSys:   %d KB\n", m.StackSys/1024)
}

func demoGCObservation() {
	fmt.Println("\n--- Before allocation ---")
	printGCStats()

	// Allocate a lot of objects
	var data [][]byte
	for range 1000 {
		data = append(data, make([]byte, 1024)) // 1KB each = ~1MB total
	}

	fmt.Printf("\n--- After allocating %d KB ---\n", len(data))
	printGCStats()

	// Let GC collect
	data = nil
	runtime.GC()

	fmt.Println("\n--- After GC ---")
	printGCStats()
}

func main() {
	fmt.Println("=== 1. GMP Model ===")
	fmt.Println("  G = Goroutine (user-level thread, ~2KB initial stack, grows on demand)")
	fmt.Println("  M = Machine (OS thread, expensive — runtime manages pool)")
	fmt.Println("  P = Processor (logical CPU context, controls scheduling)")
	fmt.Println()
	fmt.Println("  Scheduling rules:")
	fmt.Println("  - Goroutine preemption: Go 1.14+ = asynchronous preemption")
	fmt.Println("  - Work stealing: idle P takes goroutines from busy P")
	fmt.Println("  - Goroutine park/unpark: channel ops, syscalls")

	fmt.Println("\n--- Runtime info ---")
	printRuntimeInfo()

	fmt.Println("\n--- Change GOMAXPROCS ---")
	old := runtime.GOMAXPROCS(2)
	fmt.Printf("  Set GOMAXPROCS=2 (was %d)\n", old)
	fmt.Printf("  Now: %d\n", runtime.GOMAXPROCS(0))
	runtime.GOMAXPROCS(old) // restore

	fmt.Println("\n=== 2. Goroutine Stacks ===")
	fmt.Println("  Initial stack: 2KB-8KB (varies by version)")
	fmt.Println("  Grows dynamically to 1GB (64-bit), 250MB (32-bit)")
	fmt.Println("  Shrinks after GC scan")

	// Print goroutine stacks
	fmt.Println("\n--- Current goroutines ---")
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false) // false = only current goroutine
	fmt.Printf("  %s\n", buf[:n])

	fmt.Println("\n=== 3. GC — Tri-color Mark-and-Sweep ===")
	fmt.Println("  Go GC phases:")
	fmt.Println("  1. Mark Setup (STW): enable write barriers")
	fmt.Println("  2. Marking (concurrent): mark reachable objects")
	fmt.Println("  3. Mark Termination (STW): finalize marking")
	fmt.Println("  4. Sweeping (concurrent): reclaim unmarked memory")
	fmt.Println()
	fmt.Println("  Write barrier: ensures concurrent marking correctness")
	fmt.Println("  Tri-color: White=unreachable, Grey=reachable, Black=scanned")
	fmt.Println()
	fmt.Println("  Go 1.25+: Green Tea GC experiment")
	fmt.Println("  Enable: GOEXPERIMENT=greenteagc go run .")
	fmt.Println("  Goal: reduce GC overhead for large heaps")

	fmt.Println("\n=== 4. GC Tuning ===")
	fmt.Printf("  GOGC (default=%d): GC runs when heap grows GOGC%% above previous\n",
		readGOGC())
	fmt.Println("  GOGC=off: disable GC (batch jobs)")
	fmt.Println("  GOGC=100 (default): GC when heap doubles")
	fmt.Println("  GOGC=200: GC when heap triples (less frequent, more memory)")

	// debug.SetGCPercent — equivalent to GOGC
	oldGOGC := debug.SetGCPercent(200)
	fmt.Printf("  Changed GOGC: %d → 200\n", oldGOGC)
	debug.SetGCPercent(oldGOGC)

	// debug.SetMemoryLimit (Go 1.19+)
	fmt.Println()
	fmt.Println("  debug.SetMemoryLimit (Go 1.19+):")
	fmt.Println("  Sets soft memory limit — GC runs more aggressively")
	fmt.Println("  Useful for containers (instead of GOGC)")
	fmt.Println("  GOMEMLIMIT=500MiB go run .  (equivalent)")

	// demo.SetMemoryLimit(500 << 20) // 500MB — commented to not affect demo

	fmt.Println("\n=== 5. GC Observation ===")
	demoGCObservation()

	fmt.Println("\n=== 6. GODEBUG ===")
	fmt.Println("  Useful GODEBUG settings:")
	fmt.Println("  GODEBUG=schedtrace=1000  : log scheduler state every 1000ms")
	fmt.Println("  GODEBUG=gctrace=1        : log GC stats each cycle")
	fmt.Println("  GODEBUG=madvdontneed=1   : return freed memory to OS faster")
	fmt.Println("  GODEBUG=gccheckmark=1    : verify GC mark phase")

	fmt.Println("\n=== 7. Container-aware GOMAXPROCS (Go 1.25+) ===")
	fmt.Println("  Go 1.25+ automatically detects CPU quotas in containers")
	fmt.Println("  Previously needed: uber-go/automaxprocs")
	fmt.Println("  Now: runtime.GOMAXPROCS respects cgroup CPU limits")
	fmt.Printf("  Current GOMAXPROCS = %d (container-aware)\n", runtime.GOMAXPROCS(0))

	fmt.Println("\n=== 8. Goroutine Count Demo ===")
	var wg sync.WaitGroup
	fmt.Printf("  Before: %d goroutines\n", runtime.NumGoroutine())

	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(50 * time.Millisecond)
		}()
	}

	fmt.Printf("  During: %d goroutines\n", runtime.NumGoroutine())
	wg.Wait()
	time.Sleep(10 * time.Millisecond) // let goroutines fully exit
	fmt.Printf("  After: %d goroutines\n", runtime.NumGoroutine())
}

func readGOGC() int {
	v := debug.SetGCPercent(-1) // read current value
	debug.SetGCPercent(v)       // restore
	return v
}
