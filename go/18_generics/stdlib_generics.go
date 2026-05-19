package main

import (
	"fmt"
	"iter"
	"maps"
	"slices"
)

// === Stdlib Generic Packages (Go 1.21+) ===

// Custom iterator using iter.Seq (Go 1.23+)
type Student struct {
	Name string
	GPA  float64
}

func ScholarshipStudents(students []Student) iter.Seq[Student] {
	return func(yield func(Student) bool) {
		for _, s := range students {
			if s.GPA > 8.0 {
				if !yield(s) {
					return
				}
			}
		}
	}
}

// iter.Seq2 — iterator with index
func Enumerate[T any](slice []T) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i, v := range slice {
			if !yield(i, v) {
				return
			}
		}
	}
}

// Backward iterator
func Backwards[T any](slice []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := len(slice) - 1; i >= 0; i-- {
			if !yield(slice[i]) {
				return
			}
		}
	}
}

func demoStdlibGenerics() {
	fmt.Println("\n--- slices package (Go 1.21+) ---")
	nums := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3}
	fmt.Printf("  Original: %v\n", nums)

	// Sort
	sorted := slices.Clone(nums) // deep copy
	slices.Sort(sorted)
	fmt.Printf("  Sorted: %v\n", sorted)

	// SortFunc — custom comparator
	students := []Student{
		{"Alice", 9.0},
		{"Bob", 7.5},
		{"Carol", 8.5},
	}
	slices.SortFunc(students, func(a, b Student) int {
		if a.GPA > b.GPA {
			return -1
		}
		if a.GPA < b.GPA {
			return 1
		}
		return 0
	})
	for _, s := range students {
		fmt.Printf("  Student: %s (%.1f)\n", s.Name, s.GPA)
	}

	// Search & check
	fmt.Printf("  Contains(5): %t\n", slices.Contains(nums, 5))
	fmt.Printf("  Index(5): %d\n", slices.Index(nums, 5))
	idx, found := slices.BinarySearch(sorted, 5)
	fmt.Printf("  BinarySearch(5): idx=%d, found=%t\n", idx, found)

	// Max/Min (Go 1.21+)
	fmt.Printf("  Max: %d\n", slices.Max(nums))
	fmt.Printf("  Min: %d\n", slices.Min(nums))

	// Compact: remove consecutive duplicates (needs to be sorted first)
	compact := slices.Compact(sorted)
	fmt.Printf("  Compact: %v\n", compact)

	// Reverse
	rev := slices.Clone(nums)
	slices.Reverse(rev)
	fmt.Printf("  Reversed: %v\n", rev[:5])

	fmt.Println("\n--- maps package (Go 1.21+) ---")
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	// Clone
	m2 := maps.Clone(m)
	m2["d"] = 4
	fmt.Printf("  Original len=%d, Clone len=%d\n", len(m), len(m2))

	// DeleteFunc
	maps.DeleteFunc(m2, func(k string, v int) bool { return v%2 == 0 })
	fmt.Printf("  After DeleteFunc(even): %v\n", m2)

	// Equal
	fmt.Printf("  Equal: %t\n", maps.Equal(m, maps.Clone(m)))

	// Keys/Values (Go 1.23+ — returns iter.Seq)
	fmt.Println("  Keys:")
	for k := range maps.Keys(m) {
		fmt.Printf("    %s\n", k)
	}

	fmt.Println("\n--- iter.Seq (Go 1.23+) ---")
	allStudents := []Student{
		{"Alice", 9.0},
		{"Bob", 7.5},
		{"Carol", 8.5},
		{"Dave", 9.5},
	}

	fmt.Println("  Scholarship students (GPA > 8.0):")
	for s := range ScholarshipStudents(allStudents) {
		fmt.Printf("    %s (%.1f)\n", s.Name, s.GPA)
	}

	fmt.Println("  Enumerate:")
	words := []string{"go", "is", "great"}
	for i, w := range Enumerate(words) {
		fmt.Printf("    [%d] %s\n", i, w)
	}

	fmt.Println("  Backwards:")
	for v := range Backwards(words) {
		fmt.Printf("    %s\n", v)
	}
}

func demoWhenNotToUse() {
	fmt.Println("\n--- Khi nào KHÔNG dùng generics ---")
	fmt.Println("  1. Chỉ có 1 type cụ thể — dùng concrete type")
	fmt.Println("     Bad:  func ProcessInt[T int](v T) T")
	fmt.Println("     Good: func ProcessInt(v int) int")
	fmt.Println()
	fmt.Println("  2. Behavior khác nhau theo type — dùng interface")
	fmt.Println("     Bad:  func Log[T any](v T) — bạn cần fmt.Stringer")
	fmt.Println("     Good: func Log(v fmt.Stringer)")
	fmt.Println()
	fmt.Println("  3. Code mới, chưa rõ pattern — YAGNI")
	fmt.Println("     → Viết concrete types trước, generic hóa khi thực sự cần")
	fmt.Println()
	fmt.Println("  4. Simple functions không benefit from generics:")
	fmt.Println("     reverse := func(s string) string { ... }  ← đủ rồi")
	fmt.Println()
	fmt.Println("  NGUYÊN TẮC: \"If in doubt, leave it out\"")
	fmt.Println("  Generics tốt cho: container types, algorithms trên slices/maps")
}
