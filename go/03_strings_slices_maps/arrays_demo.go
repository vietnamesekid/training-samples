package main

import "fmt"

func demoArrays() {
	// Array has a fixed size and is a VALUE TYPE — copied when assigned or passed to a function
	arr1 := [3]int{1, 2, 3}
	arr2 := [...]int{4, 5, 6} // compiler counts the size from the literal

	fmt.Printf("arr1: %v (type: %T)\n", arr1, arr1)
	fmt.Printf("arr2: %v (type: %T)\n", arr2, arr2)

	// GOTCHA: assigning an array is a FULL COPY
	a1 := [3]int{1, 2, 3}
	a2 := a1     // copy
	a2[0] = 99
	fmt.Printf("\nSau khi copy: a1=%v, a2=%v\n", a1, a2)
	fmt.Println("  a1[0] vẫn là 1 — a2 là bản sao độc lập")

	// Arrays can be compared if the element type is comparable
	fmt.Printf("\n[3]int{1,2,3} == [3]int{1,2,3}: %t\n", [3]int{1, 2, 3} == [3]int{1, 2, 3})
	fmt.Printf("[3]int{1,2,3} == [3]int{1,2,4}: %t\n", [3]int{1, 2, 3} == [3]int{1, 2, 4})

	// 2D array
	matrix := [2][3]int{{1, 2, 3}, {4, 5, 6}}
	fmt.Printf("\n2D array: %v\n", matrix)

	// PRINCIPLE: In practice, use slices instead of arrays
	// Arrays are suitable when: the size is fixed at compile time, value semantics are needed,
	// used as a map key, or stack allocation is required
	fmt.Println("\n→ Thực tế: dùng slice ([]int) thay vì array ([n]int)")
	fmt.Printf("  [3]int là array (value type, size cố định)\n")
	fmt.Printf("  []int là slice (reference type, dynamic size)\n")
}
