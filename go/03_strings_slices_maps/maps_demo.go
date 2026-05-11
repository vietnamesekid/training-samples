package main

import (
	"fmt"
	"maps"
	"sort"
)

func demoMaps() {
	// Map là hash table — KHÔNG đảm bảo thứ tự khi iterate
	m := map[string]int{
		"alice": 30,
		"bob":   25,
		"carol": 35,
	}
	fmt.Printf("map: %v\n", m)

	// Tạo với make — dùng khi biết trước số lượng để giảm reallocation
	m2 := make(map[string]int, 100) // hint: ~100 entries
	m2["key"] = 100
	fmt.Printf("make map: len=%d\n", len(m2))

	// CRUD
	m["dave"] = 28                // Create/Update
	fmt.Printf("after add dave: len=%d\n", len(m))
	delete(m, "bob")              // Delete
	fmt.Printf("after delete bob: len=%d\n", len(m))

	// NGUYÊN TẮC: Luôn dùng 2-value form khi check key tồn tại
	val, ok := m["alice"]
	fmt.Printf("\nalice: val=%d, ok=%t\n", val, ok)

	val2, ok2 := m["nonexistent"]
	fmt.Printf("nonexistent: val=%d, ok=%t ← val là zero value!\n", val2, ok2)

	// GOTCHA: Đọc key không tồn tại trả về zero value, KHÔNG panic
	// Nhưng ghi vào nil map THÌ panic!
	// var nilMap map[string]int
	// nilMap["key"] = 1 // ← PANIC: assignment to entry in nil map

	// Iterate — thứ tự RANDOM mỗi lần chạy
	fmt.Println("\nIterate (thứ tự ngẫu nhiên):")
	for k, v := range m {
		fmt.Printf("  %s: %d\n", k, v)
	}

	// Iterate theo thứ tự: sort keys trước
	fmt.Println("\nIterate theo thứ tự alphabet:")
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("  %s: %d\n", k, m[k])
	}

	// GOTCHA: Không thể set field trực tiếp trong map[string]Struct
	type Point struct{ X, Y int }
	points := map[string]Point{"a": {1, 2}}
	// points["a"].X = 10  // ERROR: cannot assign to struct field in map

	// Fix 1: copy, modify, assign back
	p := points["a"]
	p.X = 10
	points["a"] = p
	fmt.Printf("\nFix struct field: %v\n", points["a"])

	// Fix 2: dùng pointer *Point
	points2 := map[string]*Point{"b": {3, 4}}
	points2["b"].X = 30  // OK — modify qua pointer
	fmt.Printf("Pointer fix: %v\n", points2["b"])

	// maps package (Go 1.21+)
	fmt.Println("\nmaps package (Go 1.21+):")
	original := map[string]int{"a": 1, "b": 2, "c": 3}
	cloned := maps.Clone(original) // deep copy
	cloned["d"] = 4
	fmt.Printf("  original len=%d, clone len=%d\n", len(original), len(cloned))

	// maps.Keys và maps.Values (Go 1.23+ via iter.Seq)
	fmt.Println("  Keys:")
	for k := range maps.Keys(original) {
		fmt.Printf("    %s\n", k)
	}

	// maps.Equal
	m3 := map[string]int{"x": 1, "y": 2}
	m4 := map[string]int{"x": 1, "y": 2}
	m5 := map[string]int{"x": 1, "y": 3}
	fmt.Printf("  Equal(m3,m4)=%t, Equal(m3,m5)=%t\n", maps.Equal(m3, m4), maps.Equal(m3, m5))

	// Set pattern — dùng map[T]struct{} (tốn ít bộ nhớ nhất)
	fmt.Println("\nSet pattern (map[T]struct{}):")
	set := map[string]struct{}{
		"apple":  {},
		"banana": {},
		"cherry": {},
	}
	set["date"] = struct{}{}
	_, exists := set["apple"]
	fmt.Printf("  apple exists: %t\n", exists)
	delete(set, "banana")
	fmt.Printf("  after delete banana, len=%d\n", len(set))

	// Counting/frequency map
	fmt.Println("\nFrequency map:")
	words := []string{"go", "is", "great", "go", "is", "awesome", "go"}
	freq := make(map[string]int)
	for _, w := range words {
		freq[w]++ // zero value của int là 0, nên không cần check tồn tại
	}
	fmt.Printf("  %v\n", freq)
}
