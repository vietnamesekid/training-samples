package main

import (
	"strings"
	"testing"
)

// === Benchmarks ===
// Run: go test -bench=. -benchmem -benchtime=3s ./...

// BenchmarkReverse measures the performance of the Reverse function
func BenchmarkReverse(b *testing.B) {
	input := "Hello, World! 🌍"
	b.ResetTimer() // ignore setup time
	for range b.N {
		Reverse(input)
	}
}

// Compare string concat with += vs strings.Builder
func BenchmarkStringConcat_Plus(b *testing.B) {
	for range b.N {
		s := ""
		for range 100 {
			s += "x" // O(n²): creates a new string each time
		}
		_ = s
	}
}

func BenchmarkStringConcat_Builder(b *testing.B) {
	for range b.N {
		var sb strings.Builder
		sb.Grow(100)
		for range 100 {
			sb.WriteByte('x') // O(n): amortized
		}
		_ = sb.String()
	}
}

// Go 1.24+: b.Loop() — runs the loop exactly N times
// Ensures the benchmark is not optimized away by the compiler
func BenchmarkReverse_Loop(b *testing.B) {
	input := "Hello, World! 🌍"
	for b.Loop() {
		Reverse(input)
	}
}

// BenchmarkCalculator with b.ReportAllocs()
func BenchmarkCalculator_Add(b *testing.B) {
	c := NewCalculator()
	b.ReportAllocs() // print allocations per iteration
	b.ResetTimer()
	for range b.N {
		c.Add(1, 2)
	}
}

// Sub-benchmark — compare different implementations
func BenchmarkSliceAppend(b *testing.B) {
	b.Run("no_prealloc", func(b *testing.B) {
		for range b.N {
			var s []int
			for i := range 1000 {
				s = append(s, i) // many reallocations
			}
			_ = s
		}
	})

	b.Run("prealloc", func(b *testing.B) {
		for range b.N {
			s := make([]int, 0, 1000) // allocate once
			for i := range 1000 {
				s = append(s, i)
			}
			_ = s
		}
	})
}
