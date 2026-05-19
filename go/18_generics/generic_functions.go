package main

import (
	"cmp"
	"fmt"
)

// === Constraints ===

// Number: all numeric types
// ~int: all types whose underlying type is int (including type MyInt int)
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Ordered: types that can be compared with <, >, ==
// cmp.Ordered from Go 1.21 (replaces constraints.Ordered)
type Ordered = cmp.Ordered

// === Generic Functions ===

// Map: transform a slice from type T to type U
func Map[T, U any](slice []T, f func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = f(v)
	}
	return result
}

// Filter: filter a slice by predicate
func Filter[T any](slice []T, predicate func(T) bool) []T {
	var result []T
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Reduce: fold a slice into a single value
func Reduce[T, U any](slice []T, initial U, f func(U, T) U) U {
	acc := initial
	for _, v := range slice {
		acc = f(acc, v)
	}
	return acc
}

// Sum: sum of a numeric slice
func Sum[T Number](slice []T) T {
	var total T
	for _, v := range slice {
		total += v
	}
	return total
}

// Max: find the max in a slice
func Max[T Ordered](slice []T) T {
	if len(slice) == 0 {
		var zero T
		return zero
	}
	max := slice[0]
	for _, v := range slice[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

// Contains: check if an element is in a slice
func Contains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// Keys: get keys from a map
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Ptr: return a pointer to a value (utility function)
func Ptr[T any](v T) *T {
	return &v
}

func demoGenericFunctions() {
	fmt.Println("\n--- Map ---")
	nums := []int{1, 2, 3, 4, 5}
	doubled := Map(nums, func(n int) int { return n * 2 })
	fmt.Printf("  Map(double): %v\n", doubled)

	strs := Map(nums, func(n int) string { return fmt.Sprintf("item%d", n) })
	fmt.Printf("  Map(to string): %v\n", strs)

	fmt.Println("\n--- Filter ---")
	evens := Filter(nums, func(n int) bool { return n%2 == 0 })
	fmt.Printf("  Filter(even): %v\n", evens)

	words := []string{"go", "python", "rust", "java", "kotlin"}
	short := Filter(words, func(s string) bool { return len(s) <= 3 })
	fmt.Printf("  Filter(len<=3): %v\n", short)

	fmt.Println("\n--- Reduce ---")
	sum := Reduce(nums, 0, func(acc, n int) int { return acc + n })
	fmt.Printf("  Reduce(sum): %d\n", sum)

	product := Reduce(nums, 1, func(acc, n int) int { return acc * n })
	fmt.Printf("  Reduce(product): %d\n", product)

	fmt.Println("\n--- Sum & Max ---")
	fmt.Printf("  Sum(ints): %d\n", Sum(nums))
	fmt.Printf("  Sum(floats): %f\n", Sum([]float64{1.1, 2.2, 3.3}))
	fmt.Printf("  Max(ints): %d\n", Max(nums))
	fmt.Printf("  Max(strings): %s\n", Max([]string{"banana", "apple", "cherry"}))

	fmt.Println("\n--- Contains ---")
	fmt.Printf("  Contains([1,2,3], 2): %t\n", Contains(nums[:3], 2))
	fmt.Printf("  Contains([1,2,3], 5): %t\n", Contains(nums[:3], 5))
	fmt.Printf("  Contains(strings, \"go\"): %t\n", Contains(words, "go"))

	fmt.Println("\n--- Ptr (generic pointer helper) ---")
	p := Ptr(42)
	fmt.Printf("  Ptr(42) = %p, *p = %d\n", p, *p)
}
