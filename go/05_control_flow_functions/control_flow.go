package main

import (
	"fmt"
	"time"
)

func demoControlFlow() {
	fmt.Println("\n--- For Loops ---")

	// Cách 1: C-style for loop
	for i := 0; i < 3; i++ {
		fmt.Printf("  C-style: %d\n", i)
	}

	// Cách 2: while-style (chỉ có condition)
	n := 0
	for n < 3 {
		fmt.Printf("  while-style: %d\n", n)
		n++
	}

	// Cách 3: range over integer (Go 1.22+)
	for i := range 3 {
		fmt.Printf("  range int: %d\n", i)
	}

	// Cách 4: range over slice
	nums := []int{10, 20, 30}
	for i, v := range nums {
		fmt.Printf("  range slice: index=%d, value=%d\n", i, v)
	}

	// Cách 5: range chỉ lấy index
	for i := range nums {
		fmt.Printf("  range index only: %d\n", i)
	}

	// Cách 6: range over map (thứ tự ngẫu nhiên)
	m := map[string]int{"a": 1, "b": 2}
	for k, v := range m {
		fmt.Printf("  range map: %s=%d\n", k, v)
	}

	// Cách 7: range over string (trả về rune)
	for i, r := range "Go🎯" {
		fmt.Printf("  range string: index=%d, rune=%c\n", i, r)
	}

	// Cách 8: infinite loop với break
	count := 0
	for {
		if count >= 3 {
			break
		}
		fmt.Printf("  infinite+break: %d\n", count)
		count++
	}

	// continue
	fmt.Println("  continue (skip chẵn):")
	for i := range 6 {
		if i%2 == 0 {
			continue
		}
		fmt.Printf("    %d\n", i)
	}

	// Labeled break/continue — break/continue outer loop từ inner loop
	fmt.Println("  Labeled break:")
outer:
	for i := range 3 {
		for j := range 3 {
			if i+j >= 3 {
				break outer // break cả 2 vòng lặp
			}
			fmt.Printf("    i=%d, j=%d\n", i, j)
		}
	}

	fmt.Println("\n--- Switch ---")

	// Expression switch
	day := time.Thursday
	switch day {
	case time.Saturday, time.Sunday:
		fmt.Println("  Weekend!")
	case time.Monday:
		fmt.Println("  Thứ Hai blues...")
	default:
		fmt.Printf("  Ngày thường: %v\n", day)
	}

	// Switch không cần expression (như if-else chain)
	score := 85
	switch {
	case score >= 90:
		fmt.Println("  Grade: A")
	case score >= 80:
		fmt.Println("  Grade: B")
	case score >= 70:
		fmt.Println("  Grade: C")
	default:
		fmt.Println("  Grade: F")
	}

	// Switch với initialization statement
	switch x := 42; {
	case x > 100:
		fmt.Println("  > 100")
	case x > 50:
		fmt.Println("  > 50")
	default:
		fmt.Printf("  x = %d\n", x)
	}

	// Type switch
	fmt.Println("  Type switch:")
	values := []any{42, "hello", 3.14, true, []int{1, 2, 3}}
	for _, v := range values {
		switch t := v.(type) {
		case int:
			fmt.Printf("    int: %d\n", t)
		case string:
			fmt.Printf("    string: %q\n", t)
		case float64:
			fmt.Printf("    float64: %.2f\n", t)
		case bool:
			fmt.Printf("    bool: %t\n", t)
		default:
			fmt.Printf("    unknown type: %T\n", t)
		}
	}

	// fallthrough — hiếm dùng, thực thi case tiếp theo
	fmt.Println("  fallthrough:")
	switch 1 {
	case 1:
		fmt.Println("    case 1")
		fallthrough
	case 2:
		fmt.Println("    case 2 (do fallthrough)")
	case 3:
		fmt.Println("    case 3 (không chạy)")
	}

	fmt.Println("\n--- if với initialization ---")
	// if statement có thể có initialization trước condition
	if err := riskyOp(); err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Println("  OK")
	}

	fmt.Println("\n--- defer ---")
	// defer: trì hoãn thực thi đến khi function return
	// LIFO order: defer cuối cùng chạy trước
	demonstrateDefer()
}

func riskyOp() error {
	return nil // thành công
}

func demonstrateDefer() {
	fmt.Println("  defer LIFO order:")
	defer fmt.Println("    defer 1 (chạy cuối cùng)")
	defer fmt.Println("    defer 2")
	defer fmt.Println("    defer 3 (chạy đầu tiên)")
	fmt.Println("    main function body")

	// GOTCHA: Argument của defer được evaluate ngay khi defer được gọi
	x := 10
	defer fmt.Printf("    defer value of x: %d (evaluated tại thời điểm defer)\n", x)
	x = 20
	fmt.Printf("    current x: %d\n", x)

	// Use case phổ biến: cleanup resources
	// f, _ := os.Open("file.txt")
	// defer f.Close()  // ← đảm bảo luôn close dù function return ở đâu

	// GOTCHA: defer trong loop tích lũy, không nên dùng trong loop lớn
	// for rows.Next() {
	//     row := rows.Scan(...)
	//     defer row.Close()  // ← TẤT CẢ defer chạy khi function return, không phải end of loop!
	// }
}
