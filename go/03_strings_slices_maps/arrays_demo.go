package main

import "fmt"

func demoArrays() {
	// Array có fixed size, là VALUE TYPE — copy khi gán hoặc truyền vào function
	arr1 := [3]int{1, 2, 3}
	arr2 := [...]int{4, 5, 6} // compiler tự đếm size từ literal

	fmt.Printf("arr1: %v (type: %T)\n", arr1, arr1)
	fmt.Printf("arr2: %v (type: %T)\n", arr2, arr2)

	// GOTCHA: gán array là COPY TOÀN BỘ
	a1 := [3]int{1, 2, 3}
	a2 := a1     // copy
	a2[0] = 99
	fmt.Printf("\nSau khi copy: a1=%v, a2=%v\n", a1, a2)
	fmt.Println("  a1[0] vẫn là 1 — a2 là bản sao độc lập")

	// Array có thể so sánh nếu element type comparable
	fmt.Printf("\n[3]int{1,2,3} == [3]int{1,2,3}: %t\n", [3]int{1, 2, 3} == [3]int{1, 2, 3})
	fmt.Printf("[3]int{1,2,3} == [3]int{1,2,4}: %t\n", [3]int{1, 2, 3} == [3]int{1, 2, 4})

	// Array 2 chiều
	matrix := [2][3]int{{1, 2, 3}, {4, 5, 6}}
	fmt.Printf("\n2D array: %v\n", matrix)

	// NGUYÊN TẮC: Trong thực tế, dùng slice thay vì array
	// Array phù hợp khi: size cố định tại compile time, cần value semantics,
	// dùng làm map key, hoặc cần stack allocation guarantee
	fmt.Println("\n→ Thực tế: dùng slice ([]int) thay vì array ([n]int)")
	fmt.Printf("  [3]int là array (value type, size cố định)\n")
	fmt.Printf("  []int là slice (reference type, dynamic size)\n")
}
