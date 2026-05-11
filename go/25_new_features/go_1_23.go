package main

import (
	"fmt"
	"iter"
	"maps"
	"slices"
)

func demoGo123() {
	fmt.Println("\n--- 1. iter.Seq & iter.Seq2 (Go 1.23+) ---")
	// iter.Seq[V] = func(yield func(V) bool)
	// iter.Seq2[K, V] = func(yield func(K, V) bool)

	// Custom iterator
	fibonacci := func(n int) iter.Seq[int] {
		return func(yield func(int) bool) {
			a, b := 0, 1
			for range n {
				if !yield(a) {
					return
				}
				a, b = b, a+b
			}
		}
	}

	fmt.Print("  Fibonacci(8): ")
	for v := range fibonacci(8) {
		fmt.Printf("%d ", v)
	}
	fmt.Println()

	// Enumerate với iter.Seq2
	words := []string{"go", "is", "great"}
	withIndex := func(s []string) iter.Seq2[int, string] {
		return func(yield func(int, string) bool) {
			for i, v := range s {
				if !yield(i, v) {
					return
				}
			}
		}
	}

	fmt.Println("  Enumerate:")
	for i, w := range withIndex(words) {
		fmt.Printf("    [%d] %s\n", i, w)
	}

	// slices.All, slices.Values, slices.Backward (Go 1.23+)
	nums := []int{10, 20, 30}
	fmt.Println("  slices.All:")
	for i, v := range slices.All(nums) {
		fmt.Printf("    [%d]=%d\n", i, v)
	}
	fmt.Println("  slices.Values:")
	for v := range slices.Values(nums) {
		fmt.Printf("    %d\n", v)
	}
	fmt.Println("  slices.Backward:")
	for i, v := range slices.Backward(nums) {
		fmt.Printf("    [%d]=%d\n", i, v)
	}

	// maps.All, maps.Keys, maps.Values (Go 1.23+)
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	fmt.Println("  maps.Keys:")
	for k := range maps.Keys(m) {
		fmt.Printf("    %s\n", k)
	}

	fmt.Println("\n--- 2. Timer fix (Go 1.23+) ---")
	fmt.Println("  time.NewTimer: sending goroutine no longer blocked")
	fmt.Println("  Before: Reset/Stop could race with drain")
	fmt.Println("  Now: safe to call Stop()/Reset() without draining")
	fmt.Println("  Old pattern needed:")
	fmt.Println("    if !t.Stop() { <-t.C }")
	fmt.Println("  New (Go 1.23+): just call t.Reset(d) directly")
}
