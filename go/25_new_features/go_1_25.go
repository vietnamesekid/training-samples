package main

import (
	"fmt"
	"runtime"
	"sync"
)

func demoGo125() {
	fmt.Println("\n--- 1. sync.WaitGroup.Go() (Go 1.25+) ---")
	// wg.Go() = wg.Add(1) + go func()
	// More convenient, fewer mistakes (forgetting Add(1))

	var wg sync.WaitGroup
	results := make([]int, 5)

	for i := range 5 {
		wg.Go(func() { // automatically calls Add(1)
			results[i] = i * i
		})
	}
	wg.Wait()
	fmt.Printf("  WaitGroup.Go results: %v\n", results)

	fmt.Println("\n--- 2. Container-aware GOMAXPROCS (Go 1.25+) ---")
	fmt.Printf("  GOMAXPROCS = %d (automatically set based on cgroup limits)\n",
		runtime.GOMAXPROCS(0))
	fmt.Println("  Before Go 1.25: needed github.com/uber-go/automaxprocs")
	fmt.Println("  Go 1.25+: runtime automatically reads /sys/fs/cgroup/cpu.max")

	fmt.Println("\n--- 3. Green Tea GC (Go 1.25 experiment) ---")
	fmt.Println("  Enable: GOEXPERIMENT=greenteagc go run .")
	fmt.Println("  Goal: reduce GC overhead for large heaps (>1GB)")
	fmt.Println("  Uses generational + regional GC techniques")
	fmt.Println("  Status: opt-in experiment in 1.25, not default yet")
	fmt.Printf("  Current GC: %s\n", runtime.Version())

	fmt.Println("\n--- 4. testing.T.Attr() (Go 1.25+) ---")
	fmt.Println("  t.Attr(key, value) — add structured attributes to test output")
	fmt.Println("  Useful for CI systems and test reporting tools")
	fmt.Println("  Example:")
	fmt.Println("  func TestSomething(t *testing.T) {")
	fmt.Println("      t.Attr(\"env\", \"staging\")")
	fmt.Println("      t.Attr(\"build\", \"12345\")")
	fmt.Println("      // ...")
	fmt.Println("  }")

	fmt.Println("\n--- 5. trace.NewFlightRecorder (Go 1.25+) ---")
	fmt.Println("  runtime/trace: always-on flight recorder")
	fmt.Println("  Captures recent trace data in a circular buffer")
	fmt.Println("  Low overhead (<1% CPU), writes on demand")
	fmt.Println("  Use case: capture trace data when anomaly detected")
	fmt.Println()
	fmt.Println("  fr := trace.NewFlightRecorder()")
	fmt.Println("  fr.Start()")
	fmt.Println("  // ... application code ...")
	fmt.Println("  if anomalyDetected {")
	fmt.Println("      fr.WriteTo(w) // dump recent traces")
	fmt.Println("  }")

	fmt.Println("\n--- Summary: Evolution from 1.22 → 1.25 ---")
	features := []struct {
		version string
		feature string
	}{
		{"1.22", "Loop variable per-iteration fix"},
		{"1.22", "range over integer (for i := range N)"},
		{"1.22", "Enhanced HTTP mux (method+path routing)"},
		{"1.23", "iter.Seq / iter.Seq2 iterators"},
		{"1.23", "slices.All/Values/Backward, maps.Keys/Values"},
		{"1.23", "Timer fix (no drain needed)"},
		{"1.24", "Generic type aliases complete"},
		{"1.24", "strings.Lines, strings.SplitSeq"},
		{"1.24", "testing.B.Loop()"},
		{"1.24", "Swiss Tables map implementation"},
		{"1.25", "sync.WaitGroup.Go()"},
		{"1.25", "Container-aware GOMAXPROCS"},
		{"1.25", "Green Tea GC (experiment)"},
		{"1.25", "testing.T.Attr()"},
		{"1.25", "trace.NewFlightRecorder"},
	}

	for _, f := range features {
		fmt.Printf("  Go %s: %s\n", f.version, f.feature)
	}
}
