package main

import (
	"strings"
	"testing"
)

// === Benchmarks ===
// Chạy: go test -bench=. -benchmem -benchtime=3s ./...

// BenchmarkReverse đo hiệu suất Reverse function
func BenchmarkReverse(b *testing.B) {
	input := "Hello, World! 🌍"
	b.ResetTimer() // bỏ qua thời gian setup
	for range b.N {
		Reverse(input)
	}
}

// So sánh string concat với += vs strings.Builder
func BenchmarkStringConcat_Plus(b *testing.B) {
	for range b.N {
		s := ""
		for range 100 {
			s += "x" // O(n²): tạo string mới mỗi lần
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

// Go 1.24+: b.Loop() — chạy loop chính xác N lần
// Đảm bảo benchmark không bị compiler optimize away
func BenchmarkReverse_Loop(b *testing.B) {
	input := "Hello, World! 🌍"
	for b.Loop() {
		Reverse(input)
	}
}

// BenchmarkCalculator với b.ReportAllocs()
func BenchmarkCalculator_Add(b *testing.B) {
	c := NewCalculator()
	b.ReportAllocs() // in số allocations per iteration
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
				s = append(s, i) // nhiều reallocation
			}
			_ = s
		}
	})

	b.Run("prealloc", func(b *testing.B) {
		for range b.N {
			s := make([]int, 0, 1000) // 1 lần alloc
			for i := range 1000 {
				s = append(s, i)
			}
			_ = s
		}
	})
}
