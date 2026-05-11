package main

import (
	"fmt"
	"strings"
	"testing"
)

func demoGo124() {
	fmt.Println("\n--- 1. Generic Type Aliases (Go 1.24+) ---")
	// Go 1.24 hoàn thiện generic type aliases
	// type MySlice[T any] = []T  (= là alias, không phải new type)

	type Predicate[T any] = func(T) bool

	isEven := Predicate[int](func(n int) bool { return n%2 == 0 })
	isLong := Predicate[string](func(s string) bool { return len(s) > 5 })

	fmt.Printf("  isEven(4): %t\n", isEven(4))
	fmt.Printf("  isLong(\"hello\"): %t\n", isLong("hello"))
	fmt.Printf("  isLong(\"hi\"): %t\n", isLong("hi"))

	fmt.Println("\n--- 2. strings.Lines & strings.SplitSeq (Go 1.24+) ---")
	text := "line1\nline2\nline3\n"
	fmt.Println("  strings.Lines:")
	for line := range strings.Lines(text) {
		fmt.Printf("    %q\n", line)
	}

	csv := "a,b,c,d,e"
	fmt.Println("  strings.SplitSeq:")
	for part := range strings.SplitSeq(csv, ",") {
		fmt.Printf("    %q\n", part)
	}

	fmt.Println("  strings.FieldsSeq:")
	for word := range strings.FieldsSeq("  foo   bar  baz  ") {
		fmt.Printf("    %q\n", word)
	}

	fmt.Println("\n--- 3. testing.B.Loop() (Go 1.24+) ---")
	// b.Loop() là replacement cho for range b.N
	// Ưu điểm: chính xác hơn, prevent compiler optimizations
	fmt.Println("  Old: for range b.N { ... }")
	fmt.Println("  New: for b.Loop() { ... }")
	fmt.Println()
	fmt.Println("  Example benchmark:")
	fmt.Println("  func BenchmarkFoo(b *testing.B) {")
	fmt.Println("      for b.Loop() {")
	fmt.Println("          result = doWork()")
	fmt.Println("      }")
	fmt.Println("  }")

	// Demo b.Loop() equivalent behavior
	b := &testing.B{}
	_ = b // chỉ để import

	fmt.Println("\n--- 4. Swiss Tables Map (Go 1.24+) ---")
	fmt.Println("  Go 1.24 rewrites map implementation with Swiss Tables")
	fmt.Println("  ~30% faster lookups for common workloads")
	fmt.Println("  Transparent: no code changes needed")
	fmt.Println("  Uses SIMD instructions when available")
}
