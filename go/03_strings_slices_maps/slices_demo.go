package main

import (
	"fmt"
	"slices"
)

func demoSlices() {
	// Slice internal: struct{ ptr *T, len int, cap int }
	// Là REFERENCE TYPE — share underlying array!

	// Tạo slice
	s1 := []int{1, 2, 3}               // slice literal
	s2 := make([]int, 5)               // len=5, cap=5, zero values
	s3 := make([]int, 0, 10)           // len=0, cap=10 — pre-allocate tránh realloc
	var s4 []int                        // nil slice (len=0, cap=0, ptr=nil)

	fmt.Printf("s1=%v len=%d cap=%d\n", s1, len(s1), cap(s1))
	fmt.Printf("s2=%v len=%d cap=%d\n", s2, len(s2), cap(s2))
	fmt.Printf("s3=%v len=%d cap=%d\n", s3, len(s3), cap(s3))
	fmt.Printf("s4=%v len=%d cap=%d nil=%t\n", s4, len(s4), cap(s4), s4 == nil)

	// append — có thể tạo underlying array mới nếu cap không đủ
	fmt.Println("\nappend:")
	s := []int{1, 2, 3}
	fmt.Printf("  before: %v cap=%d\n", s, cap(s))
	s = append(s, 4)
	fmt.Printf("  after append(4): %v cap=%d\n", s, cap(s))
	s = append(s, 5, 6, 7) // append nhiều phần tử cùng lúc
	fmt.Printf("  after append(5,6,7): %v cap=%d\n", s, cap(s))

	// Append slice vào slice với ...
	a := []int{1, 2, 3}
	b := []int{4, 5, 6}
	c := append(a, b...) // unpack b
	fmt.Printf("  append(a, b...): %v\n", c)

	// GOTCHA: slicing SHARE underlying array
	fmt.Println("\nSlicing — SHARE underlying array:")
	original := []int{1, 2, 3, 4, 5}
	shared := original[1:3]  // [2, 3], share memory với original
	fmt.Printf("  original=%v\n", original)
	fmt.Printf("  shared=original[1:3]=%v\n", shared)
	shared[0] = 99           // NGUY HIỂM: thay đổi original!
	fmt.Printf("  sau shared[0]=99, original=%v ← original bị thay đổi!\n", original)

	// Fix 1: dùng full slice expression để giới hạn capacity
	original2 := []int{1, 2, 3, 4, 5}
	safe := original2[1:3:3] // cap bị giới hạn ở index 3
	safe = append(safe, 99)   // append TẠO array mới (không ảnh hưởng original)
	fmt.Printf("\n  original2=%v (không bị ảnh hưởng)\n", original2)
	fmt.Printf("  safe (sau append)=%v\n", safe)

	// Fix 2: copy — tạo deep copy
	dst := make([]int, len(original))
	n := copy(dst, original)
	fmt.Printf("  copy: dst=%v, %d elements copied\n", dst, n)

	// Xóa element tại index i
	fmt.Println("\nXóa phần tử:")
	arr := []int{1, 2, 3, 4, 5}
	i := 2 // xóa index 2

	// Cách 1: không giữ thứ tự (hoán đổi với phần tử cuối — O(1))
	arr[i] = arr[len(arr)-1]
	arr = arr[:len(arr)-1]
	fmt.Printf("  Không giữ thứ tự: %v\n", arr)

	// Cách 2: slices.Delete (Go 1.21+) — giữ thứ tự O(n)
	arr2 := []int{1, 2, 3, 4, 5}
	arr2 = slices.Delete(arr2, 2, 3) // xóa index 2 đến 3 (exclusive)
	fmt.Printf("  slices.Delete: %v\n", arr2)

	// slices package (Go 1.21+) — các hàm tiện ích
	fmt.Println("\nslices package (Go 1.21+):")
	nums := []int{3, 1, 4, 1, 5, 9, 2, 6}
	fmt.Printf("  Before sort: %v\n", nums)
	slices.Sort(nums)
	fmt.Printf("  After Sort: %v\n", nums)
	fmt.Printf("  Contains(5): %t\n", slices.Contains(nums, 5))
	idx, found := slices.BinarySearch(nums, 5)
	fmt.Printf("  BinarySearch(5): idx=%d, found=%t\n", idx, found)
	slices.Reverse(nums)
	fmt.Printf("  After Reverse: %v\n", nums)

	// 2D slice
	fmt.Println("\n2D slice:")
	rows, cols := 3, 4
	matrix := make([][]int, rows)
	for i := range matrix {
		matrix[i] = make([]int, cols)
		for j := range matrix[i] {
			matrix[i][j] = i*cols + j
		}
	}
	for _, row := range matrix {
		fmt.Printf("  %v\n", row)
	}

	// NGUYÊN TẮC: Pre-allocate slice khi biết trước kích thước
	fmt.Println("\nPre-allocate tránh reallocation:")
	n2 := 1000
	good := make([]int, 0, n2) // ← tốt: chỉ 1 lần allocate
	for i := range n2 {
		good = append(good, i)
	}
	fmt.Printf("  cap=%d (không realloc)\n", cap(good))
}
